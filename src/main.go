package main

import (
	"encoding/json"
	"github.com/axm/apollo-utils"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"net/http"
	"time"
)

func main() {
	router := mux.NewRouter()
	app := &apollo.DefaultApp{Router: router}
	registerHandlers(app)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	server.ListenAndServe()
}

const (
	RootUrl = "http://localhost:9090/competitions/premier-league/{host}/{guest}"
)

type leagueHandler func(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp)

func badRequest(w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(reason))
}

func registerHandlers(app *apollo.DefaultApp) {
	wrapper := func(handler leagueHandler, app *apollo.DefaultApp) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, app)
		}
	}
	app.Router.HandleFunc(RootUrl+"/newGoal", wrapper(NewGoal, app)).Methods("POST")
	app.Router.HandleFunc(RootUrl+"/newOffside", wrapper(NewOffside, app)).Methods("POST")
	app.Router.HandleFunc(RootUrl+"/newSubstitution", wrapper(NewSubstitution, app)).Methods("POST")
}

type KafkaProducerConfig struct {
	bootstrapServers string
}

//
func NewGoal(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp) {
	type request struct {
		team   string
		minute uint8
		player uint8
	}
	req := &request{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		badRequest(w, "Unable to read contents")
		return
	}

	config := &KafkaProducerConfig{}

	// publish to kafka
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.bootstrapServers})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer p.Close()

	topic := "mytopic"
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          nil,
	}, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

//
func NewOffside(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp) {
	type request struct {
		team   string
		minute uint8
		player uint8
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		badRequest(w, "Unable to read contents")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func NewSubstitution(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp) {
	type request struct {
		team      string
		inPlayer  uint8
		outPlayer uint8
		minute    uint8
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "Cannot read contents")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// If second yellow card, this will be processed by the event handler
func NewYellowCard(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp) {
	type request struct {
		team   string
		minute uint8
		player uint8
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "Unable to read contents")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func NewRedCard(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp) {
	type request struct {
		team   string
		minute uint8
		player uint8
	}
}

//
func NewPenalty(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp) {
	type request struct {
		team   string
		minute uint8
		player uint8
		reason string
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		badRequest(w, "Unable to read contents")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
