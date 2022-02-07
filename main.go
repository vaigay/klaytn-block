package main

import (
	"context"
	"fmt"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/client"
	"github.com/klaytn/klaytn/common"
	"klaytn-block/kabi"
	"log"
	"time"
)

func main() {
	url := "wss://api.baobab.klaytn.net:8652"
	contract := common.HexToAddress("0xD0CAA66319290197746b92e42BaF3338c67Ff745")
	ctx := context.Background()
	c, err := client.DialContext(ctx, url)
	if err != nil {
		log.Fatal("Can not dial!")
	}
	f, err := kabi.NewKabiFilterer(contract, c)
	if err != nil {
		log.Fatal("could not create wormhole contract filter: %w", err)
	}
	timeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	messageC := make(chan *kabi.KabiLogMessagePublished)
	messageSub, err := f.WatchLogMessagePublished(&bind.WatchOpts{Context: timeout}, messageC, nil)
	if err != nil {
		log.Fatalf("could not subscribe new header: %w", err)
	}
	errC := make(chan error)
	t := 0
	i := 0
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-messageSub.Err():
				fmt.Printf("Message error: %v\n", err)
				errC <- err
				fmt.Printf("Times: %v\n", t)
				return
			case ev := <-messageC:
				t++
				i++
				log.Printf("Transaction %v: %v\n", i, ev.Raw.TxHash.String())
			}
		}
	}()

	headSink := make(chan *types.Header, 2)
	headerSubscription, err := c.SubscribeNewHead(ctx, headSink)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-headerSubscription.Err():
				fmt.Printf("Message error: %v\n", err)
				errC <- err
				fmt.Printf("Times: %v\n", t)
			case ev := <-headSink:
				t++
				log.Printf("New Header: %v\n", ev.Number)
			}
		}
	}()
	select {
	case <-errC:
		log.Fatal("End!!!")
	case <-ctx.Done():
		log.Fatal("End!!!")
	}
}
