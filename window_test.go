package window

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/araddon/dateparse"
)

func Test(t *testing.T) {
	now, _ := dateparse.ParseStrict("1 May 2022 00:00:00")

	type test struct {
		text       string
		windowFunc func() Window
	}
	tests := []test{
		// 1. Abs-Abs
		// 2. Abs-Rel
		// 3. Abs-RelN
		// 4. Rel-Abs
		// 5. Rel-Rel (Sliding window)
		// 6. Rel-RelN
		// 7. RelN-Abs
		// 8. RelN-Rel
		// 9. RelN-RelN
		{"1 April 2022 to tomorrow", func() Window {
			d1 := dateparse.MustParse("1 Apr 2022 00:00:00")
			d2, _ := dateparse.ParseStrict("2 May 2022 00:00:00")
			return Window{from: &d1, to: &d2}
		}},
		{"1 April 2022 to tomorrow", func() Window {
			d1 := dateparse.MustParse("1 Apr 2022 00:00:00")
			d2, _ := dateparse.ParseStrict("2 May 2022 00:00:00")
			return Window{from: &d1, to: &d2}
		}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			winSpec, err := Start(tt.text)
			if err != nil {
				t.Error(err)
			}
			win := winSpec.ResolveAt(now)
			expectedWin := tt.windowFunc()
			if !reflect.DeepEqual(win, &expectedWin) {
				t.Errorf("window [%v] should be [%v]", win, tt.windowFunc())
			}
		})

	}
}
