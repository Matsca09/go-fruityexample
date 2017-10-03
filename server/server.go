package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/matsca09/go-fruityexample/fruit"
	grpc "google.golang.org/grpc"
)

// FruitServer mah fruit server
type FruitServer struct {
	Apples  int32
	Bananas int32
	Oranges int32
}

// AddApples does what its name means
func (s *FruitServer) AddApples(ctx context.Context, protoInteger *pb.Integer) (*pb.Empty, error) {
	log.Printf("Increasing apples -- %d -> %d", s.Apples, s.Apples+protoInteger.Value)
	s.Apples += protoInteger.Value
	return new(pb.Empty), nil
}

// AddBananas does what its name means
func (s *FruitServer) AddBananas(ctx context.Context, protoInteger *pb.Integer) (*pb.Empty, error) {
	log.Printf("Increasing bananas -- %d -> %d", s.Bananas, s.Bananas+protoInteger.Value)
	s.Bananas += protoInteger.Value
	return new(pb.Empty), nil
}

// AddOranges does what its name means
func (s *FruitServer) AddOranges(ctx context.Context, protoInteger *pb.Integer) (*pb.Empty, error) {
	log.Printf("Increasing oranges -- %d -> %d", s.Oranges, s.Oranges+protoInteger.Value)
	s.Oranges += protoInteger.Value
	return new(pb.Empty), nil
}

// GetAllFruits returns a populated AvailableFruits
func (s *FruitServer) GetAllFruits(ctx context.Context, blackhole *pb.Empty) (*pb.AvailableFruits, error) {
	log.Printf("Client requested full data: sending -- %d apples, %d bananas and %d oranges", s.Apples, s.Bananas, s.Oranges)
	return &pb.AvailableFruits{Apple: s.Apples, Banana: s.Bananas, Orange: s.Oranges}, nil
}

// GetLiveFruits starts a stream
func (s *FruitServer) GetLiveFruits(blackhole *pb.Empty, stream pb.Fruit_GetLiveFruitsServer) error {
	for {
		strErr := stream.Send(&pb.AvailableFruits{Apple: s.Apples, Banana: s.Bananas, Orange: s.Oranges})
		if strErr != nil {
			log.Println("Error in stream: " + strErr.Error())
			return strErr
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	socket, errTCP := net.Listen("tcp", ":50000")
	if errTCP != nil {
		log.Fatal("Cannot bind port because " + errTCP.Error())
	}
	server := grpc.NewServer()
	fs := new(FruitServer)
	pb.RegisterFruitServer(server, fs)
	if err := server.Serve(socket); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
