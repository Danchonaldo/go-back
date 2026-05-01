package handlers

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

func SendNotification(message string) (string, error) {
	client := resty.New()

	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		fmt.Println("[Resty] Sending request:", req.URL)
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		fmt.Println("[Resty] Response:", resp.Status())
		return nil
	})

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"message": message,
		}).
		Post("http://notification-service:8083/notify")

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}
