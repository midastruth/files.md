package userconfig

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_strDuration_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		expect  time.Duration
		value   string
		wantErr bool
	}{
		{10 * time.Minute, `10`, false},
		{15 * time.Second, `0.25`, false},
		{DefaultConfig.PomodoroDuration(), `0`, false},
		{0, `"1"`, true},
		{0, `"x"`, true},
		{0, `x`, true},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			r := require.New(t)
			c := &Config{}
			data := []byte(fmt.Sprintf(`{"pomodoroDurationMinute": %v}`, tt.value))
			err := json.Unmarshal(data, c)
			if tt.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
				r.Equal(tt.expect, c.PomodoroDuration())
			}
		})
	}
}

func TestConfig_SetPomodoroDuration(t *testing.T) {

	tests := []struct {
		value   time.Duration
		wantErr bool
	}{
		{10 * time.Minute, false},
		{100 * time.Millisecond, false},
		{4 * time.Hour, false},
		{0, true},
		{24 * time.Hour, false},
		{24*time.Hour + time.Second, true},
	}
	for _, tt := range tests {
		t.Run(tt.value.String(), func(t *testing.T) {
			c := &Config{}
			if err := c.SetPomodoroDuration(tt.value); (err != nil) != tt.wantErr {
				t.Errorf("SetPomodoroDuration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
