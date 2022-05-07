# Window Language Parser (in Go)

A window is a term from the time series domain. It is a period of time.
This package allows to define windows in plain text and then convert them to Go structures.

Examples:
- `A`
- `B`
- `C`

## How To Define A Window

A windows can be defined by its left and right bounds. `W: (a,b)`. Bounds are time points (as precise as we want). This
package does not distinguish included/excluded bounds. A time point can be absolute or relative. Relative to what? There
are two options: relative to the other point or to the NOW point.
Summarize of window bounds:
- absolute bound (ex: `9:00 am, 22 June, 2022`)
- relative to the other bound (ex: `6 days BEFORE 1 September 2022`)
- relative to now (ex: `1 may 1992 TO today`)

Keywords allowed before every type:

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
      <code>x AGO</code> or <code>x LATER</code> where "x" is a combination of <code>number unit (and number unit)*</code>
      <br> units: nanosecond, microsecond, millisecond, second, minute, hour, day, week, month
      <br> Also possible more sophisitcated queries: <code>last X</code> (last year)
    </td>
  </tr>
  <tr>
    <td>Relative To Another Bound</td>
    <td><code>WITHIN</code></td>
    <td><code>WITHIN</code></td>
    <td><code>number unit (and number unit)*</code></td>
  </tr>
</table>

Combining these types we can specify a window in handful of ways:
<table>
    <thead>
    <tr>
        <th>Left Bound</th>        
        <th>Right Bound</th>
        <th>Pattern</th>
        <th>Example</th>        
    </tr>
    </thead>
    <tr>
        <td>Absolute</td>
        <td>Absolute</td>
        <td><code>FROM a TO b</code></td>
        <td>FROM 1 January 1991 TO 31 December 1991</td>
    </tr>
    <tr>
        <td>Absolute</td>
        <td>Relative To Another Bound</td>
        <td><code>FROM a WITHIN b</code></td>
        <td>FROM 1 January 1991 WITHIN 365 days</td>
    </tr>
    <tr>
        <td>Absolute</td>
        <td>Relative To Now</td>
        <td><code>FROM a UNTIL b</code></td>
        <td>FROM 1 January 1991 UNTIL last week</td>
    </tr>
    <tr>
        <td>Relative To Another Bound</td>
        <td>Absolute</td>
        <td><code>WITHIN a BEFORE b</code></td>
        <td>WITHIN 365 days BEFORE 31 December 1991</td>
    </tr>
    <tr>
        <td>Relative To Another Bound</td>
        <td>Relative To Another Bound</td>
        <td><code>WITHIN a</code></td>
        <td>A special case for a window that only has a duration, but starts at any point in time. Called "Sliding Window".
            <br>
            WITHIN 60 DAYS
        </td>
    </tr>
    <tr>
        <td>Relative To Another Bound</td>
        <td>Relative To Now</td>
        <td><code>WITHIN a UNTIL b</code></td>
        <td>WITHIN 7 days UNTIL yesterday</td>
    </tr>
    <tr>
        <td>Relative To Now</td>
        <td>Absolute</td>
        <td><code>FROM a TO b</code></td>
        <td>FROM 7 years ago TO 31 December 1991</td>
    </tr>
    <tr>
        <td>Relative To Now</td>
        <td>Relative To Another Bound</td>
        <td><code>FROM a UNTIL b</code></td>
        <td>FROM 7 years ago UNTIL last week</td>
    </tr>
    <tr>
        <td>Relative To Now</td>
        <td>Relative To Now</td>
        <td><code>FROM a TO b</code></td>
        <td>FROM last year TO this year</td>
    </tr>
</table>

## Types

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
## How To Use

This is a go package so the usage if straightforward:
TODO