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

package messages

import (
	"github.com/golang/protobuf/proto"
)

// ClientMessage represents any message generated by a client
//
// ClientID() returns ID of the client created the message
type ClientMessage interface {
	ClientID() uint32
}

// MessageWithUI represents any message with UI attached
//
// ReplicaID returns ID of the replica created the message
//
// Payload returns serialized message data, excluding attached UI
//
// UIBytes returns serialized UI attached to the message payload
//
// AttachUI attaches a serialized UI to the message
type MessageWithUI interface {
	ReplicaID() uint32
	Payload() []byte
	UIBytes() []byte
	AttachUI(ui []byte)
}

// MessageWithSignature represents any message with normal signature
// attached
//
// Payload returns serialized message data, excluding attached
// signature
//
// SignatureBytes returns serialized signature attached to the message
// payload
//
// AttachSignature attaches a serialized signature to the message
type MessageWithSignature interface {
	Payload() []byte
	SignatureBytes() []byte
	AttachSignature(signature []byte)
}

var (
	_ ClientMessage        = (*Request)(nil)
	_ MessageWithUI        = (*Prepare)(nil)
	_ MessageWithUI        = (*Commit)(nil)
	_ MessageWithSignature = (*Request)(nil)
	_ MessageWithSignature = (*Reply)(nil)
)

// ClientID returns ID of the client created the message
func (m *Request) ClientID() uint32 {
	return m.GetMsg().GetClientId()
}

// Payload returns serialized message data, except signature
func (m *Request) Payload() []byte {
	mBytes, err := proto.Marshal(m.GetMsg())
	if err != nil {
		panic(err)
	}
	return mBytes
}

// SignatureBytes returns serialized signature attached to the message
// payload
func (m *Request) SignatureBytes() []byte {
	return m.GetSignature()
}

// AttachSignature attaches a serialized signature to the message
func (m *Request) AttachSignature(signature []byte) {
	m.Signature = signature
}

// Payload returns serialized message data, except signature
func (m *Reply) Payload() []byte {
	mBytes, err := proto.Marshal(m.GetMsg())
	if err != nil {
		panic(err)
	}
	return mBytes
}

// SignatureBytes returns serialized signature attached to the message
// payload
func (m *Reply) SignatureBytes() []byte {
	return m.GetSignature()
}

// AttachSignature attaches a serialized signature to the message
func (m *Reply) AttachSignature(signature []byte) {
	m.Signature = signature
}

// ReplicaID returns ID of the replica created the message
func (m *Prepare) ReplicaID() uint32 {
	return m.Msg.GetReplicaId()
}

// Payload returns serialized message data expect UI
func (m *Prepare) Payload() []byte {
	mBytes, err := proto.Marshal(m.GetMsg())
	if err != nil {
		panic(err)
	}
	return mBytes
}

// UIBytes returns serialized UIBytes attached to the message payload
func (m *Prepare) UIBytes() []byte {
	return m.GetReplicaUi()
}

// AttachUI attaches a serialized UI to the message
func (m *Prepare) AttachUI(ui []byte) {
	m.ReplicaUi = ui
}

// ReplicaID returns ID of the replica created the message
func (m *Commit) ReplicaID() uint32 {
	return m.Msg.GetReplicaId()
}

// Payload returns serialized message data expect UI
func (m *Commit) Payload() []byte {
	mBytes, err := proto.Marshal(m.GetMsg())
	if err != nil {
		panic(err)
	}
	return mBytes
}

// UIBytes returns serialized UIBytes attached to the message payload
func (m *Commit) UIBytes() []byte {
	return m.GetReplicaUi()
}

// AttachUI attaches a serialized UI to the message
func (m *Commit) AttachUI(ui []byte) {
	m.ReplicaUi = ui
}

// Prepare extracts the corresponding Prepare message
func (m *Commit) Prepare() *Prepare {
	msg := m.GetMsg()
	return &Prepare{
		Msg: &Prepare_M{
			View:      msg.GetView(),
			ReplicaId: msg.GetPrimaryId(),
			Request:   msg.GetRequest(),
		},
		ReplicaUi: msg.GetPrimaryUi(),
	}
}

// Request extract the corresponding Request message
func (m *Commit) Request() *Request {
	return m.GetMsg().GetRequest()
}

// WrapMessage wraps a concrete message to a generic wrapper Message.
func WrapMessage(m interface{}) *Message {
	switch m := m.(type) {
	case *Request:
		return &Message{
			Type: &Message_Request{
				Request: m,
			},
		}
	case *Reply:
		return &Message{
			Type: &Message_Reply{
				Reply: m,
			},
		}
	case *Prepare:
		return &Message{
			Type: &Message_Prepare{
				Prepare: m,
			}}
	case *Commit:
		return &Message{
			Type: &Message_Commit{
				Commit: m,
			}}
	default:
		panic("Unknown message type")
	}
}

// UnwrapMessage unwraps a generic Message to a concrete message type.
func UnwrapMessage(m *Message) interface{} {
	switch t := m.Type.(type) {
	case *Message_Request:
		return t.Request
	case *Message_Reply:
		return t.Reply
	case *Message_Prepare:
		return t.Prepare
	case *Message_Commit:
		return t.Commit
	default:
		panic("Unknown message type")
	}
}
