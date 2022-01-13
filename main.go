package main

import (
	"context"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"log"
)

func main() {
	url := "wss://KASKUBOW3SRFSLSUY1DUOL7G:6rVa_u8c9oYjb2A1BQ3yvw5a_2cAC5K1SiofBO8H@node-api.klaytnapi.com/v1/ws/open?chain-id=1001"
	ctx := context.Background()
	c, err := client.Dial(url)
	if err != nil {
		log.Fatal("Error when dial: " + err.Error())
	}
	headSink := make(chan *types.Header, 2)
	headerSubscription, err := c.SubscribeNewHead(ctx, headSink)
	if err != nil {
		log.Fatal("Error when subscribe Header: " + err.Error())
	}
	select {
	case ev := <-headSink:
		log.Printf("New block header: %v", ev.Number)
	case err := <-headerSubscription.Err():
		log.Printf("Error while subcription: %v", err)
	}
}
