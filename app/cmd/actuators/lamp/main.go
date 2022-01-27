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

type LampServer struct {
	pkg.Actuator
	Luminosity int32
	*pb.UnimplementedActuatorServer
}

// Create gRPC server
func (la LampServer) Listen(port int) {
	log.Printf("Serving LampActuator on port %d", port)
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

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("No port argument provided")
	}
	server := LampServer{}
	server.Name = "Lamp"

	port, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Could not parse port")
	}
	server.Listen(port)
}

func (s *LampServer) GetAvailableCommands(ctx context.Context, in *emptypb.Empty) (*pb.AvailableCommandsResponse, error) {
	var commands = [3]*pb.Command{
		{Id: 1, Key: "TurnOn"},
		{Id: 2, Key: "TurnOff"},
	}

	return &pb.AvailableCommandsResponse{Commands: commands[:]}, nil
}

func (s *LampServer) IssueCommand(ctx context.Context, in *pb.IssueCommandRequest) (*pb.IssueCommandResponse, error) {
	log.Printf("RECEIVED ISSUE COMMAND %v", in.Key)
	s.EnvironmentConn = pkg.ConnectToEnviroment()
	env := pkg.ReadEnviromentData(s.EnvironmentConn)
	switch in.Key {
	case "TurnOn":
		s.TurnOn(env)
	case "TurnOff":
		s.TurnOff(env)
	}
	return &pb.IssueCommandResponse{Status: "OK"}, nil
}

func (s *LampServer) GetProperties(ctx context.Context, in *emptypb.Empty) (*pb.PropertiesResponse, error) {
	return &pb.PropertiesResponse{}, nil
}

func (la *LampServer) TurnOn(env pkg.Environment) {
	fmt.Println("Turning on the lamp")
	la.Status = true
	la.Luminosity = 100.0
	pkg.WriteEnviromentData(la.EnvironmentConn, "Luminosity", env.Luminosity+100)
}

func (la *LampServer) TurnOff(env pkg.Environment) {
	la.Status = false
	pkg.WriteEnviromentData(la.EnvironmentConn, "Luminosity", env.Luminosity-100)
}
