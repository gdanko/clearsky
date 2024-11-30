package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func FetchUrl(method string, url string, logger *logrus.Logger, payload []byte) (body []byte, err error) {
	logger.Debugf("Fetching %s", url)
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return body, err
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}

	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}
