package mr

import "fmt"
import "log"
import "net/rpc"
import "hash/fnv"
import "time"

// Map functions return a slice of KeyValue.
type KeyValue struct {
	Key   string
	Value string
}

// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

// main/mrworker.go calls this function.
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	for {
		args := AskForTaskArgs{}
		reply := AskForTaskReply{}
		ok := call("Coordinator.AskForTask", &args, &reply)
		if !ok {
			log.Printf("Worker: call failed! Coordinator likely exited. Exiting worker.\n")
			break // Coordinator exited, so worker should exit
		}

		switch reply.TaskType {
		case MapTask:
			log.Printf("Worker: Received Map task %d for file %s. NReduce: %d, NMap: %d\n",
				reply.TaskNumber, reply.FileName, reply.NReduce, reply.NMap)
			// We'll add the actual Map task execution logic here later.

		case ReduceTask:
			log.Printf("Worker: Received Reduce task %d.\n", reply.TaskNumber)
			// We'll add the actual Reduce task execution logic here later.

		case WaitTask:
			log.Println("Worker: No tasks available, waiting...")
			time.Sleep(500 * time.Millisecond) // Wait a bit before asking again

		case ExitTask:
			log.Println("Worker: Coordinator signaled exit. Exiting worker.")
			return // Exit the worker function

		}

	}
}


// send an RPC request to the coordinator, wait for the response.
// usually returns true.
// returns false if something goes wrong.
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := coordinatorSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		return false
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
