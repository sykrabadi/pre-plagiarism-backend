package prefalsification

import "go-nsq/store"

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
