package window

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_recognitionFail(t *testing.T) {
	type test struct {
		text string
		err  string
	}
	tests := []test{
		{"", "unexpected eof"},
		{"a", "unexpected character found at 0"},
		{"WITHIN", "unexpected eof"},
		{"WITHIN one day", "unexpected character found at 7"},
		{"WITHIN 1 2d", "unexpected character found at 9"},
		{"WITHIN 1", "unexpected eof"},
		{"3 days max", "unexpected character found at 7"},
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
		text   string
		window Specification
	}
	tests := []test{
		{"3 days", MakeSpecification(3*24*time.Hour, nil)},
		{"Within 3 days", MakeSpecification(3*24*time.Hour, nil)},
		{"3 days until yesterday", MakeSpecification(3*24*time.Hour, boundRelativeToNow{verbal: "yesterday"})},
		{"3 days until last year", MakeSpecification(3*24*time.Hour, boundRelativeToNow{inFuture: false, verbal: "year"})},
		{"30 days until 2 days ago", MakeSpecification(
			30*24*time.Hour,
			boundRelativeToNow{inFuture: false, duration: 2 * 24 * time.Hour},
		)},
		{"from last month until 2 hours later", MakeSpecification(
			boundRelativeToNow{inFuture: false, verbal: "month"},
			boundRelativeToNow{inFuture: true, duration: 2 * time.Hour},
		)},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			win, err := Start(tt.text)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(win, tt.window) {
				t.Errorf("window [%v] should be [%v]", win, tt.window)
			}
		})

	}
}
