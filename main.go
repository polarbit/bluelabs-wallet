package main

import (
	"fmt"
	"os"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type WalletCreator interface {
}

type Wallet struct {
	wallet   walletState
	lastTran walletTran
}

type walletState struct {
	wid     string
	balance float64
	ver     int32
	created time.Time
}

type walletTran struct {
	wid         string
	ver         float32
	old_balance float32
	new_balance float32
	amount      float32
	created     time.Time
	by          string
	ttype       string
	// wid + wwer should be unique in db.
	// op: transfer, win, play, deposit,
}

type hello struct{ Who string }
type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
		fmt.Println("Roger!")
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		time.Sleep(3 * time.Second)
		done <- true
	}()

	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &helloActor{} })
	pid := system.Root.Spawn(props)

	// Send operations are non-blocking
	system.Root.Send(pid, &hello{Who: "You"})
	system.Root.Send(pid, &hello{Who: "You Too"})
	system.Root.StopFuture(pid).Wait()

	<-done

	fmt.Println("Bitti")
}
