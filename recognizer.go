package window

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

const (
	STATE_LEFT_BOUND  = iota // parse left bound
	STATE_RIGHT_BOUND        // parse right bound
	STATE_VALIDATE           // parsing is over, validate the result
	STATE_FINISH             // all is good, stop the parsing
)

type Recognizer struct {
	p    *Parser
	spec Specification
}

func (r *Recognizer) state(state int) (nextState int, err error) {
	r.p.eatWs()
	switch state {
	case STATE_LEFT_BOUND: // start here
		// skip keywords
		r.p.expectAny([]string{"from", "since", "within"})
		r.p.eatWs()

		// Try 1: RelN spec
		oldPos := r.p.pos
		bound, parseErr := r.parseRelnBound()
		if parseErr == nil {
			r.spec.leftBoundRelN = &bound
			nextState = STATE_RIGHT_BOUND
			return
		}
		r.p.rollbackAt(oldPos)

		// Try 2: Rel spec
		oldPos = r.p.pos
		duration, parseErr := r.parseRelBound()
		if parseErr == nil {
			r.spec.leftBoundRel = &duration
			nextState = STATE_RIGHT_BOUND
			return
		}
		r.p.rollbackAt(oldPos)

		// Try 3: anything else should be treated as Abs spec
		oldPos = r.p.pos
		leftBoundText, _ := r.p.consumeUntil([]string{" to", "until", "within"})
		leftBoundText = strings.Trim(leftBoundText, " \n\t")
		absTime, absErr := dateparse.ParseStrict(leftBoundText)
		if absErr == nil {
			r.spec.leftBoundAbs = &absTime
			nextState = STATE_RIGHT_BOUND
			return
		}
		r.p.rollbackAt(oldPos)

		err = fmt.Errorf("failed to recognize the left bound")
		return

	case STATE_RIGHT_BOUND:
		if r.p.isEof() { // sliding window case
			nextState = STATE_VALIDATE
			return
		}

		r.p.expectAny([]string{"until", "to ", "within"})
		r.p.eatWs()

		// Try 1: RelN spec
		oldPos := r.p.pos
		bound, parseErr := r.parseRelnBound()
		if parseErr == nil {
			r.spec.rightBoundRelN = &bound
			nextState = STATE_VALIDATE
			return
		}
		r.p.rollbackAt(oldPos)

		// Try 2: Rel spec
		oldPos = r.p.pos
		duration, parseErr := r.parseRelBound()
		if parseErr == nil {
			r.spec.rightBoundRel = &duration
			nextState = STATE_VALIDATE
			return
		}
		r.p.rollbackAt(oldPos)

		// Try 3: anything else should be treated as Abs spec
		remainingText, _ := r.p.consumeUntil([]string{})
		remainingText = strings.Trim(remainingText, " \n\t")
		absTime, absErr := dateparse.ParseStrict(remainingText)
		if absErr == nil {
			r.spec.rightBoundAbs = &absTime
			nextState = STATE_VALIDATE
			return
		}

		err = fmt.Errorf("failed to recognize the right bound")
		return

	case STATE_VALIDATE:
		if !r.p.isEof() { // at this point there should be nothing left in the string
			err = r.fail("")
			return
		}
		r.spec.validate()
		nextState = STATE_FINISH
	default:
		err = fmt.Errorf("unexpected state %d", state)

	}

	return
}

func (r *Recognizer) fail(customMsg string) error {

	// Normal case:
	msg := fmt.Sprintf("unexpected character found at %d", r.p.pos)
	if customMsg != "" {
		msg += ": " + customMsg
	}
	if len(r.p.text) > 0 {
		visualPlacement := fmt.Sprintf("\n%s\n%s^", r.p.text, strings.Repeat(" ", r.p.pos))
		msg += visualPlacement
	}

	return fmt.Errorf(msg)
}

func (r *Recognizer) mapDurationUnit(unit string) (d time.Duration, err error) {
	d = MapUnitToDuration(unit)
	if d == 0 {
		err = fmt.Errorf("unsupported unit %s in relative bound", unit)
	}
	return
}

// parseRelBound check the current text and parses strings like "1 day" or "2 minutes and 3 seconds"
func (r *Recognizer) parseRelBound() (d time.Duration, err error) {
	r.p.eatWs()
	// parse num
	num := r.p.consumeRE(`\d+`)
	if num == "" {
		err = r.fail("")
		return
	}
	n, intErr := strconv.ParseInt(num, 10, 64)
	if intErr != nil {
		r.p.rollback(len(num))
		err = r.fail(intErr.Error())
		return
	}

	// parse units
	r.p.eatWs()
	unit := r.p.consumeRE(`\w+`)
	if unit == "" {
		err = r.fail("")
		return
	}
	unitDuration, durationErr := r.mapDurationUnit(unit)
	if durationErr != nil {
		r.p.rollback(len(unit))
		err = r.fail(durationErr.Error())
		return
	}
	d = unitDuration * time.Duration(n)

	// check for more "and X Y..."
	r.p.eatWs()
	if r.p.expect("and") {
		curPos := r.p.pos
		extraDuration, extraErr := r.parseRelBound()
		if extraErr != nil {
			r.p.rollbackAt(curPos)
			err = r.fail(extraErr.Error())
			return
		}
		d = d + extraDuration
	}

	return
}

// parseRelnBound checks that text contains relative specification like "next month" or an interval like "2 days ago"
func (r *Recognizer) parseRelnBound() (bound boundRelativeToNow, err error) {
	// check one-word onewords
	verbalKeyword := r.p.expectAny(GetShortWords())
	if verbalKeyword != "" {
		bound.verbal = verbalKeyword
		return
	}

	// check verbals "last X" or next Y"
	verbalPrefix := r.p.expectAny([]string{"last", "next"})
	if verbalPrefix != "" {
		r.p.eatWs()
		inFuture := false
		if verbalPrefix == "next" {
			inFuture = true
		}

		// check interval keywords
		verbal := r.p.expectAny(GetPeriodWords())
		if verbal != "" {
			bound.inFuture = inFuture
			bound.verbal = verbal
			return
		}
	}

	// check intervals "X Y ago" or "X Y after"
	duration, durationErr := r.parseRelBound()
	if durationErr == nil {
		r.p.eatWs()
		keywords := []string{"ago", "before", "after", "later", "ahead"}
		keyword := r.p.expectAny(keywords)
		if keyword == "" {
			err = r.fail("expected ago or after at this point")
			return
		}
		inFuture := false
		if keyword == "after" || keyword == "later" || keyword == "ahead" {
			inFuture = true
		}

		bound.inFuture = inFuture
		bound.duration = duration
		return
	}

	err = r.fail(durationErr.Error())
	return
}

func Start(text string) (s Specification, e error) {
	p := StartParsing(text)
	r := &Recognizer{
		p: p,
	}

	nextState := STATE_LEFT_BOUND // start here
	for {
		nextState, e = r.state(nextState)
		if nextState == STATE_FINISH {
			s = r.spec
			return
		}
		if e != nil {
			return
		}
	}
}
