package window

import (
	"time"
)

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
	return nil
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
