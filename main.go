package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	port = 8080
)

var rieEndpoint = getEnv("RIE_ENDPOINT", "http://localhost:9000/2015-03-31/functions/function/invocations")

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

type RequestContext struct {
	AccountID    string `json:"accountId"`
	APIID        string `json:"apiId"`
	DomainName   string `json:"domainName"`
	DomainPrefix string `json:"domainPrefix"`
	HTTP         struct {
		Method    string `json:"method"`
		Path      string `json:"path"`
		Protocol  string `json:"protocol"`
		SourceIP  string `json:"sourceIp"`
		UserAgent string `json:"userAgent"`
	} `json:"http"`
	RequestID string `json:"requestId"`
	RouteKey  string `json:"routeKey"`
	Stage     string `json:"stage"`
	Time      string `json:"time"`
	TimeEpoch int64  `json:"timeEpoch"`
}

type LambdaEvent struct {
	Version               string            `json:"version"`
	RouteKey              string            `json:"routeKey"`
	RawPath               string            `json:"rawPath"`
	RawQueryString        string            `json:"rawQueryString"`
	Headers               map[string]string `json:"headers"`
	QueryStringParameters map[string]string `json:"queryStringParameters,omitempty"`
	RequestContext        RequestContext    `json:"requestContext"`
	Body                  string            `json:"body,omitempty"`
	IsBase64Encoded       bool              `json:"isBase64Encoded"`
}

type LambdaResponse struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("[Lambda URL Proxy] %s %s\n", r.Method, r.URL.String())

		// Build the Lambda event
		lambdaEvent, err := buildLambdaEvent(r)
		if err != nil {
			http.Error(w, "Failed to build Lambda event: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Forward the event to RIE
		response, err := invokeLambda(lambdaEvent)
		if err != nil {
			http.Error(w, "Failed to invoke Lambda: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		// Map Lambda response to HTTP response
		err = mapLambdaResponseToHTTP(w, response)
		if err != nil {
			http.Error(w, "Failed to process Lambda response: "+err.Error(), http.StatusInternalServerError)
		}
	})

	log.Printf("[Lambda URL Proxy] Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func buildLambdaEvent(r *http.Request) (*LambdaEvent, error) {
	rawQuery := r.URL.RawQuery

	// Query string parameters
	queryParams := map[string]string{}
	for key, values := range r.URL.Query() {
		queryParams[key] = strings.Join(values, ",")
	}

	// Headers
	headers := map[string]string{}
	for key, values := range r.Header {
		headers[strings.ToLower(key)] = values[0]
	}

	// Body and Base64 encoding
	var body string
	var isBase64Encoded bool
	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		// Detect if the body needs Base64 encoding
		contentType := r.Header.Get("Content-Type")
		if !isTextContent(contentType) {
			body = base64.StdEncoding.EncodeToString(bodyBytes)
			isBase64Encoded = true
		} else {
			body = string(bodyBytes)
			isBase64Encoded = false
		}
	}

	// Request context
	now := time.Now().UTC()
	requestContext := RequestContext{
		AccountID:    "anonymous",
		APIID:        "mock-api-id",
		DomainName:   r.Host,
		DomainPrefix: "mock-api-id",
		RequestID:    generateRequestID(),
		RouteKey:     "$default",
		Stage:        "$default",
		Time:         now.Format("2006-01-02T15:04:05.000Z"),
		TimeEpoch:    now.UnixMilli(),
	}
	requestContext.HTTP.Method = r.Method
	requestContext.HTTP.Path = r.URL.Path
	requestContext.HTTP.Protocol = r.Proto
	requestContext.HTTP.SourceIP = r.RemoteAddr
	requestContext.HTTP.UserAgent = r.UserAgent()

	return &LambdaEvent{
		Version:               "2.0",
		RouteKey:              "$default",
		RawPath:               r.URL.Path,
		RawQueryString:        rawQuery,
		Headers:               headers,
		QueryStringParameters: queryParams,
		RequestContext:        requestContext,
		Body:                  body,
		IsBase64Encoded:       isBase64Encoded,
	}, nil
}

func invokeLambda(event *LambdaEvent) (*http.Response, error) {
	eventData, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", rieEndpoint, bytes.NewBuffer(eventData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func mapLambdaResponseToHTTP(w http.ResponseWriter, response *http.Response) error {
	var lambdaResponse LambdaResponse
	err := json.NewDecoder(response.Body).Decode(&lambdaResponse)
	if err != nil {
		return err
	}

	// Set HTTP status code
	w.WriteHeader(lambdaResponse.StatusCode)

	// Set HTTP headers
	for key, value := range lambdaResponse.Headers {
		w.Header().Set(key, value)
	}

	// Decode the body if Base64 encoded
	var body string
	if lambdaResponse.IsBase64Encoded {
		bodyBytes, err := base64.StdEncoding.DecodeString(lambdaResponse.Body)
		if err != nil {
			return err
		}
		body = string(bodyBytes)
	} else {
		body = lambdaResponse.Body
	}

	// Write the body to the response
	_, err = w.Write([]byte(body))
	return err
}

func isTextContent(contentType string) bool {
	return strings.Contains(contentType, "text") || strings.Contains(contentType, "json") || strings.Contains(contentType, "xml")
}

func generateRequestID() string {
	return fmt.Sprintf("req-%s", strings.ReplaceAll(fmt.Sprintf("%x", time.Now().UnixNano()), " ", ""))
}
