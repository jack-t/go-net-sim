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

func NewSimulator() Simulator {
	return Simulator {
		handles: make(map[NodeId]NodeHandle),
		queue: make(<-chan Message),
	}
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
	for _, v := range s.handles {
		go v.handler()
	}
}

// implemented by client, provided somehow
type NodeMaker interface {
	Make(string) Node
}

func (s* Simulator) HandleCommand(cmd string, noder NodeMaker) {
	tokens := strings.SplitN(cmd, " ", 2)
	if len(tokens) != 2 {
		
		panic(fmt.Sprintf("Syntax error on line: %s", cmd))
	}
	switch strings.ToUpper(tokens[0]) {
	case "ADD":
		args := strings.SplitN(tokens[1], " ", 2)
		id, _ := strconv.ParseInt(args[0], 0, 32)
		node := noder.Make(args[1])
		s.AddNode(node, NodeId(id))
	case "NODE":
		args := strings.Split(tokens[1], " ") // NODE {node id} CMD
		if len(args) < 2 {
			panic(fmt.Sprintf("Not enough arguments to command %s", cmd))
		}
		id, _ := strconv.ParseInt(args[0], 0, 32)
		s.handles[NodeId(id)].HandleCommand(args[0], args[1:])
	}
}