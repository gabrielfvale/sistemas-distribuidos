package actuators

import (
	"context"

	"github.com/gabrielfvale/ti0151-sistemas/grpc/impl"
	pb "github.com/gabrielfvale/ti0151-sistemas/grpc/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FireAlarmActuator struct {
	impl.Actuator
	Smoke bool
}

type FireAlarmServer struct {
	*pb.UnimplementedActuatorServer
}

func (s *FireAlarmServer) GetAvailableCommands(ctx context.Context, in *emptypb.Empty) (*pb.AvailableCommandsResponse, error) {
	var commands = [3]*pb.Command{
		{Id: 1, Key: "command1"},
	}

	return &pb.AvailableCommandsResponse{Commands: commands[:]}, nil
}

func (s *FireAlarmServer) IssueCommand(ctx context.Context, in *pb.IssueCommandRequest) (*pb.IssueCommandResponse, error) {
	return &pb.IssueCommandResponse{}, nil
}

func (s *FireAlarmServer) GetProperties(ctx context.Context, in *emptypb.Empty) (*pb.PropertiesResponse, error) {
	return &pb.PropertiesResponse{}, nil
}

func (fa *FireAlarmActuator) TurnOn() {
	fa.Status = true
}

func (fa *FireAlarmActuator) TurnOff() {
	fa.Status = false
}

func (fa *FireAlarmActuator) SetFireSmoke() {
	fa.Smoke = true
}

func (fa *FireAlarmActuator) ClearFireSmoke() {
	fa.Smoke = false
}
