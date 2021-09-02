package worker

import (
	"log"
	"strconv"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/automanaged"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/polarbit/bluelabs-wallet/worker/messages"
)

var cnt int = 0

type walletClientActor struct {
	cluster *cluster.Cluster
}

func (p *walletClientActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		cnt += 1
		req := &messages.CreateRequest{
			Id: strconv.Itoa(cnt),
		}

		client := messages.GetWalletGrainClient(p.cluster, "wallet-1")
		option := cluster.NewGrainCallOptions(p.cluster).WithRetry(3)
		res, err := client.Create(req, option)
		if err != nil {
			log.Print(err.Error())
			return
		}
		log.Printf("Received %v", res)

	case *messages.CreateReply:
		// Never comes here.
		// When the pong grain responds to the sender's gRPC call,
		// the sender is not a ping actor but a future process.
		log.Print("Received pong message")

	}
}

func StartClient(finish <-chan bool) {
	// Setup actor system
	system := actor.NewActorSystem()

	// Prepare remote env that listens to 8081
	remoteConfig := remote.Configure("localhost", 8081)

	// Configure cluster on top of the above remote env
	// This member uses port 6330 for cluster provider, and add ponger member -- localhost:6331 -- as member.
	// With automanaged implementation, one must list up all known members at first place to ping each other.
	// Note that this member itself is not registered as a member member because this only works as a client.
	cp := automanaged.NewWithConfig(1*time.Second, 6330, "localhost:6331")
	clusterConfig := cluster.Configure("cluster-example", cp, remoteConfig)
	c := cluster.New(system, clusterConfig)
	// Start as a client, not as a cluster member.
	c.StartClient()

	// Start ping actor that periodically send "ping" payload to "Ponger" cluster grain
	clientProps := actor.PropsFromProducer(func() actor.Actor {
		return &walletClientActor{
			cluster: c,
		}
	})
	clientPid := system.Root.Spawn(clientProps)

	// Periodically send ping payload till signal comes
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			system.Root.Send(clientPid, struct{}{})

		case <-finish:
			log.Print("Finish")
			return

		}
	}
}
