package api

import (
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func FetchUrl(url string, logger *logrus.Logger) (body []byte, err error) {
	logger.Debug(url)
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
