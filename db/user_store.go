package db

import (
	"context"
	"fmt"

	"github.com/vihaan404/hotel-go/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dburi    = "mongodb+srv://vihaan2005yadav:fU1LyD9y2QJKZEDF@cluster0.mxzwc.mongodb.net/"
	dbName   = "hotel-reservation"
	userColl = "users"
)

type Dropper interface {
	Drop(ctx context.Context) error
}
type UserStore interface {
	Dropper
	UpdateUser(ctx context.Context, id string, params types.UpdatedUserParams) error
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, dbName string) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbName).Collection(userColl),
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	users := []*types.User{}
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id string, params types.UpdatedUserParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	values := params.MappingBson()

	update := bson.D{{"$set", values}}

	_, err = s.coll.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("---------Dropping-----------")
	err := s.coll.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}
