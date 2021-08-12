package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	pb "example.com/cache-service/z_generated/user"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)


func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	arr := []string{"I","III","IV"}
	for roll := 1569; roll < 1599; roll += 10 {
		for _,class := range arr {
			if _,err := client.SetUser(context.Background(), &pb.User{Name: "John", Class: class, Rollnum: int64(roll)}); err != nil {
				log.Fatalf("Failed to save due to following error: %v", err)
			}
		}
	}

	key := &pb.UserKey{Name: "John", Rollnum: 1579}
	stream, err := client.GetUser(context.Background(), key)
	if err != nil {
		log.Fatalf("Couldn't fetch proto due to following error: %v",err)
	}
	fmt.Println("The following is the proto list:")
	for {
		user, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("The following error has occured: %v", err)
		}

		fmt.Println(user)
	}
}