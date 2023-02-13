package main

import (
	"log"
	"math/rand"
	"time"
)

const TimeOutTime = 3
const MeanArrivalTime = 4

const (
	LEADER = iota
	FOLLOWER
	CANDIDATE
)

var currentState int

func handleCandidate() {
	log.Println("THREAD 3 : Entered candidate state!")
	// go startElection()
	// Here, send election requests to every other node
}

func timeout(inCh chan struct{}, outCh chan struct{}, interval time.Duration) {
	for {
		// On each iteration new timer is created
		select {
		case <-time.After(interval):
			log.Println("THREAD 1 : Sending a timeout to the other thread!")
			outCh <- struct{}{}
		case <-inCh:
			log.Println("THREAD 1 : Recieved a heartbeat from the other thread!")
			if currentState == FOLLOWER {
				currentState = CANDIDATE
				go handleCandidate()
			}
			if currentState == CANDIDATE {
				currentState = FOLLOWER
			}
		}
	}
}

func listener(inCh chan struct{}, outCh chan struct{}, MeanArrivalTime int) {
	for {
		select {
		case <-time.After(time.Duration(rand.Intn(MeanArrivalTime)) * time.Second):
			log.Println("THREAD 2 : Sending a heartbeat to the other thread!")
			inCh <- struct{}{}
		case <-outCh:
			log.Println("THREAD 2 : Recieved a timeout from the other thread!")
			// handle timeout response
		}
	}
}

func main() {
	currentState = FOLLOWER
	const interval = time.Second * TimeOutTime
	// channel for incoming messages
	var inCh = make(chan struct{}, 1)
	var outCh = make(chan struct{}, 1)
	log.Println("Start")
	go timeout(inCh, outCh, interval)
	go listener(inCh, outCh, MeanArrivalTime)

	// prevent main to stop for a while
	<-time.After(30 * time.Second)
}
