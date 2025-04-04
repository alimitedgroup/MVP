package controller

import (
	"cmp"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/dto/response"
	"github.com/alimitedgroup/MVP/common/lib"
	"github.com/alimitedgroup/MVP/common/lib/broker"
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
	"golang.org/x/exp/slices"
)

func TestTransferControllerCreateTransfer(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	createTransferUseCaseMock := NewMockICreateTransferUseCase(ctrl)
	createTransferUseCaseMock.EXPECT().CreateTransfer(gomock.Any(), gomock.Any()).Return(port.CreateTransferResponse{TransferID: "1"}, nil)

	getTransferUseCaseMock := NewMockIGetTransferUseCase(ctrl)

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t),
		fx.Supply(fx.Annotate(createTransferUseCaseMock, fx.As(new(port.ICreateTransferUseCase)))),
		fx.Supply(fx.Annotate(getTransferUseCaseMock, fx.As(new(port.IGetTransferUseCase)))),
		fx.Provide(NewTransferController),
		fx.Provide(NewTransferRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *TransferRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					dto := request.CreateTransferRequestDTO{
						SenderID:   "1",
						ReceiverID: "2",
						Goods: []request.TransferGood{
							{
								GoodID:   "1",
								Quantity: 10,
							},
						},
					}
					payload, err := json.Marshal(dto)
					require.NoError(t, err)

					resp, err := ns.Request("transfer.create", payload, 1*time.Second)
					require.NoError(t, err)

					var respDto response.TransferCreateResponseDTO
					err = json.Unmarshal(resp.Data, &respDto)
					require.NoError(t, err)

					require.Empty(t, respDto.Error)
					require.Equal(t, respDto.Message.TransferID, "1")

					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err = app.Stop(ctx)
		require.NoError(t, err)
	}()
}

func TestTransferControllerGetTransfer(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)

	createTransferUseCaseMock := NewMockICreateTransferUseCase(ctrl)

	getTransferUseCaseMock := NewMockIGetTransferUseCase(ctrl)
	getTransferUseCaseMock.EXPECT().GetTransfer(gomock.Any(), gomock.Any()).Return(model.Transfer{
		ID:                "1",
		SenderID:          "1",
		ReceiverID:        "2",
		Status:            "Created",
		UpdateTime:        time.Now().UnixMilli(),
		CreationTime:      time.Now().UnixMilli(),
		ReservationID:     "",
		LinkedStockUpdate: 0,
		Goods: []model.GoodStock{
			{
				GoodID:   "1",
				Quantity: 10,
			},
		},
	}, nil)
	getTransferUseCaseMock.EXPECT().GetAllTransfers(gomock.Any()).Return([]model.Transfer{
		{
			ID:                "1",
			SenderID:          "1",
			ReceiverID:        "2",
			Status:            "Created",
			UpdateTime:        time.Now().UnixMilli(),
			CreationTime:      time.Now().UnixMilli(),
			ReservationID:     "",
			LinkedStockUpdate: 0,
			Goods: []model.GoodStock{
				{
					GoodID:   "1",
					Quantity: 10,
				},
			},
		},
		{
			ID:                "2",
			SenderID:          "3",
			ReceiverID:        "1",
			Status:            "Created",
			UpdateTime:        time.Now().UnixMilli(),
			CreationTime:      time.Now().UnixMilli(),
			ReservationID:     "",
			LinkedStockUpdate: 0,
			Goods: []model.GoodStock{
				{
					GoodID:   "2",
					Quantity: 10,
				},
			},
		},
	})

	ns, _ := broker.NewInProcessNATSServer(t)

	app := fx.New(
		lib.ModuleTest,
		fx.Supply(ns, t),
		fx.Supply(fx.Annotate(createTransferUseCaseMock, fx.As(new(port.ICreateTransferUseCase)))),
		fx.Supply(fx.Annotate(getTransferUseCaseMock, fx.As(new(port.IGetTransferUseCase)))),
		fx.Provide(NewTransferController),
		fx.Provide(NewTransferRouter),
		fx.Invoke(func(lc fx.Lifecycle, r *TransferRouter) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := r.Setup(ctx)
					require.NoError(t, err)

					dto := request.GetTransferRequestDTO{
						TransferID: "1",
					}
					payload, err := json.Marshal(dto)
					require.NoError(t, err)

					resp, err := ns.Request("transfer.get", payload, 1*time.Second)
					require.NoError(t, err)

					var respDto response.GetTransferResponseDTO
					err = json.Unmarshal(resp.Data, &respDto)
					require.NoError(t, err)

					require.Empty(t, respDto.Error)
					require.Equal(t, respDto.Message.TransferID, "1")

					resp, err = ns.Request("transfer.get.all", []byte{}, 1*time.Second)
					require.NoError(t, err)

					var respAllDto response.GetAllTransferResponseDTO
					err = json.Unmarshal(resp.Data, &respAllDto)
					require.NoError(t, err)

					require.Empty(t, respAllDto.Error)
					slices.SortFunc(respAllDto.Message, func(a response.TransferInfo, b response.TransferInfo) int {
						return cmp.Compare(a.TransferID, b.TransferID)
					})
					require.Equal(t, respAllDto.Message[0].TransferID, "1")
					require.Equal(t, respAllDto.Message[1].TransferID, "2")

					return nil
				},
			})
		}),
	)

	err := app.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err = app.Stop(ctx)
		require.NoError(t, err)
	}()
}
