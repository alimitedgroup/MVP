package persistence

import (
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransferPersistenceAdapterApplyTransferUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockITransferRepository(ctrl)

	mock.EXPECT().SetTransfer(gomock.Any(), gomock.Any()).Return(false)

	adapter := NewTransferPersistanceAdapter(mock)

	cmd := port.ApplyTransferUpdateCmd{
		Id:            "1",
		Status:        "Created",
		SenderId:      "1",
		ReceiverId:    "2",
		ReservationId: "",
		CreationTime:  time.Now().UnixMilli(),
		UpdateTime:    time.Now().UnixMilli(),
		Goods: []model.GoodStock{
			{
				ID:       "1",
				Quantity: 10,
			},
			{
				ID:       "2",
				Quantity: 10,
			},
		},
	}
	err := adapter.ApplyTransferUpdate(cmd)
	require.NoError(t, err)
}

func TestTransferPersistenceAdapterGetTransferExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockITransferRepository(ctrl)

	mock.EXPECT().GetTransfer(gomock.Any()).Return(Transfer{
		ID:                "1",
		Status:            "Created",
		SenderID:          "1",
		ReceiverID:        "2",
		LinkedStockUpdate: 0,
		Goods: []TransferUpdateGood{
			{
				GoodID:   "1",
				Quantity: 10,
			},
			{
				GoodID:   "2",
				Quantity: 10,
			},
		},
		ReservationId: "",
		UpdateTime:    time.Now().UnixMilli(),
		CreationTime:  time.Now().UnixMilli(),
	}, nil)

	adapter := NewTransferPersistanceAdapter(mock)

	transfer, err := adapter.GetTransfer("1")
	require.NoError(t, err)
	require.Equal(t, transfer.Id, model.TransferID("1"))
}

func TestTransferPersistenceAdapterGetTransferNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockITransferRepository(ctrl)

	mock.EXPECT().GetTransfer(gomock.Any()).Return(Transfer{}, ErrTransferNotFound)

	adapter := NewTransferPersistanceAdapter(mock)

	_, err := adapter.GetTransfer("1")
	require.Error(t, err, ErrTransferNotFound)
}
