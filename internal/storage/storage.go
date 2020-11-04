package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
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

	//file upload
	r := bytes.NewReader(data)
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := storageClient.Bucket("quicpos-images").Object(imageName).NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		log.Println("io.Copy: ", err)
		return "error"
	}
	if err := wc.Close(); err != nil {
		log.Println("Writer.Close: ", err)
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
