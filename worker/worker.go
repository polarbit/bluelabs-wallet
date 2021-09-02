package worker

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/automanaged"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/polarbit/bluelabs-wallet/worker/messages"
)

type wallet struct {
	cluster.Grain
}

var _ messages.Wallet = (*wallet)(nil)

// Terminate takes care of the finalization.
func (p *wallet) Terminate() {
	// Do finalization if required. e.g. Store the current state to storage and switch its behavior to reject further messages.
	// This method is called when a pre-configured idle interval passes from the last message reception.
	// The actor will be re-initialized when a message comes for the next time.
	// Terminating the idle actor is effective to free unused server resource.
	//
	// A poison pill message is enqueued right after this method execution and the actor eventually stops.
	log.Printf("Terminating ponger: %s", p.ID())
}

// ReceiveDefault is a default method to receive and handle incoming messages.
func (p *wallet) ReceiveDefault(ctx actor.Context) {
	log.Printf("A plain message is sent from sender: %+v", ctx.Sender())

	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *messages.CreateRequest:
		log.Printf("Received CreateRequest message %v", msg)
		pong := &messages.CreateReply{Wallet: &messages.WalletMessage{Id: "1", Balance: 100.}}
		ctx.Respond(pong)
	case *messages.GetRequest:
		log.Printf("Received GetRequest message %v", msg)
		pong := &messages.GetReply{Wallet: &messages.WalletMessage{Id: "1", Balance: 100.}}
		ctx.Respond(pong)
	default:
	}
}

// Ping is called when gRPC-based request is sent against Ponger service.
func (p *wallet) Create(ping *messages.CreateRequest, ctx cluster.GrainContext) (*messages.CreateReply, error) {
	// The sender process is not a sending actor, but a future process
	log.Printf("Received Ping call from sender: %+v", ctx.Sender())

	pong := &messages.CreateReply{
		Wallet: &messages.WalletMessage{Id: "1", Balance: 99},
	}
	return pong, nil
}

func (p *wallet) Get(ping *messages.GetRequest, ctx cluster.GrainContext) (*messages.GetReply, error) {
	// The sender process is not a sending actor, but a future process
	log.Printf("Received Ping call from sender: %+v", ctx.Sender())

	pong := &messages.GetReply{
		Wallet: &messages.WalletMessage{Id: "1", Balance: 1},
	}
	return pong, nil
}

func StartWorker() {
	// Setup actor system
	system := actor.NewActorSystem()

	// Register ponger constructor.
	// This is called when the wrapping PongerActor is initialized.
	// PongerActor proxies messages to ponger's corresponding methods.
	messages.WalletFactory(func() messages.Wallet {
		return &wallet{}
	})

	// Prepare remote env that listens to 8080
	// Messages are sent to this port.
	remoteConfig := remote.Configure("localhost", 8080)

	// Configure cluster provider to work as a cluster member.
	// This member uses port 6331 for cluster provider, and register itself -- localhost:6331" -- as cluster member.
	cp := automanaged.NewWithConfig(1*time.Second, 6331, "localhost:6331")

	// Register an actor constructor for the Ponger kind.
	// With this registration, the message sender and other cluster members know this member is capable of providing Ponger.
	// PongerActor will implicitly be initialized when the first message comes.
	clusterKind := cluster.NewKind(
		"Ponger",
		actor.PropsFromProducer(func() actor.Actor {
			return &messages.WalletActor{
				// The actor stops when 10 seconds passed since the last message reception.
				// When the next
				Timeout: 10 * time.Second,
			}
		}))
	clusterConfig := cluster.Configure("cluster-example", cp, remoteConfig, clusterKind)
	c := cluster.New(system, clusterConfig)

	// Start as a cluster member.
	// Use StartClient() when this process is not a member of cluster members but required to send messages to cluster grains.
	c.Start()

	// Run till signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-finish
}
