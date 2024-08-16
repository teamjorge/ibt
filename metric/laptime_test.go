package metric

import "testing"

func TestLapTime_ToString(t *testing.T) {
	tests := []struct {
		name string
		l    LapTime
		want string
	}{
		{
			name: "regular",
			l:    76.932,
			want: "01:16.932",
		},
		{
			name: "less than a minute",
			l:    48.001,
			want: "00:48.001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.ToString(); got != tt.want {
				t.Errorf("LapTime.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
