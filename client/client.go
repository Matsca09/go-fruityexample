package main

import (
	"context"
	"log"
	"time"

	pb "github.com/matsca09/go-fruityexample/fruit"
	grpc "google.golang.org/grpc"
)

func printLiveStream(c chan *pb.AvailableFruits) {
	for {
		fruitBasket := <-c
		log.Println("Received data from server!")
		log.Printf(" -- %d apples, %d bananas and %d oranges --", fruitBasket.Apple, fruitBasket.Banana, fruitBasket.Orange)
	}
}

func streamRoutine(stream pb.Fruit_GetLiveFruitsClient, c chan *pb.AvailableFruits) {
	for {
		fruitBasket, err := stream.Recv()
		if err != nil {
			log.Fatal("ERROR while receiving data: " + err.Error())
		}
		c <- fruitBasket

	}
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:50000", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannot connect to the gRPC server. Is it running? " + err.Error())
	}
	defer conn.Close()

	log.Println("Start goroutine for print live stream")
	dataChan := make(chan *pb.AvailableFruits, 1)
	go printLiveStream(dataChan)

	fruitClient := pb.NewFruitClient(conn)
	log.Println("Start stream fruit status from server...")
	stream, errAct := fruitClient.GetLiveFruits(context.Background(), &pb.Empty{})
	if errAct != nil {
		log.Fatal("Cannot get fruit basket from server! " + errAct.Error())
	}
	go streamRoutine(stream, dataChan)

	for {
		time.Sleep(11 * time.Second)
		log.Println("Add fruits!")
		fruitClient.AddApples(context.Background(), &pb.Integer{Value: 10})
		fruitClient.AddBananas(context.Background(), &pb.Integer{Value: 10})
		fruitClient.AddOranges(context.Background(), &pb.Integer{Value: 10})
	}
}
