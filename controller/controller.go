package controller

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	model "github.com/anudeep652/mongoapi/modal"
// 	"github.com/gorilla/mux"
// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// const dbName = "netflix"
// const colName = "watchlist"

// // * MOST IMPORTANT
// var collection *mongo.Collection

// // connect with mongodb
// func init() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	connectionString := os.Getenv("MONGO_URI")
// 	// client options
// 	clientOptions := options.Client().ApplyURI(connectionString)

// 	// connect to mongodb
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Mongodb connection success")

// 	collection = client.Database(dbName).Collection(colName)
// 	// collection instance
// 	fmt.Println("Collection instance is ready")
// }

// // insert one record
// func insertOneMovie(movie model.NetFlix) {
// 	inserted, err := collection.InsertOne(context.Background(), movie)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
// }

// func getAMovie(movieId string) (model.NetFlix, error) {
// 	id, err := primitive.ObjectIDFromHex(movieId)
// 	if err != nil {
// 		return model.NetFlix{ID: primitive.NewObjectID(), Movie: "", Watched: false}, err
// 	}
// 	var movie model.NetFlix
// 	filter := bson.M{"_id": id}

// 	err = collection.FindOne(context.Background(), filter).Decode(&movie)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return movie, nil
// }

// // update one record
// func updateOneMovie(movieId string) {
// 	id, _ := primitive.ObjectIDFromHex(movieId)
// 	filter := bson.M{"_id": id}
// 	update := bson.M{"$set": bson.M{"watched": true}}

// 	//result contains how many value are updated
// 	result, err := collection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("modified count:", result.ModifiedCount)
// }

// // delete one record
// func deleteOneMovie(movieId string) {
// 	_id, _ := primitive.ObjectIDFromHex(movieId)
// 	filter := bson.M{"_id": _id}

// 	deleteCount, err := collection.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Movie got deleted with delete count", deleteCount)
// }

// // delete all records
// func deleteAllMovies() int64 {
// 	// empty parenthesis in bson.D means delete all records

// 	// deleteResult.DeletedCount return how many records deleted
// 	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("No of movies deleted: ", deleteResult.DeletedCount)
// 	return deleteResult.DeletedCount
// }

// // get all movies
// func getAllMovies() []bson.M {
// 	cursor, err := collection.Find(context.Background(), bson.D{{}})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var movies []bson.M

// 	// 1st way
// 	if err = cursor.All(context.Background(), &movies); err != nil {
// 		log.Fatal(err)
// 	}

// 	// 2nd way to get all records from cursor
// 	// for cursor.Next(context.Background()) {
// 	// 	var movie bson.M
// 	// 	err := cursor.Decode(&movie)
// 	// 	if err != nil {
// 	// 		log.Fatal(err)
// 	// 	}
// 	// 	movies = append(movies, movie)
// 	// }

// 	// * important to close
// 	defer cursor.Close(context.Background())
// 	return movies
// }

// // actual controller

// func GetAMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	movie, err := getAMovie(params["id"])
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode("No documents found")
// 		return
// 	}
// 	// fmt.Print(movie)
// 	json.NewEncoder(w).Encode(movie)
// }
// func GetAllMovies(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
// 	allMovies := getAllMovies()
// 	json.NewEncoder(w).Encode(allMovies)
// }

// func CreateMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
// 	w.Header().Set("Allow-Control-Allow-Methods", "POST")
// 	var movie model.NetFlix
// 	_ = json.NewDecoder(r.Body).Decode(&movie)
// 	insertOneMovie(movie)
// 	json.NewEncoder(w).Encode(movie)
// }

// func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
// 	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

// 	params := mux.Vars(r)
// 	updateOneMovie(params["id"])
// 	json.NewEncoder(w).Encode(params["id"])
// }

// func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
// 	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

// 	params := mux.Vars(r)
// 	deleteOneMovie(params["id"])
// 	json.NewEncoder(w).Encode(params["id"])
// }

// func DeleteAllMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
// 	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

// 	count := deleteAllMovies()
// 	json.NewEncoder(w).Encode(count)

// }
