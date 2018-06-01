package simulator

import (
	"fmt"
)

type NodeId uint32

// bridge between Simulator and Node: both sides have access to the handle
type NodeHandle interface {
	// put a message on the simulator's queue
	Send(message string, recipient NodeId)
	// take a message off the queue's hands, get ready to deliver to the node (call HandleMessage)
	Recv(message Message)
	HandleCommand(cmd string, args []string)
	Log(format string, a ...interface{}) (n int, err error) // basically, passes everything on to fmt
	handler()
}

type defaultNodeHandle struct {
	me NodeId
	sim Simulator
	pendingMessages chan Message
	node Node
}

func (d defaultNodeHandle) Send(message string, recipient NodeId) {
	msg := Message{
		Payload: message,
		Recipients: []NodeId{ recipient },
		Sender: d.me,
	}
	d.sim.DeliverMessage(msg)
}

func (d defaultNodeHandle) Recv(message Message) {
	d.pendingMessages <- message
}

func (d defaultNodeHandle) handler() {
	for ;; {
		msg := <- d.pendingMessages
		d.node.HandleMessage(msg)
	}
}

func (d defaultNodeHandle) HandleCommand(cmd string, args []string) {
	d.node.HandleCommand(cmd, args)
}
func (d defaultNodeHandle) Log(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf("%d says: %s", d.me, fmt.Sprintf(format, a...))
}

func NewNodeHandle(s *Simulator, n Node, i NodeId) NodeHandle {
	ret := defaultNodeHandle{
		me: i,
		sim: *s,
		pendingMessages: make(chan Message, 128), // buffered channel; things'll get fucky if you've got a bunch of shit in the queue
		node: n,
	}
	n.SetHandle(ret)

	return ret
}

type Node interface {
	SetHandle(handle NodeHandle)
	HandleMessage(message Message)
	HandleCommand(cmd string, args []string)
}
