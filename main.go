package main

import (
	"context"
	"net/http"
	"time"

	"github.com/backendbro/go-mongodb/controllers"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	r := httprouter.New()

	client := getSession()
	uc := controllers.NewUserController(client)

	r.GET("/user/:id", uc.GetUsers)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	http.ListenAndServe("localhost:8080", r)
}

var urlString = "mongodb+srv://zubiverse:naruto123@cluster0.wihkseb.mongodb.net/"

func getSession() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(urlString))
	if err != nil {
		panic(err)
	}

	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return client
}
