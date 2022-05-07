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

// resolveAt map the relN bound to time. It uses isFuture/isLeftBound to understand which bound of the interval to pick.
// -----[last year]------[NOW]------[next year]----
//      ^		  ^					^		  ^   <---- possible picks depending on isLeftBound and isFuture
func (b *boundRelativeToNow) resolveAt(n time.Time, isLeftBound bool) time.Time {
	layout := "2006-01-02 15:04:05.000000000 MST"
	tz := n.Format("MST")
	var leftBoundString, rightBoundString string

	// verbal map
	if b.verbal != "" {
		sign := -1
		if b.inFuture {
			sign = 1
		}

		switch b.verbal {
		case "now":
			return n
		case "today":
			tomorrowString := n.Format("2006-01-02")
			leftBoundString = fmt.Sprintf("%s  00:00:00.000000000 %s", tomorrowString, tz)
			rightBoundString = fmt.Sprintf("%s 23:59:59.999999999 %s", tomorrowString, tz)
		case "tomorrow":
			tomorrowString := n.AddDate(0, 0, 1).Format("2006-01-02")
			leftBoundString = fmt.Sprintf("%s  00:00:00.000000000 %s", tomorrowString, tz)
			rightBoundString = fmt.Sprintf("%s 23:59:59.999999999 %s", tomorrowString, tz)
		case "yesterday":
			tomorrowString := n.AddDate(0, 0, -1).Format("2006-01-02")
			leftBoundString = fmt.Sprintf("%s  00:00:00.000000000 %s", tomorrowString, tz)
			rightBoundString = fmt.Sprintf("%s 23:59:59.999999999 %s", tomorrowString, tz)
		case "nanosecond", "nanoseconds":
			nanosecondString := n.Add(time.Duration(sign) * time.Nanosecond).Format("2006-01-02 15:04:05.999999999")
			leftBoundString = fmt.Sprintf("%s %s", nanosecondString, tz)
			rightBoundString = fmt.Sprintf("%s %s", nanosecondString, tz)
		case "microsecond", "microseconds":
			microsecondString := n.Add(time.Duration(sign) * time.Microsecond).Format("2006-01-02 15:04:05.999999")
			leftBoundString = fmt.Sprintf("%s000 %s", microsecondString, tz)
			rightBoundString = fmt.Sprintf("%s999 %s", microsecondString, tz)
		case "millisecond", "milliseconds":
			millisecondString := n.Add(time.Duration(sign) * time.Millisecond).Format("2006-01-02 15:04:05.999")
			leftBoundString = fmt.Sprintf("%s000000 %s", millisecondString, tz)
			rightBoundString = fmt.Sprintf("%s999999 %s", millisecondString, tz)
		case "second", "seconds":
			secondString := n.Add(time.Duration(sign) * time.Second).Format("2006-01-02 15:04:05")
			leftBoundString = fmt.Sprintf("%s.000000000 %s", secondString, tz)
			rightBoundString = fmt.Sprintf("%s.999999999 %s", secondString, tz)
		case "minute", "minutes":
			minuteString := n.Add(time.Duration(sign) * time.Minute).Format("2006-01-02 15:04")
			leftBoundString = fmt.Sprintf("%s:00.000000000 %s", minuteString, tz)
			rightBoundString = fmt.Sprintf("%s:59.999999999 %s", minuteString, tz)
		case "hour", "hours":
			hourString := n.Add(time.Duration(sign) * time.Hour).Format("2006-01-02 15")
			leftBoundString = fmt.Sprintf("%s:00:00.000000000 %s", hourString, tz)
			rightBoundString = fmt.Sprintf("%s:59:59.999999999 %s", hourString, tz)
		case "day", "days":
			dayString := n.AddDate(0, 0, sign*1).Format("2006-01-02")
			leftBoundString = fmt.Sprintf("%s  00:00:00.000000000 %s", dayString, tz)
			rightBoundString = fmt.Sprintf("%s 23:59:59.999999999 %s", dayString, tz)
		case "week", "weeks":
			var nextMonday time.Time
			switch n.Weekday() {
			case time.Monday:
				nextMonday = n.AddDate(0, 0, 7)
			case time.Tuesday:
				nextMonday = n.AddDate(0, 0, 6)
			case time.Wednesday:
				nextMonday = n.AddDate(0, 0, 5)
			case time.Thursday:
				nextMonday = n.AddDate(0, 0, 4)
			case time.Friday:
				nextMonday = n.AddDate(0, 0, 3)
			case time.Saturday:
				nextMonday = n.AddDate(0, 0, 2)
			case time.Sunday:
				nextMonday = n.AddDate(0, 0, 1)
			}
			mondayString := nextMonday.Format("2006-01-02")
			leftBoundString = fmt.Sprintf("%s  00:00:00.000000000 %s", mondayString, tz)
			rightBoundString = fmt.Sprintf("%s 23:59:59.999999999 %s", mondayString, tz)
		case "month", "months":
			daysInMonth := map[string]uint8{
				"1":  30,
				"2":  31,
				"3":  30,
				"4":  31,
				"5":  30,
				"6":  31,
				"7":  30,
				"8":  31,
				"9":  30,
				"10": 31,
				"11": 30,
				"12": 31,
			}

			d := n.AddDate(0, sign*1, 0)
			days := daysInMonth[d.Format("M")]
			monthString := n.AddDate(0, sign*1, 0).Format("2006-01")

			leftBoundString = fmt.Sprintf("%s-01  00:00:00.000000000 %s", monthString, tz)
			rightBoundString = fmt.Sprintf("%s-%d 23:59:59.999999999 %s", monthString, days, tz)
		case "year", "years":
			yearString := n.AddDate(sign*1, 0, 0).Format("2006")
			leftBoundString = fmt.Sprintf("%s-01-01  00:00:00.000000000 %s", yearString, tz)
			rightBoundString = fmt.Sprintf("%s-12-31 23:59:59.999999999 %s", yearString, tz)
		default:
			panic(fmt.Errorf("verbal [%s] not recognized", b.verbal))
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
	} else if s.rightBoundRelN != nil {
		rt := s.rightBoundRelN.resolveAt(t, false)
		w.to = &rt
	}

	// edge-case: left bound is a Rel and the right is an Abs, so calculate the left bound abs value relative to the right bound abs value
	if s.leftBoundRel != nil && w.to != nil {
		lt := w.to.Add(-*s.leftBoundRel)
		w.from = &lt
		w.slide = 0 // reset the slide
	}

	w.validate()
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
	if w.from == nil || w.to == nil {
		panic("absolute bound are not defined on this window")
	}
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

func (w *Window) validate() {
	if w.from != nil && w.to != nil && w.from.After(*w.to) {
		panic(fmt.Errorf("window bounds are in wrong order"))
	}

	if w.from == nil && w.to == nil && w.slide == 0 {
		panic("empty window")
	}
}
