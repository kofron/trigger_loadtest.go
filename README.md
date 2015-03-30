### trigger_loadtest

#### Dependencies
The golang runtime.  The program is written in standard library Go,
so only `go` is required:

```
go get github.com/kofron/trigger_loadtest
go build github.com/kofron/trigger_loadtest && ./trigger_loadtest
```


#### Introduction and running the tool
This tool is designed to test (and stress) the data handling systems
for the Project 8 experiment by writing analyzable data to a working
directory at some user-specified rate.  The main loop of the program
is simple:

* A random number is thrown from a exponential distribution according
to the user-specified rate parameter lambda, and is interpreted as the
time-to-next-event (i.e. the sequence of triggers is a Poisson process).
* Another random number on [0,1) is thrown.  If this number is less than
the probability-of-false-alarm (pfa) specified by the user, the next
event will be noise.  Otherwise, it will be signal.
* A random filename is generated, and either the signal data or the noise
data is written to disk using this random filename.
* This process repeats until the simulated run length has elapsed or the
user sends CTRL-C (or SIGTERM).

The best documentation of the usage of the utility comes from flag.Usage -
to see it, just run this program at the command line with no arguments.  

#### Limitations and assumptions
* First of all, this tool does *not* generate random data - it uses an existing
  noise event and an existing signal event to generate load.  Therefore any 
  spectrum calculated from these data will not be representative of any real
  physics.  

* The duration of the event itself is not taken into account when calculating
  the elapsed time. 
