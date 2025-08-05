package tests

import (
	"context"
	"time"

	"github.com/silvasuite/silva/trunk"
	"github.com/stretchr/testify/mock"
)

type BookingProviderMock struct {
	mock.Mock

	Slots    []trunk.Slot
	Bookings []trunk.Booking
}

func (b *BookingProviderMock) GenerateSlots(from, to time.Time, days int, slotLength time.Duration) {
	b.Slots = []trunk.Slot{}
	for c := 0; c < days; c++ {
		start := from.Add(time.Duration(c) * 24 * time.Hour)
		end := to.Add(time.Duration(c) * 24 * time.Hour)

		for t := start; t.Before(end); t = t.Add(slotLength) {
			slot := trunk.Slot{
				ID:    t.Format(time.RFC3339),
				Start: t,
				End:   t.Add(slotLength),
			}
			b.Slots = append(b.Slots, slot)
		}
	}
}

func (b *BookingProviderMock) ListAvailabeSlots(ctx context.Context, start, end time.Time) ([]trunk.Slot, error) {
	return b.Slots, nil
}

func (b *BookingProviderMock) BookSlot(ctx context.Context, slotID string) error {
	var err = trunk.ErrCannotBookEvent
	for i := range b.Slots {
		if b.Slots[i].ID == slotID {
			b.Bookings = append(b.Bookings, trunk.Booking(b.Slots[i]))
		}
		err = nil
	}

	return err

}

func (b *BookingProviderMock) CanBookSlot(ctx context.Context, slotID string) (bool, error) {
	args := b.Called(ctx, slotID)
	return args.Bool(0), args.Error(1)
}

func (b *BookingProviderMock) CancelBooking(ctx context.Context, bookingID string) error {
	for i := range b.Bookings {
		if b.Bookings[i].ID == bookingID {
			b.Bookings = append(b.Bookings[:i], b.Bookings[i+1:]...)
			return nil
		}
	}
	return trunk.ErrBookingNotFound
}

func (b *BookingProviderMock) ListBookings(ctx context.Context) ([]trunk.Booking, error) {
	return b.Bookings, nil
}
