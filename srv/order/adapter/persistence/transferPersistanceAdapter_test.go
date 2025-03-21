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
				GoodID:   "1",
				Quantity: 10,
			},
			{
				GoodID:   "2",
				Quantity: 10,
			},
		},
	}
	adapter.ApplyTransferUpdate(cmd)
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
	require.Equal(t, transfer.ID, "1")
}

func TestTransferPersistenceAdapterGetTransferNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockITransferRepository(ctrl)

	mock.EXPECT().GetTransfer(gomock.Any()).Return(Transfer{}, ErrTransferNotFound)

	adapter := NewTransferPersistanceAdapter(mock)

	_, err := adapter.GetTransfer("1")
	require.Error(t, err, ErrTransferNotFound)
}

func TestTransferPersistenceAdapterGetAllTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockITransferRepository(ctrl)

	mock.EXPECT().GetTransfers().Return([]Transfer{{
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
	}})

	adapter := NewTransferPersistanceAdapter(mock)

	transfers := adapter.GetAllTransfer()
	require.Len(t, transfers, 1)
	require.Equal(t, transfers[0].ID, "1")
}

func TestTransferPersistenceAdapterSetComplete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockITransferRepository(ctrl)

	mock.EXPECT().SetComplete(gomock.Any()).Return(nil)

	adapter := NewTransferPersistanceAdapter(mock)

	err := adapter.SetComplete("1")
	require.NoError(t, err)
}

func TestTransferPersistenceAdapterSetCompleteErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockITransferRepository(ctrl)

	mock.EXPECT().SetComplete(gomock.Any()).Return(ErrTransferNotFound)

	adapter := NewTransferPersistanceAdapter(mock)

	err := adapter.SetComplete("1")
	require.Error(t, err, ErrTransferNotFound)
}

func TestTransferPersistenceAdapterIncrementLinkedStockUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockITransferRepository(ctrl)

	mock.EXPECT().IncrementLinkedStockUpdate(gomock.Any()).Return(nil)

	adapter := NewTransferPersistanceAdapter(mock)

	err := adapter.IncrementLinkedStockUpdate("1")
	require.NoError(t, err)
}

func TestTransferPersistenceAdapterIncrementLinkedStockUpdateErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := NewMockITransferRepository(ctrl)

	mock.EXPECT().IncrementLinkedStockUpdate(gomock.Any()).Return(ErrTransferNotFound)

	adapter := NewTransferPersistanceAdapter(mock)

	err := adapter.IncrementLinkedStockUpdate("1")
	require.Error(t, err, ErrTransferNotFound)
}
