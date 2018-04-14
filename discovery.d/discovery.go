package main

import (
	"log"
	"runtime"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/spf13/viper"

	pb "../imager/imagegrpc"
)

const natsAddress = "nats://192.168.99.100:4222"

var orderServiceUri string

func init() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found")
	}
	orderServiceUri = viper.GetString("discovery.imageservice")
}
func main() {
	// Create server connection
	natsConnection, _ := nats.Connect(natsAddress)
	log.Println("Connected to " + natsAddress)

	natsConnection.Subscribe("Discovery.OrderService", func(m *nats.Msg) {
		orderServiceDiscovery := pb.ServiceDiscovery{OrderServiceUri: orderServiceUri}
		data, err := proto.Marshal(&orderServiceDiscovery)
		if err == nil {
			natsConnection.Publish(m.Reply, data)
		}
	})
	// Keep the connection alive
	runtime.Goexit()
}
