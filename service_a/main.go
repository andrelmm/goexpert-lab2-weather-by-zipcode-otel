package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/andrelmm/goexpert-lab2-weather-by-zipcode-otel/shared"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

var tracer trace.Tracer

func validateCep(cep interface{}) bool {
	cepStr, ok := cep.(string)
	if !ok {
		return false
	}
	match, _ := regexp.MatchString(`^\d{8}$`, cepStr)
	return match
}

func forwardToServiceB(ctx context.Context, cep string) (int, string, error) {
	serviceBBaseURL := os.Getenv("SERVICE_B_BASE_URL")
	if serviceBBaseURL == "" {
		serviceBBaseURL = "http://localhost:8081"
	}
	serviceBURL := fmt.Sprintf("%s/weather?zip=%s", serviceBBaseURL, cep)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serviceBURL, nil)
	if err != nil {
		return 0, "", err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}
	return resp.StatusCode, string(body), nil
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {

	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := tracer.Start(ctx, "service_a_HandleRequest")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var cepRequest map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&cepRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cep, exists := cepRequest["cep"]
	if !exists || !validateCep(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	cepStr := cep.(string)
	statusCode, response, err := forwardToServiceB(ctx, cepStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(response))
}

func main() {
	ctx := context.Background()
	shutdown, err := shared.InitProvider("service_a", "otel-collector:4317")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer = otel.Tracer("microservice-tracer")
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http.HandleFunc("/weather", HandleRequest)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
