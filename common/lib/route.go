package lib

import "context"

type APIRoute interface {
	Setup(context.Context)
}

type BrokerRoute interface {
	Setup(context.Context)
}
