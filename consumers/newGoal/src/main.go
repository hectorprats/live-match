package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/axm/apollo-utils"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type newGoal struct {
	Host   string `json:Host`
	Guest  string `json:Guest`
	Team   string `json:Team`
	Minute uint8  `json:Minute`
	Player uint8  `json:Player`
}

const (
	InsertNewGoal = `
INSERT INTO MatchEvents(matchcode, eventtype, payload)
VALUES($1, $2, $3)
`
)

func main() {
	dc := &apollo.DatabaseConnection{
		Server:     "livematch_postgres",
		Port:       5432,
		UserId:     "postgres",
		Password:   "postgres",
		Database:   "livematch",
		DriverName: "postgres",
	}

	cs := &apollo.RabbitConsumerSettings{
		Queue:     "LiveMatch.NewGoal",
		Consumer:  "",
		AutoAck:   true,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}
	rc := &apollo.RabbitConnection{
		User:     "guest",
		Password: "guest",
		Host:     "livematch_rabbit",
		Port:     5672,
	}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d", rc.User, rc.Password, rc.Host, rc.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		cs.Queue,
		cs.AutoAck,
		cs.Exclusive,
		cs.NoLocal,
		cs.NoWait,
		nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := ch.Consume(
		cs.Queue,
		cs.Consumer,
		cs.AutoAck,
		cs.Exclusive,
		cs.NoLocal,
		cs.NoWait,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	consume := func() {
		for m := range msgs {
			fmt.Println(string(m.Body))
			fmt.Println(fmt.Sprintf("Body length: %d", len(m.Body)))
			var newGoal newGoal
			err := json.Unmarshal(m.Body, &newGoal)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = storeEvent(&newGoal, dc)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
	go consume()
	whileTrue := make(chan bool)
	<-whileTrue
}

func storeEvent(goal *newGoal, dc *apollo.DatabaseConnection) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dc.Server, dc.Port, dc.UserId, dc.Password, dc.Database)
	fmt.Println(fmt.Sprintf("Connection string: %s", psqlInfo))
	db, err := sql.Open(dc.DriverName, psqlInfo)
	if err != nil {
		return errors.Wrap(err, "Unable to open database.")
	}
	defer db.Close()

	payload, err := json.Marshal(goal)
	_, err = db.Exec(InsertNewGoal, fmt.Sprintf("%s%s", goal.Host, goal.Guest), 1, string(payload))
	if err != nil {
		return errors.Wrap(err, "Failed to exec sql.")
	}

	return nil
}
