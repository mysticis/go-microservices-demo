package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	pb "github.com/mysticis/logistic-service-consignment/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultFilePath = "consignment.json"
	address         = "localhost:50051"
)

func parseFile(file string) (*pb.Consignment, error) {

	var consignment *pb.Consignment

	data, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &consignment)

	return consignment, err
}

func main() {
	//set up a connection to the server
	clientConn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}

	defer clientConn.Close()

	client := pb.NewShippingServiceClient(clientConn)

	//contact server and print out contents
	file := defaultFilePath

	if (len(os.Args)) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)

	if err != nil {
		log.Fatalf("could not parse file: %v", err)
	}

	resp, err := client.CreateConsignment(context.Background(), consignment)

	if err != nil {
		log.Fatalf("could not create consignment: %v", err)
	}

	log.Printf("Created Consignment: %v", resp.Created)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})

	if err != nil {
		log.Fatalf("could not list consignments: %v", err)
	}

	for _, item := range getAll.Consignments {
		log.Println(item)
	}
}
