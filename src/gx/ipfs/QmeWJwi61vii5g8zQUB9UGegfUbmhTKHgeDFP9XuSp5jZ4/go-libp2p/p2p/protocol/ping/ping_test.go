package ping

import (
	"context"
	"testing"
	"time"

	netutil "gx/ipfs/QmNqvnxGtJBaKQnenD6uboNGdjSjHGmZGRxMHEevKJe5Pk/go-libp2p-netutil"
	peer "gx/ipfs/QmWUswjn261LSyVxWAEpMVtPdy8zmKBJJfBpG3Qdpa8ZsE/go-libp2p-peer"
	pstore "gx/ipfs/Qme1g4e3m2SmdiSGGU3vSWmUStwUjc5oECnEriaK9Xa1HU/go-libp2p-peerstore"
	bhost "gx/ipfs/QmeWJwi61vii5g8zQUB9UGegfUbmhTKHgeDFP9XuSp5jZ4/go-libp2p/p2p/host/basic"
)

func TestPing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h1 := bhost.New(netutil.GenSwarmNetwork(t, ctx))
	h2 := bhost.New(netutil.GenSwarmNetwork(t, ctx))

	err := h1.Connect(ctx, pstore.PeerInfo{
		ID:    h2.ID(),
		Addrs: h2.Addrs(),
	})

	if err != nil {
		t.Fatal(err)
	}

	ps1 := NewPingService(h1)
	ps2 := NewPingService(h2)

	testPing(t, ps1, h2.ID())
	testPing(t, ps2, h1.ID())
}

func testPing(t *testing.T, ps *PingService, p peer.ID) {
	pctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ts, err := ps.Ping(pctx, p)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		select {
		case took := <-ts:
			t.Log("ping took: ", took)
		case <-time.After(time.Second * 4):
			t.Fatal("failed to receive ping")
		}
	}

}
