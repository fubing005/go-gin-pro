package bootstrap

import (
	"context"
	"shalabing-gin/global"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func InitializeMongoDB() *mongo.Client { //*mongo.Database
	clientOptions := mongoOptions.Client().ApplyURI(global.App.Config.MongoDB.Uri)
	if global.App.Config.MongoDB.Username != "" && global.App.Config.MongoDB.Password != "" {
		credential := mongoOptions.Credential{
			Username: global.App.Config.MongoDB.Username,
			Password: global.App.Config.MongoDB.Password,
		}
		clientOptions.SetAuth(credential)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions) //不会马上建立与 MongoDB 服务器的实际连接，它仅仅是创建并初始化客户端对象
	if err != nil {
		global.App.Log.Error("Failed to create MongoDB client: ", zap.Any("err", err))
	}

	// if err := client.Connect(ctx); err != nil { //用于建立与 MongoDB 服务器的实际连接
	// 	global.App.Log.Error("Failed to connect to MongoDB: ", zap.Any("err", err))
	// }

	// Ping the MongoDB server to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		global.App.Log.Error("Failed to ping MongoDB: ", zap.Any("err", err))
	}

	MongoClient = client
	// MongoDB = client.Database(global.App.Config.MongoDB.Database)
	global.App.Log.Info("Connected to MongoDB successfully")
	return MongoClient //MongoDB
}

// Cleanup MongoDB connection
func CloseMongoDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			global.App.Log.Error("Error disconnecting from MongoDB: ", zap.Any("err", err))
		}
		// global.App.Log.Info("Disconnected from MongoDB")
	}
}
