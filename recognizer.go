package window

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	STATE_LEFT_BOUND       = iota // parse left bound
	STATE_LEFT_BOUND_ABS          // expect absolute value
	STATE_LEFT_BOUND_REL          // expect value relative to right bound
	STATE_LEFT_BOUND_RELN         // expect bound relative to NOW
	STATE_RIGHT_BOUND             // parse right bound
	STATE_RIGHT_BOUND_ABS         // expect absolute value
	STATE_RIGHT_BOUND_REL         // expect value relative to left bound
	STATE_RIGHT_BOUND_RELN        // expect bound relative to NOW
	STATE_VALIDATE                // parsing is over, validate the result
	STATE_FINISH                  // all is good, stop the parsing
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
		leftBoundDelim := r.p.expectAny([]string{"from", "since", "within"})
		// find left bound text
		leftBoundText, rightBoundDelim := r.p.consumeUntil([]string{"to", "until", "within"})

		// try
		if r.p.expect("from") {
			nextState = STATE_LEFT_BOUND_ABS
			return
		}
		if r.p.expect("since") {
			nextState = STATE_LEFT_BOUND_RELN
			return
		}
		nextState = STATE_LEFT_BOUND_REL

	case STATE_LEFT_BOUND_REL:
		r.p.expect("within") // skip if present
		r.p.eatWs()
		duration, durationErr := r.parseDuration()
		if durationErr != nil {
			err = durationErr
			return
		}
		r.spec.leftBoundRel = &duration
		nextState = STATE_RIGHT_BOUND

	case STATE_LEFT_BOUND_RELN:
		relnBound, boundErr := r.parseRelnBound()
		if boundErr != nil {
			err = r.fail(boundErr.Error())
			return
		}
		r.spec.leftBoundRelN = &relnBound
		nextState = STATE_VALIDATE

	case STATE_RIGHT_BOUND:
		if r.p.isEof() {
			// case sliding-window
			nextState = STATE_VALIDATE
			return
		}
		if r.p.expect("until") {
			nextState = STATE_RIGHT_BOUND_RELN
			return
		}
		err = r.fail("") // should not get here
	case STATE_RIGHT_BOUND_RELN:
		relnBound, boundErr := r.parseRelnBound()
		if boundErr != nil {
			err = r.fail(boundErr.Error())
			return
		}
		r.spec.rightBoundRelN = &relnBound
		nextState = STATE_VALIDATE
	case STATE_VALIDATE:
		if !r.p.isEof() { // at this point there should be nothing left in the string
			err = r.fail("")
			return
		}
		nextState = STATE_FINISH
	default:
		err = fmt.Errorf("unexpected state %d", state)

	}

	return
}

func (r *Recognizer) fail(customMsg string) error {

	// Edge-case: eof found
	if r.p.isEof() {
		return fmt.Errorf("unexpected eof")
	}

	// Normal case:
	msg := fmt.Sprintf("unexpected character found at %d", r.p.pos)
	if customMsg != "" {
		msg += ":" + customMsg
	}
	if len(r.p.text) > 0 {
		visualPlacement := fmt.Sprintf("\n%s\n%s^", r.p.text, strings.Repeat(" ", r.p.pos))
		msg += visualPlacement
	}

	return fmt.Errorf(msg)
}

func (r *Recognizer) mapDurationUnit(unit string) (d time.Duration, err error) {
	switch unit {
	case "nanoseconds":
		d = time.Nanosecond
	case "microseconds":
		d = time.Microsecond
	case "milliseconds":
		d = time.Millisecond
	case "seconds":
		d = time.Second
	case "minutes":
		d = time.Minute
	case "hours":
		d = time.Hour
	case "days":
		d = time.Hour * 24
	case "weeks":
		d = time.Hour * 24 * 7
	default:
		err = fmt.Errorf("unsupported unit %s in relative bound", unit)
	}
	return
}

// parseDuration check the current text and parses strings like "1 day" or "2 minutes and 3 seconds"
func (r *Recognizer) parseDuration() (d time.Duration, err error) {
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
	duration, durationErr := r.parseDuration()
	if durationErr == nil {
		r.p.eatWs()
		keywords := []string{"ago", "after"}
		keyword := r.p.expectAny(keywords)
		if keyword == "" {
			err = r.fail("expected ago or after at this point")
			return
		}
		inFuture := false
		if keyword == "after" {
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
