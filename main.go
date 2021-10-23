package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Kick struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Company   string             `json:"company,omitempty" bson:"company,omitempty"`
	Timestamp string             `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	Longitude float64            `json:"longitude,omitempty" bson:"longitude,omitempty"`
	Attitude  float64            `json:"attitude,omitempty" bson:"attitude,omitempty"`
	Velocity  float64            `json:"velocity,omitempty" bson:"velocity,omitempty"`
}

func GetKick(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var kick Kick
	collection := client.Database("mongodb").Collection("kicks")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Kick{ID: id}).Decode(&kick)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(kick)
}

func CreateKick(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var kick Kick

	err := json.NewDecoder(request.Body).Decode(&kick)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Println("Bad Json")
		return
	}

	collection := client.Database("mongodb").Collection("kicks")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	now := time.Now()

	kick.Timestamp = now.Format("2006-01-02 15:04:05")

	result, err := collection.InsertOne(ctx, kick)

	if err != nil {
		response.WriteHeader(http.StatusNotAcceptable)
	}

	err = json.NewEncoder(response).Encode(result)
	if err != nil {
		response.WriteHeader(http.StatusConflict)
	}
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/kick", CreateKick).Methods("POST")
	http.ListenAndServe(":8080", router)
}
