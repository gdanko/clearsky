package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gdanko/clearsky/globals"
)

func FetchUrl(url string) (body []byte, err error) {
	if globals.GetDebugFlag() {
		fmt.Println(url)
	}

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
