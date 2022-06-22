package webhook

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func PostToWebHook(HTTPMethod, URL, err_key, err_msg string, location string) {

	webhook_url := "https://discordapp.com/api/webhooks/988372011258613791/wCOTksCb0yPIw1XWQlpLF94QOtVb1Nj2TfWy2L9wIybSQJR-3A0YINlgqiqq-u_kDZSN"


	err_message := strings.Replace(err_msg, `"`, " ", -1)
	// time := time.Now().Format("Mon, 02 Jan 2006, 15:04:05")
	temp := map[string]string{
		"content": "```" + "error_key : " + err_key + "\nerror_message : " + err_message + "\nlocation : " + location + "```",
	}

	body, _ := json.Marshal(temp)
	fetch, _ := http.NewRequest("POST", webhook_url, bytes.NewBuffer(body))
	fetch.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(fetch)

	log.Println("Sent To Webhook")
	defer resp.Body.Close()
}
