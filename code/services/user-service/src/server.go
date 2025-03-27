package server

import (
	"user-service/src/pb"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	DBClient pb.DBServiceClient
}

// TODO: implement user crud operations
