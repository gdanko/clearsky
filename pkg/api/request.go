package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func FetchUrl(method string, url string, headers map[string]string, payload []byte, logger *logrus.Logger) (body []byte, err error) {
	var (
		key   string
		value string
	)
	logger.Debugf("Fetching %s", url)
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return body, err
	}

	for key, value = range headers {
		req.Header.Set(key, value)
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
