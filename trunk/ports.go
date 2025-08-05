package trunk

import (
	"context"
	"errors"
	"time"
)

var (
	ErrCannotBookEvent = errors.New("cannot book event")
	ErrBookingNotFound = errors.New("booking not found")
)

// BookingProvider books a class/slot/event with an external service.
type BookingProvider interface {
	ListAvailabeSlots(ctx context.Context, start, end time.Time) ([]Slot, error)
	BookSlot(ctx context.Context, slotID string) error
	CanBookSlot(ctx context.Context, slotID string) (bool, error)
	CancelBooking(ctx context.Context, bookingID string) error
	ListBookings(ctx context.Context) ([]Booking, error)
}

// CalendarProvider exposes read/write operations on a calendar.
type CalendarProvider interface {
	IsAvailable(ctx context.Context, from, to time.Time) (bool, error)
	SaveBooked(ctx context.Context, e Slot) error
}
