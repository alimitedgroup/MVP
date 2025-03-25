package lib

import "context"

type BrokerRoute interface {
	Setup(context.Context) error
}
