package tensorflow

import (
	"fmt"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

var recommenderModel *tf.SavedModel
var detectorModel *tf.SavedModel

//InitModels for recommender and detector
func InitModels() {

	model, err := tf.LoadSavedModel("/home/kuba/Documents/gitfolders/QuicPos-Microservice/out/recommender", []string{"serve"}, nil)
	if err != nil {
		fmt.Printf("Error loading saved model: %s\n", err.Error())
		return
	}
	recommenderModel = model

	model, err = tf.LoadSavedModel("/home/kuba/Documents/gitfolders/QuicPos-Microservice/out/detector", []string{"serve"}, nil)
	if err != nil {
		fmt.Printf("Error loading saved model: %s\n", err.Error())
		return
	}
	detectorModel = model

	//defer model.Session.Close()
}

//Recommend post
func Recommend() {

	text, _ := tf.NewTensor([1][100]float32{})
	user, _ := tf.NewTensor([1][1]float32{})
	reports, _ := tf.NewTensor([1][100]float32{})
	creation, _ := tf.NewTensor([1][1]float32{})
	image, _ := tf.NewTensor([1][224][224][3]float32{})
	views, _ := tf.NewTensor([1][100][6]float32{})
	shares, _ := tf.NewTensor([1][100]float32{})
	requestingUser, _ := tf.NewTensor([1][1]float32{})
	requestingLat, _ := tf.NewTensor([1][1]float32{})
	requestingLong, _ := tf.NewTensor([1][1]float32{})
	requestingTime, _ := tf.NewTensor([1][1]float32{})

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
		fmt.Printf("Error running the session with input, err: %s\n", err.Error())
		return
	}

	fmt.Printf("Result value: %v \n", result[0].Value())

}
