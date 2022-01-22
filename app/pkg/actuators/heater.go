package actuators

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/gabrielfvale/ti0151-sistemas/app/grpc/proto"
	"github.com/gabrielfvale/ti0151-sistemas/app/pkg"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HeaterActuator struct {
	pkg.Actuator
	Temperature uint
}

type HeaterServer struct {
	*pb.UnimplementedActuatorServer
}

func (ha HeaterActuator) Listen(port int) {
	log.Printf("Serving HeaterActuator...")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("HeaterActuator failed to listen: %v", err)
	}
	ha.Server = grpc.NewServer()
	pb.RegisterActuatorServer(ha.Server, &HeaterServer{})
	// log.Printf("server listening at %v", lis.Addr())
	if err := ha.Server.Serve(lis); err != nil {
		log.Fatalf("HeaterActuator failed to serve: %v", err)
	}
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
