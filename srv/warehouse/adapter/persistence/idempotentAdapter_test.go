package persistence

import (
	"testing"

	"github.com/alimitedgroup/MVP/srv/warehouse/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestIdempotentAdapterGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIIdempotentRepository(ctrl)

	mock.EXPECT().IsAlreadyProcessed(gomock.Eq("event"), gomock.Eq("id")).Return(false)

	adapter := NewIDempotentAdapter(mock)
	ret := adapter.IsAlreadyProcessed(port.IdempotentCmd{
		Event: "event",
		ID:    "id",
	})
	require.False(t, ret)
}

func TestIdempotentAdapterSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockIIdempotentRepository(ctrl)

	mock.EXPECT().SaveEventID(gomock.Eq("event"), gomock.Eq("id")).Return()

	adapter := NewIDempotentAdapter(mock)
	adapter.SaveEventID(port.IdempotentCmd{
		Event: "event",
		ID:    "id",
	})
}
