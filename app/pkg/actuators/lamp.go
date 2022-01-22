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
	*pb.UnimplementedActuatorServer
}

type LampActuator struct {
	pkg.Actuator
	Luminosity int32
}

func (la LampActuator) Listen(port int) {
	log.Printf("Serving LampActuator...")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("LampActuator failed to listen: %v", err)
	}
	la.Server = grpc.NewServer()
	pb.RegisterActuatorServer(la.Server, &LampServer{})
	// log.Printf("server listening at %v", lis.Addr())
	if err := la.Server.Serve(lis); err != nil {
		log.Fatalf("LampActuator failed to serve: %v", err)
	}
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
	la.Environment.Luminosity += 100.0
}

func (la *LampActuator) TurnOff() {
	la.Status = false
	// la.Luminosity = 0.0
	la.Environment.Luminosity -= 100.0
}
