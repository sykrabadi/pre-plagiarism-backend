package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	Db     *mongo.Database
}

// Commented code ode below is used to connect to mongodb atlas
func loadMongoDBConfig() (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	// if os.Getenv("MONGODB_USERNAME") == "" {
	// 	return "", "", fmt.Errorf("Environment Variable MONGODB_USERNAME must be set")
	// }
	// if os.Getenv("MONGODB_PASSWORD") == "" {
	// 	return "", "", fmt.Errorf("Environment Variable MONGODB_PASSWORD must be set")
	// }
	// if os.Getenv("MONGODB_CLUSTER") == "" {
	// 	return "", "", fmt.Errorf("Environment Variable MONGODB_CLUSTER must be set")
	// }
	if os.Getenv("MONGODB_DB_NAME") == "" {
		return "", fmt.Errorf("Environment Variable MONGODB_DB_NAME must be set")
	}

	// connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s",
	// 	os.Getenv("MONGODB_USERNAME"),
	// 	os.Getenv("MONGODB_PASSWORD"),
	// 	os.Getenv("MONGODB_CLUSTER"),
	// )

	dbName := fmt.Sprintf("%s", os.Getenv("MONGODB_DB_NAME"))

	return dbName, nil
}

func InitMongoDB(ctx context.Context) (*Mongo, error) {
	db, err := loadMongoDBConfig()
	if db == "" {
		log.Fatal("You must set your 'MONGODB_DB_NAME' environmental variable.")
	}
	if err != nil {
		log.Fatalf("An error encountered with error message %s", err)
	}
	err = godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("[InitMongoDB] unable to load env file with error %v \n", err)
		return nil, err
	}
	mongodbContainerAddr := os.Getenv("MONGODB_CONTAINER_ADDRESS")
	var mongodbAddr string
	if mongodbContainerAddr == ""{
		mongodbAddr = "localhost"
	}else{
		mongodbAddr = mongodbContainerAddr
	}
	uri := fmt.Sprintf("mongodb://%v:27017", mongodbAddr)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("[InitMongoDB] unable to connect to mongodb with error : %s", err)
	}

	dbName := client.Database(db)

	// check connectivity via ping
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("[InitMongoDB] unable to ping to mongodb with error : %s", err)
	}
	log.Printf("Successfully connected with number of client : %d", client.NumberSessionsInProgress())

	return &Mongo{
		Client: client,
		Db:     dbName}, nil
}
