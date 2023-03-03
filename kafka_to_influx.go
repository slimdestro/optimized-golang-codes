/**
	@ Read message from kafka topic and write to InfluxDB
**/
package main

import (
	"log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/influxdata/influxdb/client/v2"
)

func main() {
	// Create a new consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Panicln(err)
	}

	// Subscribe to the topic
	c.SubscribeTopics([]string{"myTopic"}, nil)

	// Create a new InfluxDB client
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})

	if err != nil {
		log.Panicln(err)
	}

	// Read messages from kafka topic
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			/**
				// #TODO: parse the message as per your need
			**/
			point, err := client.NewPoint("myMeasurement", nil, map[string]interface{}{
				"value": msg.Value,
			})

			if err != nil {
				log.Panicln(err)
			}

			//new batch
			batch, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  "myDB",
				Precision: "s",
			})

			if err != nil {
				log.Panicln(err)
			}

			batch.AddPoint(point)
			if err := influxClient.Write(batch); err != nil {
				log.Panicln(err)
			}
		} else {
			log.Fatal(err)
		}
	}
}