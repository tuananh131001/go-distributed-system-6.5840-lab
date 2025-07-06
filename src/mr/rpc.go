package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"

//
// example to show how to declare the arguments
// and reply for an RPC.
//

// type ExampleArgs struct {
// 	X int
// }
//
// type ExampleReply struct {
// 	Y int
// }

type AskForTaskArgs struct {
	// add Worker ID here
}

type TaskType int

const (
	MapTask TaskType = iota
	ReduceTask
	WaitTask
	ExitTask
)

type AskForTaskReply struct {
	TaskType TaskType
	FileName string // for MapTask
	TaskNumber int    // for ReduceTask
	NReduce int		// for ReduceTask
	NMap int		// for MapTask
}


// Add your RPC definitions here.


// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
