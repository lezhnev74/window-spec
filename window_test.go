package window

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/araddon/dateparse"
)

func Test(t *testing.T) {
	now := dateparse.MustParse("1 May 2022 00:00:00")

	type test struct {
		text       string
		windowFunc func() Window
	}
	tests := []test{
		// 1. Abs-Abs
		{"1 April 2022 to 2 April 2022", func() Window {
			d1 := dateparse.MustParse("01 Apr 2022 00:00:00.000000000")
			d2 := dateparse.MustParse("2 Apr 2022 00:00:00.000000000")
			return Window{from: &d1, to: &d2}
		}},
		// 2. Abs-Rel
		{"1 April 2022 within 1 day", func() Window {
			d1 := dateparse.MustParse("01 Apr 2022 00:00:00.000000000")
			d2 := dateparse.MustParse("2 Apr 2022 00:00:00.000000000")
			return Window{from: &d1, to: &d2}
		}},
		// 3. Abs-RelN
		{"1 April 2022 to tomorrow", func() Window {
			d1 := dateparse.MustParse("1 Apr 2022 00:00:00")
			d2 := dateparse.MustParse("2 May 2022 00:00:00")
			return Window{from: &d1, to: &d2}
		}},
		// 4. Rel-Abs
		{"1 day to 2 April 2022", func() Window {
			d1 := dateparse.MustParse("01 Apr 2022 00:00:00.000000000")
			d2 := dateparse.MustParse("2 Apr 2022 00:00:00.000000000")
			return Window{from: &d1, to: &d2}
		}},
		// 5. Rel-Rel (Sliding window)
		{"30 days", func() Window {
			return Window{slide: 30 * 24 * time.Hour}
		}},
		// 6. Rel-RelN
		{"2 days to next week", func() Window {
			d1 := dateparse.MustParse("30 Apr 2022 00:00:00.000000000")
			d2 := dateparse.MustParse("2 May 2022 00:00:00.000000000")
			return Window{from: &d1, to: &d2}
		}},
		{"1 days to next day", func() Window {
			d1 := dateparse.MustParse("1 May 2022 00:00:00.000000000")
			d2 := dateparse.MustParse("2 May 2022 00:00:00.000000000")
			return Window{from: &d1, to: &d2}
		}},
		{"1 days to next month", func() Window {
			d1 := dateparse.MustParse("31 May 2022 00:00:00.000000000")
			d2 := dateparse.MustParse("1 Jun 2022 00:00:00.000000000")
			return Window{from: &d1, to: &d2}
		}},
		{"30 days to next year", func() Window {
			d1 := dateparse.MustParse("2 Dec 2022 00:00:00.000000000")
			d2 := dateparse.MustParse("1 Jan 2023 00:00:00.000000000")
			return Window{from: &d1, to: &d2}
		}},
		// 7. RelN-Abs
		{"next year to 20 May 2024", func() Window {
			d1 := dateparse.MustParse("31 Dec 2023 23:59:59.999999999")
			d2 := dateparse.MustParse("20 May 2024 00:00:00.000000000")
			return Window{from: &d1, to: &d2}
		}},
		// 8. RelN-Rel
		{"next year within 3 days and 2 hours", func() Window {
			d1 := dateparse.MustParse("31 Dec 2023 23:59:59.999999999")
			d2 := dateparse.MustParse("4 Jan 2024 01:59:59.999999999")
			return Window{from: &d1, to: &d2}
		}},
		// 9. RelN-RelN
		{"last hour to next minute", func() Window {
			d1 := dateparse.MustParse("30 Apr 2022 23:59:59.999999999")
			d2 := dateparse.MustParse("1 May 2022 00:01:00.000000000")
			return Window{from: &d1, to: &d2}
		}},
		{"last second to next millisecond", func() Window {
			d1 := dateparse.MustParse("30 Apr 2022 23:59:59.999999999")
			d2 := dateparse.MustParse("1 May 2022 00:00:00.001000000")
			return Window{from: &d1, to: &d2}
		}},
		{"last microsecond to next nanosecond", func() Window {
			d1 := dateparse.MustParse("30 Apr 2022 23:59:59.999999999")
			d2 := dateparse.MustParse("1 May 2022 00:00:00.000000001")
			return Window{from: &d1, to: &d2}
		}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			winSpec, err := Start(tt.text)
			if err != nil {
				t.Fatal(err)
			}
			win := winSpec.ResolveAt(now)
			expectedWin := tt.windowFunc()
			if !reflect.DeepEqual(win, &expectedWin) {
				t.Errorf("window [%v] should be [%v]", win, tt.windowFunc())
			}
		})

	}
}
