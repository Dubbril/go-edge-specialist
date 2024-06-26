package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
	"time"
)

func LOGGER() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()
}

// LogHandler is a middleware function to log information about incoming requests and outgoing responses.
func LogHandler() gin.HandlerFunc {
	logger := LOGGER()

	return func(c *gin.Context) {
		// Capture request details
		startTime := time.Now()
		var requestBytes []byte
		if c.Request.Body != nil {
			requestBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBytes)) // Reset the request body for further use
		}

		// Get the value of the X-Correlation-ID header
		correlationID := c.GetHeader("X-Correlation-ID")
		requestBody := string(requestBytes)
		reqMarking, err := MaskMiddleThird(requestBody)
		if err != nil {
			reqMarking = requestBody
		}

		logger.Info().
			Str("SERVICE", "EDGE-USER-SERVICE").
			Str("CORRELATION_ID", correlationID).
			Str("METHOD", c.Request.Method).
			Str("URL", c.Request.URL.RequestURI()).
			//Str("USER_AGENT", c.Request.UserAgent()).
			Str("CLIENT_IP", c.ClientIP()).
			Msg(reqMarking)

		// Create a custom response writer
		w := &responseLogger{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}

		// Continue with processing the request
		c.Writer = w
		c.Next()

		// Capture response details
		responseBytes := w.body.Bytes()
		responseStatus := w.status
		duration := time.Since(startTime)

		responseBody := string(responseBytes)
		respMarking, err := MaskMiddleThird(responseBody)
		if err != nil {
			respMarking = responseBody
		}

		logger.Info().
			Str("SERVICE", "EDGE-USER-SERVICE").
			Str("CORRELATION_ID", correlationID).
			Str("RESPONSE_STATUS", fmt.Sprintf("%d", responseStatus)).
			Str("FULL_REQUEST_TIME", fmt.Sprintf("%d", duration)).
			Msg(respMarking)

	}
}

// responseLogger is a custom response writer to capture the response body
type responseLogger struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

// Write is called to write the response body
func (w *responseLogger) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteHeader is called to set the response status code
func (w *responseLogger) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// WriteString is a helper function to write a string to the response body
func (w *responseLogger) WriteString(s string) (int, error) {
	return io.WriteString(w, s)
}

// MaskMiddleThird masks the specified keys by replacing the middle section with asterisks
func MaskMiddleThird(jsonString string) (string, error) {
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return "", err
	}

	// Retrieve keys to be masked from viper
	keysToMask := viper.GetStringSlice("keysToMask")

	// Loop through each key and apply the marking logic
	for _, keyToMask := range keysToMask {
		// Check if the current key exists in the JSON
		if keyVal, ok := data[keyToMask]; ok {
			// Convert the current key to a string
			key, ok := keyVal.(string)
			if !ok {
				return "", fmt.Errorf("value of '%s' key is not a string", keyToMask)
			}

			// If the current key has more than 45 characters, apply new logic
			if len(key) > 45 {
				// Select the first 15 characters
				firstChars := key[:15]

				// Create a mask with 15 asterisks
				mask := strings.Repeat("*", 15)

				// Select the last 15 characters from the back
				lastChars := key[len(key)-15:]

				// Concatenate the parts
				data[keyToMask] = firstChars + mask + lastChars
			} else {
				// If the current key is 45 characters or shorter, use the existing logic
				// Calculate the start and end indices of the middle third
				length := len(key)
				if length < 3 {
					// If the current key is too short, just return the original JSON
					continue
				}

				start := length / 3
				end := 2 * length / 3

				// Extract the middle third of the current key
				middleThird := key[start:end]

				// Create a mask with asterisks
				mask := strings.Repeat("*", len(middleThird))

				// Replace the middle third of the current key with the mask
				data[keyToMask] = key[:start] + mask + key[end:]
			}
		}
	}

	// Marshal the modified map back to JSON
	updatedJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(updatedJSON), nil
}
