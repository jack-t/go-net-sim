package simulator

import (
	"fmt"
)

type Message struct {
	Payload string
	Recipients []NodeId
	Sender NodeId
}

type Simulator struct {
	handles map[NodeId]NodeHandle
	queue <-chan Message
}

func (s *Simulator) DeliverMessage(m Message) {
	for _, d := range m.Recipients {
		if val, ok := s.handles[d]; ok {
			val.Recv(m)
		} else {
			fmt.Printf("No such node id = %d\n", d)
		}
	}
}

func (s *Simulator) getNextId() NodeId {
	return NodeId(len(s.handles))
}

func (s *Simulator) AddNode(n Node, i NodeId) {
	s.handles[i] = NewNodeHandle(s, n, i)
}

func (s *Simulator) Launch() {
	for k, v := range s.handles {
		go v.handler()
	}
}