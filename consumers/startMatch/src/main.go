package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/axm/apollo-utils"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	StartMatch = ``
)

type startMatch struct {
}

func main() {
	dc := &apollo.DatabaseConnection{
		Server:     "",
		Port:       0,
		UserId:     "",
		Password:   "",
		Database:   "",
		DriverName: "",
	}

	cs := &apollo.RabbitConsumerSettings{
		Queue:     "",
		Consumer:  "",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}

	rc := &apollo.RabbitConnection{
		User:     "",
		Password: "",
		Host:     "",
		Port:     0,
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

	whileTrue := make(chan bool)
	<-whileTrue
}

func storeEvent(goal *startMatch, dc *apollo.DatabaseConnection) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dc.Server, dc.Port, dc.UserId, dc.Password, dc.Database)
	fmt.Println(fmt.Sprintf("Connection string: %s", psqlInfo))
	db, err := sql.Open(dc.DriverName, psqlInfo)
	if err != nil {
		return errors.Wrap(err, "Unable to open database.")
	}
	defer db.Close()

	//payload, err := json.Marshal(goal)
	//_, err = db.Exec(StartMatch, fmt.Sprintf("%s%s", goal.Host, goal.Guest), 1, string(payload))
	if err != nil {
		return errors.Wrap(err, "Failed to exec sql.")
	}

	return nil
}
