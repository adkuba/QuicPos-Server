package tensorflow

import (
	"QuicPos/internal/data"
	"QuicPos/internal/storage"
	"log"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"bytes"
	"image"
	"image/jpeg"
	"io"
	"math/rand"
)

type netData struct {
	Text     [1][400]float32
	User     [1][1]float32
	Reports  [1][100]float32
	Creation [1][1]float32
	Image    [1][224][224][3]float32
	Views    [1][100][6]float32
	Shares   [1][100]float32
}

var recommenderModel *tf.SavedModel
var detectorModel *tf.SavedModel
var imageModel *tf.SavedModel

//InitModels for recommender and detector
func InitModels() error {

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

func convertPost(post data.Post) netData {

	var netPost netData

	//text convert
	for i := 0; i < 400; i++ {
		char := float32(0)
		if i < len(post.Text) {
			char = float32(int(post.Text[i]))
		}
		netPost.Text[0][i] = char
	}

	//user convert
	netPost.User[0][0] = float32(0)

	//reports convert
	for i := 0; i < 100; i++ {
		if len(post.Reports) > 0 {
			randomIndex := rand.Intn(len(post.Reports))
			netPost.Reports[0][i] = float32(0)
			post.Reports = remove(post.Reports, randomIndex)
		} else {
			netPost.Reports[0][i] = 0
		}
	}

	//creation convert
	netPost.Creation[0][0] = float32(float64(post.CreationTime.Unix()) / float64(100000))

	//image convert
	if post.Image != "" {
		image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
		imageData := storage.ReadFile(post.Image + "_small")
		imageReader := bytes.NewReader(imageData)
		netPost.Image, _ = getPixels(imageReader)
	}

	//views convert
	for i := 0; i < 100; i++ {
		if len(post.Views) > 0 {
			randomIndex := rand.Intn(len(post.Views))
			netPost.Views[0][i][0] = float32(0)
			netPost.Views[0][i][1] = float32(0)
			netPost.Views[0][i][2] = float32(0)
			netPost.Views[0][i][3] = float32(0)
			netPost.Views[0][i][4] = float32(float64(post.Views[randomIndex].Date.Unix()) / float64(100000))
			netPost.Views[0][i][5] = float32(post.Views[randomIndex].Time)
			post.Views = removeView(post.Views, randomIndex)
		} else {
			netPost.Views[0][i][0] = 0
			netPost.Views[0][i][1] = 0
			netPost.Views[0][i][2] = 0
			netPost.Views[0][i][3] = 0
			netPost.Views[0][i][4] = 0
			netPost.Views[0][i][5] = 0
		}
	}

	//shares convert
	for i := 0; i < 100; i++ {
		if len(post.Shares) > 0 {
			randomIndex := rand.Intn(len(post.Shares))
			netPost.Shares[0][i] = float32(0)
			post.Shares = remove(post.Shares, randomIndex)
		} else {
			netPost.Shares[0][i] = 0
		}
	}

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

	netPost := convertPost(post)

	text, _ := tf.NewTensor(netPost.Text)
	user, _ := tf.NewTensor(netPost.User)
	reports, _ := tf.NewTensor(netPost.Reports)
	creation, _ := tf.NewTensor(netPost.Creation)
	image, _ := tf.NewTensor(netPost.Image)
	views, _ := tf.NewTensor(netPost.Views)
	shares, _ := tf.NewTensor(netPost.Shares)

	result, err := detectorModel.Session.Run(
		map[tf.Output]*tf.Tensor{
			recommenderModel.Graph.Operation("serving_default_input_1").Output(0): text,
			recommenderModel.Graph.Operation("serving_default_input_2").Output(0): user,
			recommenderModel.Graph.Operation("serving_default_input_3").Output(0): reports,
			recommenderModel.Graph.Operation("serving_default_input_4").Output(0): creation,
			recommenderModel.Graph.Operation("serving_default_input_5").Output(0): image,
			recommenderModel.Graph.Operation("serving_default_input_6").Output(0): views,
			recommenderModel.Graph.Operation("serving_default_input_7").Output(0): shares,
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

	netPost := convertPost(post)

	//requesting user convert
	var requestingUserArray [1][1]float32
	requestingUserArray[0][0] = float32(0)

	//requesting lat convert
	var requestingLatArray [1][1]float32
	requestingLatArray[0][0] = float32(0)

	//requesting long convert
	var requestingLongArray [1][1]float32
	requestingLongArray[0][0] = float32(0)

	//requesting time convert
	var requestingTimeArray [1][1]float32
	requestingTimeArray[0][0] = float32(float64(requesting.Date.Unix()) / float64(100000))

	text, _ := tf.NewTensor(netPost.Text)
	user, _ := tf.NewTensor(netPost.User)
	reports, _ := tf.NewTensor(netPost.Reports)
	creation, _ := tf.NewTensor(netPost.Creation)
	image, _ := tf.NewTensor(netPost.Image)
	views, _ := tf.NewTensor(netPost.Views)
	shares, _ := tf.NewTensor(netPost.Shares)
	requestingUser, _ := tf.NewTensor(requestingUserArray)
	requestingLat, _ := tf.NewTensor(requestingLatArray)
	requestingLong, _ := tf.NewTensor(requestingLongArray)
	requestingTime, _ := tf.NewTensor(requestingTimeArray)

	//log.Println(netPost, requestingLatArray, requestingLongArray, requestingUserArray, requestingTimeArray)

	result, err := recommenderModel.Session.Run(
		map[tf.Output]*tf.Tensor{
			recommenderModel.Graph.Operation("serving_default_input_1").Output(0):  text,
			recommenderModel.Graph.Operation("serving_default_input_2").Output(0):  user,
			recommenderModel.Graph.Operation("serving_default_input_3").Output(0):  reports,
			recommenderModel.Graph.Operation("serving_default_input_4").Output(0):  creation,
			recommenderModel.Graph.Operation("serving_default_input_5").Output(0):  image,
			recommenderModel.Graph.Operation("serving_default_input_6").Output(0):  views,
			recommenderModel.Graph.Operation("serving_default_input_7").Output(0):  shares,
			recommenderModel.Graph.Operation("serving_default_input_8").Output(0):  requestingUser,
			recommenderModel.Graph.Operation("serving_default_input_9").Output(0):  requestingLat,
			recommenderModel.Graph.Operation("serving_default_input_10").Output(0): requestingLong,
			recommenderModel.Graph.Operation("serving_default_input_11").Output(0): requestingTime,
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
