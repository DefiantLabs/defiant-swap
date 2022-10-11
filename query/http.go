package query

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func PostJson(url string, request interface{}, response interface{}, jwt *JWT) (error, int) {
	b, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err, -1
	}
	req.Header.Set("content-type", "application/json")
	if jwt != nil {
		req.Header.Set("Authorization", jwt.Token)
	}

	r, err := myClient.Do(req)
	if err != nil {
		return err, -1
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(response), r.StatusCode
}
