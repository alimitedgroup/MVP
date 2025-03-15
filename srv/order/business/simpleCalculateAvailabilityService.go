package business

import (
	"context"

	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type SimpleCalculateAvailabilityService struct {
}

func NewSimpleCalculateAvailabilityService() *SimpleCalculateAvailabilityService {
	return &SimpleCalculateAvailabilityService{}
}

func (s *SimpleCalculateAvailabilityService) GetAvailable(ctx context.Context, cmd port.CalculateAvailabilityCmd) (port.CalculateAvailabilityResponse, error) {

	return port.CalculateAvailabilityResponse{}, nil
}
