package course

import (
	"time"
)

var NoRepetition Repetition = Repetition{}

type Repetition struct {
	interval Interval
	end      Date
	on       bool
}

type Interval struct {
	Freq Frequency
	Skip int
}

type Frequency interface {
	Next(Date, int) Date
	IsON(Date) bool
}

type Daily struct{}
type Monthly MonthDays
type Weekly Weekdays

// make sure they implement Frequency interface
var _ Frequency = Daily{}
var _ Frequency = Monthly(0)
var _ Frequency = Weekly(0)

func NewDaily() Daily {
	return Daily{}
}

func NewWeekly(on ...time.Weekday) Weekly {
	w := Weekdays(0)
	w.SetON(on...)
	return Weekly(w)
}

func NewMonthly(on ...int) Monthly {
	w := MonthDays(0)
	w.SetON(on...)
	return Monthly(w)
}

func NewRepetition(end Date, inv Interval) Repetition {
	if inv.Skip > 100 || inv.Skip < 0 {
		panic("Interval too long or negative")
	}
	return Repetition{
		on:       true,
		end:      end,
		interval: inv,
	}
}

func (d Daily) Next(start Date, n int) Date {
	return start.AddDate(0, 0, n)
}

func (d Daily) IsON(start Date) bool {
	return true
}

func (d Monthly) Next(start Date, n int) Date {
	m := MonthDays(d)

	if d.IsON(start) {
		month := start.month
		start = start.AddDate(0, 0, 1)
		if month != start.month {
			start.day = m.First()
			n--
		}
		for i := start.day; i <= 31 && month == start.month; i++ {
			if m.IsON(i) {
				max := daysInMonth(start.month, isLeapYear(start.year))
				if i > max {
					return start.AddDate(0, 0, int(max-start.day))
				}
				return start.AddDate(0, 0, int(i-start.day))
			}
		}
	}
	for i := 0; i < n; i++ {
		start = m.Next(start)
	}
	return start
}

func (d Monthly) IsON(start Date) bool {
	current := start.day
	w := MonthDays(d)
	if w.IsON(current) {
		return true
	}

	if start.month == time.February {
		if isLeapYear(start.year) && start.day == 29 && (w.IsON(30) || w.IsON(31)) {
			return true
		} else if start.day == 28 && (w.IsON(30) || w.IsON(31) || w.IsON(29)) {
			return true
		}
	}
	return false
}

func (d Weekly) Next(start Date, n int) Date {
	w := Weekdays(d)
	if d.IsON(start) {
		weekday := start.ToTime().Weekday()
		for i := weekday + 1; i <= time.Saturday; i++ {
			if w.IsON(i) {
				return start.AddDate(0, 0, int(i-weekday))
			}
		}
	}
	for i := 0; i < n; i++ {
		start = w.Next(start)
	}
	return start
}

func (d Weekly) IsON(start Date) bool {
	w := Weekdays(d)
	return w.IsON(start.ToTime().Weekday())
}
