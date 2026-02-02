package emailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	c "{{ cookiecutter.module_path }}/utils"
)

type TermiiSender struct {
	APIKey string
	Config c.Config
}

type TermiiPayload struct {
	ApiKey         string `json:"api_key"`
	MessageType    string `json:"message_type"`
	To             string `json:"to"`
	From           string `json:"from"`
	Channel        string `json:"channel"`
	PinAttempts    int    `json:"pin_attempts"`
	PinTimeToLive  int    `json:"pin_time_to_live"`
	PinLength      int    `json:"pin_length"`
	PinPlaceholder string `json:"pin_placeholder"`
	MessageText    string `json:"message_text"`
	PinType        string `json:"pin_type"`
}

type TermiiResData struct {
	SmsStatus    string `json:"smsStatus"`
	PhoneNumber  string `json:"phone_number"`
	To           string `json:"to"`
	PinId        string `json:"pinId"`
	PinId1       string `json:"pin_id"`
	MessageIdStr string `json:"message_id_str"`
	Status       string `json:"status"`
}

func (t *TermiiSender) SendPhoneOTP(identifier string, token string) (string, error) {
	url := t.Config.MailtrapURL
	MessageData := fmt.Sprintf("Your verification pin is  < %s >", token)
	method := "POST"

	payload := TermiiPayload{
		ApiKey:         t.Config.TermiiApiKey,
		MessageType:    "NUMERIC",
		To:             identifier,
		From:           t.Config.SENDER,
		Channel:        "generic",
		PinAttempts:    1,
		PinTimeToLive:  1,
		PinLength:      4,
		PinPlaceholder: token,
		MessageText:    MessageData,
		PinType:        "NUMERIC",
	}
	log.Println(payload)
	// Marshal payload to json
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(string(body))
	msg := fmt.Sprintf(" OTP token sent to the identifier %s", identifier)
	return msg, nil
}

func (t *TermiiSender) SendPhoneWelcomeMessage(identifier string) (string, error) {
	// TODO: Implement later
	return "Welcome message sending not implemented yet", nil
}
