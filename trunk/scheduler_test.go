package trunk_test

import (
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/silvasuite/silva/internal/tests"
	"github.com/silvasuite/silva/trunk"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SchedulerTestSuite struct {
	suite.Suite

	scheduler       *trunk.Scheduler
	bookingProvider *tests.BookingProviderMock
}

func (s *SchedulerTestSuite) SetupTest() {
	s.bookingProvider = &tests.BookingProviderMock{}

	preferred := time.Date(2020, 1, 1, 14, 30, 0, 0, time.UTC)
	// Generate mock slots for testing
	start := time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC)
	end := start.Add(10 * time.Hour)

	s.bookingProvider.GenerateSlots(start, end, 5, 1*time.Hour)

	s.Len(s.bookingProvider.Slots, 50, "Expected 50 slots to be generated")

	s.scheduler = &trunk.Scheduler{
		Calendar:          nil, // Mock or real calendar provider can be set here
		Booking:           s.bookingProvider,
		PreferredTime:     trunk.NewPreferredBookTime(preferred, time.Hour, nil),
		LookAhead:         7 * 24 * time.Hour,
		MaxBookingsPerDay: 1,
	}
}

func (s *SchedulerTestSuite) TestRunOneCycle() {
	s.Len(s.bookingProvider.Bookings, 0)
	s.bookingProvider.On("CanBookSlot", mock.Anything, mock.Anything).Return(true, nil)

	err := s.scheduler.Run(s.T().Context())
	s.NoError(err, "Expected no error during scheduler run")

	s.Len(s.bookingProvider.Bookings, 5)

	// s.bookingProvider.Bookings[0].Start
	for i := range s.bookingProvider.Bookings {
		s.Equal(s.bookingProvider.Bookings[i].Start.Format("15:04:05"), "15:00:00")
	}
}

func (s *SchedulerTestSuite) TestCanBookOnlyTwoDays() {
	s.Len(s.bookingProvider.Bookings, 0)
	s.bookingProvider.On("CanBookSlot", mock.Anything, "2020-01-01T15:00:00Z").Return(true, nil)
	s.bookingProvider.On("CanBookSlot", mock.Anything, "2020-01-02T15:00:00Z").Return(true, nil)
	s.bookingProvider.On("CanBookSlot", mock.Anything, mock.Anything).Return(false, nil)

	err := s.scheduler.Run(s.T().Context())
	s.NoError(err, "Expected no error during scheduler run")

	s.Len(s.bookingProvider.Bookings, 2)

	// s.bookingProvider.Bookings[0].Start
	for i := range s.bookingProvider.Bookings {
		s.Equal(s.bookingProvider.Bookings[i].Start.Format("15:04:05"), "15:00:00")
	}
}

func (s *SchedulerTestSuite) TestErrorArePropagatedFromBookingProvider() {
	s.Len(s.bookingProvider.Bookings, 0)
	s.bookingProvider.On("CanBookSlot", mock.Anything, mock.Anything).Return(false, trunk.ErrCannotBookEvent)

	err := s.scheduler.Run(s.T().Context())
	s.Error(err, "Expected error when booking is not possible")
	s.True(errors.Is(err, trunk.ErrCannotBookEvent), "Expected specific booking error")

	s.Len(s.bookingProvider.Bookings, 0)
}

func (s *SchedulerTestSuite) TestCustomizePreferredTime() {
	s.Len(s.bookingProvider.Bookings, 0)

	preferred := time.Date(2020, 1, 1, 16, 0, 0, 0, time.UTC)
	byDay := map[time.Weekday]time.Time{
		time.Wednesday: preferred.Add(30 * time.Minute),
		time.Thursday:  preferred.Add(2 * time.Hour),
		time.Friday:    preferred.Add(3 * time.Hour),
	}

	s.scheduler.PreferredTime = trunk.NewPreferredBookTime(preferred, time.Hour, byDay)

	s.bookingProvider.On("CanBookSlot", mock.Anything, mock.Anything).Return(true, nil)

	err := s.scheduler.Run(s.T().Context())
	s.NoError(err, "Expected no error during scheduler run")

	s.Len(s.bookingProvider.Bookings, 5)

	sort.Slice(s.bookingProvider.Bookings, func(i, j int) bool {
		return s.bookingProvider.Bookings[i].Start.Before(s.bookingProvider.Bookings[j].Start)
	})

	expectedTimes := []string{
		"17:00:00", // Wednesday
		"18:00:00", // Thursday
		"19:00:00", // Friday
		"16:00:00", // Saturday
		"16:00:00", // Sunday
	}

	for i := range s.bookingProvider.Bookings {
		s.Equal(s.bookingProvider.Bookings[i].Start.Format("15:04:05"), expectedTimes[i])
	}
}

func (s *SchedulerTestSuite) TestWithCalendarProvider() {
	s.Len(s.bookingProvider.Bookings, 0)

	t1 := time.Date(2020, 1, 2, 15, 0, 0, 0, time.UTC)
	t2 := time.Date(2020, 1, 2, 16, 0, 0, 0, time.UTC)
	t3 := time.Date(2020, 1, 4, 15, 0, 0, 0, time.UTC)
	t4 := time.Date(2020, 1, 4, 16, 0, 0, 0, time.UTC)

	calendarMock := &tests.CalendarProviderMock{}
	calendarMock.On("IsAvailable", mock.Anything, t1, t2).Return(false, nil)
	calendarMock.On("IsAvailable", mock.Anything, t3, t4).Return(false, nil)
	calendarMock.On("IsAvailable", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	calendarMock.On("SaveBooked", mock.Anything, mock.Anything).Return(nil)

	s.scheduler.Calendar = calendarMock
	s.bookingProvider.On("CanBookSlot", mock.Anything, mock.Anything).Return(true, nil)

	err := s.scheduler.Run(s.T().Context())
	s.NoError(err, "Expected no error during scheduler run")

	s.Len(s.bookingProvider.Bookings, 3)

	sort.Slice(s.bookingProvider.Bookings, func(i, j int) bool {
		return s.bookingProvider.Bookings[i].Start.Before(s.bookingProvider.Bookings[j].Start)
	})

	expectedWeekDays := []time.Weekday{
		time.Wednesday, // 2020-01-01
		time.Friday,    // 2020-01-03
		time.Sunday,    // 2020-01-05
	}

	for i := range s.bookingProvider.Bookings {
		s.Equal(s.bookingProvider.Bookings[i].Start.Format("15:04:05"), "15:00:00")
		s.Equal(s.bookingProvider.Bookings[i].Start.Weekday(), expectedWeekDays[i])
	}

	calendarMock.AssertExpectations(s.T())
}

func TestScheduler(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}
