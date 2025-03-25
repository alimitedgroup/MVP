package business

import (
	"errors"
	"fmt"
	"github.com/alimitedgroup/MVP/common/dto"
	"github.com/alimitedgroup/MVP/srv/api_gateway/business/types"
	"github.com/alimitedgroup/MVP/srv/api_gateway/portin"
	"go.uber.org/fx"
	"time"

	"github.com/alimitedgroup/MVP/srv/api_gateway/portout"
)

var (
	ErrorGetToken           = errors.New("error getting token for given credentials")
	ErrorGetRole            = errors.New("error getting role for given token")
	ErrorGetUsername        = errors.New("error getting username")
	ErrorGetWarehouses      = errors.New("error getting warehouses")
	ErrorGetGoods           = errors.New("error getting goods")
	ErrorGetStock           = errors.New("error getting global stock")
	ErrorInvalidCredentials = errors.New("invalid credentials")
	ErrorTokenInvalid       = errors.New("this token is invalid")
	ErrorTokenExpired       = errors.New("this token is expired")
)

var Module = fx.Options(
	fx.Provide(fx.Annotate(
		NewBusiness,
		fx.As(new(portin.Auth)),
		fx.As(new(portin.Warehouses)),
	)),
)

func NewBusiness(auth portout.AuthenticationPortOut, catalog portout.CatalogPortOut) *Business {
	return &Business{auth: auth, catalog: catalog}
}

type Business struct {
	auth    portout.AuthenticationPortOut
	catalog portout.CatalogPortOut
}

func (b *Business) GetWarehouseByID(_ int64) (dto.Warehouse, error) {
	//TODO da implementare quando catalog supporta questa query
	panic("implement me")
}

func (b *Business) GetWarehouses() ([]portin.WarehouseOverview, error) {
	warehouses, err := b.catalog.ListWarehouses()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorGetWarehouses, err)
	}

	result := make([]portin.WarehouseOverview, 0, len(warehouses))
	for _, warehouse := range warehouses {
		result = append(result, portin.WarehouseOverview{ID: warehouse.ID})
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

	result := make([]dto.GoodAndAmount, 0, len(goods))
	for _, good := range goods {
		amount, ok := amounts[good.ID]
		if !ok {
			amount = 0
		}

		result = append(result, dto.GoodAndAmount{
			Name:        good.Name,
			Description: good.Description,
			ID:          good.ID,
			Amount:      amount,
		})
	}

	return result, nil
}

func (b *Business) Login(username string) (portin.LoginResult, error) {
	token, err := b.auth.GetToken(username)
	if err != nil {
		return portin.LoginResult{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
	}
	if token == "" {
		return portin.LoginResult{}, ErrorInvalidCredentials
	}

	parsed, err := b.auth.VerifyToken(token)
	if err != nil {
		return portin.LoginResult{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
	}

	role, err := b.auth.GetRole(parsed)
	if err != nil {
		return portin.LoginResult{}, fmt.Errorf("%w: %w", ErrorGetRole, err)
	}

	// TODO: bisognerebbe prendere la scadenza dall'output del servizio di Authentication
	expiration := time.Now().Add(7 * 24 * time.Hour)

	return portin.LoginResult{
		Token:           token,
		TokenExpiration: expiration,
		Role:            role,
	}, nil
}

func (b *Business) ValidateToken(token string) (portin.UserData, error) {
	tok, err := b.auth.VerifyToken(types.UserToken(token))
	if err != nil {
		if errors.Is(err, portout.ErrTokenExpired) {
			return portin.UserData{}, ErrorTokenExpired
		} else if errors.Is(err, portout.ErrTokenInvalid) {
			return portin.UserData{}, ErrorTokenInvalid
		} else {
			return portin.UserData{}, fmt.Errorf("%w: %w", ErrorGetToken, err)
		}
	}

	username, err := b.auth.GetUsername(tok)
	if err != nil {
		if errors.Is(err, portout.ErrTokenInvalid) {
			return portin.UserData{}, ErrorTokenInvalid
		} else if errors.Is(err, portout.ErrTokenExpired) {
			return portin.UserData{}, ErrorTokenExpired
		} else {
			return portin.UserData{}, fmt.Errorf("%w: %w", ErrorGetUsername, err)
		}
	}

	role, err := b.auth.GetRole(tok)
	if err != nil {
		if errors.Is(err, portout.ErrTokenInvalid) {
			return portin.UserData{}, ErrorTokenInvalid
		} else if errors.Is(err, portout.ErrTokenExpired) {
			return portin.UserData{}, ErrorTokenExpired
		} else {
			return portin.UserData{}, fmt.Errorf("%w: %w", ErrorGetRole, err)
		}
	}

	return portin.UserData{Username: username, Role: role}, err
}

// Asserzione a compile time che Business implementi le interfaccie delle porte di input
var _ portin.Auth = (*Business)(nil)
var _ portin.Warehouses = (*Business)(nil)
