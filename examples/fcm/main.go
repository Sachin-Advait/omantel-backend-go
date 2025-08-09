package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

	"google.golang.org/api/option"
)

func main() {

	//https://medium.com/readytowork-org/push-notification-with-firebase-go-and-react-b52910923988

	app, err := connection()
	if err != nil {
		log.Fatal("failed fcm connection: " + err.Error())
	}

	_, err = app.Messaging(context.TODO())
	if err != nil {
		log.Fatalf("messaging: %s", err)
	}

}

func connection() (*firebase.App, error) {

	opt := option.WithCredentialsFile("./incentive-test-bdb73-firebase-adminsdk-ekmc6-454deae12a.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	// app, err := firebase.NewApp(ctx, nil, opts...)
	// if err != nil {
	// 	return nil, err
	// }

	return app, nil
}

func sendMessageToSpecificDevice(fcmClient *messaging.Client) {

	registrationToken := "BPEkuW9iLU_WDgGKzY5HalO0zvcvY7RQ9G0rPrDzYvlHsYA5-iOe68M_lWcDNkfbrd4Gk6vvTBk7Aui0upw35WI"

	message := &messaging.Message{
		Data: map[string]string{
			"desc":  "testing incentive fcm",
			"score": "850",
			"time":  "2:45",
		},
		Token: registrationToken,
	}

	response, err := fcmClient.Send(context.TODO(), message)
	if err != nil {
		log.Fatalln("sendMessageToSpecificDevice: send: " + err.Error())
	}
	// Response is a message ID string.
	fmt.Println("sendMessageToSpecificDevice: Successfully sent message: ", response)
	//Successfully sent message: projects/incentive-test-bdb73/messages/7795480998665223986
}

func sendMessageToMultipleDevice(fcmClient *messaging.Client) {

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: "Congratulations!!",
			Body:  "You have just implemented push notification",
		},
		Data: map[string]string{
			"desc":  "testing incentive fcm",
			"score": "850",
			"time":  "2:45",
		},
		Tokens: []string{},
	}

	response, err := fcmClient.SendMulticast(context.TODO(), message)
	if err != nil {
		log.Fatalln("sendMessageToMultipleDevice: send: " + err.Error())
	}
	// Response is a message ID string.
	fmt.Println("sendMessageToMultipleDevice: Successfully sent message:", response)
}

func sendViaTopicMessage(fcmClient *messaging.Client) {

	topic := "test-incentive1"

	message := &messaging.Message{
		Data: map[string]string{
			"desc":  "testing incentive fcm",
			"score": "850",
			"time":  "2:45",
		},
		Topic: topic,
	}

	response, err := fcmClient.Send(context.TODO(), message)
	if err != nil {
		log.Fatalln("sendViaTopicMessage: send: " + err.Error())
	}
	// Response is a message ID string.
	fmt.Println("sendViaTopicMessage: Successfully sent message:", response)
}

func fetchPushNotification(fcmClient *messaging.Client) {

}

func addingSDK() {
	/*
		npm install firebase
	*/

	/*
			// Import the functions you need from the SDKs you need
		import { initializeApp } from "firebase/app";
		import { getAnalytics } from "firebase/analytics";
		// TODO: Add SDKs for Firebase products that you want to use
		// https://firebase.google.com/docs/web/setup#available-libraries

		// Your web app's Firebase configuration
		// For Firebase JS SDK v7.20.0 and later, measurementId is optional
		const firebaseConfig = {
		  apiKey: "AIzaSyD0KpI0SIHLl5AYW-sLzdf3y4ZbMY5qlZI",
		  authDomain: "incentive-test-bdb73.firebaseapp.com",
		  projectId: "incentive-test-bdb73",
		  storageBucket: "incentive-test-bdb73.appspot.com",
		  messagingSenderId: "1043387257847",
		  appId: "1:1043387257847:web:71c72d626d242d3b8d94f2",
		  measurementId: "G-KVLSMKBSDW"
		};

		// Initialize Firebase
		const app = initializeApp(firebaseConfig);
		const analytics = getAnalytics(app);
	*/
}
