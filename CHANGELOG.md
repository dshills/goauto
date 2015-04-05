# Change Log

## [0.1.4](https://github.com/dshills/goauto/tree/0.1.4) (2015-04-05)

**Enhancements:**
* Simpler and more consistent API
* Concurrent Workflow support
* Built in task as sub packages (gotask, shelltask, webtask)
* RestartTask for starting and stopping processes i.e. Go server or Web server
* NewSassTask for creating a task that runs the command line sass tool
* Pipeline recursive directory watches will add new directories as they are added

**Fixes:**
* Better handling of long running Workflows by queueing Pipeline notifications 

**Closed issues:**

- Handle directory add after Watching has begun [\#2](https://github.com/dshills/goauto/issues/2)

- Multiple GOPATH support for Transformers [\#1](https://github.com/dshills/goauto/issues/1)
