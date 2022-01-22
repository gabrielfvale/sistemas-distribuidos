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

type HeaterServer struct {
	pkg.Actuator
	Temperature uint
	*pb.UnimplementedActuatorServer
}

// Create gRPC server
func (ha HeaterServer) Listen(port int) {
	log.Printf("Serving HeaterActuator...")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("HeaterActuator failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterActuatorServer(s, &HeaterServer{})
	// log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("HeaterActuator failed to serve: %v", err)
	}
}

func (s *HeaterServer) GetAvailableCommands(ctx context.Context, in *emptypb.Empty) (*pb.AvailableCommandsResponse, error) {
	var commands = [4]*pb.Command{
		{Id: 1, Key: "TurnOn"},
		{Id: 2, Key: "TurnOff"},
		{Id: 3, Key: "RaiseTemp"},
		{Id: 4, Key: "LowerTemp"},
	}

	return &pb.AvailableCommandsResponse{Commands: commands[:]}, nil
}

func (s *HeaterServer) IssueCommand(ctx context.Context, in *pb.IssueCommandRequest) (*pb.IssueCommandResponse, error) {
	switch in.Key {
	case "TurnOn":
		s.TurnOn()
	case "TurnOff":
		s.TurnOff()
	case "RaiseTemp":
		s.RaiseTemp()
	case "LowerTemp":
		s.LowerTemp()
	}
	return &pb.IssueCommandResponse{Status: "OK"}, nil
}

func (s *HeaterServer) GetProperties(ctx context.Context, in *emptypb.Empty) (*pb.PropertiesResponse, error) {
	return &pb.PropertiesResponse{}, nil
}

func (ha *HeaterServer) TurnOn() {
	ha.Status = true
	ha.Temperature = 28 // arbitrary Celsius temperature
	ha.Environment.Temperature = 28
}

func (ha *HeaterServer) TurnOff() {
	ha.Status = false
	ha.Temperature = 0
}

func (ha *HeaterServer) RaiseTemp() {
	ha.Temperature += 1
	if ha.Status {
		ha.Environment.Temperature += 1
	}
}

func (ha *HeaterServer) LowerTemp() {
	ha.Temperature -= 1
	if ha.Status {
		ha.Environment.Temperature -= 1
	}
}

func (ha *HeaterServer) SetTemp(temp int32) {
	// ha.Temperature = temp
	if ha.Status {
		ha.Environment.Temperature = temp
	}
}
