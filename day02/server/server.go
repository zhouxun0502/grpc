package main

import (
	"flag"
	"fmt"
	pb "grpc/day02/proto"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

type StreamService struct {
	pb.UnimplementedStreamServiceServer
}

// 服务端流式RPC
func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n < 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			}})

		if err != nil {
			return err
		}
	}

	return nil
}

// 客户端流式RPC
func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			// 全部接受客户端发送来的消息后，通知客户端
			return stream.SendAndClose(&pb.StreamResponse{
				Pt: &pb.StreamPoint{
					Name:  "gRPC Stream Server: Record",
					Value: 1,
				}})
		}
		if err != nil {
			return err
		}

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}
}

// 双向流式RPC
func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "gRPC Stream Client: Route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("stream.Recv pt.Name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value+int32(n))
		n++
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamServiceServer(s, &StreamService{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
