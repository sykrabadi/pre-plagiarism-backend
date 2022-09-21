package entrypoint

import (
	"context"
	"go-nsq/store"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewEntryPointService(
	store store.Store,
) IEntryPointService {
	return &EntryPointService{
		DBStore: store,
	}
}

func (c *EntryPointService) SendData() error {
	err := c.DBStore.DocumentStore().SendData()

	if err != nil {
		return err
	}

	return nil
}

func (c *EntryPointService) UpdateData(ctx context.Context, objectID string) error {
	fromHexID, _ := primitive.ObjectIDFromHex(objectID)
	id := primitive.ObjectID.String(fromHexID)
	err := c.DBStore.DocumentStore().UpdateData(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
