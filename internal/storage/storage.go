package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

var storageClient *storage.Client
var ctx context.Context

//InitStorage client
func InitStorage() {
	ctx = context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("QuicPos-key.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	} else {
		storageClient = client
	}
}

//ReadFile from storage
func ReadFile(fileName string) (data []byte) {

	rc, err := storageClient.Bucket("quicpos-images").Object(fileName).NewReader(ctx)
	if err != nil {
		log.Println("readFile: unable to open file from bucket, file", fileName, err)
		return
	}
	defer rc.Close()

	data, err = ioutil.ReadAll(rc)
	if err != nil {
		log.Println("readFile: unable to read data from bucket, file", fileName, err)
		return
	}

	return data
}

func saveToStorage(imageName string, data []byte) error {

	//file upload
	r := bytes.NewReader(data)
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := storageClient.Bucket("quicpos-images").Object(imageName).NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

//UploadFile uploads file to google cloud storage client returns file name or error
func UploadFile(data string) (string, error) {

	//check
	if data == "" {
		return "", nil
	}

	//file decoding
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		return "", errors.New("InvalidImage")
	}
	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		return "", errors.New("Cannot decode b64")
	}

	imageName := uuid.New().String()
	result := saveToStorage(imageName, unbased)
	if result != nil {
		return "", errors.New("Cannot send to storage")
	}

	return imageName, nil
}
