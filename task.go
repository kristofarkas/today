package today

import (
	"fmt"
	"time"
)

// Task is a unit based task with an end date.
type Task struct {
	Name    string
	Total   int
	Current int
	Changes []Mod
	EndDate time.Time
}

type Mod struct {
	Date  time.Time
	Value int
}

// SetCurrent sets the current unit index of the task.
func (t *Task) SetCurrent(current int) error {
	if current > t.Total {
		return fmt.Errorf("current exceeds total")
	}

	t.Current = current
	t.Changes = append([]Mod{Mod{time.Now(), current}}, t.Changes...)
	return nil
}

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func (t *Task) UnitsPerDay() int {

	daysToGo := int(t.EndDate.Sub(beginningOfDay(time.Now())).Hours()) / 24

	if daysToGo < 1 {
		return t.Total - t.EffectiveCurrent()
	}

	return (t.Total - t.EffectiveCurrent()) / daysToGo
}

func (t *Task) Progress() int {
	return int(100 * float64(t.Current) / float64(t.Total))
}

func (t *Task) ToGo() int {
	if toGo := t.Today() - t.Current; toGo > 0 {
		return toGo
	}
	return 0
}

func (t *Task) EffectiveCurrent() int {
	bot := beginningOfDay(time.Now())

	if bot.After(t.EndDate) {
		return t.Total
	}

	// If there is no change made before today, we set the effective current value to zero.
	latestCurrent := 0

	for _, mod := range t.Changes {
		if mod.Date.Before(bot) {
			latestCurrent = mod.Value
			break
		}
	}

	return latestCurrent
}

// Today is the desired unit index to be achieved today.
func (t *Task) Today() int {
	return t.EffectiveCurrent() + t.UnitsPerDay()
}
