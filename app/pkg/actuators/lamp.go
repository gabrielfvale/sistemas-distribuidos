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

type LampServer struct {
	pkg.Actuator
	Luminosity int32
	*pb.UnimplementedActuatorServer
}

// Create gRPC server
func (la LampServer) Listen(port int) {
	log.Printf("Serving LampActuator...")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("LampActuator failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterActuatorServer(s, &LampServer{})
	// log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("LampActuator failed to serve: %v", err)
	}
}

func (s *LampServer) GetAvailableCommands(ctx context.Context, in *emptypb.Empty) (*pb.AvailableCommandsResponse, error) {
	var commands = [3]*pb.Command{
		{Id: 1, Key: "TurnOn"},
		{Id: 2, Key: "TurnOff"},
	}

	return &pb.AvailableCommandsResponse{Commands: commands[:]}, nil
}

func (s *LampServer) IssueCommand(ctx context.Context, in *pb.IssueCommandRequest) (*pb.IssueCommandResponse, error) {
	fmt.Println("ISSUING COMMAND FOR LAMP")
	switch in.Key {
	case "TurnOn":
		s.TurnOn()
	case "TurnOff":
		s.TurnOff()
	}
	return &pb.IssueCommandResponse{Status: "OK"}, nil
}

func (s *LampServer) GetProperties(ctx context.Context, in *emptypb.Empty) (*pb.PropertiesResponse, error) {
	return &pb.PropertiesResponse{}, nil
}

func (la *LampServer) TurnOn() {
	fmt.Println("Turning on the lamp")

	la.Status = true
	la.Luminosity = 100.0
	la.Environment.Luminosity += 100.0
}

func (la *LampServer) TurnOff() {
	la.Status = false
	// la.Luminosity = 0.0
	la.Environment.Luminosity -= 100.0
}
