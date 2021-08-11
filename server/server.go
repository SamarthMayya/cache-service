package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
	//"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	"github.com/go-redis/redis/v8"

	pb "example.com/cache-service/cache"
)

var (
	tls	   	   = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file","","The TLS cert file")
	keyFile    = flag.String("key_file","","The TLS key file")
	port       = flag.Int("port",10000,"The server port")
	rdb	   	   = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: 0, // use default DB
	})
)

const (
	dur time.Duration = 300 * time.Second
	name = "Samarth"
)

type cacheServiceServer struct {
	pb.UnimplementedCacheServiceServer
}

func (s *cacheServiceServer) Set(ctx context.Context, cache_entry *pb.CacheEntry) (*pb.Empty, error) {
	err := rdb.Set(ctx, name + ":" + cache_entry.Key.Key, cache_entry.Val.Value, dur).Err()
	obj := new(pb.Empty)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

func (s *cacheServiceServer) Get(ctx context.Context, key *pb.Key) (*pb.Value, error) {
	val, err := rdb.Get(ctx, name + ":" + key.Key).Result()
	if err != nil {
		return nil, err
	}

	return &pb.Value{Value: []byte(val)}, nil
}

func main() {
	flag.Parse()
        lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
        if err != nil {
                log.Fatalf("failed to listen: %v", err)
        }
        var opts []grpc.ServerOption
        if *tls {
                if *certFile == "" {
                        *certFile = data.Path("x509/server_cert.pem")
                }
                if *keyFile == "" {
                        *keyFile = data.Path("x509/server_key.pem")
                }
                creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
                if err != nil {
                        log.Fatalf("Failed to generate credentials %v", err)
                }
                opts = []grpc.ServerOption{grpc.Creds(creds)}
        }
        grpcServer := grpc.NewServer(opts...)
        pb.RegisterCacheServiceServer(grpcServer, &cacheServiceServer{})
        grpcServer.Serve(lis)

}
