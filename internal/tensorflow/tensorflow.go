package tensorflow

import (
	"QuicPos/internal/data"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"image"
	"io"
)

type netData struct {
	Text   [1][200]float32
	Image  [1][1280]float32
	Views  [1][1]float32
	Shares [1][1]float32
}

var recommenderModel *tf.SavedModel
var detectorModel *tf.SavedModel
var imageModel *tf.SavedModel
var recommenderDictionary []string
var detectorDictionary []string

//InitModels for recommender and detector and dictionaries
func InitModels() error {

	jsonFile, err := os.Open("./out/recommenderDictionary.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &recommenderDictionary)
	//log.Println(recommenderDictionary)

	jsonFile, err = os.Open("./out/detectorDictionary.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ = ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &detectorDictionary)

	model, err := tf.LoadSavedModel("./out/recommender", []string{"serve"}, nil)
	if err != nil {
		return err
	}
	recommenderModel = model

	model, err = tf.LoadSavedModel("./out/detector", []string{"serve"}, nil)
	if err != nil {
		return err
	}
	detectorModel = model

	model, err = tf.LoadSavedModel("./out/image", []string{"serve"}, nil)
	if err != nil {
		return err
	}
	imageModel = model
	return nil
	//defer model.Session.Close()
}

func getPixels(data io.Reader) ([1][224][224][3]float32, error) {
	img, _, err := image.Decode(data)

	if err != nil {
		return [1][224][224][3]float32{}, err
	}

	width, height := 224, 224

	var converted [1][224][224][3]float32
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			converted[0][y][x][0], converted[0][y][x][1], converted[0][y][x][2] = normalize(int(r), int(g), int(b))
			//check
			if converted[0][y][x][0] > 1 || converted[0][y][x][0] < -1 {
				log.Println("Bad r value")
			}
		}
	}
	return converted, nil
}

//value 0-255
func normalize(r int, g int, b int) (float32, float32, float32) {
	rValue := float32(r / 257)
	gValue := float32(g / 257)
	bValue := float32(b / 257)
	return rValue/float32(127.5) - 1, gValue/float32(127.5) - 1, bValue/float32(127.5) - 1
}

func removeView(s []*data.View, i int) []*data.View {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func remove(s []*string, i int) []*string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func indexOf(tab []string, elem string) int {
	for idx, value := range tab {
		if value == elem {
			return idx
		}
	}
	return -1
}

//from keras documanetation
func textToWordSequence(s string) []string {
	s = strings.ToLower(s)
	for _, filter := range "!\"#$%&()*+,-./:;<=>?@[\\]^_`{|}~\t\n" {
		s = strings.ReplaceAll(s, string(filter), " ")
	}
	splitted := strings.Split(s, " ")
	var ready []string
	for _, value := range splitted {
		value = strings.ReplaceAll(value, " ", "")
		if value != "" {
			ready = append(ready, value)
		}
	}
	return ready
}

func convertPost(post data.Post, recommender bool) netData {

	var netPost netData

	//text convert
	tokens := textToWordSequence(post.Text)
	for i := 0; i < 200; i++ {
		if i < len(tokens) {
			if recommender {
				netPost.Text[0][i] = float32(indexOf(recommenderDictionary, tokens[i]))
			} else {
				netPost.Text[0][i] = float32(indexOf(detectorDictionary, tokens[i]))
			}
		} else {
			netPost.Text[0][i] = -1
		}
	}
	//log.Println(netPost.Text)

	//image convert
	if post.ImageFeatures != nil {
		for i := 0; i < 1280; i++ {
			netPost.Image[0][i] = post.ImageFeatures[i]
		}
	} else {
		for i := 0; i < 1280; i++ {
			netPost.Image[0][i] = 0
		}
	}
	//log.Println(netPost.Image)

	//views convert
	netPost.Views[0][0] = float32(len(post.Views))

	//shares convert
	netPost.Shares[0][0] = float32(len(post.Shares))

	return netPost
}

//InitialReview of post
func InitialReview(post data.Post) (bool, error) {
	result, err := Spam(post)
	if err != nil {
		return false, err
	}

	if result.([][]float32)[0][0] > float32(0.5) {
		return true, nil
	}
	return false, nil
}

//Spam detection
func Spam(post data.Post) (interface{}, error) {

	netPost := convertPost(post, false)

	text, _ := tf.NewTensor(netPost.Text)
	image, _ := tf.NewTensor(netPost.Image)

	result, err := detectorModel.Session.Run(
		map[tf.Output]*tf.Tensor{
			recommenderModel.Graph.Operation("serving_default_input_1").Output(0): text,
			recommenderModel.Graph.Operation("serving_default_input_2").Output(0): image,
		},
		[]tf.Output{
			recommenderModel.Graph.Operation("StatefulPartitionedCall").Output(0),
		},
		nil,
	)

	if err != nil {
		return nil, err
	}

	return result[0].Value(), nil
}

//Recommend post
func Recommend(post data.Post, requesting data.Requesting) (interface{}, error) {

	netPost := convertPost(post, true)

	//requesting user convert
	var requestingUserArray [1][4]float32
	userBytes, _ := hex.DecodeString(requesting.User)
	requestingUserArray[0][0] = float32(binary.BigEndian.Uint32(userBytes[0:4]))
	requestingUserArray[0][1] = float32(binary.BigEndian.Uint32(userBytes[4:8]))
	requestingUserArray[0][2] = float32(binary.BigEndian.Uint32(userBytes[8:12]))
	requestingUserArray[0][3] = float32(binary.BigEndian.Uint32(userBytes[12:16]))
	//log.Println(requestingUserArray[0][0])

	text, _ := tf.NewTensor(netPost.Text)
	image, _ := tf.NewTensor(netPost.Image)
	views, _ := tf.NewTensor(netPost.Views)
	shares, _ := tf.NewTensor(netPost.Shares)
	requestingUser, _ := tf.NewTensor(requestingUserArray)

	result, err := recommenderModel.Session.Run(
		map[tf.Output]*tf.Tensor{
			recommenderModel.Graph.Operation("serving_default_input_1").Output(0): text,
			recommenderModel.Graph.Operation("serving_default_input_2").Output(0): image,
			recommenderModel.Graph.Operation("serving_default_input_3").Output(0): requestingUser,
			recommenderModel.Graph.Operation("serving_default_input_4").Output(0): views,
			recommenderModel.Graph.Operation("serving_default_input_5").Output(0): shares,
		},
		[]tf.Output{
			recommenderModel.Graph.Operation("StatefulPartitionedCall").Output(0),
		},
		nil,
	)

	if err != nil {
		return nil, err
	}

	return result[0].Value(), nil

}
