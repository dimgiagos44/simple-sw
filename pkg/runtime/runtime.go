package runtime

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/serverlessworkflow/sdk-go/model"
)

func (r *Runtime) Start() {
	r.exec()
}

func (r *Runtime) exec() {
	initState := r.Workflow.States[0]
	fmt.Println("The input file is: ", r.InputFile)
	if r.InputFile != "" {
		jsonFile, _ := os.Open(r.InputFile)
		byteValue, _ := ioutil.ReadAll(jsonFile)
		r.lastOutput = byteValue
	}

	r.begin(initState)

}

func (r *Runtime) begin(st model.State) error {
	switch st.(type) {
	case *model.EventState:
		//fmt.Println("event")
		handleEventState(st.(*model.EventState), r)
	case *model.OperationState:
		//fmt.Println("operation")
		handleOperationState(st.(*model.OperationState), r)
	case *model.EventBasedSwitchState:
		//fmt.Println("event based switch")
	case *model.DataBasedSwitchState:
		//fmt.Println("data based switch")
		HandleDataBasedSwitch(st.(*model.DataBasedSwitchState), r.lastOutput, r)
	case *model.ForEachState:
		//fmt.Println("foreach")
	case *model.ParallelState:
		//fmt.Println("parallel")
	case *model.InjectState:
		//fmt.Println("inject")
		handleInjectState(st.(*model.InjectState), r)
	}
	return nil
}
