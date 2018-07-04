// Copyright (c) 2018 NEC Laboratories Europe GmbH.
//
// Authors: Wenting Li <wenting.li@neclab.eu>
//          Sergey Fedorov <sergey.fedorov@neclab.eu>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"fmt"
	"time"
)

//======= Interface for module 'config' =======

// Configer defines the interface to obtain the protocol parameters from
// the configuration
type Configer interface {
	// n: number of nodes in the network
	N() uint32
	// f: number of byzantine nodes the network can tolerate
	F() uint32

	// cp: checkpoint period
	CheckpointPeriod() uint32
	// L: must be larger than CheckpointPeriod
	Logsize() uint32

	// starts when receives a request and stops when request is accepted
	TimeoutRequest() time.Duration
	// starts when sends VIEW-CHANGE and stops when receives a valid NEW-VIEW
	TimeoutViewChange() time.Duration
}

//======= Interface for module 'network' =======

// ReplicaConnector establishes connections to replicas
//
// ReplicaMessageStreamHandler provides a local representation of the
// specified replica for the purpose of message exchange. The local
// representation is returned as MessageStreamHandler interface. It
// guarantees reliable and secure message delivery to the replica.
// Message delivery delay is assumed not to grow indefinitely. This
// assumption has to be satisfied to ensure the liveness of the system.
type ReplicaConnector interface {
	ReplicaMessageStreamHandler(replicaID uint32) (MessageStreamHandler, error)
}

// MessageStreamHandler handles streams of messages
//
// HandleMessageStream initiates asynchronous handling of an incoming
// message stream and returns another stream of messages that might be
// produced in reply. Each value sent/received through a channel is a
// single complete serialized message. Once a message is received from
// any of the channels, it is the receiver's responsibility to finish
// handling of the message.
type MessageStreamHandler interface {
	HandleMessageStream(in <-chan []byte) (out <-chan []byte, err error)
}

//======= Interface for module 'authentication' ========

// AuthenticationRole defines the authentication roles
type AuthenticationRole int

const (
	// ReplicaAuthen specifies authentication of replica messages
	// signed by using a normal replica node key without utilizing
	// the tamper-proof component
	ReplicaAuthen AuthenticationRole = 1 + iota

	// USIGAuthen specifies authentication of replica messages
	// signed by means of a USIG certificate produced in the
	// tamper-proof component of a replica node
	USIGAuthen

	// ClientAuthen specifies authentication of client messages
	ClientAuthen
)

func (r AuthenticationRole) String() string {
	switch r {
	case ReplicaAuthen:
		return "replica"
	case USIGAuthen:
		return "usig"
	case ClientAuthen:
		return "client"
	}
	return fmt.Sprintf("AuthenticationRole(%d)", r)
}

// Authenticator manages the identities of the replicas and clients
// and provides an interface to authenticate the message senders as
// well as to generate authentication tags for the message to send.
// Methods of this interface may be invoked from spawned goroutines.
type Authenticator interface {
	// VerifyMessageAuthenTag verifies authenticity of a message,
	// given an authentication tag, ID of replica/client that
	// signed the message, and the authentication role used to
	// generate the tag.
	VerifyMessageAuthenTag(role AuthenticationRole, id uint32, msg []byte, tag []byte) error

	// GenerateMessageAuthenTag generates an authentication tag
	// for the message using the credentials selected by the
	// specified authentication role
	GenerateMessageAuthenTag(role AuthenticationRole, msg []byte) ([]byte, error)
}

//======= Interface for module 'requestconsumer' ========

// RequestConsumer defines the interface for the replicated state machine. It
// accepts the *accepted* requests from the BFT core and triggers further
// actions such as execution of the message
type RequestConsumer interface {
	// Deliver delivers an accepted message and returns the result once the
	// delivered message has been executed
	Deliver(msg []byte) []byte

	// State returns the digest of the current system state
	StateDigest() []byte
}
