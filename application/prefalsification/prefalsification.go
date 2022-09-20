package prefalsification

import (
	"context"
	"go-nsq/store"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewPrefalsificationService(
	store store.Store,
) IPrefalsificationService {
	return &PrefalsificationService{
		DBStore: store,
	}
}

func (c *PrefalsificationService) SendData() error {
	err := c.DBStore.DocumentStore().SendData()

	if err != nil {
		return err
	}

	return nil
}

func (c *PrefalsificationService) UpdateData(ctx context.Context, objectID string) error {
	fromHexID, _ := primitive.ObjectIDFromHex(objectID)
	id := primitive.ObjectID.String(fromHexID)
	err := c.DBStore.DocumentStore().UpdateData(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
