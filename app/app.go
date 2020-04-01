package app

import (
	"strconv"
	"twist-supervisor/app/eventbus"
	app "twist-supervisor/app/interface"
	"twist-supervisor/app/signalbus"

	log "github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
	"github.com/spf13/viper"
)

type App struct {
	id        uint64
	flake     *sonyflake.Sonyflake
	eventbus  *eventbus.EventBus
	signalbus *signalbus.SignalBus
}

func CreateApp() *App {

	// Genereate a unique ID for instance
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return nil
	}

	idStr := strconv.FormatUint(id, 16)

	return &App{
		id:    id,
		flake: flake,
		eventbus: eventbus.CreateConnector(
			viper.GetString("event_store.host"),
			viper.GetString("event_store.cluster_id"),
			idStr,
		),
		signalbus: signalbus.CreateConnector(
			viper.GetString("signal_server.host"),
			idStr,
		),
	}
}

func (a *App) Init() error {

	log.WithFields(log.Fields{
		"a_id": a.id,
	}).Info("Starting application")

	// Connect to event server
	err := a.eventbus.Connect()
	if err != nil {
		return err
	}

	// Connect to signal server
	err = a.signalbus.Connect()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Uninit() {
}

func (a *App) Run() error {

	port := strconv.Itoa(viper.GetInt("service.port"))
	err := a.InitGRPCServer(":" + port)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) GetEventBus() app.EventBusImpl {
	return app.EventBusImpl(a.eventbus)
}

func (a *App) GetSignalBus() app.SignalBusImpl {
	return app.SignalBusImpl(a.signalbus)
}
