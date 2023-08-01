package utils

import (
	"reflect"
	"testing"
	"time"

	"github.com/joaopegoraro/ahpsico-go/utils"
)

func TestGetStartOfDay(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{
		{
			name: "Should return 0 hour, minute and seconds of the same day of provided time",
			t:    now,
			want: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.GetStartOfDay(tt.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStartOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEndOfDay(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{
		{
			name: "Should return 0 hour, minute and seconds of the next day of provided time",
			t:    now,
			want: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(time.Hour * 24),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.GetEndOfDay(tt.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEndOfDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
