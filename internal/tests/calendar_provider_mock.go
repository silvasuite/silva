package tests

import (
	"context"
	"time"

	"github.com/silvasuite/silva/trunk"
	"github.com/stretchr/testify/mock"
)

type CalendarProviderMock struct {
	mock.Mock
}

func (c *CalendarProviderMock) IsAvailable(ctx context.Context, from, to time.Time) (bool, error) {
	args := c.Called(ctx, from, to)
	return args.Bool(0), args.Error(1)
}

func (c *CalendarProviderMock) SaveBooked(ctx context.Context, e trunk.Slot) error {
	args := c.Called(ctx, e)
	return args.Error(0)
}
