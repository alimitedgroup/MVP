package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/common/dto/request"
	"github.com/alimitedgroup/MVP/common/lib/observability"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	ErrorGetToken           = errors.New("error getting token for given credentials")
	ErrorGetRole            = errors.New("error getting role for given token")
	ErrorGetUsername        = errors.New("error getting username")
	ErrorGetWarehouses      = errors.New("error getting warehouses")
	ErrorGetGoods           = errors.New("error getting goods")
	ErrorGetStock           = errors.New("error getting global stock")
	ErrorGetTransfers       = errors.New("error getting transfers")
	ErrorGetOrders          = errors.New("error getting orders")
	ErrorGetQueries         = errors.New("error getting queries")
	ErrorAddStock           = errors.New("error adding stock")
	ErrorRemoveStock        = errors.New("error adding stock")
	ErrorCreateOrder        = errors.New("error creating order")
	ErrorCreateTransfer     = errors.New("error creating transfer")
	ErrorCreateQuery        = errors.New("error creating query")
	ErrorCreateGood         = errors.New("error creating good")
	ErrorUpdateGood         = errors.New("error updating good")
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorTokenInvalid       = errors.New("this token is invalid")
	ErrorTokenExpired       = errors.New("this token is expired")
)

var Module = fx.Module(
	"business",
	fx.Provide(fx.Annotate(
		NewBusiness,
		fx.As(new(portin.Auth)),
		fx.As(new(portin.Warehouses)),
		fx.As(new(portin.Order)),
		fx.As(new(portin.Notifications)),
	)),
	fx.Decorate(observability.WrapLogger("business")),
)

type BusinessParams struct {
	fx.In

	Auth         portout.AuthenticationPortOut
	Catalog      portout.CatalogPortOut
	Order        portout.OrderPortOut
	Notification portout.NotificationPortOut
	Logger       *zap.Logger
}

func NewBusiness(p BusinessParams) *Business {
	return &Business{auth: p.Auth, catalog: p.Catalog, order: p.Order, notification: p.Notification, Logger: p.Logger}
}

type Business struct {
	auth         portout.AuthenticationPortOut
	catalog      portout.CatalogPortOut
	order        portout.OrderPortOut
	notification portout.NotificationPortOut
	*zap.Logger
}

func (b *Business) CreateQuery(goodId string, operator string, threshold int) (string, error) {
	queryId, err := b.notification.CreateQuery(dto.Rule{
		GoodId:    goodId,
		Operator:  operator,
		Threshold: threshold,
	})
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrorCreateQuery, err)
	}
	return queryId, nil

}

func (b *Business) GetQueries() ([]dto.Query, error) {
	queries, err := b.notification.GetQueries()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorGetQueries, err)
	}
	resp := make([]dto.Query, 0, len(queries))
	for _, query := range queries {
		resp = append(resp, dto.Query{
			QueryID:   query.RuleId.String(),
			GoodID:    query.GoodId,
			Operator:  query.Operator,
			Threshold: query.Threshold,
		})
	}
	return resp, nil
}

func (b *Business) AddStock(warehouseId string, goodId string, quantity int64) error {
	err := b.catalog.AddStock(warehouseId, goodId, quantity)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorAddStock, err)
	}
	return nil
}

func (b *Business) RemoveStock(warehouseId string, goodId string, quantity int64) error {
	err := b.catalog.RemoveStock(warehouseId, goodId, quantity)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorRemoveStock, err)
	}
	return nil
}

func (b *Business) CreateOrder(name string, fullname string, address string, goods map[string]int64) (string, error) {
	goodList := make([]request.CreateOrderGood, 0, len(goods))
	for goodID, quantity := range goods {
		goodList = append(goodList, request.CreateOrderGood{
			GoodID:   goodID,
			Quantity: quantity,
		})
	}
	createDto := request.CreateOrderRequestDTO{
		Name:     name,
		FullName: fullname,
		Address:  address,
		Goods:    goodList,
	}
	orderId, err := b.order.CreateOrder(createDto)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrorCreateOrder, err)
	}
	return orderId.OrderID, nil

}

func (b *Business) GetOrders() ([]dto.Order, error) {
	orders, err := b.order.GetAllOrders()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorGetOrders, err)
	}

	resp := make([]dto.Order, 0, len(orders))
	for _, order := range orders {
		goods := make(map[string]int64, len(order.Goods))
		for _, good := range order.Goods {
			goods[good.GoodID] = good.Quantity
		}

		resp = append(resp, dto.Order{
			Status:       order.Status,
			CreationTime: order.CreationTime,
			UpdateTime:   order.UpdateTime,
			OrderID:      order.OrderID,
			Name:         order.Name,
			FullName:     order.FullName,
			Address:      order.Address,
			Reservations: order.Reservations,
			Goods:        goods,
		})
	}

	return resp, nil
}

func (b *Business) CreateTransfer(senderID string, receiverID string, goods map[string]int64) (string, error) {
	goodList := make([]request.TransferGood, 0, len(goods))
	for goodID, quantity := range goods {
		goodList = append(goodList, request.TransferGood{
			GoodID:   goodID,
			Quantity: quantity,
		})
	}
	createDto := request.CreateTransferRequestDTO{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Goods:      goodList,
	}

	transferId, err := b.order.CreateTransfer(createDto)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrorCreateTransfer, err)
	}
	return transferId.TransferID, nil
}

func (b *Business) GetTransfers() ([]dto.Transfer, error) {
	transfers, err := b.order.GetAllTransfers()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorGetTransfers, err)
	}

	resp := make([]dto.Transfer, 0, len(transfers))
	for _, transfer := range transfers {
		goods := make(map[string]int64, len(transfer.Goods))
		for _, good := range transfer.Goods {
			goods[good.GoodID] = good.Quantity
		}

		resp = append(resp, dto.Transfer{
			Status:       transfer.Status,
			CreationTime: transfer.CreationTime,
			UpdateTime:   transfer.UpdateTime,
			TransferID:   transfer.TransferID,
			SenderID:     transfer.SenderID,
			ReceiverID:   transfer.ReceiverID,
			Goods:        goods,
		})
	}

	return resp, nil
}

func (b *Business) CreateGood(ctx context.Context, name string, description string) (string, error) {
	goodId, err := b.catalog.CreateGood(ctx, name, description)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrorCreateGood, err)
	}
	return goodId, nil
}

func (b *Business) UpdateGood(ctx context.Context, goodId string, name string, description string) error {
	err := b.catalog.UpdateGood(ctx, goodId, name, description)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUpdateGood, err)
	}
	return nil
}

func (b *Business) GetWarehouses() ([]types.WarehouseOverview, error) {
	warehouses, err := b.catalog.ListWarehouses()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorGetWarehouses, err)
	}

	result := make([]types.WarehouseOverview, 0, len(warehouses))
	for _, warehouse := range warehouses {
		result = append(result, types.WarehouseOverview{ID: warehouse.ID})
	}
	return result, nil
}

func (b *Business) GetGoods() ([]dto.GoodAndAmount, error) {
	goods, err := b.catalog.ListGoods()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorGetGoods, err)
	}

	amounts, err := b.catalog.ListStock()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorGetStock, err)
	}

	warehouses, err := b.catalog.ListWarehouses()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorGetWarehouses, err)
	}

	result := make([]dto.GoodAndAmount, 0, len(goods))
	for _, good := range goods {
		amount, ok := amounts[good.ID]
		if !ok {
			amount = 0
		}

		// amount of stock per warehouse
		amounts := make(map[string]int64, len(warehouses))
		for _, warehouse := range warehouses {
			if stock, ok := warehouse.Stock[good.ID]; ok {
				amounts[warehouse.ID] = stock
			}
		}

		result = append(result, dto.GoodAndAmount{
			Name:        good.Name,
			Description: good.Description,
			ID:          good.ID,
			Amount:      amount,
			Amounts:     amounts,
		})
	}

	return result, nil
}

func (b *Business) Login(username string) (types.LoginResult, error) {
	token, err := b.auth.GetToken(username)
	if err != nil {
		b.Error("Failed to get JWT token for given username", zap.Error(err))
		return types.LoginResult{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
	}
	if token == "" {
		return types.LoginResult{}, ErrorInvalidCredentials
	}

	parsed, err := b.auth.VerifyToken(token)
	if err != nil {
		b.Error("Failed to parse JWT returned by authentication microservice", zap.Error(err))
		return types.LoginResult{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
	}

	role, err := b.auth.GetRole(parsed)
	if err != nil {
		b.Error("Failed to get role from JWT returned by authentication microservice", zap.Error(err))
		return types.LoginResult{}, fmt.Errorf("%w: %w", ErrorGetRole, err)
	}

	return types.LoginResult{
		Token: token,
		Role:  role,
	}, nil
}

func (b *Business) ValidateToken(token string) (types.UserData, error) {
	tok, err := b.auth.VerifyToken(types.UserToken(token))
	if err != nil {
		if errors.Is(err, portout.ErrTokenExpired) {
			return types.UserData{}, ErrorTokenExpired
		} else if errors.Is(err, portout.ErrTokenInvalid) {
			return types.UserData{}, ErrorTokenInvalid
		} else {
			b.Error("Failed to validate JWT token", zap.Error(err))
			return types.UserData{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
		}
	}

	username, err := b.auth.GetUsername(tok)
	if err != nil {
		if errors.Is(err, portout.ErrTokenInvalid) {
			return types.UserData{}, ErrorTokenInvalid
		} else if errors.Is(err, portout.ErrTokenExpired) {
			return types.UserData{}, ErrorTokenExpired
		} else {
			b.Error("Failed to get username from valid JWT token", zap.Error(err))
			return types.UserData{}, fmt.Errorf("%w: %w", ErrorGetUsername, err)
		}
	}

	role, err := b.auth.GetRole(tok)
	if err != nil {
		if errors.Is(err, portout.ErrTokenInvalid) {
			return types.UserData{}, ErrorTokenInvalid
		} else if errors.Is(err, portout.ErrTokenExpired) {
			return types.UserData{}, ErrorTokenExpired
		} else {
			b.Error("Failed to get role from valid JWT token", zap.Error(err))
			return types.UserData{}, fmt.Errorf("%w: %w", ErrorGetRole, err)
		}
	}

	return types.UserData{Username: username, Role: role}, err
}

// Asserzione a compile time che Business implementi le interfaccie delle porte di input
var _ portin.Auth = (*Business)(nil)
var _ portin.Warehouses = (*Business)(nil)
var _ portin.Order = (*Business)(nil)
var _ portin.Notifications = (*Business)(nil)
