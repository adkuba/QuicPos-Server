package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

//UploadFile uploads file to google cloud storage client returns file name or error
func UploadFile(data string) string {

	//check
	if data == "" {
		return ""
	}

	// Creates a client.
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("QuicPos-key.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	//file decoding
	idx := strings.Index(data, ";base64,")
	if idx < 0 {
		panic("InvalidImage")
	}
	unbased, err := base64.StdEncoding.DecodeString(data[idx+8:])
	if err != nil {
		panic("Cannot decode b64")
	}
	r := bytes.NewReader(unbased)
	imageName := uuid.New().String()

	//file upload
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := storageClient.Bucket("quicpos-images").Object(imageName).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		log.Fatalf("io.Copy: %v", err)
		return "error"
	}
	if err := wc.Close(); err != nil {
		log.Fatalf("Writer.Close: %v", err)
		return "error"
	}

	return imageName
}
