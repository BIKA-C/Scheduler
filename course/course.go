package course

import (
	"io"
	"scheduler/util"
	"time"
)

type Course struct {
	ID          util.ID
	Title       string
	Description string
	Sections    []Section
}

// currently doesn't support overnight sections
type Section struct {
	Meta
	Title        string
	Description  string
	UnitPrice    int
	InstructorID util.UUID
	repeat       Repetition
	start        Date
	duration     time.Duration
	at           Time
}

func NewSection(title, des string, p int, start Date, at Time, dur time.Duration, rep Repetition) Section {
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
		UnitPrice:   p,
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

func (c *Section) Next(from Date, n int) ([]Class, Date, error) {

	if ok, list, err := c.cleanup(from, n); !ok {
		return list, from, err
	}

	var list []Class = nil
	var current Date = from
	if !c.repeat.interval.Freq.IsON(current) {
		current = c.repeat.interval.Freq.Next(current, 1)
	}
	var at Time = c.at
	var class Class = Class{
		Index:        0,
		UnitPrice:    c.UnitPrice,
		Title:        c.Title,
		Duration:     c.duration,
		InstructorID: c.InstructorID,
	}

	for (current.Same(c.repeat.end) || current.Before(c.repeat.end)) && n > 0 {
		// lazy init
		if list == nil {
			list = make([]Class, 0, n)
		}
		class.Index++
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
	if !c.repeat.on || c.repeat.end.Same(c.start) || c.repeat.interval.Freq.IsON(c.start) {
		return Class{
			Title:        c.Title,
			Index:        1,
			UnitPrice:    c.UnitPrice,
			Time:         c.start.ToTimeAt(c.at),
			Duration:     c.duration,
			InstructorID: c.InstructorID,
		}
	} else {
		date := c.repeat.interval.Freq.Next(c.start, 1)
		return Class{
			Title:        c.Title,
			Index:        1,
			UnitPrice:    c.UnitPrice,
			Time:         date.ToTimeAt(c.at),
			Duration:     c.duration,
			InstructorID: c.InstructorID,
		}

	}
}

func (c *Section) All() []Class {

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
	class.Index = 0

	for current.Before(c.repeat.end) {
		class.Time = at.ToTime(current)
		class.Index++
		list = append(list, class)
		current = c.repeat.interval.Freq.Next(current, c.repeat.interval.Skip+1)
	}
	return list
}
