package runtime

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/itchyny/gojq"
	"github.com/serverlessworkflow/sdk-go/model"
)

func handleEventState(state *model.EventState, r *Runtime) error {
	fmt.Println("Event:", state.GetName())
	//TODO 
	if (state.GetTransition() != nil) {
		newStateName := state.Transition.NextState
		ns := findNewStateObject(newStateName, r)
		fmt.Println("New State = ", ns)
		r.begin(ns)
		return nil
	}
	fmt.Println("This is the end..")
	return nil
}

func handleOperationState(state *model.OperationState, r *Runtime) error {
	fmt.Println("Operation:", state.GetName())
	// TODO
	//dirty ifs here for now
	if (state.GetTransition() != nil) {
                newStateName := state.Transition.NextState
                ns := findNewStateObject(newStateName, r)
                fmt.Println("New State = ", ns)
        }
	if (state.GetTransition() == nil) {
		fmt.Println("this is the end..")
		return nil
	}
	// Check for the action Mode (default: sequential)
	if (state.ActionMode == "sequential") {
		fmt.Println("This OperationState is sequential")
		return nil
	}
	fmt.Println("This OperationState is parallel")
	return nil
}

func handleInjectState(state *model.InjectState, r *Runtime) error {
	fmt.Println("--> Current state Inject: ", state.GetName())
	//injectFilter := state.GetStateDataFilter()
	injectData := state.Data
	fmt.Println("Data of inject state: ", injectData)
	//fmt.Println("Input filter: ", injectFilter.Input, " Output filter: ", injectFilter.Output)
	//outFilter := strings.Split(injectFilter.Output, " ")[1]
	//outFilter = strings.Split(outFilter, ".")[1]
	if (state.GetTransition() != nil) {
		newStateName  := state.Transition.NextState
		ns := findNewStateObject(newStateName, r) //type of model.State
		ns2 := ns.(*model.InjectState) //typecasting so as to be compatible
		ns2.Data = injectData
		fmt.Println("ns2 : ", ns2)
		r.begin(ns2)
		return nil
	}
	fmt.Println("This is the end..")
	return nil
}


func findNewStateObject(name string, r *Runtime) model.State {
	fmt.Println("Searching the next State by name.. ")
	states := r.Workflow.States
	for _, state := range states {
		if  (name == state.GetName()){
			return state
		}
	}
	fmt.Println("Next state not found")
	return nil
}

func HandleDataBasedSwitch(state *model.DataBasedSwitchState, in []byte, r *Runtime) error {
	for _, cond := range state.DataConditions {
		fmt.Println(cond.GetCondition())
		switch cond.(type) {
		case *model.TransitionDataCondition:
			var result map[string]interface{}
			json.Unmarshal(in, &result)
			op, _ := gojq.Parse(cond.GetCondition())
			iter := op.Run(result)
			v, _ := iter.Next()
			if err, ok := v.(error); ok {
				log.Fatalln(err)
			}
			// fmt.Printf("%v\n", v)
			if v.(bool) {
				fmt.Println("GOTO", cond.(*model.TransitionDataCondition).Transition.NextState)
				newStateName  := cond.(*model.TransitionDataCondition).Transition.NextState
				ns := findNewStateObject(newStateName, r)
				r.begin(ns)

			} else {
				fmt.Println("Not True")
			}
			// test := map[string]interface{}{"foo": []interface{}{"age", 2, 3}}

			// fmt.Println("Result is:", string(res))

			// return cond.(*model.TransitionDataCondition).Transition.NextState
			// if this condition is true
			// HandleTransition(state, ns)
			//find next state object
			// InferType()
		case *model.EndDataCondition:
			fmt.Println(cond.(*model.EndDataCondition).End)
			// this is the end, you know
			fmt.Println("This is the end..")
		}

	}
	return nil
}
