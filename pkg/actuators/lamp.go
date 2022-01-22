package actuators

import (
	"context"

	"github.com/gabrielfvale/ti0151-sistemas/grpc/impl"
	pb "github.com/gabrielfvale/ti0151-sistemas/grpc/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type LampServer struct {
	*pb.UnimplementedActuatorServer
}

type LampActuator struct {
	impl.Actuator
	Server     *LampServer
	Luminosity int32
}

func (s *LampServer) GetAvailableCommands(ctx context.Context, in *emptypb.Empty) (*pb.AvailableCommandsResponse, error) {
	var commands = [3]*pb.Command{
		{Id: 1, Key: "TurnOn"},
	}

	return &pb.AvailableCommandsResponse{Commands: commands[:]}, nil
}

func (s *LampServer) IssueCommand(ctx context.Context, in *pb.IssueCommandRequest) (*pb.IssueCommandResponse, error) {
	return &pb.IssueCommandResponse{}, nil
}

func (s *LampServer) GetProperties(ctx context.Context, in *emptypb.Empty) (*pb.PropertiesResponse, error) {
	return &pb.PropertiesResponse{}, nil
}

func (la *LampActuator) TurnOn() {
	la.Status = true
	la.Luminosity = 100.0
}

func (la *LampActuator) TurnOff() {
	la.Status = false
	la.Luminosity = 0.0
}
