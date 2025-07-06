package mr

// class responsible for coordinating the MapReduce job. In design, it will manage the tasks, track their status, and communicate with workers.

import "log"
import "net"
import "os"
import "net/rpc"
import "net/http"
import "sync"

type Coordinator struct {
	mu         sync.Mutex // Mutex to protect shared state
	InputFiles []string   // List of input files to be processed
	MapTasks   []Task     // Slice of map tasks
	NReduce    int        // Number of reduce tasks
	NMap       int        // Number of map tasks (equal to the number of input files)
}

type Task struct {
	FileName string     // Name of the file for map tasks
	TaskType   TaskType   // Type of the task (MapTask, ReduceTask, etc.)
	TaskNumber int        // Task number (index for map tasks, reduce task number)
	Status     TaskStatus // Status of the task (Pending, InProgress, Completed, etc.)
	WorkerID   string     // ID of the worker assigned to this task (optional)
	StartTime  int64      // Start time of the task (for tracking progress)
}

type TaskStatus int

const (
	Pending TaskStatus = iota
	InProgress
	Completed
)

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.

func (c *Coordinator) AskForTask(args *AskForTaskArgs, reply *AskForTaskReply) error {
	c.mu.Lock() // Lock the mutex to protect shared state
	defer c.mu.Unlock() // Ensure the mutex is unlocked when the function returns

	// Get tasks and process them
	for i, task := range c.MapTasks {
		if task.Status == Pending {
			reply.TaskType = MapTask
			reply.FileName = task.FileName
			reply.TaskNumber = task.TaskNumber
			reply.NReduce = c.NReduce
			reply.NMap = c.NMap

			c.MapTasks[i].Status = InProgress // Update task status to InProgress
			log.Printf("Assigned task %d of type %v to worker", task.TaskNumber, reply.TaskType)
			return nil // Return the assigned task
		}
	}
	reply.TaskType = WaitTask // No pending tasks, ask worker to wait
	log.Println("No pending tasks, worker should wait")
	return nil // Return without an error, worker will wait

}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	c.mu.Lock() // Lock the mutex to protect shared state
	defer c.mu.Unlock() // Ensure the mutex is unlocked when the function returns

	return true
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Initialize the Coordinator with the list of input files and number of reduce tasks
	c.InputFiles = files // List of input files
	c.NReduce = nReduce  // number of reduce tasks
	c.NMap = len(files)  // Number of len files

	// Allocate a slice of tasks for reduce tasks
	c.MapTasks = make([]Task, c.NMap) // wat da hMapTask, len(files))eo is make? Allocate a slice of tasks for map tasks
	log.Printf("Coordinator initialized with %d map tasks and %d reduce tasks\n", c.NMap, c.NReduce)

	// Loop through the input files and create a map task for each file
	for i, files := range files {
		c.MapTasks[i] = Task{
			FileName:   files,       // Set the filename for each map task
			TaskType:   MapTask,     // Set the task type to MapTask
			TaskNumber: i,           // Set the task number for each map task
			Status:     Pending, // Set the initial status of the task to pending
		}
	}

	c.server()
	return &c
}
