package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/axm/apollo-utils"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	StartMatch = ``
)

type startMatch struct {
}

type Config map[string]*json.RawMessage

func readFile(relativePath string) ([]byte, error) {
	cwd, err := os.Getwd()
	fmt.Println(fmt.Sprintf("cwd: %s", cwd))
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get cwd.")
	}
	path := path.Join(cwd, relativePath)
	fmt.Println(fmt.Sprintf("Path: %s", path))
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read file.")
	}

	return contents, nil
}

func readDatabaseConnection(config *Config) (*apollo.DatabaseConnection, error) {
	var dc apollo.DatabaseConnection
	bytes, err := (*config)["Database"].MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read database settings section")
	}
	err = json.Unmarshal(bytes, &dc)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to parse database settings")
	}

	return &dc, nil
}

func readRabbitSettings(config *Config) (*apollo.RabbitConnection, *apollo.RabbitConsumerSettings, error) {
	buffer, err := (*config)["Rabbit"].MarshalJSON()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to read rabbit settings section.")
	}

	rabbitMap := make(Config)
	err = json.Unmarshal(buffer, &rabbitMap)

	buffer, err = rabbitMap["Connection"].MarshalJSON()
	var rc apollo.RabbitConnection
	err = json.Unmarshal(buffer, &rc)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to parse rabbit connection settings.")
	}

	var cs apollo.RabbitConsumerSettings
	buffer, err = rabbitMap["Consumer"].MarshalJSON()
	err = json.Unmarshal(buffer, &cs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to parse rabbit consumer settings.")
	}

	return &rc, &cs, nil
}

func main() {
	// configContents, err := readFile("consumers/startMatch/src/config.json")
	configContents, err := readFile("config.json")
	if err != nil {
		log.Fatal("Error reading config file. Aborting.")
		return
	}

	jsonConfig := make(Config)
	err = json.Unmarshal(configContents, &jsonConfig)
	if err != nil {
		log.Fatal("Error parsing config. Aborting.")
		return
	}

	dc, err := readDatabaseConnection(&jsonConfig)
	if err != nil || dc == nil {
		log.Fatal("Error reading database connection config")
		return
	}

	rc, cs, err := readRabbitSettings(&jsonConfig)
	if err != nil {
		log.Fatal("Error reading rabbit connection config")
	}

	rabbitConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%d", rc.User, rc.Password, rc.Host, rc.Port)
	conn, err := amqp.Dial(rabbitConnectionString)
	if err != nil {
		fmt.Println(fmt.Sprintf("Rabbit connection string: %s", rabbitConnectionString))
		fmt.Println(fmt.Sprintf("Error connecting to rabbit: %s", err))
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

	closeFlag := make(chan bool)
	consume := func(closeFlag chan bool) {
		for m := range msgs {
			fmt.Println(string(m.Body))
		}
	}

	go consume(closeFlag)
	_ = <-closeFlag
	fmt.Println("Received close flag")
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
