package server

import (
	dbpb "gen/db-service/pb"
	userpb "gen/user-service/pb"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	DBClient dbpb.DBServiceClient
}

// TODO: implement user crud operations
