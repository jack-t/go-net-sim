package simulator

import (
	"fmt"
	"strings"
	"strconv"
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
// implemented by client, provided somehow
type Noder interface {
	func Node(type string) Node
}

func (s* Simulator) HandleCommand(cmd string, noder Noder) {
	tokens := strings.SplitN(" ", 2)
	if len(tokens) != 2 {
		fmt.Printf("Syntax error on line: %s", cmd)
		panic()
	}
	switch strings.ToUpper(tokens[0]) {
	case "ADD":
		args := strings.SplitN(tokens[1], 2)
		id := strconv.ParseInt(args[0], 0, 32)
		node := noder.Node(args[1])
		s.AddNode(node, id)
	case "NODE":
		// NODE 1 JOIN 0 -- in chord, have node 1 join through 0; JOIN 0 is passed to node 1
	}
}