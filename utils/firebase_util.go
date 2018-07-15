package utils

import (
	"firebase.google.com/go"
	"os"
	"google.golang.org/appengine"
	"google.golang.org/api/option"
	"context"
	"fmt"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/db"
)

type FirebaseUtil interface {
	InitAuthClient(ctx context.Context) *auth.Client
	InitRTDBClient(ctx context.Context) *db.Client
}

type firebaseUtil struct {
}

func NewFirebaseUtil() FirebaseUtil {
	return &firebaseUtil{}
}

func (u *firebaseUtil) initApp(ctx context.Context) *firebase.App {
	var app *firebase.App
	var err error

	conf := &firebase.Config{
		DatabaseURL: os.Getenv("RTDB_URL"),
	}

	if appengine.IsDevAppServer() {
		opt := option.WithCredentialsFile(os.Getenv("FIREBASE_ADMIN_KEY"))
		app, err = firebase.NewApp(ctx, conf, opt)
	} else {
		app, err = firebase.NewApp(ctx, conf)
	}

	if err != nil {
		panic(fmt.Errorf("failed firebase app: %v", err))
	}

	return app
}

func (u *firebaseUtil) InitAuthClient(ctx context.Context) *auth.Client {
	client, err := u.initApp(ctx).Auth(ctx)
	if err != nil {
		panic(fmt.Errorf("error getting Auth client: %v", err))
	}

	return client
}

func (u *firebaseUtil) InitRTDBClient(ctx context.Context) *db.Client {
	client, err := u.initApp(ctx).Database(ctx)
	if err != nil {
		panic(fmt.Errorf("error getting RTDB client: %v", err))
	}

	return client
}
