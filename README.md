TaskQ - Multidimensional Task Queue
===================


**TaskQ** allows you to define a bunch of tasks and then execute them in a previously defined order. It can be used to 'glue' together different kind of tools and then execute them all at once. Each task can be associated with a stage. When a stage gets loaded all tasks belonging to the stage will be executed as go routines. Stages are executed from a low stage numbers to high stage numbers.  Due to its flexibility and simplicity TaskQ can be used for a whole range of automation tasks.

----------

### Installation
```
git clone https://github.com/Schreibmaschine/taskq.git
go build
./taskq
```
----------


### Usage

```
./taskq 

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
		
[>] Error: Not enough commandline arguments.
Usage of ./taskq:
  -p string
    	Path to .json task file (required)
  -q	Disable console output
  -t int
    	Number of threads (default 3)
```
----------
### Task definition
All tasks are defined in **.json format**. In the following you can see an example **.json file** containing four tasks ready to be executed via **TaskQ**, as you can see from the file it defines the stages **0**, **1** and **2**. All tasks belonging to stage **0** are executed first, after this tasks from stage **1**, and then all remaining tasks from stage **2**. As you can see from stage **1** it is possible to associate more then one task with a stage. The stage numbers are also variable so that you could have used the stages **101**, **102** and **103** instead of **0**, **1** and **2**.
```
[
 {
        "stage": 0,
        "name": "touch-stage0",
        "description": "create a file named stage0.txt",
        "executable": "/bin/bash",
        "arguments": ["-c", "touch stage0.txt"]
 },

 {
        "stage": 1,
        "name": "echo-stage1",
        "description": "write 'stage1' into the file created in stage 0",
        "executable": "/bin/bash",
        "arguments": ["-c", "echo 'stage1' > stage0.txt"]
 },

 {
        "stage": 1,
        "name": "sleep-stage1",
        "description": "sleep for 15 seconds",
        "executable": "/bin/bash",
        "arguments": ["-c", "sleep 15s"]
 },
 
 {
        "stage": 2,
        "name": "rm-stage2",
        "description": "delete the file created in stage 0",
        "executable": "/bin/bash",
        "arguments": ["-c", "rm stage0.txt"]
 }
]
```
----------

### Execute tasks
```
$ ./taskq -p /home/user/tasks.json

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
		
[*] Start
[*] Opening .json task file
[*] Parsing .json task file
[*] Arranging tasks into stages
[*] Spawning workers
[*] Processing tasks from stage [1/3]
[*] Processing Task [1/1]
[*] Processing tasks from stage [2/3]
[*] Processing Task [2/2]
[*] Processing tasks from stage [3/3]
[*] Processing Task [1/1]
[*] Finished

```
----------
EOF
