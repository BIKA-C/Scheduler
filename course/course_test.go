package course

import (
	"reflect"
	"testing"
	"time"
)

func TestSection_Next_From_Start(t *testing.T) {
	var EightAM = NewTime(8, 0, 0)
	tests := []struct {
		name    string
		section Section
		next    int
		want    []Class
		wantErr bool
	}{
		{
			name: "Base Test: Single Section",
			section: Section{
				start:     NewDate(2023, 1, 1),
				repeat:    Repetition{},
				UnitPrice: 1,
				at:        EightAM,
			},
			next: 1,
			want: []Class{
				{Index: 1, Time: NewDate(2023, 1, 1).ToTimeAt(EightAM), UnitPrice: 1},
			},
			wantErr: true,
		}, {
			name: "Test: Repeat Daily",
			section: Section{
				start: NewDate(2023, 1, 1),
				repeat: Repetition{
					end:      NewDate(2023, 1, 1),
					on:       true,
					interval: Interval{Skip: 0, Freq: Daily{}},
				},
			},
			next: 1,
			want: []Class{
				{Index: 1, Time: NewDate(2023, 1, 1).ToTime()},
			},
			wantErr: true,
		}, {
			name: "Test: Repeat Daily Advanced",
			section: Section{
				start: NewDate(2023, 1, 31),
				repeat: Repetition{
					end:      NewDate(2023, 2, 5),
					on:       true,
					interval: Interval{Skip: 1, Freq: Daily{}},
				},
			},
			next: 3,
			want: []Class{
				{Index: 1, Time: NewDate(2023, 1, 31).ToTime()},
				{Index: 2, Time: NewDate(2023, 2, 2).ToTime()},
				{Index: 3, Time: NewDate(2023, 2, 4).ToTime()},
			},
			wantErr: true,
		}, {
			name: "Test: Repeat Weekly",
			section: Section{
				start: NewDate(2023, 6, 25),
				repeat: Repetition{
					end:      NewDate(2023, 8, 3),
					on:       true,
					interval: Interval{Skip: 0, Freq: NewWeekly(time.Monday)},
				},
			},
			next: 3,
			want: []Class{
				{Index: 1, Time: NewDate(2023, 6, 26).ToTime()},
				{Index: 2, Time: NewDate(2023, 7, 3).ToTime()},
				{Index: 3, Time: NewDate(2023, 7, 10).ToTime()},
			},
			wantErr: false,
		}, {
			name: "Test: Repeat Weekly Advanced",
			section: Section{
				start: NewDate(2023, 6, 25),
				repeat: Repetition{
					end:      NewDate(2023, 7, 21),
					on:       true,
					interval: Interval{Skip: 2, Freq: NewWeekly(time.Monday, time.Friday)},
				},
			},
			next: 4,
			want: []Class{
				{Index: 1, Time: NewDate(2023, 6, 26).ToTime()},
				{Index: 2, Time: NewDate(2023, 6, 30).ToTime()},
				{Index: 3, Time: NewDate(2023, 7, 17).ToTime()},
				{Index: 4, Time: NewDate(2023, 7, 21).ToTime()},
			},
			wantErr: true,
		}, {
			name: "Test: Repeat Monthly",
			section: Section{
				start: NewDate(2023, 6, 2),
				repeat: Repetition{
					end:      NewDate(2023, 9, 21),
					on:       true,
					interval: Interval{Skip: 0, Freq: NewMonthly(1)},
				},
			},
			next: 2,
			want: []Class{
				{Index: 1, Time: NewDate(2023, 7, 1).ToTime()},
				{Index: 2, Time: NewDate(2023, 8, 1).ToTime()},
			},
			wantErr: false,
		}, {
			name: "Test: Repeat Monthly Advanced 1",
			section: Section{
				start: NewDate(2023, 6, 2),
				repeat: Repetition{
					end:      NewDate(2024, 3, 1),
					on:       true,
					interval: Interval{Skip: 3, Freq: NewMonthly(1, 31)},
				},
			},
			next: 5,
			want: []Class{
				{Index: 1, Time: NewDate(2023, 7, 1).ToTime()},
				{Index: 2, Time: NewDate(2023, 7, 31).ToTime()},
				{Index: 3, Time: NewDate(2023, 11, 1).ToTime()},
				{Index: 4, Time: NewDate(2023, 11, 30).ToTime()},
				{Index: 5, Time: NewDate(2024, 3, 1).ToTime()},
			},
			wantErr: true,
		}, {
			name: "Test: Repeat Monthly Advanced 2",
			section: Section{
				start: NewDate(2023, 12, 2),
				repeat: Repetition{
					end:      NewDate(2025, 3, 1),
					on:       true,
					interval: Interval{Skip: 1, Freq: NewMonthly(31)},
				},
			},
			next: 3,
			want: []Class{
				{Index: 1, Time: NewDate(2023, 12, 31).ToTime()},
				{Index: 2, Time: NewDate(2024, 2, 29).ToTime()},
				{Index: 3, Time: NewDate(2024, 4, 30).ToTime()},
			},
			wantErr: false,
		},
	}
	// t.Error(unsafe.Sizeof(tests[0].section))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.section
			got, _, err := c.Next(tt.section.start, tt.next)
			if (err != nil) != tt.wantErr {
				t.Errorf("Section.Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Section.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}
