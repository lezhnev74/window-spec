package window

import (
	"fmt"
	"time"
)

func MapUnitToDuration(unit string) (d time.Duration) {
	switch unit {
	case "nanosecond", "nanoseconds":
		d = time.Nanosecond
	case "microsecond", "microseconds":
		d = time.Microsecond
	case "millisecond", "milliseconds":
		d = time.Millisecond
	case "second", "seconds":
		d = time.Second
	case "minute", "minutes":
		d = time.Minute
	case "hour", "hours":
		d = time.Hour
	case "day", "days":
		d = time.Hour * 24
	case "week", "weeks":
		d = time.Hour * 24 * 7
	}
	return d
}

// GetPeriodWords returns a list of possible predefined words that can be used in the bound definition relative to now
// ex: "last X" or "next Y"
func GetPeriodWords() []string {
	return []string{
		"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday",
		"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december",
		"nanosecond", "microsecond", "millisecond", "second", "minute", "hour", "day", "week", "month", "year",
		"nanoseconds", "microseconds", "milliseconds", "seconds", "minutes", "hour", "days", "weeks", "months", "years",
	}
}

// GetShortWords returns a list of possible predefined words that can be used in the bound definition relative to now as is
func GetShortWords() []string {
	return []string{"today", "yesterday", "now", "tomorrow"}
}

// boundRelativeToNow contains a time specification relative to another point in time
// ex: "yesterday", "last june", "next week", "2 days after"
type boundRelativeToNow struct {
	inFuture bool          // direction
	verbal   string        // "june", "year", "week", "today", "yesterday"
	duration time.Duration // "2 days", "1 second"
}

// resolveAt map the relN bound to time. It uses isFuture flag to understand which bound of the interval to pick.
// -----[last year]------[NOW]------[next year]----
//      ^		  ^					^		  ^   <---- possible picks depending on isLeftBound and isFuture
func (b *boundRelativeToNow) resolveAt(n time.Time, isLeftBound bool) time.Time {
	layout := "2006-01-02 15:04:05.000000000 MST"
	locationOffset := n.Format("MST")
	var leftBoundString, rightBoundString string

	// verbal map
	if b.verbal != "" {
		switch b.verbal {
		case "now":
			return n
		case "tomorrow":
			tomorrowString := n.AddDate(0, 0, 1).Format("2006-01-02")
			leftBoundString = fmt.Sprintf("%s  00:00:00.000000000 %s", tomorrowString, locationOffset)
			rightBoundString = fmt.Sprintf("%s 23:59:59.999999999 %s", tomorrowString, locationOffset)
		default:
			panic(fmt.Errorf("vrbal [%s] not recognized", b.verbal))
		}
	}

	leftBoundTime, _ := time.Parse(layout, leftBoundString)
	rightBoundTime, _ := time.Parse(layout, rightBoundString)

	// If the period is in the left bound (from yesterday to ...) then we use the right bound of the period
	// ----[period]-----NOW---
	//            ^				<-- this bound is used if the period is met in the left window bound
	// ----NOW------[period]--
	//              ^			<-- otherwise the left bound is used
	if isLeftBound {
		return rightBoundTime
	}
	return leftBoundTime
}

// Specification contains left/right bounds for a window that can be resolved to absolute time when needed
type Specification struct {
	leftBoundAbs, rightBoundAbs   *time.Time          // "2 April 2022"
	leftBoundRel, rightBoundRel   *time.Duration      // "3 days"
	leftBoundRelN, rightBoundRelN *boundRelativeToNow // "2 days ago" or "last june"
}

func MakeSpecification(leftBound, rightBound any) Specification {
	s := Specification{}

	switch v := leftBound.(type) {
	case time.Time:
		s.leftBoundAbs = &v
	case time.Duration:
		s.leftBoundRel = &v
	case boundRelativeToNow:
		s.leftBoundRelN = &v
	}

	switch v := rightBound.(type) {
	case time.Time:
		s.rightBoundAbs = &v
	case time.Duration:
		s.rightBoundRel = &v
	case boundRelativeToNow:
		s.rightBoundRelN = &v
	}

	return s
}

// ResolveAt will generate a new Window instance
// It resolves all relative time points to absolute ones relatively to the given time point
func (s *Specification) ResolveAt(t time.Time) *Window {
	w := Window{}

	// left bound
	if s.leftBoundAbs != nil {
		w.from = s.leftBoundAbs
	} else if s.leftBoundRel != nil {
		w.slide = *s.leftBoundRel
	} else {
		rt := s.leftBoundRelN.resolveAt(t, true)
		w.from = &rt
	}

	// right bound
	if s.rightBoundAbs != nil {
		w.to = s.rightBoundAbs
	} else if s.rightBoundRel != nil {
		rt := w.from.Add(*s.rightBoundRel)
		w.to = &rt
	} else {
		rt := s.rightBoundRelN.resolveAt(t, false)
		w.to = &rt
	}

	return &w
}

func (s *Specification) validate() {
	if s.rightBoundRel != nil && s.leftBoundRel != nil {
		panic(fmt.Errorf("two rel bound are not allowed"))
	}
}

type Window struct {
	slide    time.Duration
	from, to *time.Time
}

// GetBounds return absolute times as left and right bound of the window
func (w *Window) GetBounds() (from, to time.Time) {
	return *w.from, *w.to
}

// IsSliding return true if the window has no absolute bounds, only the duration
func (w *Window) IsSliding() bool {
	return w.slide != 0
}

// GetSlide return duration of a sliding window
func (w *Window) GetSlide() time.Duration {
	return w.slide
}
