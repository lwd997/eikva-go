package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	envvars "eikva.ru/eikva/env_vars"
)

type PostConfig struct {
	Url      string
	ReqBody  interface{}
	RespBody interface{}
	Headers  *map[string]string
}

func Post(config *PostConfig) error {
	jsonData, err := json.Marshal(config.ReqBody)
	if err != nil {
		return err
	}

	client := &http.Client{}
	noSSL := envvars.Get(envvars.NoSSLVerify) == "1"

	if noSSL {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	req, err := http.NewRequest("POST", config.Url, bytes.NewBuffer(jsonData));
	if err != nil {
		return err
	}

	for name, value := range *config.Headers {
		req.Header.Set(name, value)
	}

	resp, err := client.Do(req)
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

	if err := json.Unmarshal(body, &config.RespBody); err != nil {
		return err
	}

	return nil
}
