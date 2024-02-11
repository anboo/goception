package goception

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type Actor struct {
	client        *http.Client
	baseURL       string
	lastReq       *http.Request
	lastRes       *http.Response
	customHeaders map[string]string
	t             *testing.T
}

func NewActor(t *testing.T) *Actor {
	return &Actor{
		client: &http.Client{},
		t:      t,
	}
}

func (a *Actor) SetBaseURL(baseUrl string) *Actor {
	a.baseURL = baseUrl
	return a
}

func (a *Actor) HaveHeader(header, value string) *Actor {
	if a.customHeaders == nil {
		a.customHeaders = make(map[string]string)
	}
	a.customHeaders[header] = value
	return a
}

func (a *Actor) SendGet(path string) *Actor {
	return a.sendRequest("GET", path, nil)
}

func (a *Actor) SendPatch(path string, body interface{}) *Actor {
	return a.sendRequest("PATCH", path, body)
}

func (a *Actor) SendPost(path string, body interface{}) *Actor {
	return a.sendRequest("POST", path, body)
}

func (a *Actor) SendDelete(path string) *Actor {
	return a.sendRequest("DELETE", path, nil)
}

func (a *Actor) SendPut(path string, body interface{}) *Actor {
	return a.sendRequest("PUT", path, body)
}

func (a *Actor) sendRequest(method, path string, body interface{}) *Actor {
	url := a.baseURL + path

	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			a.t.Fatalf("failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		a.t.Fatalf("failed to create request: %v", err)
	}

	for header, value := range a.customHeaders {
		req.Header.Set(header, value)
	}

	resp, err := a.client.Do(req)
	if err != nil {
		a.t.Fatalf("failed to send request: %v", err)
	}

	a.lastReq = req
	a.lastRes = resp
	return a
}

func (a *Actor) ExpectResponse(expectedStatus int) *Actor {
	if a.lastRes.StatusCode != expectedStatus {
		a.t.Fatalf("unexpected status code: expected %d, got %d", expectedStatus, a.lastRes.StatusCode)
	}
	return a
}

func (a *Actor) ParseJSON(v interface{}) *Actor {
	defer a.lastRes.Body.Close()

	body, err := ioutil.ReadAll(a.lastRes.Body)
	if err != nil {
		a.t.Fatalf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		a.t.Fatalf("failed to parse JSON: %v", err)
	}

	return a
}

func (a *Actor) ParseFieldFromJSONPath(jsonPath string, to interface{}) *Actor {
	jsonData := make(map[string]interface{})
	err := json.NewDecoder(a.lastRes.Body).Decode(&jsonData)
	if err != nil {
		a.t.Fatalf("failed to decode JSON: %v", err)
	}

	parts := strings.Split(jsonPath, ".")
	current := jsonData
	for _, part := range parts {
		val, ok := current[part].(map[string]interface{})
		if !ok {
			a.t.Fatalf("field %s not found or is not an object in JSON response", part)
		}
		current = val
	}

	err = getDataByPath(jsonData, jsonPath, to)
	if err != nil {
		a.t.Fatalf("failed get %s json path: %w", jsonPath, err)
	}

	return a
}

func getDataByPath(s map[string]interface{}, path string, to interface{}) error {
	parts := strings.Split(path, ".")
	current := s
	for _, part := range parts {
		val, ok := current[part]
		if !ok {
			return fmt.Errorf("field %s not found in data", part)
		}
		if m, ok := val.(map[string]interface{}); ok {
			current = m
		} else {
			switch t := to.(type) {
			case *string:
				if s, ok := val.(string); ok {
					*t = s
				} else {
					return fmt.Errorf("failed to convert value to string")
				}
			case *int:
				if f, ok := val.(float64); ok {
					*t = int(f)
				} else {
					return fmt.Errorf("failed to convert value to int")
				}
			case *float64:
				if f, ok := val.(float64); ok {
					*t = f
				} else {
					return fmt.Errorf("failed to convert value to float64")
				}
			case *bool:
				if b, ok := val.(bool); ok {
					*t = b
				} else {
					return fmt.Errorf("failed to convert value to bool")
				}
			default:
				return fmt.Errorf("unsupported type for destination variable: %T", to)
			}
			return nil
		}
	}
	return fmt.Errorf("path %s does not point to a leaf node", path)
}
