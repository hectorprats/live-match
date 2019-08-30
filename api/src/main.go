package main

import (
	"encoding/json"
	"fmt"
	"github.com/axm/apollo-utils"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting live match...")
	router := mux.NewRouter()
	app := &apollo.DefaultApp{Router: router}
	app.Rabbit = &apollo.RabbitConnection{
		User:     "guest",
		Password: "guest",
		Host:     "livematch_rabbit",
		Port:     5672,
	}
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
	RootUrl = "/competitions/premier-league/matches/{Host}/{Guest}"
)

type leagueHandler func(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings)

func badRequest(w http.ResponseWriter, reason string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(reason))
}

type wrapper func(handlerMap map[string]leagueHandler, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) http.HandlerFunc

func registerHandlers(app *apollo.DefaultApp) {
	wrapper := func(handlerMap map[string]leagueHandler, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			version := r.Header.Get("version")
			if version == "" {
				badRequest(w, "Missing version header.")
				return
			}
			handler := handlerMap[version]
			if handler == nil {
				badRequest(w, "Invalid version.")
				return
			}
			handler(w, r, app, pub)
		}
	}

	registerV1Handlers(app, wrapper)
}
func registerV1Handlers(app *apollo.DefaultApp, wrapper wrapper) {
	registerNewGoalV1(app, wrapper)
	registerNewOffisdeV1(app, wrapper)
	registerNewSubstitutionV1(app, wrapper)
	registerNewYellowCardV1(app, wrapper)
	registerNewRedCardV1(app, wrapper)
	registerNewPenaltyV1(app, wrapper)
	registerStartMatchV1(app, wrapper)
}

func registerNewPenaltyV1(app *apollo.DefaultApp, wrapper wrapper) {
	newPenalty := &apollo.RabbitPublisherSettings{
		Queue:      "LiveMatch.NewPenalty",
		Exchange:   "LiveMatch.NewPenalty",
		RoutingKey: "LiveMatch.NewPenalty",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	newPenaltyMap := make(map[string]leagueHandler)
	newPenaltyMap["1.0"] = NewPenalty
	app.Router.HandleFunc(RootUrl+"/NewPenalty", wrapper(newPenaltyMap, app, newPenalty)).Methods("POST")
}

func registerNewSubstitutionV1(app *apollo.DefaultApp, wrapper wrapper) {
	newSubstitution := &apollo.RabbitPublisherSettings{
		Queue:      "LiveMatch.NewSubstitution",
		Exchange:   "LiveMatch.NewSubstitution",
		RoutingKey: "LiveMatch.NewSubstitution",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	newSubstitutionMap := make(map[string]leagueHandler)
	newSubstitutionMap["1.0"] = NewSubstitution
	app.Router.HandleFunc(RootUrl+"/newSubstitution", wrapper(newSubstitutionMap, app, newSubstitution)).Methods("POST")
}

func registerNewRedCardV1(app *apollo.DefaultApp, wrapper wrapper) {
	newRedCard := &apollo.RabbitPublisherSettings{
		Queue:      "LiveMatch.NewRedCard",
		Exchange:   "LiveMatch.NewRedCard",
		RoutingKey: "LiveMatch.NewRedCard",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	newRedCardMap := make(map[string]leagueHandler)
	newRedCardMap["1.0"] = NewRedCard
	app.Router.HandleFunc(RootUrl+"/NewRedCard", wrapper(newRedCardMap, app, newRedCard)).Methods("POST")
}

func registerNewYellowCardV1(app *apollo.DefaultApp, wrapper wrapper) {
	newYellowCard := &apollo.RabbitPublisherSettings{
		Queue:      "LiveMatch.NewYellowCard",
		Exchange:   "LiveMatch.NewYellowCard",
		RoutingKey: "LiveMatch.NewYellowCard",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	newYellowCardMap := make(map[string]leagueHandler)
	newYellowCardMap["1.0"] = NewYellowCard
	app.Router.HandleFunc(RootUrl+"/NewYellowCard", wrapper(newYellowCardMap, app, newYellowCard)).Methods("POST")
}

func registerNewGoalV1(app *apollo.DefaultApp, wrapper wrapper) {
	newGoal := &apollo.RabbitPublisherSettings{
		Queue:      "LiveMatch.NewGoal",
		Exchange:   "LiveMatch.NewGoal",
		RoutingKey: "LiveMatch.NewGoal",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	newGoalVersionMap := make(map[string]leagueHandler)
	newGoalVersionMap["1.0"] = NewGoal
	app.Router.HandleFunc(RootUrl+"/newGoal", wrapper(newGoalVersionMap, app, newGoal)).Methods("POST")
}

//
func NewGoal(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) {
	fmt.Println("Received new goal")
	type request struct {
		Team   string `json:"Team"`
		Minute uint8  `json:"Minute"`
		Player uint8  `json:"Player"`
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		badRequest(w, "Unable to read contents")
		return
	}
	var req request
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		badRequest(w, "Unable to read contents")
		return
	}

	type newGoalEvent struct {
		Host   string
		Guest  string
		Team   string
		Minute uint8
		Player uint8
	}
	vars := mux.Vars(r)
	newGoal := &newGoalEvent{
		Host:   vars["Host"],
		Guest:  vars["Guest"],
		Team:   req.Team,
		Minute: req.Minute,
		Player: req.Player,
	}
	fmt.Println(fmt.Sprintf("Host: %s, Guest: %s", newGoal.Host, newGoal.Guest))
	bytes, err = json.Marshal(newGoal)
	fmt.Println("Sending payload: " + string(bytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rc := app.Rabbit
	addr := fmt.Sprintf("amqp://%s:%s@%s:%d/", rc.User, rc.Password, rc.Host, rc.Port)
	fmt.Println(addr)

	err = apollo.PublishMessage(app.Rabbit, pub, &bytes)
	if err != nil {
		fmt.Println("Unable to publish message")
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func registerNewOffisdeV1(app *apollo.DefaultApp, wrapper wrapper) {
	newOffside := &apollo.RabbitPublisherSettings{
		Queue:      "LiveMatch.NewOffside",
		Exchange:   "LiveMatch.NewOffside",
		RoutingKey: "LiveMatch.NewOffside",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	newOffsideMap := make(map[string]leagueHandler)
	newOffsideMap["1.0"] = NewOffside
	app.Router.HandleFunc(RootUrl+"/newOffside", wrapper(newOffsideMap, app, newOffside)).Methods("POST")
}

//
func NewOffside(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) {
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

	bytes, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = apollo.PublishMessage(app.Rabbit, pub, &bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func NewSubstitution(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) {
	type request struct {
		team      string
		inPlayer  uint8
		outPlayer uint8
		minute    uint8
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		badRequest(w, "Unable to read contents")
		return
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = apollo.PublishMessage(app.Rabbit, pub, &bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// If second yellow card, this will be processed by the event handler
func NewYellowCard(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) {
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

	bytes, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = apollo.PublishMessage(app.Rabbit, pub, &bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func NewRedCard(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) {
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

	bytes, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = apollo.PublishMessage(app.Rabbit, pub, &bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

//
func NewPenalty(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) {
	type request struct {
		team   string
		minute uint8
		player uint8
		reason string
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		badRequest(w, fmt.Sprintf("Unable to read contents: %s", err))
		return
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = apollo.PublishMessage(app.Rabbit, pub, &bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func registerStartMatchV1(app *apollo.DefaultApp, wrapper wrapper) {
	startMatch := &apollo.RabbitPublisherSettings{
		Queue:      "LiveMatch.StartMatch",
		Exchange:   "LiveMatch.StartMatch",
		RoutingKey: "LiveMatch.StartMatch",
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	startMatchMap := make(map[string]leagueHandler)
	startMatchMap["1.0"] = StartMatch
	app.Router.HandleFunc(RootUrl+"/StartMatch", wrapper(startMatchMap, app, startMatch)).Methods("POST")
}

//
func StartMatch(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) {

}

//
func EndMatch(w http.ResponseWriter, r *http.Request, app *apollo.DefaultApp, pub *apollo.RabbitPublisherSettings) {
}
