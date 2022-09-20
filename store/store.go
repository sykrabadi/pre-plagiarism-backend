package store

type Store interface {
	DocumentStore() DocumentStore
}

type DocumentStore interface {
	SendData() error
}

// type IMongoStore interface {
// 	SendData() error
// }

// type MongoStore struct {
// 	conn *db.Mongo
// }

// func NewMongoStore(db *db.Mongo) *MongoStore {
// 	return &MongoStore{
// 		conn: db,
// 	}
// }

// func (c *MongoStore) SendData() error {
// 	documentCollection := c.conn.Db.Collection("docs")
// 	_, err := documentCollection.InsertOne(context.Background(), bson.D{
// 		{Key: "name", Value: "TestInsertFromGo"},
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
