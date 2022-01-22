package actuators

import (
	"context"

	pb "github.com/gabrielfvale/ti0151-sistemas/app/grpc/proto"
	"github.com/gabrielfvale/ti0151-sistemas/app/pkg"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FireAlarmActuator struct {
	pkg.Actuator
	Smoke bool
}

type FireAlarmServer struct {
	*pb.UnimplementedActuatorServer
}

func (s *FireAlarmServer) GetAvailableCommands(ctx context.Context, in *emptypb.Empty) (*pb.AvailableCommandsResponse, error) {
	var commands = [4]*pb.Command{
		{Id: 1, Key: "TurnOn"},
		{Id: 2, Key: "TurnOff"},
		{Id: 3, Key: "SetFireSmoke"},
		{Id: 4, Key: "ClearFireSmoke"},
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
	if fa.Status {
		fa.Smoke = true
	}
	fa.Environment.Smoke = true
}

func (fa *FireAlarmActuator) ClearFireSmoke() {
	if fa.Status {
		fa.Smoke = false
	}
	fa.Environment.Smoke = false
}
