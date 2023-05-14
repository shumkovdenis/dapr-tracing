package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	pb.UnimplementedGreeterServer
}

func runGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &grpcServer{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *grpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())

	out := fmt.Sprintf("%s %s", in.GetName(), serviceName)

	if !disableCall {
		out = callGRPCService(ctx, out)
	}

	return &pb.HelloReply{Message: out}, nil
}

func callGRPCService(ctx context.Context, msg string) string {
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%s", daprGrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx = metadata.AppendToOutgoingContext(ctx, "dapr-app-id", calledService)

	md, _ := metadata.FromIncomingContext(ctx)
	traceContext := md["grpc-trace-bin"]
	log.Println("traceContext", traceContext)

	if len(traceContext) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", traceContext[0])
	}

	client := pb.NewGreeterClient(conn)

	res, err := client.SayHello(ctx, &pb.HelloRequest{Name: msg})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return res.GetMessage()
}
