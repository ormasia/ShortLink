package main

import "shortLink/logservice/consumer"

func main() {
	consumer.StartLogConsumer([]string{"localhost:9092"}, "log-consumer-group", []string{"shortlink-log"})
}
