package main

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

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
		ch <- true
		fmt.Println("Roger!")
	}
}

var ch chan bool

func main() {
	ch = make(chan bool, 1)

	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &helloActor{} })
	pid := system.Root.Spawn(props)

	// Send operations are non-blocking
	system.Root.Send(pid, &hello{Who: "You"})
	system.Root.Send(pid, &hello{Who: "You Too"})
	<-ch
	system.Root.StopFuture(pid).Wait()
}
