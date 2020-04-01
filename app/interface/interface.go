package app

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type EventBusImpl interface {
	Emit(string, []byte) error
	On(string, func(*stan.Msg)) error
}

type SignalBusImpl interface {
	Emit(string, []byte) error
	Watch(string, func(*nats.Msg)) (*nats.Subscription, error)
}

type AppImpl interface {
	GetEventBus() EventBusImpl
	GetSignalBus() SignalBusImpl
}
