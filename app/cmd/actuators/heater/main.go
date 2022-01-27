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

type HeaterServer struct {
	pkg.Actuator
	Temperature uint
	*pb.UnimplementedActuatorServer
}

// Create gRPC server
func (ha HeaterServer) Listen(port int) {
	ha.EnvironmentConn = pkg.ConnectToEnviroment()
	log.Printf("Serving HeaterActuator on port %d", port)
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

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("No port argument provided")
	}
	server := HeaterServer{}
	server.Name = "Heater"
	server.EnvironmentConn = pkg.ConnectToEnviroment()

	port, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("Could not parse port")
	}
	server.Listen(port)
}

func (s *HeaterServer) GetAvailableCommands(ctx context.Context, in *emptypb.Empty) (*pb.AvailableCommandsResponse, error) {
	fmt.Println("ISSUING COMMAND FOR HEATER")
	var commands = [4]*pb.Command{
		{Id: 1, Key: "TurnOn"},
		{Id: 2, Key: "TurnOff"},
		{Id: 3, Key: "RaiseTemp"},
		{Id: 4, Key: "LowerTemp"},
	}

	return &pb.AvailableCommandsResponse{Commands: commands[:]}, nil
}

func (s *HeaterServer) IssueCommand(ctx context.Context, in *pb.IssueCommandRequest) (*pb.IssueCommandResponse, error) {
	log.Printf("RECEIVED ISSUE COMMAND %v", in.Key)
	s.EnvironmentConn = pkg.ConnectToEnviroment()
	env := pkg.ReadEnviromentData(s.EnvironmentConn)

	switch in.Key {
	case "TurnOn":
		s.TurnOn(env)
	case "TurnOff":
		s.TurnOff()
	case "RaiseTemp":
		s.RaiseTemp(env)
	case "LowerTemp":
		s.LowerTemp(env)
	}
	return &pb.IssueCommandResponse{Status: "OK"}, nil
}

func (s *HeaterServer) GetProperties(ctx context.Context, in *emptypb.Empty) (*pb.PropertiesResponse, error) {
	return &pb.PropertiesResponse{}, nil
}

func (ha *HeaterServer) TurnOn(env pkg.Environment) {
	ha.Status = true
	ha.Temperature = 28 // arbitrary Celsius temperature
	pkg.WriteEnviromentData(ha.EnvironmentConn, "Temperature", 28)
}

func (ha *HeaterServer) TurnOff() {
	ha.Status = false
	ha.Temperature = 0
}

func (ha *HeaterServer) RaiseTemp(env pkg.Environment) {
	ha.Temperature += 1
	if ha.Status {
		pkg.WriteEnviromentData(ha.EnvironmentConn, "Temperature", env.Temperature+1)
	}
}

func (ha *HeaterServer) LowerTemp(env pkg.Environment) {
	ha.Temperature -= 1
	if ha.Status {
		pkg.WriteEnviromentData(ha.EnvironmentConn, "Temperature", env.Temperature-1)
	}
}

func (ha *HeaterServer) SetTemp(temp int32) {

}
