package appixLogger_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/Travix-International/appix/appixLogger"
	"github.com/alecthomas/assert"
)

// TestAppixLogger_StartStop tests the start & stop of the logger without logging statements
func TestAppixLogger_StartStop(t *testing.T) {
	url := "http://localhost:1/nope"
	sut := appixLogger.NewAppixLogger(url)

	sut.Start()
	sut.Stop()
}

// TestAppixLogger_LogToHttp makes sure that the HTTP target is called ad the default context is added
func TestAppixLogger_LogToHttp_InfoIncludesDefaultContext(t *testing.T) {
	// Set us up to receive logging events
	logServer := &mockHttpLogServer{}
	server := httptest.NewServer(logServer)
	sut := appixLogger.NewAppixLogger(server.URL)

	sut.Start()
	sut.AddMessageToQueue(appixLogger.LoggerNotification{
		Level:    appixLogger.LevelInfo,
		LogEvent: "SomeEvent",
		Message:  "SomeMessage",
	})
	sut.Stop()
	server.Close()

	// Check details of the message that got logged
	assert.Equal(t, 1, logServer.LogCount)
	logMessage := make(map[string]interface{})
	err := json.Unmarshal(logServer.Bodies[0], &logMessage)
	assert.Nil(t, err)
	assert.Equal(t, "Info", logMessage["level"].(string))
	assert.Equal(t, "SomeEvent", logMessage["event"].(string))
	assert.Equal(t, "SomeMessage", logMessage["message"].(string))
	assert.Equal(t, "SomeEvent", logMessage["messageType"].(string))
	assert.Equal(t, "core", logMessage["applicationgroup"].(string))
	assert.Equal(t, "appix", logMessage["applicationname"].(string))
}

// TestAppixLogger_LogToHttp makes sure that the HTTP target is called ad the default context is added
func TestAppixLogger_LogToHttp_ErrorIncludesDefaultContext(t *testing.T) {
	// Set us up to receive logging events
	logServer := &mockHttpLogServer{}
	server := httptest.NewServer(logServer)
	sut := appixLogger.NewAppixLogger(server.URL)

	sut.Start()
	sut.AddMessageToQueue(appixLogger.LoggerNotification{
		Level:    appixLogger.LevelError,
		LogEvent: "SomeEvent",
		Message:  "SomeMessage",
	})
	sut.Stop()
	server.Close()

	// Check details of the message that got logged
	assert.Equal(t, 1, logServer.LogCount)
	logMessage := make(map[string]interface{})
	err := json.Unmarshal(logServer.Bodies[0], &logMessage)
	assert.Nil(t, err)
	assert.Equal(t, "Error", logMessage["level"].(string))
	assert.Equal(t, "SomeEvent", logMessage["event"].(string))
	assert.Equal(t, "SomeMessage", logMessage["message"].(string))
	assert.Equal(t, "SomeEvent", logMessage["messageType"].(string))
	assert.Equal(t, "core", logMessage["applicationgroup"].(string))
	assert.Equal(t, "appix", logMessage["applicationname"].(string))
}

type mockHttpLogServer struct {
	LogCount int
	Bodies   [][]byte
}

func (m *mockHttpLogServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	_, readErr := buf.ReadFrom(r.Body)
	if readErr != nil {
		panic("Missing content")
	}
	body := buf.Bytes()

	m.Bodies = append(m.Bodies, body)
	m.LogCount++
}
