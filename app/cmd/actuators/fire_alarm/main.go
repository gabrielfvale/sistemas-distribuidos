package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/gabrielfvale/ti0151-sistemas/app/grpc/proto"
	"github.com/gabrielfvale/ti0151-sistemas/app/pkg"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type FireAlarmServer struct {
	pkg.Actuator
	*pb.UnimplementedActuatorServer
}

// Create gRPC server
func (fa FireAlarmServer) Listen(port int) {

	log.Printf("Serving FireAlarmActuator on port %d", port)
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

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("No port argument provided")
	}
	server := FireAlarmServer{}
	server.Name = "Fire alarm"

	port, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Could not parse port")
	}
	server.Listen(port)
}

func (s *FireAlarmServer) GetAvailableCommands(ctx context.Context, in *emptypb.Empty) (*pb.AvailableCommandsResponse, error) {
	fmt.Println("ISSUING COMMAND FOR FIRE ALARM")
	var commands = [4]*pb.Command{
		{Id: 1, Key: "TurnOn"},
		{Id: 2, Key: "TurnOff"},
		{Id: 3, Key: "SetFireSmoke"},
		{Id: 4, Key: "ClearFireSmoke"},
	}

	return &pb.AvailableCommandsResponse{Commands: commands[:]}, nil
}

func (s *FireAlarmServer) IssueCommand(ctx context.Context, in *pb.IssueCommandRequest) (*pb.IssueCommandResponse, error) {
	log.Printf("RECEIVED ISSUE COMMAND %v", in.Key)
	s.EnvironmentConn = pkg.ConnectToEnviroment()
	env := pkg.ReadEnviromentData(s.EnvironmentConn)
	switch in.Key {
	case "TurnOn":
		s.TurnOn()
	case "TurnOff":
		s.TurnOff()
	case "SetFireSmoke":
		s.SetFireSmoke(env)
	case "ClearFireSmoke":
		s.ClearFireSmoke(env)
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

func (fa *FireAlarmServer) SetFireSmoke(env pkg.Environment) {
	if fa.Status {
		pkg.WriteEnviromentData(fa.EnvironmentConn, "Temperature", 1)
	}
}

func (fa *FireAlarmServer) ClearFireSmoke(env pkg.Environment) {
	if fa.Status {
		pkg.WriteEnviromentData(fa.EnvironmentConn, "Temperature", 0)
	}
}
