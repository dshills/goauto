# GoAuto [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/dshills/goauto)
>"What makes you so ashamed of being a grownup?" - The War Doctor

## Overview
**UPDATE** Pipeline now includes an experimental flag OSX. If you are using OS X and have received the "To many open files" warning this is an attempt to fix it. The watcher code has now been extracted into it's own interface and can use the experimental OSX events package https://github.com/go-fsnotify/fsevents. I have done heavy testing locally with no problems but your mileage may vary. This should have no affect on the current usage of GoAuto.

Task automation for grownups. GoAuto is a package that makes building a native executable tailored to a specific work flow, simple. 

Here is a complete example of a Go build and test process triggered by any source file in a project changing. The workflows are run concurrently.

```go
package main

import (
	"path/filepath"

	"github.com/dshills/goauto"
	"github.com/dshills/goauto/gotask"
)

func main() {
	// Create a pipeline (Develop using Verbose, Change to Silent after testing)
	p := goauto.NewPipeline("Go Pipeline", goauto.Verbose)
	defer p.Stop()

	// Watch directories recursively, ignoring hidden directories
	wd := filepath.Join("src", "github.com", "me", "myproject")
	if err := p.WatchRecursive(wd, goauto.IgnoreHidden); err != nil {
		panic(err)
	}

	// Create an install workflow, don't wait for the tests to run
	wf := goauto.NewWorkflow(gotask.NewGoVetTask(), gotask.NewGoLintTask(), gotask.NewGoInstallTask())
	wf.Name = "Install"
	wf.Concurrent = true

	// Add a file pattern to match
	if err := wf.WatchPattern(".*\\.go$"); err != nil {
		panic(err)
	}

	// Create a Test workflow
	wf2 := goauto.NewWorkflow(gotask.NewGoTestTask())
	wf2.Name = "Test"
	wf2.Concurrent = true

	// Add a file pattern to match
	if err := wf2.WatchPattern(".*\\.go$"); err != nil {
		panic(err)
	}

	// Add workflows to pipeline
	p.Add(wf, wf2)

	// start the pipeline, it will block
	p.Start()
}
```

## Features
* No config files
* Built in support for Go projects
* No new syntax just Go
* No external dependencies
* Highly customizable
* Fast
* Small
* Tool for building tools

Building a general purpose build tool with GoAuto that used config files would be a fairly trivial project. Feel free if that is your thing.


## Installation
	go get github.com/dshills/goauto
	
## Concepts

### Pipelines
A Pipeline monitors one or more file system directories for changes. When it detects a change it asks each Workflow if the specific file is a match and if it is launches the Workflow. Output and Error io can be set for a Pipeline. If not specified it will use StdIn and StdErr. One or more Pipelines can be declared. Running them concurrently is a choice left to the developer.

GoAuto follows the go tools convention of no news is good news. It will silently watch for file changes, launch workflows and tasks without any output other than the output from the task itself. If running a task like a go tool that has the same philosophy, no output will be generated at all. When first writing a set of tasks this can be a little disconcerting. Did it run? Did it work?

```go
goauto.Verbose
goauto.Silent
```

Verbose will print debug information about what events are being received by a Pipeline, starting a workflow, and tasks being run. At some point an option may be added to just show specific output. When writing your own tasks you always have the choice of outputting whatever you wish.

Watches can be absolute or $GOPATH relative.

```go
// Create a pipeline
p := goauto.NewPipeline("Go Pipeline", goauto.Verbose)

// watch directories recursively, ignoring hidden directories
wd := filepath.Join("src", "github.com", "myprojects")
if err := p.WatchRecursive(wd, goauto.IgnoreHidden); err != nil {
	panic(err)
}
```

Watch directories can be added with Watch to add a single directory. The absolute path of the added path will be returned.

	func (p *Pipeline) Watch(watchDir string) (string, error)
	
To add directories recursively use WatchRecursive optionally ignoring hidden directories

	func (p *Pipeline) WatchRecursive(watchDir string, ignoreHidden bool) error 

Adding Workflows are added using Add

	func (p *Pipeline) Add(ws ...Workflower)

After adding Workflows and Tasks to your pipeline simply tell the Pipeline to Start watching. The Pipeline will block. To stop a running Pipeline call Stop.

```go
p.Start()

Or

go p.Start()

	do stuff

p.Stop()
```

### Workflows

Workflows run a set of tasks for files matching a regular expression pattern.  Workflows only really need to know two things, what files to process and what tasks to perform. Workflow implements the Workflower interface.

Here we create a Workflow with a number of built in tasks. We then add a pattern to match using WatchPattern.

```go
// Create a workflow with some tasks
wf := goauto.NewWorkflow(goauto.NewGoVetTask(), goauto.NewGoTestTask(), goauto.NewGoLintTask(), goauto.NewGoInstallTask())

// Add a regex pattern to match
err := wf.WatchPattern(".*\\.go$")
```
Tasks can also be added using Add

	func (wf *Workflow) Add(tasks ...Tasker)

Workflows run tasks sequentially, passing the TaskInfo struct (See Tasks below) to each task on the way. Before a task is run the TaskInfo.Src is updated to the TaskInfo.Target of the previous task if it was set. TaskInfo.Src is set to the matching file name for the first task. 

Any task that returns an error will stop the Workflow. If you want the Workflow to continue even if an error occurs make sure to handle the error and not return it.

#### Advanced Options

```go
// A Workflow represents a set of tasks for files matching one or more regex patterns
type Workflow struct {
	Name       string
	Concurrent bool
	Op         Op
	Regexs     []*regexp.Regexp
	Tasks      []Tasker
}
```

Setting Concurrent to true will run a Workflow concurrently. This should be used with caution. If multiple Workflows work with the same set of files there is a potential for confusion and even data loss.

	Op = goauto.Create | goauto.Write | goauto.Remove | goauto.Rename | goauto.Chmod

By default a Workflow will check file match for Create, Write, Remove, and Rename. This can be controlled by setting the Op value.

The Workflow struct implements the Workflower interface. Most use cases will have no need for anything more than a Workflow, however, Pipelines will accept anything that implements the Workflower interface. An example might be a new Workflower that implemented the WatchPattern using glob syntax rather than a regex. 

### Tasks

Tasks are generally small, atomic pieces of work. Run tests, compile, copy a file, etc. They are what makeup a Workflow. A task implements the Tasker interface

#### Task Builtins
GoAuto includes a number of pre built tasks that can be used directly.

##### goauto/gotask

* NewGoPrjTask will run a go command with arguments
* NewGoTestTask will run tests for a project
* NewGoVetTask will run vet for a project
* NewGoBuildTask will run build for a project
* NewGoInstallTask will run install for a project
* NewGoLintTask will run golint for a project
* NewGoMetaLinter task for github.com/alecthomas/gometalinter

##### goauto/shelltask

* NewShellTask task that runs a shell command with arguments
* NewCatTask task that cats a file
* NewRemoveTask task that deleted a file
* NewMoveTask task that moves a file
* NewMkdirTask task that makes a new directory
* NewCopyTask task that copies a file
* NewRestartTask task that will restart an executable file such as a Go server or Web server

##### goauto/webtask

* NewSassTask task that runs sass command line utility with options

#### Task Generators
The built in tasks are a great way to get started with GoAuto. They do many useful things and serve as guides for building your own tasks. GoAuto also includes generator functions that will help you build your own simple tasks. NewTask, NewShellTask and NewGoPrjTask are examples of generic task generators.

NewTask is the most generic of the generators and can be used for building your own task. NewTask usage will be covered in Task Building

NewShellTask will return a new Tasker that will call a shell function with arguments on the target file. 

	func NewShellTask(cmd string, args ...string) Tasker

	st := shelltask.NewShellTask("echo", "-n") // a Task that will echo the source file name


NewGoPrjTask will produce a new Tasker that will call the go command in the directory of target file

	func NewGoPrjTask(gocmd string, args ...string) Tasker 

	gt := gotask.NewGoPrjTask('build') // A fancy new task that will build a project

#### TaskInfo
Before diving into task building we need to introduce the TaskInfo struct. TaskInfo is passed between tasks as they run. 

```go
// A TaskInfo contains the results of running a Task
type TaskInfo struct {
	Src        string       // Incoming file name for a task to process
	Target     string       // Output file name after running a task
	Buf        bytes.Buffer // Output of running a task
	Tout, Terr io.Writer    // Writers to write output and errors
	Collect    []string     // List of file names processed by a Workflow
	Verbose    bool         // output debug info
}
```

Your tasks are expected to update Target and Buf and to use Tout and Terr for output. For example a task that renames a file would set Target equal to the new file name. If your task has output useful to another task then reset the Buf and write it. User messages or error text can be written to Tout and Terr. 

As the Workflow executes each task the Src will be set to the Target of the last run task. For the first task in a Workflow Src is set to the filename matched by the Workflow.

By using Buf and Target a Workflow creates a flow similar to a using a Unix pipe 

Collect keeps a running list of file targets over the course of one run of a Workflow. This gives tasks access to run functions on all the files processed by the Workflow.

#### Task Building
The real power comes from building custom tasks. This can be done using the NewTask generator or by writing a Tasker compliant interface. Here are examples of both for calling the cat shell command.

	func NewTask(t Transformer, r Runner) Tasker 

A Transformer is the function used to convert the incoming file name to something new. A number of Transformers are built in.

	func(string)string

##### Built In Transformers

* Identity function that returns the string passed in 
* GoRelBase function that returns the file path relative to GOPATH 
* GoRelDir function that returns the directory path relative to GOPATH
* GoRelSrcDir function that returns the directory path relative to GOPATH/src
* ExtTransformer function that returns a Transformer function that returns file path with a new extension

A runner is the function called to run your task and is in the form

	func foo(i *TaskInfo)error

So lets write a task that cats a file. A cat task is already included but it is a simple example. We are using the Identity Transformer which just returns the file name passed to it. 

```go
func myCat(i *goauto.TaskInfo) (err error) {
	cmd := exec.Command("cat", i.Target)
	i.Buf.Reset()
	cmd.Stdout = &i.Buf // Write the output to the buffer
	cmd.Stderr = i.Terr // use Terr as stderr
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

```go
type myCatTask struct{}
func (t *myCatTask)Run(i *goauto.TaskInfo) (err error) {
	i.Target = i.Src // Not changing the file name so not technically required but a good habit
	cmd := exec.Command("cat", i.Target)
	i.Buf.Reset()
	cmd.Stdout = &i.Buf // Write the output to the buffer
	cmd.Stderr = i.Terr // use Terr as stderr
	defer func() {
		i.Tout.Write(i.Buf.Bytes()) // Write the buffer to Tout
	}()
	return cmd.Run()
}
```

## To Do
* More built ins for Web development LESS, Reload (Certainly can be done now but it would be nice to have built ins)
* Test large, concurrent, multi Pipeline, multi Workflow systems

## Alternatives

* [Slurp](https://github.com/omeid/slurp) Go
* [Gulp](http://gulpjs.com/) Node.js
* [Grunt](http://gruntjs.com/) Node.js

## License
Copyright 2015 Davin Hills. All rights reserved.
MIT license. License details can be found in the LICENSE file.
