# GoAuto
>"What makes you so ashamed of being a grownup?" - The War Doctor

## Overview
Task automation for grownups. GoAuto is a package that makes building a native executable tailored to a specific work flow, simple. Task runners are relatively easy things. They just need to know what tasks, in what order, and in what place. 

Here is a complete example of running tests when any source file in a project changes.

```golang
package main
import "goauto"

func main() {
	// Workflow for .go files with a go test task
	wf, err := goauto.NewWorkflow("Test Workflow", ".*.go", goauto.NewGoTestTask())
	if err != nil {
		panic(err)
	}

	// Create a pipeline for a $GOPATH directory to watch and our workflow
	p := goauto.NewPipeline("My Pipeline", "src/github.com/me/myproject", os.Stdout, os.Stderr, wf)

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
* Support for large concurrent work flows

## Instalation
	go get github.com/dshills/goauto

## Tasks

Tasks are generally small, atomic pieces of work. Run tests, compile, copy a file, etc. They are what makeup a Workflow. 

### Task Builtins
GoAuto includes a number of pre built tasks that can be used directly.

* NewGoTestTask will run tests for a project
* NewGoVetTask will run vet for a project
* NewCatTask will cat a file
* Many more 

### Task Generators
The built in tasks are a great way to get started with GoAuto. They do many useful things and serve as guides for building your own tasks. GoAuto also includes generator functions that will help you build your own simple tasks. NewShellTask and NewGoPrjTask are two examples of generic task generators. NewShellTask will return a new Task pointer that will call a shell function with arguments on the target file. NewGoPrjTask will produce a new Task pointer that will call the go command with arguments on the project the target is contained in.

	st := NewShellTask("cat") // A shiny new task that will cat a file
	gt := NewGoPrjTask('build') // A fancy new task that will build a project

### Task Building
The real power comes from building custom tasks and that process is straightforward. To write a task only requires writing one function. It will take the form:

	func(task *Task, wout, werr io.Writer) error

So lets write a task that cats a file. A cat task is already included but it is a simple example.

```golang
func MyCat(t *Task, wout, werr io.Writer) error {
	cmd := exec.Command("cat", t.Target) // Create an exec command calling cat with our task target
	cmd.Stdout = wout // use wout as stdout
	cmd.Stderr = werr // use werr as stdin
	if err := cmd.Run(); err != nil {
		return err // returning an error will stop the workflow
	}
	return nil
}
```

That's it. We have a function that will cat any file passed to it in the Task struct. Let's turn that into a Task that we can use in our Workflow.

	task := NewTaskType(MyCat)

The Task struct looks like this:

```golang
type Task struct {
	Banner     string // Prints to wout before starting a task i.e. Compiling...
	Buffer     bytes.Buffer // Store task output to pass to the next task in the workflow
	FileName   string // The file name matched by the workflow or the output filename from the previous task
	Target     string // The file name after being run through TargetFunc if defined or Filename. Use this inside your functions
	TaskFunc   func(task *Task, wout, werr io.Writer) error // The real work of the task
	TargetFunc func(string) string // a function to manipulate the FileName => Target
}
```

If your task changes the file name or outputs to a new file make sure to change the FileName so the next task can continue the work you started. 

```golang
MyRenameFile(t.Target, newName)
t.FileName = newName
```

Now we can add our task to one or more Workflows.

## Workflows
Workflows sequentially run a set of tasks for files matching a regular expression pattern.  Workflows only really need to know two things, what files to process and what tasks to perform. 

In this example we use the task that we created above and add it to a Workflow that will run on any files with the .go extension.
```golang
t := goauto.NewTaskType(MyCat)

wf := goauto.Workflow{}
wf.AddPattern(".*.go")
wf.AddTask(t) 
```

Workflows run tasks sequentially, updating the Task.FileName and Task.Buffer into each task being run. Any task that returns an error will stop the Workflow. If you want the Workflow to continue even if an error occurs within a task make sure to handle the error in your task function and return a nil. 

Multiple Workflows within a pipeline will be run concurrently. Future versions may make this an option. *Caution* should be used, having multiple Workflows working on the same file set concurrently could have unpredictable results.

Now we can add our Workflow to a Pipeline.

## Pipelines
A Pipeline monitors one or more file system directories for changes. When it detects a change it asks each Worflow if the specific file is a match and if it is launches the Workflow. Workflows are run concurrently. Output and Error io can be set for a Pipeline. If not specified it will use StdIn and StdErr. One or more Pipelines can be declared. Running them concurrently is a choice left to the developer.

This example uses the Workflow we created above and adds it to the Pipeline. Watches can be absolute or $GOPATH relative.

```golang
p := goauto.Pipeline{}
p.AddWatch("src/github.com/me/myproject")

p.AddWorkflow(&wf)

done := make(chan bool)
p.Watch(done)

```

## License
Copyright 2015 Davin Hills. All rights reserved.
MIT license. License details can be found in the LICENSE file.
