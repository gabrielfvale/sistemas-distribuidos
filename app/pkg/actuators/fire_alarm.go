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

type FireAlarmServer struct {
	pkg.Actuator
	Smoke bool
	*pb.UnimplementedActuatorServer
}

// Create gRPC server
func (fa FireAlarmServer) Listen(port int) {
	log.Printf("Serving FireAlarmActuator...")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("FireAlarmActuator failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterActuatorServer(s, &FireAlarmServer{})
	// log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("FireAlarmActuator failed to serve: %v", err)
	}
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
	fmt.Printf("Fire IssueCommand: %s\n", in.Key)
	switch in.Key {
	case "TurnOn":
		s.TurnOn()
	case "TurnOff":
		s.TurnOff()
	case "SetFireSmoke":
		s.SetFireSmoke()
	case "ClearFireSmoke":
		s.ClearFireSmoke()
	}
	return &pb.IssueCommandResponse{Status: "OK"}, nil
}

func (s *FireAlarmServer) GetProperties(ctx context.Context, in *emptypb.Empty) (*pb.PropertiesResponse, error) {
	return &pb.PropertiesResponse{}, nil
}

func (fa *FireAlarmServer) TurnOn() {
	fa.Status = true
}

func (fa *FireAlarmServer) TurnOff() {
	fa.Status = false
}

func (fa *FireAlarmServer) SetFireSmoke() {
	if fa.Status {
		fa.Smoke = true
	}
	fa.Environment.Smoke = true
}

func (fa *FireAlarmServer) ClearFireSmoke() {
	if fa.Status {
		fa.Smoke = false
	}
	fa.Environment.Smoke = false
}
