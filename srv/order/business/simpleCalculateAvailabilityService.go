package business

import (
	"context"
	"errors"

	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type SimpleCalculateAvailabilityService struct {
}

func NewSimpleCalculateAvailabilityService() *SimpleCalculateAvailabilityService {
	return &SimpleCalculateAvailabilityService{}
}

func (s *SimpleCalculateAvailabilityService) GetAvailable(ctx context.Context, cmd port.CalculateAvailabilityCmd) (port.CalculateAvailabilityResponse, error) {

	return port.CalculateAvailabilityResponse{}, errors.New("no stock")
	// return port.CalculateAvailabilityResponse{}, nil
}
