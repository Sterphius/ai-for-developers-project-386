// Package slots computes the 14-day availability grid from the booking window.
package slots

import (
	"time"

	"calendar-booking/server/internal/model"
)

// WindowDays is the fixed availability window: 14 days from now, 24/7.
const WindowDays = 14

// gridUnit is the granularity the window start is snapped to. Snapping to the
// minute gives a stable grid origin across requests (sub-second jitter in "now"
// would otherwise make slot generation and booking validation disagree).
const gridUnit = time.Minute

// WindowStart returns the grid origin: now rounded up to the next minute.
// This is the earliest bookable instant and the anchor for grid alignment.
func WindowStart(now time.Time) time.Time {
	t := now.UTC().Truncate(gridUnit)
	if t.Before(now.UTC()) {
		t = t.Add(gridUnit)
	}
	return t
}

// WindowEnd returns the end of the booking window: WindowStart + 14 days.
func WindowEnd(now time.Time) time.Time {
	return WindowStart(now).Add(WindowDays * 24 * time.Hour)
}

// IsAligned reports whether start sits on the grid: a whole number of steps
// from the window start and not before it.
func IsAligned(start, now time.Time, step time.Duration) bool {
	if step <= 0 {
		return false
	}
	ws := WindowStart(now)
	if start.Before(ws) {
		return false
	}
	delta := start.Sub(ws)
	return delta%step == 0
}

// InWindow reports whether [start, start+step) fits within [windowStart, windowEnd].
func InWindow(start, now time.Time, step time.Duration) bool {
	if start.Before(WindowStart(now)) {
		return false
	}
	return !start.Add(step).After(WindowEnd(now))
}

// Overlaps reports whether [aStart, aEnd) intersects [bStart, bEnd).
// Adjacent intervals (touching at a boundary) do not overlap.
func Overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}

// Generate builds the slot grid for an event type within [now, now+14d].
// Each slot is step=durationMinutes long; occupied slots are returned with
// Available=false. Optional from/to clamp the visible range to the window.
func Generate(et model.EventType, bookings []model.Booking, now time.Time, from, to *time.Time) []model.Slot {
	step := time.Duration(et.DurationMinutes) * time.Minute
	if step <= 0 {
		return nil
	}

	windowStart := WindowStart(now)
	windowEnd := WindowEnd(now)

	rangeStart := windowStart
	if from != nil && from.After(rangeStart) {
		rangeStart = *from
	}
	rangeEnd := windowEnd
	if to != nil && to.Before(rangeEnd) {
		rangeEnd = *to
	}

	// Align the first slot up to a grid boundary relative to windowStart.
	first := windowStart
	if rangeStart.After(windowStart) {
		delta := rangeStart.Sub(windowStart)
		steps := delta / step
		if delta%step != 0 {
			steps++
		}
		first = windowStart.Add(steps * step)
	}

	// Pre-filter bookings to those relevant to this event type's grid.
	var slots []model.Slot
	for s := first; !s.Add(step).After(rangeEnd); s = s.Add(step) {
		end := s.Add(step)
		available := true
		for _, b := range bookings {
			if Overlaps(s, end, b.Start, b.End) {
				available = false
				break
			}
		}
		slots = append(slots, model.Slot{Start: s, End: end, Available: available})
	}
	return slots
}
