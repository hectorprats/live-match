package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/axm/apollo-utils"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const (
	NewYellowCardSql = `
INSERT INTO MatchEvents(matchcode, eventtype, payload)
VALUES($1, $2, $3)
`
)

type newYellowCard struct {
	Host  string
	Guest string
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
	// configContents, err := readFile("consumers/newYellowCard/src/config.json")
	configContents, err := readFile("config.json")
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading config file: %s", err))
		return
	}

	jsonConfig := make(Config)
	err = json.Unmarshal(configContents, &jsonConfig)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error parsing config: %s", err))
		return
	}

	dc, err := readDatabaseConnection(&jsonConfig)
	if err != nil || dc == nil {
		log.Fatal(fmt.Sprintf("Error reading database connection config: %s", err))
		return
	}

	rc, cs, err := readRabbitSettings(&jsonConfig)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error reading rabbit connection config: %s", err))
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
			var newYellowCard newYellowCard
			err := json.Unmarshal(m.Body, &newYellowCard)
			fmt.Println("Deserialized stuff")
			if err != nil {
				err = errors.Wrap(err, "Unable to deserialize message contents.")
				fmt.Println(err)
				continue
			}
			fmt.Println("Storing event...")
			err = storeEvent(&newYellowCard, dc)
			fmt.Println("Past store event")
			if err != nil {
				err = errors.Wrap(err, "Unable to store NewYellowCard event")
				fmt.Println(err)
				continue
			}
			fmt.Println("Event stored")
		}

		closeFlag <- true
	}

	go consume(closeFlag)
	_ = <-closeFlag
	fmt.Println("Received close flag")
}

func storeEvent(newYellowCard *newYellowCard, dc *apollo.DatabaseConnection) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dc.Server, dc.Port, dc.UserId, dc.Password, dc.Database)
	fmt.Println(fmt.Sprintf("Connection string: %s", psqlInfo))
	db, err := sql.Open(dc.DriverName, psqlInfo)
	fmt.Println("Connection open")
	if err != nil {
		return errors.Wrap(err, "Unable to open database.")
	}
	defer db.Close()

	payload, err := json.Marshal(newYellowCard)
	_, err = db.Exec(NewYellowCardSql, fmt.Sprintf("%s%s", newYellowCard.Host, newYellowCard.Guest), 32189, string(payload))
	fmt.Println("SQL executed")
	if err != nil {
		return errors.Wrap(err, "Failed to exec sql.")
	}

	return nil
}
