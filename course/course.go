package course

import (
	"io"
	"scheduler/account"
	"time"
)

type Course struct {
	Title       string
	Description string
	Sections    []Section
}

// currently doesn't support overnight sections
type Section struct {
	Title       string
	Description string
	Instructor  *account.Instructor
	repeat      Repetition
	start       Date
	duration    time.Duration
	at          Time
}

func NewSection(title, des string, start Date, at Time, dur time.Duration, rep Repetition) Section {
	if rep.on && rep.end.Before(start) {
		panic("Repetition has an invalid end date")
	}
	// xxx check overnight section
	if at.Add(dur) < at {
		panic("Section overnight")
	}

	s := Section{
		Title:       title,
		Description: des,
		start:       start,
		at:          at,
		duration:    dur,
		repeat:      rep,
	}
	if l, _, _ := s.Next(start, 1); len(l) == 0 {
		panic("Repeated section has no actual class")
	}
	return s
}

func (c *Section) cleanup(from Date, n int) (bool, []Class, error) {
	if n == 0 {
		if from.After(c.start) {
			return false, nil, io.EOF
		} else {
			return false, nil, nil
		}
	}

	if !c.repeat.on {
		if from.After(c.start) {
			return false, nil, io.EOF
		} else {
			return false, []Class{c.First()}, io.EOF
		}
	}

	if from.After(c.repeat.end) {
		return false, nil, io.EOF
	}

	return true, nil, nil
}

func (c Section) Next(from Date, n int) ([]Class, Date, error) {

	if ok, list, err := c.cleanup(from, n); !ok {
		return list, from, err
	}

	var list []Class = nil
	var current Date = from
	if !c.repeat.interval.Freq.IsON(current) {
		current = c.repeat.interval.Freq.Next(current, 1)
	}
	ins := account.Instructor{}
	if c.Instructor != nil {
		ins = *c.Instructor
	}
	var at Time = c.at
	var class Class = Class{
		Order:      0,
		Title:      c.Title,
		Duration:   c.duration,
		Instructor: ins,
	}

	for (current.Same(c.repeat.end) || current.Before(c.repeat.end)) && n > 0 {
		// lazy init
		if list == nil {
			list = make([]Class, 0, n)
		}
		class.Order++
		class.Time = current.ToTimeAt(at)
		list = append(list, class)
		current = c.repeat.interval.Freq.Next(current, c.repeat.interval.Skip+1)
		n--
	}
	if len(list) == 0 {
		return nil, current, io.EOF
	}
	if current.After(c.repeat.end) {
		return list, current, io.EOF
	}
	return list, current, nil
}

func (c *Section) First() Class {
	var ins account.Instructor = account.Instructor{}
	if c.Instructor != nil {
		ins = *c.Instructor
	}
	if !c.repeat.on || c.repeat.end.Same(c.start) || c.repeat.interval.Freq.IsON(c.start) {
		return Class{
			Title:      c.Title,
			Order:      1,
			Time:       c.start.ToTimeAt(c.at),
			Duration:   c.duration,
			Instructor: ins,
		}
	} else {
		date := c.repeat.interval.Freq.Next(c.start, 1)
		return Class{
			Title:      c.Title,
			Order:      1,
			Time:       date.ToTimeAt(c.at),
			Duration:   c.duration,
			Instructor: ins,
		}

	}
}

func (c Section) All() []Class {

	if !c.repeat.on || c.repeat.end.Same(c.start) {
		return []Class{c.First()}
	}
	if l, _, _ := c.Next(c.start, 1); len(l) == 0 {
		return nil
	}

	var at Time = c.at
	var list []Class = make([]Class, 0, 10)
	var class Class = c.First()
	var current Date = NewDateFromTime(class.Time)
	class.Order = 0

	for current.Before(c.repeat.end) {
		class.Time = at.ToTime(current)
		class.Order++
		list = append(list, class)
		current = c.repeat.interval.Freq.Next(current, c.repeat.interval.Skip+1)
	}
	return list
}
