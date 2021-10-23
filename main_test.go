package main

import (
    "io/ioutil"
    "go.mongodb.org/mongo-driver/bson"
    "encoding/json"
    "os"
	"context"
	"log"
	 "time"
	"go.mongodb.org/mongo-driver/mongo/options"
    "net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/gorilla/mux"
	"testing"
)

func init() {
    finished := make(chan bool)
    go listen(finished)
    <-finished

    connect()
}

func listen(finished chan bool) {
    router := mux.NewRouter()
    router.HandleFunc("/kick", CreateKick).Methods("POST")

    finished <- true
    http.ListenAndServe(":8080", router)

}

func connect(){
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, _ = mongo.Connect(ctx, clientOptions)

    err := client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatalf("MongoDB connection error: %v", err)
    }
}

func sendJson(url string) {
    //reading initial struct
    var initial Kick
    file, _ := ioutil.ReadFile("test.json")
    _ = json.Unmarshal([]byte(file), &initial)

    //posting json
    r, err := os.Open("test.json")
    resp, err := http.Post(url, "application/json", r)
    r.Close()

    //checking errors
    if err != nil {
        log.Printf("Can not Post json")
    }

    if resp.StatusCode == 400 {
        log.Printf("Bad json: %v\n", r)
    }

    if resp.StatusCode == 406 {
        log.Printf("Cannot Insert document")
    }
    if resp.StatusCode == 409 {
        log.Printf("Cannot Encode json")
    }

    //reading last inserted document from BD
    var kick Kick
    collection := client.Database("mongodb").Collection("kicks")
    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

    cursor, err := collection.Find(ctx, bson.M{})
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        cursor.Decode(&kick)
    }

    //checking if the initial struct equals the readed one
    if (kick.Company != initial.Company) ||
        (kick.Longitude != initial.Longitude) ||
        (kick.Attitude != initial.Attitude) ||
        (kick.Velocity != initial.Velocity) {
            log.Println("Document Inserted Uncorrectly")
    }

    //delete testing document from BD
    collection.DeleteOne(ctx, bson.M{"_id": kick.ID})

}

func TestMain(t *testing.T) {
    sendJson("http://localhost:8080/kick")
}

