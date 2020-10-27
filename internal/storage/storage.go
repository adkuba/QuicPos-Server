package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"google.golang.org/api/option"
)

func uploadSmaller(data []byte, name string) {

	image, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic("Failed to decode image")
	}
	newImage := resize.Resize(224, 224, image, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, newImage, nil)
	if err != nil {
		panic("Failed to encode image")
	}

	result := saveToStorage(name+"_small", buf.Bytes())
	if result == "error" {
		panic("Cannot save smaller")
	}
}

func saveToStorage(imageName string, data []byte) string {

	// Creates a client.
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("QuicPos-key.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	//file upload
	r := bytes.NewReader(data)
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
	return "ok"
}

//UploadFile uploads file to google cloud storage client returns file name or error
func UploadFile(data string) string {

	//check
	if data == "" {
		return ""
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

	imageName := uuid.New().String()
	result := saveToStorage(imageName, unbased)
	if result == "error" {
		panic("Cannot send to storage")
	}

	go uploadSmaller(unbased, imageName)

	return imageName
}
