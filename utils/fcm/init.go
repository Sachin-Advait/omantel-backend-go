package fcm

import (
	"context"
	"database/sql"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"

	"google.golang.org/api/option"
)

type FCM interface {
	Send(ctx context.Context, message *messaging.Message) (string, error)
}

// dbClient is the default client
type fcmClient struct {
	*sql.DB
	*messaging.Client
}

var client *fcmClient

func InitFCM(configFile string) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	opt := option.WithCredentialsFile(configFile)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("FCMService - setup failed")
	}

	fcm, err := app.Messaging(context.TODO())
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("FCMService - message instance failed")
	}

	// now create the database
	client = &fcmClient{Client: fcm}
}

// Get is used to get the fcm instance client
func Get() FCM {
	return client
}

// Close is used to close the fcm instance
func Close() error {
	return client.Close()
}
