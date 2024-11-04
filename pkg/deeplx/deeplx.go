package deeplx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const url = "https://deeplx.mingming.dev/translate"

type request struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

type Response struct {
	Code int64  `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

func Translate(text, sourceLang, targetLang string) (string, error) {
	if len(text) == 0 {
		return "", nil
	}

	if len(sourceLang) == 0 {
		sourceLang = "auto"
	}

	if len(targetLang) == 0 {
		targetLang = "EN"
	}

	req := &request{
		Text:       text,
		SourceLang: sourceLang,
		TargetLang: targetLang,
	}
	jsonBody, _ := json.Marshal(req)

	var body []byte

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	response, err := client.Post(url, "application/json", strings.NewReader(string(jsonBody)))

	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	responseData := &Response{}

	if err = json.Unmarshal(body, responseData); err != nil {
		return "", err
	}

	if responseData.Code != 200 {
		return "", fmt.Errorf("Failed to translate: %v", responseData.Msg)
	}

	return responseData.Data, nil
}
