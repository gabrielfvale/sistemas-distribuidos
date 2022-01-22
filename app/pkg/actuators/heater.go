package actuators

import (
	"context"

	"github.com/gabrielfvale/ti0151-sistemas/app/grpc/impl"
	pb "github.com/gabrielfvale/ti0151-sistemas/app/grpc/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HeaterActuator struct {
	impl.Actuator
	Temperature uint
}

type HeaterServer struct {
	*pb.UnimplementedActuatorServer
}

func (s *HeaterServer) GetAvailableCommands(ctx context.Context, in *emptypb.Empty) (*pb.AvailableCommandsResponse, error) {
	var commands = [3]*pb.Command{
		{Id: 1, Key: "command1"},
	}

	return &pb.AvailableCommandsResponse{Commands: commands[:]}, nil
}

func (s *HeaterServer) IssueCommand(ctx context.Context, in *pb.IssueCommandRequest) (*pb.IssueCommandResponse, error) {
	return &pb.IssueCommandResponse{}, nil
}

func (s *HeaterServer) GetProperties(ctx context.Context, in *emptypb.Empty) (*pb.PropertiesResponse, error) {
	return &pb.PropertiesResponse{}, nil
}

func (ha *HeaterActuator) TurnOn() {
	ha.Status = true
	ha.Temperature = 28 // arbitrary Celsius temperature
	ha.Environment.Temperature = 28
}

func (ha *HeaterActuator) TurnOff() {
	ha.Status = false
	ha.Temperature = 0
}

func (ha *HeaterActuator) RaiseTemp() {
	ha.Temperature += 1
	if ha.Status {
		ha.Environment.Temperature += 1
	}
}

func (ha *HeaterActuator) LowerTemp() {
	ha.Temperature -= 1
	if ha.Status {
		ha.Environment.Temperature -= 1
	}
}

func (ha *HeaterActuator) SetTemp(temp uint) {
	ha.Temperature = temp
	if ha.Status {
		ha.Environment.Temperature = temp
	}
}
