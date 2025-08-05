package trunk

import "time"

// Event represents a calendar entry (class, meeting, etc.).
type Slot struct {
	ID          string
	Title       string
	Description string
	Start       time.Time
	End         time.Time
}

func (e *Slot) Overlaps(other Slot) bool {
	return e.Start.Before(other.End) && e.End.After(e.Start)
}

// Booking confirms a Slot for the user.
type Booking Slot

type PreferredBookTime struct {
	maxDelay time.Duration
	def      time.Time
	byDay    map[time.Weekday]time.Time
}

func NewPreferredBookTime(def time.Time, maxDelay time.Duration, byDay map[time.Weekday]time.Time) PreferredBookTime {
	return PreferredBookTime{
		maxDelay: maxDelay,
		def:      def,
		byDay:    byDay,
	}
}

func (p *PreferredBookTime) GetPreferredBookTime(day time.Weekday) time.Time {
	if t, ok := p.byDay[day]; ok {
		return t
	}
	return p.def
}
