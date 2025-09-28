package requests

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func Post(url string, reqBody interface{}, responseBody interface{}) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, responseBody); err != nil {
		return err
	}

	return nil
}
