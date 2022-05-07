# Window Language Parser (in Go)

A window is a term from the time series domain. It is a period of time.
This package allows to define windows in plain text and then convert them to Go structures.

## Playground

Run a tool from `window` folder to try recognition of your input, here are a few examples:

```shell
$ go run main.go "within 30 days and 2 minutes and 3 nanoseconds"
Window resolved at:     2022-05-07, 17:36:53.163478476 +05
You defined a sliding window of 30 day(s) 2 minute(s) 0.000000003 seconds(s) 

$ go run main.go --timezone="Europe/Moscow" "from yesterday to today"
Your Using time.Local set to location=Europe/Moscow MSK 
Window resolved at:     2022-05-07, 15:38:50.736614500 MSK
Left Bound:             2022-05-06, 23:59:59.999999999 MSK
Right Bound:            2022-05-07, 00:00:00.000000000 MSK
```

## Syntax

A windows can be defined by its left and right bounds: `FROM a TO b` (you can use different delimiters). Bounds are time
points (as precise as we want).
This package does not distinguish included/excluded bounds. A time point can be absolute or relative. Relative to what?
There are two options: relative to the other point or to the NOW point.

For convenience a few delimiters are supported: `FROM/SINCE`, `UNTIL/TO/BEFORE` and `WITHIN`.

## Types

Supported window bound types:

- absolute bound (ex: `9:00 am 22 June, 2022`)
- relative to the other bound (ex: `6 days BEFORE 1 September 2022`)
- relative to now (ex: `1 may 1992 UNTIL today`)

<table>
  <thead>
    <tr>
      <th>Type</th>
      <th>Left Bound Keywords</th>
      <th>Right Bound Keyword</th>
      <th>Formats</th>
    </tr>
  </thead>
  <tr>
    <td>Absolute</td>
    <td rowspan="2"><code>FROM, SINCE</code></td>
    <td rowspan="2"><code>TO, BEFORE, UNTIL</code></td>
    <td>Many formats are supported, see <a href="https://github.com/araddon/dateparse">araddon/dateparse</a></td>
  </tr>
  <tr>
    <td>Relative To Now</td>
    <td>
      <code>x AGO/BEFORE/</code> or <code>x LATER/AFTER/AHEAD</code> where "x" is a combination of <code>number unit (and number unit)*</code>
      <br> units: nanosecond, microsecond, millisecond, second, minute, hour, day, week, month
      <br> Also possible more sophisitcated queries: <code>last X</code> or <code>next Y</code>
    </td>
  </tr>
  <tr>
    <td>Relative To Another Bound</td>
    <td><code>WITHIN</code></td>
    <td><code>WITHIN</code></td>
    <td><code>number unit (and number unit)*</code></td>
  </tr>
</table>

### Absolute

Absolute bound is specified as a date-time string in one of the supported formats (which is many). Ideally it should
support anything you can type for a date-time. Of course it does not, but it supports a lot.
See https://github.com/araddon/dateparse which does the absolute datetime parsing.

### Relative To Another Bound

This bound is specified as a period that is applied to another bound. The format is simple: `number unit` (like `1 day`)
. And you can add as many as you need: `1 minute and 32 seconds`.

### Relative To Now

This window bound is defined relatively to the current point in time.
It can be expressed in 2 ways: as an interval in the past or as a point in the past. Let's review them.
Window bound relative to Now:

- **a point**. Ex: `2 days and 1 second ago` or `2 hours later`. Under the hood this specification is converted to a
  duration and subtracted from (added to)  the Now time.
- **an interval**. Ex: `last week` or `next year`. This is a window by itself. It has two bounds. Depending on where
  this specification
  is used one of the bounds is picked as a result. Consider this spec `FROM last week UNTIL yesterday`. "Last week" is
  used in the left bound definition, so the left bound of the "last week" is picked as a bound. Contrary "yesterday"
  also has two bounds and the right bound is assumed as a result of evaluation.

### Examples

Combining these types we can specify a window in 9 ways:
<table>
    <thead>
    <tr>
        <th>Left Bound</th>        
        <th>Right Bound</th>
        <th>Example</th>        
    </tr>
    </thead>
    <tr>
        <td>Absolute</td>
        <td>Absolute</td>
        <td>
          <code>FROM 1 January 1991 TO 31 December 1991</code>
        </td>
    </tr>
    <tr>
        <td>Absolute</td>
        <td>Relative To Another Bound</td>
        <td>
          <code>FROM May 8, 2009 5:57:51 PM WITHIN 365 days</code>
        </td>
    </tr>
    <tr>
        <td>Absolute</td>
        <td>Relative To Now</td>
        <td>
          <code>FROM Mon Jan 02 15:04:05 -0700 2006 TO now</code>
          <br>
          <code>FROM oct 7, 1970 UNTIL last week</code>
        </td>
    </tr>
    <tr>
        <td>Relative To Another Bound</td>
        <td>Absolute</td>
        <td>
            <code>1 day TO 1332151919</code>
            <br>
            <code>WITHIN 365 days BEFORE 12 Feb 2006 19:17</code>
        </td>
    </tr>
    <tr>
        <td>Relative To Another Bound</td>
        <td>Relative To Another Bound</td>
        <td>
            A special case for a window that only has duration, starts at any point in time. 
            Called **"Sliding Window"**.
            <br>
            <code>WITHIN 60 DAYS</code> or simply <code>1 week</code>
        </td>
    </tr>
    <tr>
        <td>Relative To Another Bound</td>
        <td>Relative To Now</td>
        <td>
          <code>WITHIN 7 days UNTIL yesterday</code>
          <br>
          <code>1 second 3 nanoseconds UNTIL today</code>
        </td>
    </tr>
    <tr>
        <td>Relative To Now</td>
        <td>Absolute</td>
        <td><code>FROM 7 years ago TO 2014-12-16 06:20:00 UTC</code></td>
    </tr>
    <tr>
        <td>Relative To Now</td>
        <td>Relative To Another Bound</td>
        <td>
          <code>FROM 7 years ago UNTIL last week</code>
          <br>
          <code>2 days ago TO next week</code>
          <br>
          <code>FROM next week TO 7 days LATER</code>
        </td>
    </tr>
    <tr>
        <td>Relative To Now</td>
        <td>Relative To Now</td>
        <td><code>FROM last year TO this year</code></td>
    </tr>
</table>

## How To Use

```go
// go get github.com/lezhnev74/window-spec

windowQuery := "30 days"
winSpec, err := Start(windowQuery) // convert the text to a specification data structure
if err != nil {
    fatal(err)
}
w := winSpec.ResolveAt(time.Now()) // resolve specification relatively to a given time 
```