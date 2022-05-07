package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"window"
)

var (
	timezone  = ""
	windowStr = ""
)

func main() {
	now := time.Now()

	flag.StringVar(&timezone, "timezone", "", "Timezone aka `America/Los_Angeles` formatted time-zone")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println(`Must pass a window, and optional location:

		./window "from yesterday to 12 Apr 2022" 

		./window --timezone="America/Denver" "within 30 days"
		`)
		return
	}

	if timezone != "" {
		// NOTE:  This is very, very important to understand
		// time-parsing in go
		l, err := time.LoadLocation(timezone)
		if err != nil {
			log.Fatal(err)
		}
		zonename, _ := time.Now().In(l).Zone()
		fmt.Printf("Your Using time.Local set to location=%s %v \n", timezone, zonename)

		now = now.In(l)
	}

	windowStr = flag.Args()[0]

	winSpec, err := window.Start(windowStr)
	if err != nil {
		log.Fatal(err)
	}
	win := winSpec.ResolveAt(now)

	fmt.Printf("Window resolved at:\t%s\n", now.Format("2006-01-02, 15:04:05.000000000 MST"))

	if win.GetSlide() != 0 {
		fmt.Printf("You defined a sliding window of %s", humanizeDuration(win.GetSlide()))
	} else {
		l, r := win.GetBounds()
		fmt.Printf("Left Bound:\t\t%s\n", l.Format("2006-01-02, 15:04:05.000000000 MST"))
		fmt.Printf("Right Bound:\t\t%s\n", r.Format("2006-01-02, 15:04:05.000000000 MST"))
	}
}

func humanizeDuration(d time.Duration) (s string) {

	if days := d / (24 * time.Hour); days != 0 {
		s = fmt.Sprintf("%s %d day(s)", s, days)
		d -= days * 24 * time.Hour
	}

	if hours := d / time.Hour; hours != 0 {
		s = fmt.Sprintf("%s %d hour(s)", s, hours)
		d -= hours * time.Hour
	}

	if minutes := d / time.Minute; minutes != 0 {
		s = fmt.Sprintf("%s %d minute(s)", s, minutes)
		d -= minutes * time.Minute
	}

	if d != 0 {
		s = fmt.Sprintf("%s %.9f seconds(s)", s, float64(d)/float64(time.Second))
	}

	s = strings.Trim(s, " ")

	return
}
