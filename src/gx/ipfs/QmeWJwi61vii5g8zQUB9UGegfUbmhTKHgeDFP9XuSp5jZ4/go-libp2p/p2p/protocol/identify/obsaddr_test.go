package identify

import (
	"testing"
	"time"

	ma "gx/ipfs/QmSWLfmj5frN9xVLMMN846dMDriy5wN5jeghUm7aTW3DAG/go-multiaddr"
)

// TestObsAddrSet
func TestObsAddrSet(t *testing.T) {
	m := func(s string) ma.Multiaddr {
		m, err := ma.NewMultiaddr(s)
		if err != nil {
			t.Error(err)
		}
		return m
	}

	addrsMarch := func(a, b []ma.Multiaddr) bool {
		for _, aa := range a {
			found := false
			for _, bb := range b {
				if aa.Equal(bb) {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		return true
	}

	a1 := m("/ip4/1.2.3.4/tcp/1231")
	a2 := m("/ip4/1.2.3.4/tcp/1232")
	a3 := m("/ip4/1.2.3.4/tcp/1233")
	a4 := m("/ip4/1.2.3.4/tcp/1234")
	a5 := m("/ip4/1.2.3.4/tcp/1235")
	a6 := m("/ip4/1.2.3.6/tcp/1236")
	a7 := m("/ip4/1.2.3.7/tcp/1237")

	oas := ObservedAddrSet{}

	if !addrsMarch(oas.Addrs(), nil) {
		t.Error("addrs should be empty")
	}

	oas.Add(a1, a4)
	oas.Add(a2, a4)
	oas.Add(a3, a4)

	// these are all different so we should not yet get them.
	if !addrsMarch(oas.Addrs(), nil) {
		t.Error("addrs should _still_ be empty (once)")
	}

	// same observer, so should not yet get them.
	oas.Add(a1, a4)
	oas.Add(a2, a4)
	oas.Add(a3, a4)
	if !addrsMarch(oas.Addrs(), nil) {
		t.Error("addrs should _still_ be empty (same obs)")
	}

	// different observer, but same observer group.
	oas.Add(a1, a5)
	oas.Add(a2, a5)
	oas.Add(a3, a5)
	if !addrsMarch(oas.Addrs(), nil) {
		t.Error("addrs should _still_ be empty (same obs group)")
	}

	oas.Add(a1, a6)
	if !addrsMarch(oas.Addrs(), []ma.Multiaddr{a1}) {
		t.Error("addrs should only have a1")
	}

	oas.Add(a2, a5)
	oas.Add(a1, a5)
	oas.Add(a1, a5)
	oas.Add(a2, a6)
	oas.Add(a1, a6)
	oas.Add(a1, a6)
	oas.Add(a2, a7)
	oas.Add(a1, a7)
	oas.Add(a1, a7)
	if !addrsMarch(oas.Addrs(), []ma.Multiaddr{a1, a2}) {
		t.Error("addrs should only have a1, a2")
	}

	// change the timeout constant so we can time it out.
	oas.SetTTL(time.Millisecond * 200)
	<-time.After(time.Millisecond * 210)
	if !addrsMarch(oas.Addrs(), []ma.Multiaddr{nil}) {
		t.Error("addrs should have timed out")
	}
}
