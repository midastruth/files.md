// Parses Garmin export JSON files and prints a 10-day journal summary.
//
// Usage: go run ./cmd/garmin/garmin.go <path-to-garmin-export-dir>
// Example: go run ./cmd/garmin/garmin.go ./garmin
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type day struct {
	Date     time.Time
	Sleep    sleepData
	HRV      int
	RHR      int
	HasHR    bool
	Workouts []workout
}

type sleepData struct {
	Score    int
	TotalMin int
	Present  bool
}

type workout struct {
	Name        string
	Activity    string
	DurationMin int
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	days := map[string]*day{}
	dayOf := func(date string) *day {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil
		}
		if d, ok := days[date]; ok {
			return d
		}
		d := &day{Date: t}
		days[date] = d
		return d
	}

	wellnessDir := filepath.Join(dir, "DI-Connect-Wellness")
	fitnessDir := filepath.Join(dir, "DI-Connect-Fitness")

	// Parse sleep data
	forEachJSON(wellnessDir, "sleepData", func(raw json.RawMessage) {
		var entries []struct {
			CalendarDate    string `json:"calendarDate"`
			DeepSeconds     int    `json:"deepSleepSeconds"`
			LightSeconds    int    `json:"lightSleepSeconds"`
			REMSeconds      int    `json:"remSleepSeconds"`
			SleepScores     *struct {
				OverallScore int `json:"overallScore"`
			} `json:"sleepScores"`
		}
		json.Unmarshal(raw, &entries)
		for _, e := range entries {
			if e.CalendarDate == "" || e.SleepScores == nil {
				continue
			}
			d := dayOf(e.CalendarDate)
			if d == nil {
				continue
			}
			totalSec := e.DeepSeconds + e.LightSeconds + e.REMSeconds
			d.Sleep = sleepData{
				Score:    e.SleepScores.OverallScore,
				TotalMin: totalSec / 60,
				Present:  true,
			}
		}
	})

	// Parse health status (HRV, RHR)
	forEachJSON(wellnessDir, "healthStatusData", func(raw json.RawMessage) {
		var entries []struct {
			CalendarDate string `json:"calendarDate"`
			Metrics      []struct {
				Type  string   `json:"type"`
				Value *float64 `json:"value"`
			} `json:"metrics"`
		}
		json.Unmarshal(raw, &entries)
		for _, e := range entries {
			d := dayOf(e.CalendarDate)
			if d == nil {
				continue
			}
			for _, m := range e.Metrics {
				if m.Value == nil {
					continue
				}
				switch m.Type {
				case "HRV":
					d.HRV = int(*m.Value)
					d.HasHR = true
				case "HR":
					d.RHR = int(*m.Value)
					d.HasHR = true
				}
			}
		}
	})

	// Parse activities
	forEachJSON(fitnessDir, "summarizedActivities", func(raw json.RawMessage) {
		var wrapper []struct {
			Activities []struct {
				Name         string  `json:"name"`
				ActivityType string  `json:"activityType"`
				StartTimeGmt float64 `json:"startTimeGmt"`
				Duration     float64 `json:"duration"` // milliseconds
			} `json:"summarizedActivitiesExport"`
		}
		json.Unmarshal(raw, &wrapper)
		for _, w := range wrapper {
			for _, a := range w.Activities {
				t := time.UnixMilli(int64(a.StartTimeGmt))
				date := t.Format("2006-01-02")
				d := dayOf(date)
				if d == nil {
					continue
				}
				d.Workouts = append(d.Workouts, workout{
					Name:        a.Name,
					Activity:    formatActivity(a.ActivityType),
					DurationMin: int(a.Duration / 60000),
				})
			}
		}
	})

	// Sort days descending, print last 10
	var sorted []*day
	for _, d := range days {
		sorted = append(sorted, d)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.After(sorted[j].Date)
	})

	n := 10
	if len(sorted) < n {
		n = len(sorted)
	}
	for _, d := range sorted[:n] {
		printDay(d)
		fmt.Println()
	}
}

func printDay(d *day) {
	fmt.Printf("#### %d %s, %s\n", d.Date.Day(), d.Date.Format("January"), d.Date.Weekday())

	if d.Sleep.Present {
		h := d.Sleep.TotalMin / 60
		m := d.Sleep.TotalMin % 60
		fmt.Printf("- Sleep: %d%%, %dh %02dm\n", d.Sleep.Score, h, m)
	}

	if d.HasHR {
		fmt.Printf("- HRV %d, RHR %d\n", d.HRV, d.RHR)
	}

	for _, w := range d.Workouts {
		fmt.Printf("- %s: %s %dm\n", w.Activity, w.Name, w.DurationMin)
	}
}

func formatActivity(activityType string) string {
	r := strings.NewReplacer("_", " ", "v2", "")
	name := r.Replace(activityType)
	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + name[1:]
	}
	return strings.TrimSpace(name)
}

// forEachJSON finds all JSON files in dir matching suffix, reads and calls fn.
func forEachJSON(dir, suffix string, fn func(json.RawMessage)) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		if !strings.Contains(e.Name(), suffix) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		fn(json.RawMessage(data))
	}
}
