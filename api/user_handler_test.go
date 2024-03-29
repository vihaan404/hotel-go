package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vihaan404/hotel-go/db"
	"github.com/vihaan404/hotel-go/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdburi = "mongodb+srv://vihaan2005yadav:fU1LyD9y2QJKZEDF@cluster0.mxzwc.mongodb.net/"
	dbName    = "hotel-reservation-test"
	userColl  = "users"
)

type testDB struct {
	UserStore db.UserStore
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func Setup(t *testing.T) *testDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}
	return &testDB{
		UserStore: db.NewMongoUserStore(client, dbName),
	}
}

func TestPostUser(t *testing.T) {
	tdb := Setup(t)
	defer tdb.teardown(t)

	app := fiber.New()

	userHandler := NewUserHandler(tdb.UserStore)

	app.Post("/", userHandler.HandlePostUser)

	params := types.UserParams{
		FirstName: "Orsted",
		LastName:  "Dragon",
		Email:     "corpo@crop.com",
		Password:  "asdfasdkfjalsdfkj",
	}

	b, err := json.Marshal(&params)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	res, err := app.Test(req, int(time.Second)*2)
	if err != nil {
		t.Error("hello", err)
	}

	fmt.Println(res.Status)
	var user types.User

	json.NewDecoder(res.Body).Decode(&user)

	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName %s but got %s ", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected Last Name %s but got %s ", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email to be  %s but got %s ", params.Email, user.Email)
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected the EncryptedPassword to be not included")
	}
	if len(user.ID) == 0 {
		t.Errorf("expected the user id to be set")
	}
}
