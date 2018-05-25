package simulator

type NodeId uint32

// bridge between Simulator and Node: both sides have access to the handle
type NodeHandle interface {
	// put a message on the simulator's queue
	Send(message string, recipient NodeId)
	// take a message off the queue's hands, get ready to deliver to the node (call HandleMessage)
	Recv(message Message)
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

func NewNodeHandle(s *Simulator, n *Node, i NodeId) NodeHandle {
	return defaultNodeHandle{
		me: i,
		sim: *s,
		pendingMessages: make(chan Message, 128), // buffered channel; things'll get fucky if you've got a bunch of shit in the queue
		node: n,
	}
}

type Node interface {
	SetHandle(handle NodeHandle)
	HandleMessage(message Message)
}
