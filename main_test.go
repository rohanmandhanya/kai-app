package main

import (
	"bytes"
	"kai-app/api/controller"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndToEndScanAndQuery(t *testing.T) {

	// Step 1: Test scan endpoint
	scanReqBody := `{"repo": "https://github.com/velancio/vulnerability_scans","files": ["vulnscan19.json","vulscan123.json"]}`
	scanReq, _ := http.NewRequest("POST", "/scan", bytes.NewBuffer([]byte(scanReqBody)))
	scanReq.Header.Set("Content-Type", "application/json")

	scanRecorder := httptest.NewRecorder()
	scanHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller.ScanHandler(w, r)
	})

	scanHandler.ServeHTTP(scanRecorder, scanReq)
	assert.Equal(t, http.StatusOK, scanRecorder.Code)

	// Step 2: Test query endpoint
	queryReqBody := `{"filters": {"severity": "HIGH"}}`
	queryReq, _ := http.NewRequest("POST", "/query", bytes.NewBuffer([]byte(queryReqBody)))
	queryReq.Header.Set("Content-Type", "application/json")

	queryRecorder := httptest.NewRecorder()
	queryHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller.QueryHandler(w, r)
	})

	queryHandler.ServeHTTP(queryRecorder, queryReq)
	assert.Equal(t, http.StatusOK, queryRecorder.Code)
}
