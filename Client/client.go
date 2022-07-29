package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vishal-swiggy/calculator/protobuf"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Starting Requests....")

	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	serivceClient := protobuf.NewCalculatorServiceClient(cc)

	Sum(serivceClient)
	PrimeNumbers(serivceClient)
	ComputeAverage(serivceClient)
	FindMaxNumber(serivceClient)

	fmt.Println("All APIs called successfully")

}

func Sum(serivceClient protobuf.CalculatorServiceClient) {

	fmt.Println("Sum API call...")

	req := protobuf.SumIn{
		In1: 213.31,
		In2: 112.34,
	}

	resp, err := serivceClient.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling Sum API: %v", err)
	}

	log.Printf("Sum API Response, Sum : %v", resp.GetOut())

}

func PrimeNumbers(serivceClient protobuf.CalculatorServiceClient) {
	fmt.Println("PrimeNumbers API call...")

	req := protobuf.PrimeNumbersIn{
		In: 111,
	}

	respStream, err := serivceClient.PrimeNumbers(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling PrimeNumbers API : %v", err)
	}

	for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			//we have reached to the end of the file
			break
		}

		if err != nil {
			log.Fatalf("error while receving server stream : %v", err)
		}

		log.Println("PrimeNumbers API Response, Prime Number : ", msg.GetOut())
	}
}

func ComputeAverage(serivceClient protobuf.CalculatorServiceClient) {

	fmt.Println("ComputeAverage API call....")

	stream, err := serivceClient.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error calling ComputeAverage API : %v", err)
	}

	requests := [] protobuf.ComputeAverageIn{
		{
			In: 234,
		},
		{
			In: 123,
		},
		{
			In: 204,
		},
		{
			In: 614,
		},
	}

	for _, req := range requests {
		log.Println("Sending Request.... : ", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server : %v", err)
	}
	log.Println("ComputeAverage API response, Average: ", resp.GetOut())
}

func FindMaxNumber(serivceClient protobuf.CalculatorServiceClient) {
	fmt.Println("FindMaxNumber API call......")

	requests := [] protobuf.FindMaxNumberIn{
		{
			In: 1123,
		},
		{
			In: 3321,
		},
		{
			In: 1235,
		},
		{
			In: 124,
		},
		{
			In: 8321,
		},
	}

	stream, err := serivceClient.FindMaxNumber(context.Background())
	if err != nil {
		log.Fatalf("Error calling FindMaxNumber API : %v", err)
	}

	//wait channel to block receiver
	waitchan := make(chan struct{})

	go func(requests [] protobuf.FindMaxNumberIn) {
		for _, req := range requests {

			log.Println("Sending Request..... : ", req.GetIn())
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("Error : %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}(requests)

	go func() {
		for {

			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}

			if err != nil {
				log.Fatalf("error receiving response from server : %v", err)
			}

			log.Printf("FindMaxNumber API response, Max : %v\n", resp.GetOut())
		}
	}()

	//block until everything is finished
	<-waitchan
}