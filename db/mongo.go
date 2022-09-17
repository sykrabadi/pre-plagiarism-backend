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
}

func loadMongoDBConfig() (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if os.Getenv("MONGODB_USERNAME") == "" {
		return "", fmt.Errorf("Environment Variable MONGODB_USERNAME must be set")
	}
	if os.Getenv("MONGODB_PASSWORD") == "" {
		return "", fmt.Errorf("Environment Variable MONGODB_PASSWORD must be set")
	}
	if os.Getenv("MONGODB_CLUSTER") == "" {
		return "", fmt.Errorf("Environment Variable MONGODB_CLUSTER must be set")
	}

	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s",
		os.Getenv("MONGODB_USERNAME"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_CLUSTER"),
	)

	return connStr, nil
}

func InitMongoDB(ctx context.Context) (*mongo.Client, error) {

	uri, err := loadMongoDBConfig()
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	if err != nil {
		log.Fatalf("An error encountered with error message %s", err)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("An error encountered : %s", err)
	}

	// check connectivity via ping
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("An error encountered : %s", err)
	}

	return client, nil
}
