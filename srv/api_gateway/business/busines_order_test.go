package business

import (
	"fmt"
	"testing"

	"github.com/alimitedgroup/MVP/common/dto"
	response "github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

func TestGetTransfers(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	orderMock.EXPECT().GetAllTransfers().Return([]response.TransferInfo{
		{
			Status:     "Created",
			TransferID: "1",
			SenderID:   "1",
			ReceiverID: "2",
			Goods: []response.TransferInfoGood{
				{
					GoodID:   "1",
					Quantity: 1,
				},
			},
		},
		{
			Status:     "Filled",
			TransferID: "2",
			SenderID:   "2",
			ReceiverID: "1",
			Goods: []response.TransferInfoGood{
				{
					GoodID:   "2",
					Quantity: 10,
				},
			},
		},
	}, nil)

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	transfers, err := business.GetTransfers()
	require.NoError(t, err)
	require.Len(t, transfers, 2)
	require.ElementsMatch(t, []dto.Transfer{
		{
			Status:     "Created",
			TransferID: "1",
			SenderID:   "1",
			ReceiverID: "2",
			Goods: map[string]int64{
				"1": 1,
			},
		},
		{
			Status:     "Filled",
			TransferID: "2",
			SenderID:   "2",
			ReceiverID: "1",
			Goods: map[string]int64{
				"2": 10,
			},
		},
	}, transfers)
}

func TestGetTransfersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	orderMock.EXPECT().GetAllTransfers().Return(nil, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	transfers, err := business.GetTransfers()
	require.Nil(t, transfers)
	require.ErrorIs(t, err, ErrorGetTransfers)
}

func TestGetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	orderMock.EXPECT().GetAllOrders().Return([]response.OrderInfo{
		{
			Status:   "Created",
			OrderID:  "1",
			Name:     "Order 1",
			FullName: "Mario Rossi",
			Address:  "Via Roma 1",
			Goods: []response.OrderInfoGood{
				{
					GoodID:   "1",
					Quantity: 1,
				},
			},
		},
		{
			Status:   "Filled",
			OrderID:  "2",
			Name:     "Order 2",
			FullName: "Luigi Verdi",
			Address:  "Via Milano 2",
			Goods: []response.OrderInfoGood{
				{
					GoodID:   "2",
					Quantity: 10,
				},
			},
		},
	}, nil)

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	orders, err := business.GetOrders()
	require.NoError(t, err)
	require.Len(t, orders, 2)
	require.ElementsMatch(t, []dto.Order{
		{
			Status:   "Created",
			OrderID:  "1",
			Name:     "Order 1",
			FullName: "Mario Rossi",
			Address:  "Via Roma 1",
			Goods: map[string]int64{
				"1": 1,
			},
		},
		{
			Status:   "Filled",
			OrderID:  "2",
			Name:     "Order 2",
			FullName: "Luigi Verdi",
			Address:  "Via Milano 2",
			Goods: map[string]int64{
				"2": 10,
			},
		},
	}, orders)
}

func TestGetOrdersError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	orderMock.EXPECT().GetAllOrders().Return(nil, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	orders, err := business.GetOrders()
	require.Nil(t, orders)
	require.ErrorIs(t, err, ErrorGetOrders)
}

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	orderMock.EXPECT().CreateOrder(gomock.Any()).Return(response.OrderCreateInfo{OrderID: "1"}, nil)

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	orderId, err := business.CreateOrder("1", "Mario Rossi", "Via Roma 1", map[string]int64{"id1": 1})
	require.NoError(t, err)
	require.Equal(t, "1", orderId)
}

func TestCreateTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	orderMock.EXPECT().CreateTransfer(gomock.Any()).Return(response.TransferCreateInfo{TransferID: "1"}, nil)

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	transferId, err := business.CreateTransfer("1", "2", map[string]int64{"id1": 1})
	require.NoError(t, err)
	require.Equal(t, "1", transferId)
}

func TestCreateTransferError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	orderMock.EXPECT().CreateTransfer(gomock.Any()).Return(response.TransferCreateInfo{}, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	transferId, err := business.CreateTransfer("1", "2", map[string]int64{"id1": 1})
	require.Empty(t, transferId)
	require.ErrorIs(t, err, ErrorCreateTransfer)
}

func TestCreateOrderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	auth := NewMockAuthenticationPortOut(ctrl)
	catalog := NewMockCatalogPortOut(ctrl)
	orderMock := NewMockOrderPortOut(ctrl)

	orderMock.EXPECT().CreateOrder(gomock.Any()).Return(response.OrderCreateInfo{}, fmt.Errorf("some error"))

	business := NewBusiness(auth, catalog, orderMock, zaptest.NewLogger(t))
	orderId, err := business.CreateOrder("1", "Mario Rossi", "Via Roma 1", map[string]int64{"id1": 1})
	require.Empty(t, orderId)
	require.ErrorIs(t, err, ErrorCreateOrder)
}
