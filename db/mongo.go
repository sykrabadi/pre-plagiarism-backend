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
	DB     *mongo.Client
	DBName *mongo.Database
}

func loadMongoDBConfig() (string, string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if os.Getenv("MONGODB_USERNAME") == "" {
		return "", "", fmt.Errorf("Environment Variable MONGODB_USERNAME must be set")
	}
	if os.Getenv("MONGODB_PASSWORD") == "" {
		return "", "", fmt.Errorf("Environment Variable MONGODB_PASSWORD must be set")
	}
	if os.Getenv("MONGODB_CLUSTER") == "" {
		return "", "", fmt.Errorf("Environment Variable MONGODB_CLUSTER must be set")
	}
	if os.Getenv("MONGODB_DB_NAME") == "" {
		return "", "", fmt.Errorf("Environment Variable MONGODB_DB_NAME must be set")
	}

	connStr := fmt.Sprintf("mongodb+srv://%s:%s@%s",
		os.Getenv("MONGODB_USERNAME"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_CLUSTER"),
	)

	dbName := fmt.Sprintf("%s", os.Getenv("MONGODB_DB_NAME"))

	return connStr, dbName, nil
}

func InitMongoDB(ctx context.Context) (*Mongo, error) {
	uri, db, err := loadMongoDBConfig()
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

	dbName := client.Database(db)

	// check connectivity via ping
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("An error encountered : %s", err)
	}
	log.Printf("Successfully connected with number of client : %d", client.NumberSessionsInProgress())

	return &Mongo{
		DB:     client,
		DBName: dbName}, nil
}
