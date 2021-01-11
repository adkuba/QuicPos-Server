package storage

import (
	"QuicPos/internal/data"
	"QuicPos/internal/mongodb"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/api/iterator"
)

//RemoveParentless remove images that don't have parent post
func RemoveParentless() (int, error) {

	result, err := mongodb.PostsCol.Find(context.TODO(), bson.M{"image": bson.M{"$ne": ""}})
	if err != nil {
		return 0, err
	}

	var posts []*data.Post
	if err = result.All(context.TODO(), &posts); err != nil {
		return 0, err
	}

	var images []string
	var toDelete []string
	imagesExceptions := []string{"index_meta.png"}

	for _, post := range posts {
		images = append(images, post.Image)
	}

	it := storageClient.Bucket("quicpos-images").Objects(context.Background(), nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}

		//check if id is in array
		name := attrs.Name
		if name[len(name)-5:] == "small" {
			toDelete = append(toDelete, name)
		}
		if stringInSlice(name, imagesExceptions) {
			continue
		}

		if !stringInSlice(name, images) {
			toDelete = append(toDelete, name)
		}
	}

	for idx, imageToDelete := range toDelete {
		err = deleteFile(imageToDelete)
		if err != nil {
			log.Println(idx, err)
		}
	}

	return len(toDelete), nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func deleteFile(name string) error {

	log.Println("Deleteing: " + name)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	o := storageClient.Bucket("quicpos-images").Object(name)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	return nil
}
