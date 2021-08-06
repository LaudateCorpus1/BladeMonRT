package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"os"
	"syscall"
	"os/signal"
	"github.com/microsoft/BladeMonRT/logging"
	"github.com/microsoft/BladeMonRT/configs"
)

type Main struct {
	workflowFactory *WorkflowFactory
	logger          *log.Logger
}

func main() {
	// Set GOMAXPROCS such that all operations execute on a single thread.
	runtime.GOMAXPROCS(1)

	var mainObj Main = NewMain()
	mainObj.logger.Println("Initialized main.")
	mainObj.logger.Println("Number of threads:", runtime.GOMAXPROCS(-1))

	// Setup main such that main does not exit unless there is a keyboard interrupt.
	quitChannel := make(chan os.Signal, 1)
    signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM) 
	<-quitChannel
}

func NewMain() Main {
	workflowsJson, err := ioutil.ReadFile(configs.WORKFLOW_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var workflowFactory WorkflowFactory = newWorkflowFactory(workflowsJson, NodeFactory{})

	schedulesJson, err := ioutil.ReadFile(configs.SCHEDULE_FILE)
	if err != nil {
		log.Fatal(err)
	}
	var workflowScheduler *WorkflowScheduler = newWorkflowScheduler(schedulesJson, workflowFactory)
	fmt.Println(workflowScheduler) // TODO: Remove print statement.

	var logger *log.Logger = logging.LoggerFactory{}.ConstructLogger("Main")
	return Main{workflowFactory: &workflowFactory, logger: logger}
}