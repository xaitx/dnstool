package internal

import (
	"bytes"
	"io"
	"net/http"
)

func queryDoH(dohURL string, msg []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", dohURL, bytes.NewReader(msg))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/dns-message")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
