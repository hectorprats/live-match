package main

import (
	"encoding/json"
	"fmt"
	"github.com/axm/apollo-utils"
	"github.com/go-redis/redis"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func Consume() {
	type newGoalCommand struct {
	}
	config := &apollo.KafkaConsumerConfig{}
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"group.id":          config.GroupId,
		"auto.offset.reset": config.AutoOffsetReset,
	})

	if err != nil {

	}

	c.SubscribeTopics([]string{"", ""}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {

		}
		command := &newGoalCommand{}
		err = json.Unmarshal(msg.Value, &command)
		if err != nil {
			continue
		}
	}
}

type NewYellowCardCommand struct {
	match  uint64
	team   string
	player uint8
	minute uint8
}

type RedisConfig struct {
	address  string
	password string
	db       uint8
}

func NewYellowCardConsumer(yellowCard *NewYellowCardCommand, app *apollo.DefaultApp) error {
	// this is a player event but also a match event
	client := redis.NewClient(&redis.Options{
		Addr:     app.Redis.Address,
		Password: app.Redis.Password,
		DB:       app.Redis.DB,
	})
}
