package ds

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/pkg/errors"
	"os"
	"time"
)

// from https://cloud.google.com/datastore/docs/concepts/entities

type Task struct {
	Category        string
	Done            bool
	Priority        float64
	Description     string `datastore:",noindex"`
	PercentComplete float64
	Created         time.Time
}

func CreateClient(ctx context.Context) (*datastore.Client, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT") // environment variable provided by app engine

	// Creates a client.
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new datastore client")
	}

	return client, nil
}

func CreateEntity(ctx context.Context, client *datastore.Client, task *Task) (*datastore.Key, error) {
	key := datastore.IncompleteKey("Task", nil)
	key, err := client.Put(ctx, key, task)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create datastore entity")
	}
	return key, nil
}

func GetEntity(ctx context.Context, client *datastore.Client, id int64) (Task, error) {
	var task Task
	taskKey := datastore.IDKey("Task", id, nil)
	err := client.Get(ctx, taskKey, &task)
	if err != nil {
		return Task{}, errors.Wrap(err, "failed to get datastore entity")
	}

	return task, nil
}

func DeleteEntity(ctx context.Context, client *datastore.Client, id int64) error {
	key := datastore.IDKey("Task", id, nil)
	err := client.Delete(ctx, key)
	if err != nil {
		return errors.Wrap(err, "failed to delete datastore entity")
	}

	return nil
}
