/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package discovery

import (
	"testing"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/common/discovery/staticdiscovery"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	fabImpl "github.com/hyperledger/fabric-sdk-go/pkg/fab"
	mocks "github.com/hyperledger/fabric-sdk-go/pkg/fab/mocks"
	"github.com/hyperledger/fabric-sdk-go/pkg/msp/test/mockmsp"
)

type mockFilter struct {
	called bool
}

// Accept returns true if this peer is to be included in the target list
func (df *mockFilter) Accept(peer fab.Peer) bool {
	df.called = true
	return true
}

func TestDiscoveryFilter(t *testing.T) {

	configBackend, err := config.FromFile("../../../../test/fixtures/config/config_test.yaml")()
	if err != nil {
		t.Fatalf(err.Error())
	}

	config1, err := fabImpl.ConfigFromBackend(configBackend...)
	if err != nil {
		t.Fatalf(err.Error())
	}

	discoveryProvider, err := staticdiscovery.New(config1)
	if err != nil {
		t.Fatalf("Failed to  setup discovery provider: %s", err)
	}
	discoveryProvider.Initialize(mocks.NewMockContext(mockmsp.NewMockSigningIdentity("user1", "Org1MSP")))

	discoveryService, err := discoveryProvider.CreateDiscoveryService("mychannel")
	if err != nil {
		t.Fatalf("Failed to setup discovery service: %s", err)
	}

	discoveryFilter := &mockFilter{called: false}

	discoveryService = NewDiscoveryFilterService(discoveryService, discoveryFilter)

	peers, err := discoveryService.GetPeers()
	if err != nil {
		t.Fatalf("Failed to get peers from discovery service: %s", err)
	}

	// One peer is configured for "mychannel"
	expectedNumOfPeers := 1
	if len(peers) != expectedNumOfPeers {
		t.Fatalf("Expecting %d, got %d peers", expectedNumOfPeers, len(peers))
	}

	if !discoveryFilter.called {
		t.Fatalf("Expecting true, got false")
	}

}