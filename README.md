# GoAuto
>"What makes you so ashamed of being a grownup?" - The War Doctor

## Overview
Full disclosure: This is ALPHA software and my first Go program. It works on my machine so mileage may vary. Suggestions, corrections, bug reports are very welcome. I have every intention of growing this into a replacement for Grunt, Gulp, Make, Rake, etc. 

Building a general purpose build tool with GoAuto that used config files would be a fairly trivial project. Feel free if that is your thing.

Task automation for grownups. GoAuto is a package that makes building a native executable tailored to a specific work flow, simple. 

Here is a complete example of a Go build process triggered by any source file in a project changing.

```golang
package main

import "goauto"

func main() {
	// Create a Pipeline
	p := goauto.Pipeline{Nmae: "Go Pipeline"}

	// Add all my project directories recursevely ignoring hidden directories
	if err := p.AddRecWatch("src/github.com/me/myprojects", true); err != nil {
		panic(err)
	}

	// Create a Workflow
	wf := goauto.Workflow{Name: "Go Build Workflow"}

	// Add a pattern to watch
	wf.AddPattern(".*\\.go$")

	// Add Tasks to the Workflow
	wf.AddTask(goauto.NewGoVetTask())
	wf.AddTask(goauto.NewGoTestTask())
	wf.AddTask(goauto.NewGoLintTask())
	wf.AddTask(goauto.NewGoInstallTask())

	p.AddWorkflow(&wf)

	// Start watching
	done := make(chan bool)
	p.Watch(done)
}
```

## Features
* No config files
* Built in support for Go projects
* No new syntax just Go
* Highly customizable
* Fast
* Small
* Tool for building tools

## Instalation
	go get github.com/dshills/goauto
	
## Concepts

### Tasks

Tasks are generally small, atomic pieces of work. Run tests, compile, copy a file, etc. They are what makeup a Workflow. A task implements the Tasker interface

#### Task Builtins
GoAuto includes a number of pre built tasks that can be used directly.

* NewGoTestTask will run tests for a project
* NewGoVetTask will run vet for a project
* NewRenameTask will rename a file
* Many more 

#### Task Generators
The built in tasks are a great way to get started with GoAuto. They do many useful things and serve as guides for building your own tasks. GoAuto also includes generator functions that will help you build your own simple tasks. NewTask, NewShellTask and NewGoPrjTask are examples of generic task generators.

NewTask is the most generic of the generators and can be used for building your own task. NewTask useage will be covered in Task Building

NewShellTask will return a new Tasker that will call a shell function with arguments on the target file. 

	func NewShellTask(cmd string, args ...string) Tasker

	st := NewShellTask("echo", "-n") // a Task that will echo the source file name


NewGoPrjTask will produce a new Tasker that will call the go command in the directory of target file

	func NewGoPrjTask(gocmd string, args ...string) Tasker 

	gt := NewGoPrjTask('build') // A fancy new task that will build a project

#### TaskInfo
Before diving into task building we need to introduce the TaskInfo struct. TaskInfo is passed between tasks as they run. 

```golang
type TaskInfo struct {
	Src        string 
	Target     string
	Buf        bytes.Buffer
	Tout, Terr io.Writer
}
```

Your tasks are expected to update Target and Buf and to use Tout and Terr for output. For example a task that renames a file would set Target equal to the new file name. If your task has output useful to another task then reset the Buf and write it. User messages or error text can be written to Tout and Terr. Src will be set to the Target, if set, of the last run task in the workflow or if it's the first task will be set to the file matched by the workflow.

#### Task Building
The real power comes from building custom tasks. This can be done using the NewTask generator or by writing a Tasker compliant interface. Here are examples of both for calling the cat shell command.

	func NewTask(t Transformer, r Runner) Tasker 

A Transformer is the function used to convert the incoming file name to something new. A number of Transformers are built in.

	func(string)string

A runner is the function called to run your task and is in the form

	func foo(i *TaskInfo)error

So lets write a task that cats a file. A cat task is already included but it is a simple example. We are using the Identity Transformer which just returns the file name passed to it. 

```golang
func myCat(i *goauto.TaskInfo) (err error) {
	cmd := exec.Command("cat", i.Target)
	i.Buf.Reset()
	cmd.Stdout = &i.Buf // Write the output to the buffer
	cmd.Stderr = i.werr // use werr as stdin
	defer func() {
		i.Tout.Write(i.Buf.Bytes()) // Write the buffer to Tout
	}()
	return cmd.Run()
}

t := goauto.NewTask(goauto.Identity, myCat)
// Identity is a built in Transformer that returns what was passed to it
// We could have written goauto.NewTask(func(f string)string {return f}, myCat)
```

Here it is written as a Tasker. In this case we don't need a Transformer because we are controlling the entire task from start to finish. In this simple example it is actually shorter to make our own Tasker

```golang
type myCatTask struct{}
func (t *myCatTask)Run(i *goauto.TaskInfo) (err error) {
	i.Target = i.Src // Not changing the file name so not technically required but a good habbit
	cmd := exec.Command("cat", i.Target)
	i.Buf.Reset()
	cmd.Stdout = &i.Buf // Write the output to the buffer
	cmd.Stderr = i.werr // use werr as stdin
	defer func() {
		i.Tout.Write(i.Buf.Bytes()) // Write the buffer to Tout
	}()
	return cmd.Run()
}
```

Now we can add our task to one or more Workflows.

### Workflows
Workflows sequentially run a set of tasks for files matching a regular expression pattern.  Workflows only really need to know two things, what files to process and what tasks to perform. 

In this example we use the task that we created above and add it to a Workflow that will run on any files with the .go extension.

```golang
wf := goauto.NewWorkflow("Cat Workflow", ".*\\.go$", new(myCatTask))

Or 

wf := &goauto.Workflow{Name:"Cat Workflow"}
wf.AddPattern(".*\\.go$")
wf.AddTask(new(myCatTask)) 
```

Workflows run tasks sequentially, passing the TaskInfo to each task on the way. Before a task is run the TaskInfo.Src is updated to the TaskInfo.Target of the previous task if it was set. TaskInfo.Src is set to the matching file name for the first task. 
Any task that returns an error will stop the Workflow. If you want the Workflow to continue even if an error occurs make sure to handle the error and not return it.

Now we can add our Workflow to a Pipeline.

### Pipelines
A Pipeline monitors one or more file system directories for changes. When it detects a change it asks each Workflow if the specific file is a match and if it is launches the Workflow. Workflows are run sequentially, however future versions may include an option for running concurrently. Output and Error io can be set for a Pipeline. If not specified it will use StdIn and StdErr. One or more Pipelines can be declared. Running them concurrently is a choice left to the developer.

This example uses the Workflow we created above and adds it to the Pipeline. Watches can be absolute or $GOPATH relative.

```golang
p := goauto.NewPipeline("My Pipeline", "src/github.com/dshills/my/myproject", os.Stdout, os.Stderr, wf)

Or

p := goauto.Pipeline{}
p.AddWatch("src/github.com/me/myproject")
p.AddWorkflow(wf)

done := make(chan bool)
p.Watch(done)

```

Watch directories can be added with AddWatch(watchDir string) to add a single directory or AddRecWatch(watchDir string, ignoreHidden bool) to add directories recursively optionally ignoring hidden directories.

## License
Copyright 2015 Davin Hills. All rights reserved.
MIT license. License details can be found in the LICENSE file.
