package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
        "google.golang.org/grpc/examples/data"
	pb "example.com/cache-service/cache"
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
	client := pb.NewCacheServiceClient(conn)
	cacheentry := &pb.CacheEntry{Val: &pb.Value{Value: []byte("Bello")}, Key: &pb.Key{Key: "Hello"}}
	_, err = client.Set(context.Background(), cacheentry)
	if err != nil {
		log.Printf("An error has been encountered: %v",err)
	}

	key := &pb.Key{Key: "Hello"}
	val, err := client.Get(context.Background(), key)
	if err != nil {
		log.Fatalf("An error has occurred while fetching: %v",err)
	}
	log.Println(val)
}


