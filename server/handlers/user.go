package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	//"go/types"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"

	pb "example.com/cache-service/z_generated/user" // proto package for the structured data
)

var (
	tls	   	   = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file","","The TLS cert file")
	keyFile    = flag.String("key_file","","The TLS key file")
	port       = flag.Int("port",10000,"The server port")
	rdb	   	   = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Configured to connect to port 6379
		DB: 0, // use default DB
	})
)

const (
	dur time.Duration = 300 * time.Second // Expiry time for all the keys in redis cache
	name string = "Samarth"
)

// Struct that implements the CacheServiceServer interface
type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

// SetUser Function to set the key-value pair in cache
func (s *userServiceServer) SetUser(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	err := rdb.Set(ctx, name + ":" + user.Name + ":" + user.Class + ":" + strconv.FormatInt(user.Rollnum, 10), user.Metadata, dur).Err()
	obj := new(pb.Empty)
	if err != nil {
		return obj, err
	}

	return obj, nil
}

// GetUser Function to get User(s) based on name and rollnum
func (s *userServiceServer) GetUser(key *pb.UserKey, stream pb.UserService_GetUserServer) error {
	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, getPattern(key), 0).Iterator()
	for iter.Next(ctx) {
		val, err := rdb.Get(ctx, iter.Val()).Result()
		if err != nil {
			return err
		} else if err == redis.Nil {
			return errors.New("key doesn't exist")
		} else {
			arr := strings.Split(iter.Val(), ":")
			if err := stream.Send(&pb.User{Name: arr[1], Class: arr[2], Rollnum: key.Rollnum, Metadata: []byte(val)}); err != nil {
				return err
			}
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func getPattern(key *pb.UserKey) string {
	return name + ":" + key.Name + ":*:" + strconv.FormatInt(key.Rollnum, 10)
	
}

// Driver code that starts the server and accepts requests.
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
	pb.RegisterUserServiceServer(grpcServer, &userServiceServer{})
	grpcServer.Serve(lis)
}
