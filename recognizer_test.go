package window

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/araddon/dateparse"
)

func Test_recognitionFail(t *testing.T) {
	type test struct {
		text string
		err  string
	}
	tests := []test{
		{"", "failed to recognize the left bound"},
		{"a", "failed to recognize the left bound"},
		{"WITHIN", "failed to recognize the left bound"},
		{"WITHIN one day", "failed to recognize the left bound"},
		{"WITHIN 1 2d", "failed to recognize the left bound"},
		{"WITHIN 1", "failed to recognize the left bound"},
		{"3 days max", "failed to recognize the right bound"},
		{"1 April 2022 to", "failed to recognize the right bound"},
		{"1 April 2022 to ", "failed to recognize the right bound"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			_, err := Start(tt.text)
			minLen := len(tt.err)
			if minLen > len(err.Error()) {
				minLen = len(err.Error())
			}
			if err.Error()[:minLen] != tt.err[:minLen] {
				t.Errorf("error [%s] should be [%s]", err, tt.err)
			}
		})

	}
}

func Test_recognitionSuccess(t *testing.T) {
	type test struct {
		text       string
		windowFunc func() Specification
	}
	tests := []test{
		// 1. Abs-Abs
		{"1 Jan 1991 to 2 Feb 1992", func() Specification {
			d1, _ := dateparse.ParseStrict("1 Jan 1991")
			d2, _ := dateparse.ParseStrict("2 Feb 1992")
			return MakeSpecification(d1, d2)
		}},
		// 2. Abs-Rel
		{"1 Jan 1991 within 1 day", func() Specification {
			d1, _ := dateparse.ParseStrict("1 Jan 1991")
			return MakeSpecification(d1, 1*24*time.Hour)
		}},
		{"1 April 2022 within 1 day", func() Specification {
			d1, _ := dateparse.ParseStrict("1 Apr 2022")
			return MakeSpecification(d1, 1*24*time.Hour)
		}},
		// 3. Abs-RelN
		{"1 April 2022 to tomorrow", func() Specification {
			d1, _ := dateparse.ParseStrict("1 Apr 2022")
			return MakeSpecification(d1, boundRelativeToNow{verbal: "tomorrow"})
		}},
		{"1 Jan 1991 to last week", func() Specification {
			d1, _ := dateparse.ParseStrict("1 Jan 1991")
			return MakeSpecification(d1, boundRelativeToNow{inFuture: false, verbal: "week"})
		}},
		{"1 Jan 1991 to 1 day ahead", func() Specification {
			d1, _ := dateparse.ParseStrict("1 Jan 1991")
			return MakeSpecification(d1, boundRelativeToNow{inFuture: true, duration: 24 * time.Hour})
		}},
		{"1 Jan 1991 until now", func() Specification {
			d1, _ := dateparse.ParseStrict("1 Jan 1991")
			return MakeSpecification(d1, boundRelativeToNow{verbal: "now"})
		}},
		// 4. Rel-Abs
		{"yesterday to 1 Apr 2022", func() Specification {
			d1, _ := dateparse.ParseStrict("1 Apr 2022")
			return MakeSpecification(boundRelativeToNow{verbal: "yesterday"}, d1)
		}},
		// 5. Rel-Rel (Sliding window)
		{"3 days", func() Specification {
			return MakeSpecification(3*24*time.Hour, nil)
		}},
		{"Within 3 days", func() Specification {
			return MakeSpecification(3*24*time.Hour, nil)
		}},
		// 6. Rel-RelN
		{"3 days until yesterday", func() Specification {
			return MakeSpecification(3*24*time.Hour, boundRelativeToNow{verbal: "yesterday"})
		}},
		{"3 days until last year", func() Specification {
			return MakeSpecification(3*24*time.Hour, boundRelativeToNow{inFuture: false, verbal: "year"})
		}},
		{"30 days until 2 days ago", func() Specification {
			return MakeSpecification(
				30*24*time.Hour,
				boundRelativeToNow{inFuture: false, duration: 2 * 24 * time.Hour},
			)
		}},
		// 7. RelN-Abs
		{"30 days to 1 Apr 2022", func() Specification {
			d1, _ := dateparse.ParseStrict("1 Apr 2022")
			return MakeSpecification(30*24*time.Hour, d1)
		}},
		// 8. RelN-Rel
		{"yesterday within 30 days", func() Specification {
			return MakeSpecification(boundRelativeToNow{verbal: "yesterday"}, 30*24*time.Hour)
		}},
		{"next year within 3 days and 2 hours", func() Specification {
			return MakeSpecification(boundRelativeToNow{inFuture: true, verbal: "year"}, 3*24*time.Hour+2*time.Hour)
		}},
		// 9. RelN-RelN
		{"from last month until 2 hours later", func() Specification {
			return MakeSpecification(
				boundRelativeToNow{inFuture: false, verbal: "month"},
				boundRelativeToNow{inFuture: true, duration: 2 * time.Hour},
			)
		}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			win, err := Start(tt.text)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(win, tt.windowFunc()) {
				t.Errorf("window [%v] should be [%v]", win, tt.windowFunc())
			}
		})

	}
}
