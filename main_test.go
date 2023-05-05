package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetRedis1(t *testing.T) {
	router := gin.Default()
	router.GET("/:tenant/:key", GetRedis)

	req, _ := http.NewRequest("GET", "/tenant1/sample1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}

func TestPutRedis1(t *testing.T) {
	router := gin.Default()
	router.POST("/:tenant", PutRedis)

	data := Data{
		Key:   "sample1",
		Value: map[string]string{"foo": "bar1"},
	}

	body, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "/tenant1", bytes.NewReader(body))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
}
func TestGetRedis(t *testing.T) {
	sourceFileName := "test/tenant1.json"
	jsonFile, err := os.Open(sourceFileName)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", sourceFileName, err)
	}
	defer jsonFile.Close()

	content, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", sourceFileName, err)
	}

	var data Data
	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal %s: %v", sourceFileName, err)
	}

	router := gin.Default()
	router.POST("/:tenant", PutRedis)
	router.GET("/:tenant/:key", GetRedis)

	// Store the data in Redis using a POST request
	body, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "/tenant1", bytes.NewReader(body))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Fetch the data from Redis using a GET request
	req, _ = http.NewRequest("GET", fmt.Sprintf("/tenant1/%s", data.Key), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// Compare the input JSON content with the actual REST request output
	responseBytes := resp.Body.Bytes()

	// Marshal the Go value back into JSON and remove possible extra whitespace
	expectedResponseBytes, _ := json.Marshal(data.Value)
	expectedResponse := string(expectedResponseBytes)

	fmt.Println("get-Expected File: ", expectedResponse)
	fmt.Println("get-Actual File: ", string(responseBytes))
	t.Logf("t-Expected File: %s", expectedResponse)
	t.Logf("t-Actual File: %s", string(responseBytes))

	assert.JSONEq(t, expectedResponse, string(responseBytes))
}

func TestPutRedis(t *testing.T) {
	router := gin.Default()
	router.POST("/:tenant", PutRedis)
	router.GET("/:tenant/:key", GetRedis)

	sourceFileName := "test/tenant1.json"
	jsonFile, err := os.Open(sourceFileName)
	if err != nil {
		t.Fatalf("Failed to open %s: %v", sourceFileName, err)
	}
	defer jsonFile.Close()

	content, err := io.ReadAll(jsonFile)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", sourceFileName, err)
	}

	var data Data
	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal %s: %v", sourceFileName, err)
	}

	body, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", "/tenant1", bytes.NewReader(body))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// Fetch the data from Redis using a GET request
	req, _ = http.NewRequest("GET", fmt.Sprintf("/tenant1/%s", data.Key), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)

	// Compare the input JSON content with the actual REST request output
	responseBytes := resp.Body.Bytes()

	// Marshal the Go value back into JSON and remove possible extra whitespace
	expectedResponseBytes, _ := json.Marshal(data.Value)
	expectedResponse := string(expectedResponseBytes)
	fmt.Println("put-Expected File: ", expectedResponse)
	fmt.Println("put-Actual File: ", string(responseBytes))
	t.Logf("t-put-Expected File: %s", expectedResponse)
	t.Logf("t-put-Actual File: %s", string(responseBytes))

	assert.JSONEq(t, expectedResponse, string(responseBytes))
}
