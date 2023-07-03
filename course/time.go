package course

import (
	"math/bits"
	"scheduler/util"
	"time"
)

type Time uint32

func NewTime(h, m, s int) Time {
	if h < 0 || h > 24 {
		panic("Invalid hour")
	}
	if m < 0 || m > 60 {
		panic("Invalid minute")
	}
	if s < 0 || s > 60 {
		panic("Invalid second")
	}
	return Time(h*3600 + m*60 + s)
}

func (t Time) Unix() int64 {
	return time.Date(0, 0, 0, int(t.Hour()), int(t.Minute()), int(t.Second()), 0, time.Local).Unix()
}

func (t Time) UnixWithDate(d Date) int64 {
	return time.Date(d.year, d.month, d.day, int(t.Hour()), int(t.Minute()), int(t.Second()), 0, time.Local).Unix()
}

func (t Time) ToTime(d Date) time.Time {
	return time.Date(d.year, d.month, d.day, int(t.Hour()), int(t.Minute()), int(t.Second()), 0, time.Local)
}

func (t Time) String() string {
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()
	var b [8]byte
	b[0] = byte(hour/10 + 48)
	b[1] = byte(hour%10 + 48)
	b[2] = ':'
	b[3] = byte(minute/10 + 48)
	b[4] = byte(minute%10 + 48)
	b[5] = ':'
	b[6] = byte(second/10 + 48)
	b[7] = byte(second%10 + 48)
	return util.B2S(b[:])
}

func (t Time) Hour() int {
	return int(t / 3600)
}

func (t Time) Minute() int {
	return int((t % 3600) / 60)
}

func (t Time) Second() int {
	return int(t % 60)
}

func (t Time) Before(o Time) bool {
	return t < o
}

func (t Time) After(o Time) bool {
	return t > o
}

const day Time = 24 * 60 * 60

func (t Time) Add(d time.Duration) Time {
	t += Time(d.Seconds())
	if t >= day {
		t %= day
	}
	return t
}

type Date struct {
	day, year int
	month     time.Month
}

func (d Date) Day() int {
	return d.day
}
func (d Date) Month() time.Month {
	return d.month
}
func (d Date) Year() int {
	return d.year
}

func (d Date) After(other Date) bool {
	if d.year > other.year {
		return true
	} else if d.year == other.year {
		if d.month > other.month {
			return true
		} else if d.month == other.month {
			return d.day > other.day
		}
	}
	return false
}

func (d Date) Before(other Date) bool {
	if d.year < other.year {
		return true
	} else if d.year == other.year {
		if d.month < other.month {
			return true
		} else if d.month == other.month {
			return d.day < other.day
		}
	}
	return false
}

func (d Date) Same(other Date) bool {
	return d == other
}

func NewDateFromTime(t time.Time) Date {
	return NewDate(t.Date())
}

func NewDate(year int, month time.Month, day int) Date {
	if month > time.December || month < time.January {
		panic("invalid month")
	}
	if daysMonth[month] < day && !(day == 29 && isLeapYear(year)) {
		panic("invalid days in month")
	}
	return Date{
		day:   day,
		month: month,
		year:  year,
	}
}

var daysMonth = [13]int{
	0,  // 0th index is not used, months start from index 1 (January) [compatible with time.Month]
	31, // January
	28, // February (non-leap year)
	31, // March
	30, // April
	31, // May
	30, // June
	31, // July
	31, // August
	30, // September
	31, // October
	30, // November
	31, // December
}

func (d Date) String() string {
	return d.ToTime().Format(time.DateOnly)
}

func (d Date) ToTime() time.Time {
	return time.Date(d.year, d.month, d.day, 0, 0, 0, 0, time.Local)
}

func (d Date) ToTimeAt(t Time) time.Time {
	return time.Date(d.year, d.month, d.day, t.Hour(), t.Minute(), t.Second(), 0, time.Local)
}

func (d Date) NextDay(n int) Date {
	return d.AddDate(0, 0, n)
}

func (d Date) NextWeek(n int) Date {
	return d.AddDate(0, 0, n*7)
}

func (d Date) NextMonth(n int) Date {
	return d.AddDate(0, n, 0)
}

func (d Date) NextYear(n int) Date {
	return d.AddDate(n, 0, 0)
}

func (d Date) AddDate(years, months, days int) Date {
	// Adjust years and months
	d.year += years
	d.month += time.Month(months)

	// Normalize months
	if d.month < 1 {
		if d.month == 0 {
			d.year--
			d.month = 12
		} else {
			d.year += int(d.month) / 12
			if d.month > -12 {
				d.year--
				d.month = 12 + d.month
			} else {
				d.month = -d.month % 12
				if d.month == 0 {
					d.year--
					d.month = 12
				}
			}
		}
	} else if d.month > 12 {
		d.year += int(d.month) / 12
		d.month = d.month % 12
	}

	// Get the maximum number of days in the new month
	maxDays := daysInMonth(d.month, isLeapYear(d.year))

	// Adjust the day if it exceeds the maximum
	if d.day > maxDays {
		d.day = maxDays
	}

	// Subtract remaining days if negative
	if days < 0 {
		d.normalizeNegativeDate(days)
	} else {
		d.normalizePositiveDate(days)
	}
	return d
}

func (d *Date) normalizePositiveDate(days int) {
	// Add remaining days
	d.day += days
	maxDays := daysInMonth(d.month, isLeapYear(d.year))

	// Handle overflow to the next month
	for d.day > maxDays {
		d.day -= maxDays
		d.month++
		if d.month > 12 {
			d.year++
			d.month = 1
		}
		maxDays = daysInMonth(d.month, isLeapYear(d.year))
	}
}

func (d *Date) normalizeNegativeDate(days int) {
	remove := -days
	for remove != 0 {
		if d.day > remove {
			d.day -= remove
			remove = 0
		} else {
			d.month--
			if d.month < 1 {
				d.year--
				d.month = 12
			}
			remove -= d.day
			d.day = daysInMonth(d.month, isLeapYear(d.year))
		}
	}
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func daysInYear(year int) int {
	if isLeapYear(year) {
		return 366
	} else {
		return 365
	}
}

func daysInMonth(month time.Month, leap bool) int {
	if leap && month == time.February {
		return 29
	} else {
		return daysMonth[month]
	}
}

type Weekdays uint8

func (w *Weekdays) SetON(c ...time.Weekday) {
	for _, d := range c {
		if d > 6 {
			continue
		}
		*w |= 1 << d
	}
}

func (w *Weekdays) SetOFF(c ...time.Weekday) {
	for _, day := range c {
		*w &^= 1 << day
	}
}

func (w Weekdays) IsON(day time.Weekday) bool {
	return (w & (1 << day)) != 0
}

func (w Weekdays) First() time.Weekday {
	a := bits.TrailingZeros8(uint8(w))
	if a >= 7 {
		return 1 << 7
	} else {
		return time.Weekday(a)
	}
}

func (w Weekdays) Next(d Date) Date {
	current := d.ToTime().Weekday()
	diff := w.First() - current
	if diff == 0 {
		return d.AddDate(0, 0, 7)
	} else if diff < 0 {
		return d.AddDate(0, 0, 7+int(diff))
	} else {
		return d.AddDate(0, 0, int(diff))
	}
}

type MonthDays uint32

func (w *MonthDays) SetON(c ...int) {
	for _, d := range c {
		if d > 31 {
			continue
		}
		*w |= 1 << d
	}
}

func (w *MonthDays) SetOFF(c ...int) {
	for _, day := range c {
		*w &^= 1 << day
	}
}

func (w MonthDays) IsON(day int) bool {
	return (w & (1 << day)) != 0
}

func (w MonthDays) First() int {
	a := bits.TrailingZeros32(uint32(w))
	if a < 1 || a > 31 {
		return 1 << 31
	} else {
		return a
	}
}

func (w MonthDays) Next(d Date) Date {
	current := d.day

	if d.month == time.February {
		if isLeapYear(d.year) && d.day == 29 && (w.IsON(30) || w.IsON(31)) {
			return d
		} else if d.day == 28 && (w.IsON(30) || w.IsON(31) || w.IsON(29)) {
			return d
		}
	}

	diff := w.First() - current
	if diff <= 0 {
		return d.AddDate(0, 1, diff)
	} else {
		return d.AddDate(0, 0, diff)
	}
}

// daysBefore[m] counts the number of days in a non-leap year
// before month m begins. There is an entry for m=12, counting
// the number of days before January of next year (365).
var daysBefore = [...]int32{
	0,
	31,
	31 + 28,
	31 + 28 + 31,
	31 + 28 + 31 + 30,
	31 + 28 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30,
	31 + 28 + 31 + 30 + 31 + 30 + 31 + 31 + 30 + 31 + 30 + 31,
}

// func (a Date) negative(b Date) Interval {
// 	// Calculate the difference in years
// 	yearsDiff := a.year - b.year

// 	// Calculate the difference in months
// 	monthsDiff := a.month - b.month

// 	// Calculate the difference in days
// 	daysDiff := a.day - b.day

// 	// Adjust for negative differences
// 	if daysDiff > 0 {
// 		monthsDiff++
// 		daysDiff = -b.day
// 		maxDays := daysInMonth(a.month, isLeapYear(b.year))
// 		daysDiff -= maxDays - a.day
// 	}
// 	if monthsDiff > 0 {
// 		yearsDiff++
// 		monthsDiff -= 12
// 	}

// 	return Interval{
// 		Year:  yearsDiff,
// 		Month: int(monthsDiff),
// 		Day:   daysDiff,
// 	}

// }

// func (a Date) Sub(b Date) Interval {
// 	// var i Interval
// 	if a.Before(b) {
// 		return a.negative(b)
// 	}
// 	// Calculate the difference in years
// 	yearsDiff := a.year - b.year

// 	// Calculate the difference in months
// 	monthsDiff := a.month - b.month

// 	// Calculate the difference in days
// 	daysDiff := a.day - b.day

// 	// Adjust for negative differences
// 	if daysDiff < 0 {
// 		monthsDiff--
// 		daysDiff = a.day
// 		maxDays := daysInMonth(a.month-1, isLeapYear(b.year))
// 		daysDiff += maxDays - b.day
// 	}
// 	if monthsDiff < 0 {
// 		yearsDiff--
// 		monthsDiff += 12
// 	}

// 	return Interval{
// 		Year:  yearsDiff,
// 		Month: int(monthsDiff),
// 		Day:   daysDiff,
// 	}
// }
