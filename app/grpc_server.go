package app

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	app "twist-supervisor/app/interface"
	pb "twist-supervisor/pb"
	supervisor "twist-supervisor/services/supervisor"
)

func (a *App) InitGRPCServer(host string) error {

	// Start to listen on port
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.WithFields(log.Fields{
		"host": host,
	}).Info("Starting gRPC server on " + host)

	// Create gRPC server
	s := grpc.NewServer()

	// Register data source adapter service
	supervisorService := supervisor.CreateService(app.AppImpl(a))
	pb.RegisterSupervisorServer(s, supervisorService)
	reflection.Register(s)

	log.WithFields(log.Fields{
		"service": "Supervisor",
	}).Info("Registered service")

	// Starting server
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
