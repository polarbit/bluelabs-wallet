package main

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type hello struct{ Who string }
type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *hello:
		fmt.Printf("Hello %v\n", msg.Who)
		ch <- true
	}
}

var ch chan bool

func main() {
	ch = make(chan bool, 1)
	system := actor.NewActorSystem()
	props := actor.PropsFromProducer(func() actor.Actor { return &helloActor{} })

	pid := system.Root.Spawn(props)
	system.Root.Send(pid, &hello{Who: "You"})
	system.Root.Send(pid, &hello{Who: "You Too"})
	<-ch
}
