package tensorflow

import (
	"QuicPos/internal/storage"
	"errors"
	"github.com/nfnt/resize"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"bytes"
	"image"
	"image/jpeg"
)

//GenerateImageFeatures based on mobilenet
func GenerateImageFeatures(imageString string) ([]float32, error) {
	imageData := storage.ReadFile(imageString)

	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	originalImage, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return []float32{}, err
	}

	newImage := resize.Resize(224, 224, originalImage, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, newImage, nil)
	if err != nil {
		return []float32{}, errors.New("Can't encode image")
	}
	pixels, _ := getPixels(buf)

	imageTensor, _ := tf.NewTensor(pixels)

	//lock
	mutex.Lock()

	result, err := imageModel.Session.Run(
		map[tf.Output]*tf.Tensor{
			imageModel.Graph.Operation("serving_default_input_1").Output(0): imageTensor,
		},
		[]tf.Output{
			imageModel.Graph.Operation("StatefulPartitionedCall").Output(0),
		},
		nil,
	)

	//unlock
	mutex.Unlock()

	if err != nil {
		return []float32{}, err
	}

	arrayFeatures := result[0].Value().([][]float32)

	return arrayFeatures[0], nil
}
