package firebaseconfig

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/storage"
	"google.golang.org/api/option"
)

var Credential = map[string]string{
	"type": "service_account",
    "project_id": "project1-9c741",
    "private_key_id": "52c3b59ceeb67cefecd6e264317ded95b34fdaa8",
    "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCgMNbPdVT8GjEF\n1TBZlxa7cais/CLSNToZW0zmQJi27eP2LBbaL7PTEjqMaTSkMK9coz3kfm1lUc/Q\nRVMvQX9KX2M6Zv5YZsR5iyaPk0B8SPij/JJyLp6wzZ5vo9ZFVxHNdzfSsvjuuWZ5\nnfAQanEpra+OCY+IRpx7ZtPOTX6i7bw9J5tQQWKjtlYj8szTRGBA9ezgJx0VpyLV\nSfRZ3/duPe8khfoVwYYxOf+JZZKoM9CYH7Fht8n/WK9sv8fBos2ZCHbMeRNJ4zKJ\nnay2kkSF/j+3se8bSaYtLsPzAuCE+UwDKvkL/OfI49nrXrsEh/H8C+KyqxU7KHCP\nTfvPqVd9AgMBAAECggEAG0se+oN5CqgyxJXlkkowLRlJPkoKa8r2XnYkJNIKQxbG\nwb2S2jOI/dRMBfrx5WzHIC9PHxSccogtMoTxTqqn88XLXBrKyMifvr7C//D8qoBr\nXtNsp0hCsQijr01Yt4FFjv66U+u7Br+vklai3cUpCUsmz3pZCoTPaiYjvfGNu1gl\ntGpjTlIdFbN8ldF2/uDf1A1HFTuyinU9uzbDIc342QrqjqxiquhEb4M9/YFHzmGJ\noZ+biNtICKwSUV/k4VKVbrFgdCsiTywpO664/iYLI60qWe3QFWRlhebWukC7itcU\nI3Laa7Ny3PZGoS0P5JcM2Mu9FllmTT/ec24PwucGMQKBgQDPlFAdQV2+V9y/iOlM\neMxcP9M63pD3SO4cpDwZAeIu+ZQ0VqdvvJsIU48yu+JAxXxlGPZ3Bvr4lPuCjYHv\ni7t8P9iJRF9bPKbKe3U4QpdqiJC15PzovoHMmMHGM6ZSKww0HLdLWTVQyR0mXanu\n9/B0OW+uadCdnnffz5qj35TZrQKBgQDFjrKw1EkwoFHr6DkMhbAPTwrKqcFkRp1z\neHPI3RfImrnYNHuQ2MsrK5hYJgH66r/+eRvYzGYIhoVL4Yj6vPxaO1Af15tK80Fo\nE3pZXmB3JmJjzbGo/xNtnQPJFu9XF3dJlKmpS5d40ZxsXmNJeZ/d0xqHypZlUzId\n8Fr/plTPEQKBgCpOcwmTYfTCWYZb9BiW4ifHBlkQNYGAxq7lti3umVeznEYShyuZ\nyAspZJ0Vd+Z1mYXNUaYusQXq71vLVXkqBHstVAed/MVOljvcb6aYw919mejIk9cn\nxLKbS5sGudYzWdhhJeZgdyJQ6vT/z/uZYtN9RUrJ3C/TtWDTIhRWJDM9AoGAO0RV\nDUGSd3kROpGfU8djiyQuW6BTuU6J+9M+ARt7AB6S9G7CRzESum/Is2ErpOf6E1Cn\n7wFa9uHMaXhOzNIXbsZugi7/kpupmpyFTvxGOliUmdZinw1u+apqMVClGt6aVcO4\npmZcGc3gfI5QgQFw2W50fhpjxiAQX/T0h8+Rj9ECgYBG12c5pnTA7GgtRs6jc7OT\nDstCRY8gUCVKHnyk6ssZifokPc4dw4YG5bInUB6uloGNbIYhQlKlOSGZqYonx20b\n2UyPI7nON+bTWreFVp7fM9LRZ2qmakiI+z9+dJFUP/J2O3zioakHU6VPWoYaoY7+\nm6MreGPFNsqJG6yrYnu+4A==\n-----END PRIVATE KEY-----\n",
    "client_email": "firebase-adminsdk-tfxf1@project1-9c741.iam.gserviceaccount.com",
    "client_id": "118304371762752988321",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-tfxf1%40project1-9c741.iam.gserviceaccount.com",
}


func FirebaseInit() *firebase.App{
	var cd []byte

	json.Marshal(Credential)
	
	// opt := option.WithCredentialsFile("D:/Go/src/Golang-Fiber/serviceAccountKey.json")
	opt := option.WithCredentialsJSON(cd)
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.SetFlags(16)
		log.Printf("error initializing app: %v\n", err)
	}
	return app
}

func CloudFirestore() *firestore.Client {
	app := FirebaseInit()
	ctx := context.Background()
	client,err := app.Firestore(ctx)
	if err !=nil{
		log.SetFlags(16)
		log.Printf("error initializing client: %v\n",err)
	}
	return client
}

func AuthClient() *auth.Client{
	app := FirebaseInit()
	client,err := app.Auth(context.Background())
	if err != nil{
		log.SetFlags(16)
		log.Printf("error getting Auth client: %v\n",err)
	}
	return client
}

func StorageClient() *storage.Client{
	var cd []byte
	var config *firebase.Config

	cd,_= json.Marshal(Credential)
	config = &firebase.Config{
		StorageBucket: "project1-9c741.appspot.com",
	}

	opt := option.WithCredentialsJSON(cd)
	app, err := firebase.NewApp(context.Background(),config,opt)
	if err !=nil{
		log.SetFlags(16)
		log.Println(err)
	}

	client,err := app.Storage(context.Background())
	if err!= nil{
		log.SetFlags(16)
		log.Println(err)
	}

	return client
}