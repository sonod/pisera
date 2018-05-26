package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sonod/pisera/config"
)

type httpError struct {
	endpoint   string
	statusCode int
}

func (e *httpError) Error() string {
	return fmt.Sprintf("http error occurred, url: %s, code: %d", e.endpoint, e.statusCode)
}

func PHPIPAMRequest(conf config.Config, method string, token string, user string, path string, param interface{}) ([]byte, error) {
	endpoint := conf.Server.Endpoint + "/api/"
	pass := conf.Server.Password

	paramJson, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	u := endpoint + path

	tls := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest(method, u, bytes.NewReader(paramJson))
	if err != nil {
		return nil, err
	}
	if len(user) != 0 {
		req.SetBasicAuth(user, pass)
	}
	if len(token) != 0 {
		req.Header.Set("Token", token)
		req.Header.Set("Content-type", "application/json")
	}

	res, err := tls.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, &httpError{
			endpoint:   req.URL.String(),
			statusCode: res.StatusCode,
		}
	}

	return ioutil.ReadAll(res.Body)
}
