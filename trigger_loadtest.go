/*
 * author: jared kofron <jared.kofron@gmail.com>
 * date: 3/25/2015
 *
 * this program is designed to produce an artificial load which is
 * similar to what could be expected from a real-time spectrum analyzer.
 * the purpose is to stress the automation facilities for data handling
 * implemented by project 8.
 *
 * there are six knobs to turn:
 *   1) the noise file (a MAT file which can be processed by katydid)
 *   2) the signal file (also a MAT file which can be processed by katydid)
 *   3) the probability of noise triggers (between 0 and 1)
 *   4) the overall rate (in hertz i.e. 10.5)
 *   5) the local working directory
 *   6) the total number of events to produce (default is infinite)
 *
 * main is simple:
 *   1) open the signal and noise files
 *   2) loop forever.  at the beginning of the loop, decide if this
 *      is to be a signal or noise event based on the probability
 *      of noise triggers.
 *   3) generate a random filename.
 *   4) save either signal or noise with a random filename to the
 *      working directory
 *   5) throw an exponential time-to-next-event according to the overall
 *      rate (which is interpreted as lambda), and sleep for that amount
 *      of time
 *   6) if an interrupt is received or the total number of events is reached,
 *      exit.  interrupt yields an exit code of 1.  "normal" exit is 0.
 *
 * on ctrl-C, will gracefully exit the loop and stop.
 */
package main

import (
	"flag"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Configuration variables that are set by command line flags.
var (
	T                                      time.Duration = 30 * time.Second
	lambda                                 float64 = 10.0
	pfa                                    float64 = 0.5
	work_dir, sig_filename, noise_filename string
)

func checkFlags() (e error) {
	if work_dir == "REQUIRED" {
		e = errors.New("Working directory is required!")
	}
	if sig_filename == "REQUIRED" {
		e = errors.New("Signal filename is required!")
	}
	if noise_filename == "REQUIRED" {
		e = errors.New("Noise filename is required!")
	}
	return
}

func main() {
	// command line flag setup
	flag.StringVar(&work_dir,
		"working-dir",
		"REQUIRED",
		"directory to save data files in")
	flag.StringVar(&sig_filename,
		"signal-file",
		"REQUIRED",
		"The signal file archetype")
	flag.StringVar(&noise_filename,
		"noise-file",
		"REQUIRED",
		"The noise file archetype.")
	flag.Float64Var(&lambda,
		"rate",
		10.0,
		"average rate of events.  equal to 1/lambda.")
	flag.Float64Var(&pfa,
		"pfa",
		0.1,
		"Probability of false alarm, i.e. percent of total events which are noise.")
	flag.DurationVar(&T,
		"run-length",
		1 * time.Second,
		"Total time of run in Duration format.")

	flag.Parse()
	if flags_err := checkFlags(); flags_err != nil {
		log.Printf("\n%s\n\n",flags_err)
		flag.Usage()
		os.Exit(-1)
	}

	// working directory
	work_dir = CanonicalizeDirName(work_dir)

	if IsAbsolutePath(work_dir) == false {
		log.Fatal("(FATAL) Working directory must be an absolute path!")
	}

	work_dir_info, wd_err := os.Stat(work_dir)
	if wd_err != nil {
		log.Fatal("(FATAL) Can't open working directory!")
	} else if work_dir_info.IsDir() == false {
		log.Fatal("(FATAL) Working directory isn't a directory!")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// open signal and noise files.  verify that they can be
	// opened and exist.
	var write_buffer []byte
	var write_handle *os.File
	var write_err error

	sig_data, sig_file_err := ioutil.ReadFile(sig_filename)
	if sig_file_err != nil {
		log.Fatal("(FATAL) couldn't open signal data file!")
	}

	noise_data, noise_file_err := ioutil.ReadFile(noise_filename)
	if noise_file_err != nil {
		log.Fatal("(FATAL) couldn't open noise data file!")
	}

	// handle SIGTERM or ctrl-c gracefully
	itr_chan := make(chan os.Signal, 1)
	signal.Notify(itr_chan, os.Interrupt)
	signal.Notify(itr_chan, syscall.SIGTERM)

	// if this thread receives the signal to stop, it will cause the
	// program to exit
	go func() {
		<-itr_chan
		log.Printf("Stop requested.  Terminating...\n")
		os.Exit(1)
	}()

	log.Printf("Starting load test in 3 seconds...\n")
	time.Sleep(3 * time.Second)

	var triggers int = 0
	var elapsed time.Duration = 0 *time.Second
	for elapsed < T {
		time_to_next := time.Duration(1000*rng.ExpFloat64()/lambda) * time.Millisecond

		if rng.Float64() < pfa {
			write_buffer = noise_data
		} else {
			write_buffer = sig_data
		}

		write_handle, write_err = os.Create(RandomFileName(work_dir, rng))
		if write_err != nil {
			log.Fatal("(FATAL) couldn't create output file!")
		}
		defer write_handle.Close()

		write_handle.Write(write_buffer)

		time.Sleep(time_to_next)
		elapsed += time_to_next
		triggers += 1
	}
	log.Printf("Exiting normally.  %d triggers recorded (average rate: %0.2f Hz).", triggers, float64(triggers)/elapsed.Seconds())

}
