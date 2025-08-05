package trunk

import (
	"context"
	"fmt"
	"sort"
	"time"
)

// Service orchestrates calendar reading, decision logic and booking.
type Service struct {
	Calendar      CalendarProvider
	Booking       BookingProvider
	PreferredTime PreferredBookTime

	// Optional: window in which to look ahead (default 7 days).
	LookAhead         time.Duration
	MaxBookingsPerDay int
}

// Run executes one scheduling cycle.
func (s *Service) Run(ctx context.Context) error {
	if s.LookAhead == 0 {
		s.LookAhead = 7 * 24 * time.Hour
	}

	now := time.Now()
	from, to := now, now.Add(s.LookAhead)

	bookableSlots, err := s.Booking.ListAvailabeSlots(ctx, from, to)
	if err != nil {
		return err
	}

	sort.Slice(bookableSlots, func(i, j int) bool {
		return bookableSlots[i].Start.Before(bookableSlots[j].Start)
	})

	bookings, err := s.Booking.ListBookings(ctx)
	if err != nil {
		return err
	}
	bookingByDay := make(map[time.Weekday]int)
	for _, booking := range bookings {
		day := booking.Start.Weekday()
		bookingByDay[day] = bookingByDay[day] + 1
	}

	for _, slot := range bookableSlots {
		if bookingByDay[slot.Start.Weekday()] >= s.MaxBookingsPerDay {
			// reached maximum bookings for this day
			continue
		}

		day := slot.Start.Weekday()
		preferredTimeForDay := s.PreferredTime.GetPreferredBookTime(day)
		maxPreferredTime := preferredTimeForDay.Add(s.PreferredTime.maxDelay)

		// Simple hack: formatted are lexicographically ordered in the same way as time
		slotTime := slot.Start.Format("15:04:05")
		preferredTime := preferredTimeForDay.Format("15:04:05")
		maxPreferredTimeStr := maxPreferredTime.Format("15:04:05")
		if slotTime >= preferredTime && slotTime <= maxPreferredTimeStr {
			canBook, err := s.Booking.CanBookSlot(ctx, slot.ID)
			if err != nil {
				return err
			}
			if s.Calendar != nil {
				available, err := s.Calendar.IsAvailable(ctx, slot.Start, slot.End)
				if err != nil {
					return err
				}
				canBook = canBook && available
			}

			if canBook {
				err := s.Booking.BookSlot(ctx, slot.ID)
				if err != nil {
					return err
				}
				if s.Calendar != nil {
					err = s.Calendar.SaveBooked(ctx, slot)
					if err != nil {
						return err
					}
				}
				bookingByDay[slot.Start.Weekday()]++

				fmt.Printf("Booked slot %s for %s\n", slot.ID, slot.Start.Format(time.RFC3339))
			}
		}
	}

	return nil
}
