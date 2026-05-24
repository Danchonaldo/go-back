package handlers

import (
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

var restyClient = resty.New()

func getNotificationURL() string {
	url := os.Getenv("NOTIFICATION_SERVICE_URL")
	if url == "" {
		url = "http://localhost:8083"
	}
	return url
}

func SendNotification(message string) (string, error) {
	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"message": message,
		}).
		SetResult(map[string]interface{}{}).
		Post(getNotificationURL() + "/notify")

	if err != nil {
		log.Printf("Notification error: %v", err)
		return "", err
	}

	log.Printf("Notification sent: status=%d, body=%s", resp.StatusCode(), resp.String())
	return resp.String(), nil
}

func GetNotificationHealth(c interface{}) {
	resp, err := restyClient.R().
		Get(getNotificationURL() + "/health")

	if err != nil {
		log.Printf("Notification service health check failed: %v", err)
		return
	}

	log.Printf("Notification service health: %s", resp.String())
}
