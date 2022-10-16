package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type process struct {
	ProcessName string `json:"processName" bson:"processName"`
}
type processDeep struct {
	ProcessName    string `json:"processName" bson:"processName"`
	IssuedQuantity int    `json:"issuedQuantity" bson:"issuedQuantity" default:"0"`
	Rejected       int    `json:"rejected" bson:"rejected" default:"0"`
	Accepted       int    `json:"accepted" bson:"accepted" default:"0"`
}

type batches struct {
	BatchName       string        `json:"batchName" bson:"batchName"`
	IssuedQuantityB int           `json:"issuedQuantityB" bson:"issuedQuantityB" default:"0"`
	Process         []processDeep `json:"process" bson:"process"`
}

type Component struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	CompanyName string             `json:"companyName" bson:"companyName"`
	Process     []process          `json:"process" bson:"process"`
	Batches     []batches          `json:"batches" bson:"batches"`
}

const dbName = "test"
const colName = "components"

// * MOST IMPORTANT
var collection *mongo.Collection

// go file execution nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run cmd/MyProgram/main.go
// connect with mongodb
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	connectionString := os.Getenv("MONGO_URI")
	// client options
	clientOptions := options.Client().ApplyURI(connectionString)
	// connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongodb connection success")
	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance is ready")
}

// get all components from db
func getAllComponents() []bson.M {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var movies []bson.M

	for cursor.Next(context.Background()) {
		var movie bson.M
		if err := cursor.Decode(&movie); err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	return movies
}

// create a component in db
func createComponent(component Component) {
	inserted, err := collection.InsertOne(context.Background(), component)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(inserted.InsertedID)
}

// delete a component from db
func deleteComponent(componentId string) {
	_id, err := primitive.ObjectIDFromHex(componentId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": _id}
	deleteCount, err2 := collection.DeleteOne(context.Background(), filter)
	if err2 != nil {
		log.Fatal(err)
	}
	fmt.Println(deleteCount)
}

// delete a batch from db
func deleteABatch(componentName string, batchName string) {
	fmt.Println(componentName, batchName)

	update := bson.M{"$pull": bson.M{"batches": bson.M{"batchName": batchName}}}
	var updatedDocument bson.M

	err := collection.FindOneAndUpdate(context.Background(), bson.M{"name": componentName}, update).Decode(&updatedDocument)
	if err != nil {
		fmt.Println(err)
		fmt.Println("error")
		log.Fatal(err)
	}
	fmt.Println(updatedDocument)
}

// update batch in db
func updateBatch(componentName string, batches []batches) {
	fmt.Println(batches)

	update := bson.M{"$set": bson.M{"batches": batches}}

	err := collection.FindOneAndUpdate(context.Background(), bson.M{"name": componentName}, update)
	if err != nil {
		log.Fatal(err)
	}
}

// update process in db
func updateProcess(componentName string, batchName string, process []processDeep) {
	identifier := []interface{}{bson.M{"batch.batchName": batchName}}

	update := bson.M{"$set": bson.M{"batches.$[batch].process": process}}
	opts := options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: identifier})

	err := collection.FindOneAndUpdate(context.Background(), bson.M{"name": componentName}, update, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(identifier)
}

// controller to get all components
func GetAll(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/x-www-form-urlencode")
	response.Header().Set("Allow-Control-Allow-Methods", "GET")

	movies := getAllComponents()
	json.NewEncoder(response).Encode(movies)
}

// controller to create a component
func CreateComponent(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/x-www-form-urlencode")
	response.Header().Set("Allow-Control-Allow-Methods", "POST")

	var component Component
	_ = json.NewDecoder(request.Body).Decode(&component)
	createComponent(component)
}

// controller to delete a component
func DeleteComponent(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(request)
	deleteComponent(params["id"])
}

func DeleteABatch(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Allow-Control-Allow-Methods", "DELETE")
	params := mux.Vars(request)
	deleteABatch(params["componentName"], params["batchName"])

}

func UpdateBatch(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(request)

	var batches []batches

	_ = json.NewDecoder(request.Body).Decode(&batches)

	updateBatch(params["componentName"], batches)
}

func UpdateProcess(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Allow-Control-Allow-Methods", "PUT")
	params := mux.Vars(request)

	var process []processDeep

	_ = json.NewDecoder(request.Body).Decode(&process)

	updateProcess(params["componentName"], params["batchName"], process)
}
