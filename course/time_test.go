package course

import (
	"testing"
	"time"
)

type addDate struct {
	years  int
	months int
	days   int
}

var addDateTests = []struct {
	name   string
	fields Date
	toAdd  addDate
	want   Date
}{
	{
		name:   "Test 1",
		fields: NewDate(2023, 1, 31),
		toAdd:  addDate{0, 1, 0},
		want:   NewDate(2023, 2, 28),
	}, {
		name:   "Test 2",
		fields: NewDate(2024, 1, 31),
		toAdd:  addDate{0, 1, 0},
		want:   NewDate(2024, 2, 29),
	}, {
		name:   "Test 3",
		fields: NewDate(2023, 1, 31),
		toAdd:  addDate{0, 2, 0},
		want:   NewDate(2023, 3, 31),
	}, {
		name:   "Test 4",
		fields: NewDate(2023, 3, 31),
		toAdd:  addDate{0, 1, 0},
		want:   NewDate(2023, 4, 30),
	}, {
		name:   "Test 5",
		fields: NewDate(2023, 1, 31),
		toAdd:  addDate{1, 1, 0},
		want:   NewDate(2024, 2, 29),
	}, {
		name:   "Test 6",
		fields: NewDate(2023, 3, 31),
		toAdd:  addDate{0, 0, 3},
		want:   NewDate(2023, 4, 3),
	}, {
		name:   "Test 7",
		fields: NewDate(2023, 12, 31),
		toAdd:  addDate{0, 0, 3},
		want:   NewDate(2024, 1, 3),
	}, {
		name:   "Test 8",
		fields: NewDate(2023, 12, 31),
		toAdd:  addDate{0, 0, 365},
		want:   NewDate(2024, 12, 30),
	}, {
		name:   "Test 9",
		fields: NewDate(1993, 12, 23),
		toAdd:  addDate{1, 1, 0},
		want:   NewDate(1995, 1, 23),
	}, {
		name:   "Test 10",
		fields: NewDate(1993, 12, 23),
		toAdd:  addDate{1, 337, 7},
		want:   NewDate(2023, 1, 30),
	}, {
		name:   "Reverse Test 1",
		fields: NewDate(2023, 2, 28),
		toAdd:  addDate{0, -1, 0},
		want:   NewDate(2023, 1, 28),
	}, {
		name:   "Reverse Test 2",
		fields: NewDate(2024, 2, 29),
		toAdd:  addDate{0, -1, 0},
		want:   NewDate(2024, 1, 29),
	}, {
		name:   "Reverse Test 3",
		fields: NewDate(2023, 3, 31),
		toAdd:  addDate{0, -2, 0},
		want:   NewDate(2023, 1, 31),
	}, {
		name:   "Reverse Test 4",
		fields: NewDate(2023, 4, 30),
		toAdd:  addDate{0, -1, 0},
		want:   NewDate(2023, 3, 30),
	}, {
		name:   "Reverse Test 5",
		fields: NewDate(2024, 2, 29),
		toAdd:  addDate{-1, -1, 0},
		want:   NewDate(2023, 1, 29),
	}, {
		name:   "Reverse Test 6",
		fields: NewDate(2023, 4, 3),
		toAdd:  addDate{0, 0, -3},
		want:   NewDate(2023, 3, 31),
	}, {
		name:   "Reverse Test 7",
		fields: NewDate(2024, 1, 3),
		toAdd:  addDate{0, 0, -3},
		want:   NewDate(2023, 12, 31),
	}, {
		name:   "Reverse Test 8",
		fields: NewDate(2024, 12, 30),
		toAdd:  addDate{0, 0, -365},
		want:   NewDate(2023, 12, 31),
	}, {
		name:   "Reverse Test 9",
		fields: NewDate(1995, 1, 23),
		toAdd:  addDate{-1, -1, 0},
		want:   NewDate(1993, 12, 23),
	}, {
		name:   "Reverse Test 10",
		fields: NewDate(2023, 1, 30),
		toAdd:  addDate{-1, -337, -7},
		want:   NewDate(1993, 12, 23),
	}, {
		name:   "Reverse Test A",
		fields: NewDate(1995, 1, 23),
		toAdd:  addDate{0, 0, -24},
		want:   NewDate(1994, 12, 30),
	}, {
		name:   "Reverse Test B",
		fields: NewDate(2024, 3, 23),
		toAdd:  addDate{0, 0, -24},
		want:   NewDate(2024, 2, 28),
	}, {
		name:   "General Test 1",
		fields: NewDate(2024, 3, 23),
		toAdd:  addDate{2, -7, -24},
		want:   NewDate(2025, 7, 30),
	},
}

func TestDate_AddDate(t *testing.T) {
	for _, tt := range addDateTests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Date{
				day:   tt.fields.day,
				year:  tt.fields.year,
				month: tt.fields.month,
			}
			got := d.AddDate(tt.toAdd.years, tt.toAdd.months, tt.toAdd.days)
			if got != tt.want {
				t.Errorf("%s: from %s got %d-%d-%d, want %s",
					tt.name, tt.fields, got.year, got.month, got.day, tt.want)
			}
		})
	}
}

func TestTime_String(t *testing.T) {
	tests := []struct {
		name string
		tr   Time
		want string
	}{
		{
			tr:   NewTime(8, 37, 45),
			want: "08:37:45",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("Time.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkTime_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := NewTime(10, 9, 8)
		if t.String() == "" {
		}
	}
}

func TestTime_Add(t *testing.T) {
	tests := []struct {
		name string
		tr   Time
		args time.Duration
		want Time
	}{
		{
			tr:   NewTime(23, 0, 0),
			args: time.Hour,
			want: NewTime(0, 0, 0),
		}, {
			tr:   NewTime(23, 0, 0),
			args: time.Minute*7 + time.Second*40,
			want: NewTime(23, 7, 40),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Add(tt.args); got != tt.want {
				t.Errorf("Time.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
