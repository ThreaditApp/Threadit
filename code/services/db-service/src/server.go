package server

import (
	dbpb "gen/db-service/pb"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBServer struct {
	dbpb.UnimplementedDBServiceServer
	Mongo *mongo.Database
}
