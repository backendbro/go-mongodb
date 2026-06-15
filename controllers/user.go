package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/backendbro/go-mongodb/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	Client *mongo.Client
}

func NewUserController(c *mongo.Client) *UserController {
	return &UserController{c}
}

// GET USER
func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	collection := uc.Client.Database("mongo-golang").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := models.User{}

	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

// CREATE USER
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u.Id = primitive.NewObjectID()

	collection := uc.Client.Database("mongo-golang").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

// DELETE USER
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	collection := uc.Client.Database("mongo-golang").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil || res.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "DELETED USER %s\n", id)
}
