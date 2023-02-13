package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	Follower = iota
	Candidate
	Leader
)

var currentState int
var votes int

type node_information struct {
	CurrentTerm       int    `json:"currentTerm"`
	VotedFor          string `json:"votedFor"`
	Log               []int  `json:"log"`
	TimeoutInterval   int    `json:"timeoutInterval"`
	HeartbeatInterval int    `json:"heartbeatInterval"`
	LeaderInfo        string `json:"leaderInfo"`
}

type append_entry struct {
	LeaderID     string `json:"leaderID"`
	Term         int    `json:"term"`
	PrevLogIndex int    `json:"prevLogIndex"`
	PrevLogTerm  int    `json:"prevLogTerm"`
	Entries      int    `json:"entries"`
}

type request_vote struct {
	CandidateID  string `json:"candidateID"`
	Term         int    `json:"term"`
	LastLogIndex int    `json:"lastLogIndex"`
	LastLogTerm  int    `json:"lastLogTerm"`
}

type controller_request struct {
	Sender_Name string `json:"sender_name"`
	Request     string `json:"request"`
	Term        int    `json:"term"`
	Key         string `json:"key"`
	Value       string `json:"value"`
}

// reads from JSON
func readNodeInformation() node_information {
	var input node_information
	jsonNodeInformation, err := os.Open("Data/node_information.json")
	if err != nil {
		fmt.Println("Error opening file")
	}
	byteValue, _ := ioutil.ReadAll(jsonNodeInformation)
	json.Unmarshal(byteValue, &input)
	return input
}

// writes to JSON
func writeNodeInformation(nodeInformation node_information) {
	jsonNodeInformation, err := json.Marshal(nodeInformation)
	if err != nil {
		fmt.Println("Error marshalling!")
	}
	err = ioutil.WriteFile("Data/node_information.json", jsonNodeInformation, 0644)
	if err != nil {
		fmt.Println("Error writing file!")
	}
}

// sends data to connection buffer
func sendUDP(payload string, node_namespace string) {
	conn, err := net.Dial("udp", node_namespace+":8080")
	if err != nil {
		if conn != nil {
			conn.Close()
		}
	} else {
		fmt.Fprintf(conn, "%s", payload)
		conn.Close()
	}
}

func handleVoteRequest(requestVote request_vote) {
	//fmt.Println("Entered handleVoteRequest!")
	nodeInformation := readNodeInformation()
	if requestVote.Term > nodeInformation.CurrentTerm && nodeInformation.VotedFor == "" {
		// nodeInformation.CurrentTerm = requestVote.Term
		nodeInformation.VotedFor = requestVote.CandidateID
		writeNodeInformation(nodeInformation)
		sendUDP("true", requestVote.CandidateID)
		fmt.Println("Voted for " + requestVote.CandidateID)
	}
}

// if node's term is less than heartbreat, persist
func handleHeartbeat(heartbeat append_entry) {
	//fmt.Println("Entered handleHeartbeat!")
	//fmt.Println(heartbeat)
	nodeInformation := readNodeInformation()
	if heartbeat.Term > nodeInformation.CurrentTerm {
		nodeInformation.CurrentTerm = heartbeat.Term
		nodeInformation.VotedFor = ""
		nodeInformation.LeaderInfo = heartbeat.LeaderID
		writeNodeInformation(nodeInformation)
	}
}

func startElection() {
	fmt.Println("Entered startElection!")
	nodeInformation := readNodeInformation()
	nodeInformation.CurrentTerm += 1
	voteStructure := &request_vote{CandidateID: os.Getenv("NODE_NAME"), Term: nodeInformation.CurrentTerm, LastLogIndex: 0, LastLogTerm: 0}
	vote_json_bytes, _ := json.Marshal(voteStructure)
	vote_json_string := string(vote_json_bytes)
	for i := 0; i < 4; i++ {
		sendUDP(vote_json_string, os.Getenv("COMPANION_"+strconv.Itoa(i)))
	}
}

func sendHeartbeats() {
	fmt.Println(os.Getenv("NODE_NAME") + " became a leader!")
	nodeInformation := readNodeInformation()
	go func() {
		for {
			time.Sleep(time.Duration(nodeInformation.HeartbeatInterval) * time.Millisecond)
			heartbeatStructure := &append_entry{LeaderID: os.Getenv("NODE_NAME"), Term: nodeInformation.CurrentTerm, PrevLogIndex: 0, PrevLogTerm: 0, Entries: 0}
			heartbeat_json_bytes, _ := json.Marshal(heartbeatStructure)
			heartbeat_json_string := string(heartbeat_json_bytes)
			for i := 0; i < 4; i++ {
				sendUDP(heartbeat_json_string, os.Getenv("COMPANION_"+strconv.Itoa(i)))
			}
		}
	}()
}

func runUDPListener(c chan string) {
	address, _ := net.ResolveUDPAddr("udp", os.Getenv("NODE_NAME")+":8080")
	connection, _ := net.ListenUDP("udp", address)
	votes = 0
	for {
		p := make([]byte, 2048)
		_, _, _ = connection.ReadFromUDP(p)
		p = bytes.Trim(p, "\x00")
		heartbeatStructure := append_entry{}
		voteRequestStructure := request_vote{}
		_ = json.Unmarshal(p, &voteRequestStructure)
		err := json.Unmarshal(p, &heartbeatStructure)
		if heartbeatStructure.LeaderID != "" {
			handleHeartbeat(heartbeatStructure)
			c <- "heartbeat"
		}
		if err != nil {
			if string(p) == "true" && currentState == Candidate {
				fmt.Println("Recieved a vote!")
				votes += 1
				if votes > 1 {
					currentState = Leader
					votes = 0
					c <- "kill"
					go sendHeartbeats()
				}
			}
		}
		if heartbeatStructure.LeaderID == "" {
			if voteRequestStructure.CandidateID == "" && string(p) != "true" {
				controllerRequestStructure := controller_request{}
				_ = json.Unmarshal(p, &controllerRequestStructure)
				// The controller can send one of four request types to the node:
				// 1. "CONVERT_FOLLOWER" - the node should convert to a follower
				// 2. "TIMEOUT" - The node should timeout immediately
				// 3. "SHUTDOWN" - The node should shutdown immediately
				// 4. "LEADER_INFO" - The node should send back its leader's info
				if controllerRequestStructure.Request == "CONVERT_FOLLOWER" {
					fmt.Println("Entered the convert follower block!")
					currentState = Follower
				} else if controllerRequestStructure.Request == "TIMEOUT" {
					fmt.Println("Entered the timeout block, timing out!")
					c <- "restart"
				} else if controllerRequestStructure.Request == "SHUTDOWN" {
					fmt.Println("Shutting down!")
					os.Exit(0)
				} else if controllerRequestStructure.Request == "LEADER_INFO" {
					fmt.Println("Sending leader info!")
					node_information := readNodeInformation()
					voteStructure := &controller_request{Sender_Name: os.Getenv("NODE_NAME"), Request: "LEADER_INFO", Term: node_information.CurrentTerm, Key: "LEADER", Value: node_information.LeaderInfo}
					vote_json_bytes, _ := json.Marshal(voteStructure)
					vote_json_string := string(vote_json_bytes)
					go sendUDP(vote_json_string, controllerRequestStructure.Sender_Name)
				}
			} else {
				if string(p) != "true" {
					go handleVoteRequest(voteRequestStructure)
				}
			}
		}
	}
}

func timeout(inCh chan string, interval time.Duration) {
	exitFlag := false
	for {
		// On each iteration new timer is created
		select {
		case <-time.After(interval):
			votes = 0
			currentState = Candidate
			go startElection()
		case listener_output := <-inCh:
			if listener_output == "heartbeat" {
				currentState = Follower
			}
			if listener_output == "kill" {
				exitFlag = true
			}
			if listener_output == "restart" {
				continue
			}
		}
		if exitFlag {
			break
		}
	}
	fmt.Println("Exited timeout!")
}

func listener(inCh chan string) {
	c := make(chan string)
	go runUDPListener(c)
	for {
		inCh <- <-c
	}
}

func main() {
	currentState = Follower
	nodeInformation := readNodeInformation()
	sleepDuration := nodeInformation.TimeoutInterval + rand.Intn(50)
	var inCh = make(chan string, 1)
	var patsy = make(chan string, 1)
	log.Println("Start")
	go timeout(inCh, time.Millisecond*time.Duration(sleepDuration))
	go listener(inCh)
	<-patsy
}
