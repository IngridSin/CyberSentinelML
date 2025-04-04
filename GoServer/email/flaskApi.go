package email

import (
	"bytes"
	"encoding/json"
	"goServer/config"
	"io"
	"log"
	"net/http"
)

func triggerMLPrediction(messageID string) {
	payload := map[string]string{"message_id": messageID}
	jsonData, _ := json.Marshal(payload)

	fullApi := config.FlaskUrl + config.PredictAPI

	println("JSOIN:", jsonData, fullApi)

	resp, err := http.Post(fullApi, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error calling Flask /predict: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Flask /predict returned non-200 status: %d", resp.StatusCode)
	}

	respBody := new(bytes.Buffer)
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		return
	}

	log.Printf("Flask response for %s: %s", messageID, respBody.String())
}
