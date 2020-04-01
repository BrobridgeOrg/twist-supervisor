package supervisor

import (
	"github.com/gogo/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	app "twist-supervisor/app/interface"
	pb "twist-supervisor/pb"
)

type Service struct {
	app app.AppImpl
}

func CreateService(app app.AppImpl) *Service {

	// Preparing service
	service := &Service{
		app: app,
	}

	return service
}

func (service *Service) PrepareTransaction(ctx context.Context, in *pb.PrepareTransactionRequest) (*pb.PrepareTransactionReply, error) {

	log.Info("Prepare a new transaction")

	// Prepare transaction information
	request := &pb.TransactionRequest{
		TransactionID: in.TransactionID,
		Mode:          in.Mode,
	}
	data, err := proto.Marshal(request)
	if err != nil {
		return &pb.PrepareTransactionReply{
			Success: false,
		}, err
	}

	// Dispatch transation request
	eb := service.app.GetEventBus()
	err = eb.Emit("twist.transactions", data)
	if err != nil {
		return &pb.PrepareTransactionReply{
			Success: false,
		}, err
	}

	return &pb.PrepareTransactionReply{
		Success:       true,
		TransactionID: request.TransactionID,
	}, nil
}

func (service *Service) UpdateAssignment(ctx context.Context, in *pb.UpdateAssignmentRequest) (*pb.UpdateAssignmentReply, error) {

	log.WithFields(log.Fields{
		"runner":      in.RunnerID,
		"transaction": in.TransactionID,
	}).Info("Update Assignment")

	// Prepare event
	event := &pb.TransactionEvent{
		TransactionID: in.TransactionID,
		RunnerID:      in.RunnerID,
		EventName:     "Assigned",
		Payload:       "",
	}

	data, err := proto.Marshal(event)
	if err != nil {
		return &pb.UpdateAssignmentReply{
			Success: false,
		}, err
	}

	// Emit transaction event
	sb := service.app.GetSignalBus()
	err = sb.Emit("twist.transaction."+in.TransactionID+".eventEmitted", data)
	if err != nil {
		return &pb.UpdateAssignmentReply{
			Success: false,
		}, err
	}

	return &pb.UpdateAssignmentReply{
		Success: true,
	}, nil
}
