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

import "testing"

func TestInsertTask(test *testing.T) {
	stages := make([]stage, 1)

	for i := 0; i < 3; i++ {
		t := task{Name: "Task", Description: "Some Task", Stage: (i + 10), Executable: "", Arguments: nil}
		stages = insertTask(stages, t)
	}

	if len(stages) != 4 {
		test.Error("error")
	}

	for i := 1; i < 4; i++ {
		if stages[i].Num != i+10-1 {
			test.Error("error")
		}
	}

	if stages[0].Num != 0 {
		test.Error("error")
	}

	t := task{Name: "Task", Description: "Some Task", Stage: 1, Executable: "", Arguments: nil}
	stages = insertTask(stages, t)

	if len(stages) != 5 {
		test.Error("erro")
	}

	if stages[4].Num != 12 || stages[1].Num != 1 {
		test.Error("error")
	}

	t = task{Name: "Task", Description: "Some Task", Stage: 100, Executable: "", Arguments: nil}
	stages = insertTask(stages, t)

	if stages[4].Num != 12 ||
		stages[1].Num != 1 ||
		stages[5].Num != 100 {
		test.Error("error")
	}

	t = task{Name: "Task", Description: "Some Task", Stage: 99, Executable: "", Arguments: nil}
	stages = insertTask(stages, t)

	if stages[4].Num != 12 ||
		stages[1].Num != 1 ||
		stages[6].Num != 100 ||
		stages[5].Num != 99 {
		test.Error("error")
	}

	t = task{Name: "Task", Description: "Some Task", Stage: -1, Executable: "", Arguments: nil}
	stages = insertTask(stages, t)

	if stages[0].Num != -1 ||
		stages[1].Num != 0 ||
		stages[2].Num != 1 {
		test.Error("error")
	}

	if len(stages) != 8 {
		test.Error("error")
	}

	stages = make([]stage, 1)

	for i := 0; i < 3; i++ {
		t := task{Name: "Task", Description: "Some Task", Stage: i, Executable: "", Arguments: nil}
		stages = insertTask(stages, t)
	}

	if len(stages) != 3 {
		test.Error("error")
	}

	t = task{Name: "a", Description: "Some Task", Stage: 0, Executable: "", Arguments: nil}
	stages = insertTask(stages, t)

	t = task{Name: "b", Description: "Some Task", Stage: 1, Executable: "", Arguments: nil}
	stages = insertTask(stages, t)

	t = task{Name: "c", Description: "Some Task", Stage: 2, Executable: "", Arguments: nil}
	stages = insertTask(stages, t)

	if stages[0].Tasks[0].Name != "Task" ||
		stages[0].Tasks[1].Name != "a" ||
		stages[1].Tasks[1].Name != "b" ||
		stages[2].Tasks[1].Name != "c" {
		test.Error("error")
	}

	t = task{Name: "d", Description: "Some Task", Stage: 2, Executable: "", Arguments: nil}
	stages = insertTask(stages, t)

	if stages[2].Tasks[2].Name != "d" {
		test.Error("error")
	}
}
