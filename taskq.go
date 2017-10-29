//-------------------------------------------------------------------------------
// TaskQ - Multidimensional Task Queue
// Copyright (c) 2017 schreibmaschine
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see [http://www.gnu.org/licenses/](http://www.gnu.org/licenses/)
//-------------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Global flag to suppress console output if needed.
var quietPtr *bool

// A stage is used to bundle tasks with the same execution priority.
// A stage consits of a num and a list of tasks.
// The smaller the num the higher the execution priority.
type stage struct {
	Num   int
	Tasks []task
}

// A task is associated with a stage and points to an executable (binary/script).
// A task has some additional fields like Name, Description that can be used to
// specify what the task is used for.
type task struct {
	Stage       int      `json:"stage"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Executable  string   `json:"executable"`
	Arguments   []string `json:"arguments"`
}

// A worker receives tasks through a channel and executes them.
// 0 : On sucess
// 1 : On failure
func worker(id int, tasks <-chan task, ret chan<- int) {
	for t := range tasks {
		err := exec.Command(t.Executable, t.Arguments...).Run()
		if err != nil {
			cPrintf("[>] Error: Task execution failed.", *quietPtr)
			ret <- 1
		}
		ret <- 0
	}
}

// The arrangeWork function receives an unsorted list of tasks that is then associated
// with the appropriate stages.
func arrangeWork(tasks []task) []stage {
	stages := make([]stage, 1)
	for _, t := range tasks {
		stages = insertTask(stages, t)
	}
	return stages
}

// The doWork function takes a sorted list of stages and appends the tasks of each stage
// to the number of workers defined by the variable threads.
// The tasks belonging to the stage with a lowest Num are executed first.
func doWork(stages []stage, threads int) {
	for i, s := range stages {
		cPrintf("[*] Processing tasks from stage [%d/%d]\n", *quietPtr, i+1, len(stages))

		jobs := make(chan task, 100)
		ret := make(chan int, 100)

		for i := 0; i < threads; i++ {
			go worker(i, jobs, ret)
		}

		for _, t := range s.Tasks {
			jobs <- t
		}

		close(jobs)

		cPrintf("\r[*] Processing Task [0/%d]", *quietPtr, len(s.Tasks))
		for i := range s.Tasks {
			<-ret
			cPrintf("\r[*] Processing Task [%d/%d]", *quietPtr, i+1, len(s.Tasks))
		}
		cPrintf("\n", *quietPtr)
	}
}

// The insertTask function takes a list of stages and a task and inserts it into the appropriate stage
// if no such stage was found a new stage is created.
func insertTask(stages []stage, t task) []stage {
	for i := 0; i < len(stages); i++ {
		switch {
		case stages[i].Num > t.Stage:
			s := stage{Num: 0, Tasks: nil}

			stages = append(stages, s)
			copy(stages[i+1:], stages[i:])

			tasks := make([]task, 1)
			tasks[0] = t

			s = stage{Num: t.Stage, Tasks: tasks}

			stages[i] = s
			return stages
		case stages[i].Num == t.Stage:
			stages[i].Tasks = append(stages[i].Tasks, t)
			return stages
		case i+1 == len(stages):
			tasks := make([]task, 1)
			tasks[0] = t

			s := stage{Num: t.Stage, Tasks: tasks}
			stages = append(stages, s)
			return stages
		}
	}
	return nil
}

// The function processTaskFile takes a path to a .json file containing tasks.
// The function reads the file and also does the parsing of the tasks.
func processTaskFile(path string) []task {
	cPrintf("[*] Opening .json task file\n", *quietPtr)
	taskFile, err := os.Open(path)

	if err != nil {
		log.Fatal("[>] Fatal Error: Unable to read task file.")
	}

	cPrintf("[*] Parsing .json task file\n", *quietPtr)
	tasks := make([]task, 1)
	jsonParser := json.NewDecoder(taskFile)

	if err = jsonParser.Decode(&tasks); err != nil {
		log.Fatal("[>] Fatal Error: Unable to parse .json task file.")
	}

	return tasks
}

// The function cPrintf is a conditional Printf. Based on the provided boolean the
// function will either print the provided output or suppress it.
func cPrintf(format string, c bool, args ...interface{}) {
	if !c {
		fmt.Printf(format, args...)
	}
}

// print fancy banner
func printBanner() {
	cPrintf(
		`
		  Multidimensional Task Queue
		*=============================*
		  _______        _     ____  
		 |__   __|      | |   / __ \ 
		    | | __ _ ___| | _| |  | |
		    | |/ _  / __| |/ / |  | |
		    | | (_| \__ \   <| |__| |
		    |_|\__,_|___/_|\_\\___\_\
		   
		       by schreibmaschine		 
		*=============================*
		`,
		*quietPtr)

	cPrintf("\n", *quietPtr)
}

// parse arguments and orchestrate stuff
func main() {
	pathPtr := flag.String("p", "", "Path to .json task file (required)")
	threadPtr := flag.Int("t", 3, "Number of threads")
	quietPtr = flag.Bool("q", false, "Disable console output")
	flag.Parse()

	printBanner()

	if len(*pathPtr) == 0 {
		cPrintf("[>] Error: Not enough commandline arguments.\n", *quietPtr)
		flag.Usage()
		os.Exit(1)
	}

	cPrintf("[*] Start\n", *quietPtr)
	tasks := processTaskFile(*pathPtr)

	cPrintf("[*] Arranging tasks into stages\n", *quietPtr)
	stages := arrangeWork(tasks)

	cPrintf("[*] Spawning workers\n", *quietPtr)
	doWork(stages, *threadPtr)

	cPrintf("[*] Finished\n", *quietPtr)
}
