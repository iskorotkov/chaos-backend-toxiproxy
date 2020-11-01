package main

import (
	"flag"
	"fmt"
	"github.com/Shopify/toxiproxy/client"
	"time"
)

var tpHost = flag.String("tp_host", "localhost", "Toxiproxy host")
var tpPort = flag.Int("tp_port", 8474, "Toxiproxy port")

var host = flag.String("host", "localhost", "Host that will be proxied")
var listenPort = flag.Int("listen", 18811, "Listen port for Toxiproxy")
var upstreamPort = flag.Int("upstream", 8811, "Upstream port for Toxiproxy")

func main() {
	flag.Parse()

	client := toxiproxy.NewClient(fmt.Sprintf("%v:%v", *tpHost, *tpPort))
	proxy, err := client.CreateProxy("server", fmt.Sprintf(":%v", *listenPort), fmt.Sprintf("%v:%v", *host, *upstreamPort))
	if err != nil {
		fmt.Println("Failed to create proxy")
		panic(err)
	}
	defer func() {
		_ = proxy.Delete()
	}()

	_, err = proxy.AddToxic("latency_down", "latency", "", 1.0, toxiproxy.Attributes{
		"latency": 1000,
	})
	if err != nil {
		fmt.Println("Error occurred while adding toxic")
		panic(err)
	}

	ch := time.After(time.Hour)
	<-ch
}
