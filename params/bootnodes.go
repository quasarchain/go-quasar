// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import "github.com/ethereum/go-ethereum/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var MainnetBootnodes = []string{
	// Ethereum Foundation Go Bootnodes
	"enode://48da1eaa78b2e5af26c3b6cbd3ae920e3c8e484de02d6824fa28f06346be45d86b58179674aa4e0420a7a58d4e0ac27fc85599999a5a780b2a8865078482352f@102.129.155.241:37742",
	"enode://3550c06d3a0214b93e08ad0aa9b1d58b3e29a063c0033985da2abf86fb44c42c0316f05ce66af273ef74ffad276ec38d824894346dd05324eea40810e0ca083b@45.32.67.14:37742",
	"enode://699ff785887742465615bb27c2a2767ff33c53093443e37eda19caadb3db47ec8ce1d10f220e09d10366c06bef8f0c46bdb27050b3604ba4b21c5a4168b5180c@209.222.30.54:37742",
	}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
var TestnetBootnodes = []string{
	"",
	"",
}

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	return ""
}
