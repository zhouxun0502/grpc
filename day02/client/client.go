package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc/day02/proto"
	"io"
	"log"
	"time"
)

var (
	addr = flag.String("addr", "localhost:50052", "the address to connect to")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := printList(ctx, client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: List ", Value: 2022}}); err != nil {
		log.Fatalf("printList.err: %v", err)
	}

	if err := printRecord(ctx, client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Record ", Value: 2022}}); err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}

	if err := printRoute(ctx, client, &pb.StreamRequest{Pt: &pb.StreamPoint{Name: "gRPC Stream Client: Route ", Value: 2022}}); err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}
}

func printList(ctx context.Context, client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.List(ctx, r)
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	return nil
}

func printRecord(ctx context.Context, client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Record(ctx)
	if err != nil {
		return err
	}

	for n := 0; n < 6; n++ {
		r.Pt.Value = r.Pt.Value + 1
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)

	return nil
}

func printRoute(ctx context.Context, client pb.StreamServiceClient, r *pb.StreamRequest) error {
	stream, err := client.Route(ctx)
	if err != nil {
		return err
	}

	for n := 0; n <= 6; n++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	if err := stream.CloseSend(); err != nil {
		return err
	}

	return nil
}
