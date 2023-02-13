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
var nextSendingIndex [4]int

type log_data struct {
	Term  int    `json:"term"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type node_information struct {
	CurrentTerm       int        `json:"currentTerm"`
	VotedFor          string     `json:"votedFor"`
	Log               []log_data `json:"log"`
	TimeoutInterval   int        `json:"timeoutInterval"`
	HeartbeatInterval int        `json:"heartbeatInterval"`
	LeaderInfo        string     `json:"leaderInfo"`
	CommitIndex       int        `json:"commitIndex"`
}

type append_entry struct {
	LeaderID     string     `json:"leaderID"`
	Term         int        `json:"term"`
	PrevLogIndex int        `json:"prevLogIndex"`
	PrevLogTerm  int        `json:"prevLogTerm"`
	Entries      []log_data `json:"entries"`
	LeaderCommit int        `json:"leaderCommit"`
}

type append_reply struct {
	Sender   string `json:"sender"`
	Response string `json:"response"`
}

type request_vote struct {
	CandidateID  string `json:"candidateID"`
	Term         int    `json:"term"`
	LastLogIndex int    `json:"lastLogIndex"`
	LastLogTerm  int    `json:"lastLogTerm"`
}

type vote_response struct {
	Response string `json:"response"`
	Term     int    `json:"term"`
}

type request_format struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type controller_request struct {
	Sender_Name string `json:"sender_name"`
	Request     string `json:"request"`
	Term        int    `json:"term"`
	Key         string `json:"key"`
	Value       string `json:"value"`
}

type controller_request_retrieve struct {
	Sender_Name string     `json:"sender_name"`
	Request     string     `json:"request"`
	Term        int        `json:"term"`
	Key         string     `json:"key"`
	Value       []log_data `json:"value"`
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

func writeVoteResponse(response string, term int) string {
	voteResponseStructure := &vote_response{Response: response, Term: term}
	jsonVoteResponse, _ := json.Marshal(voteResponseStructure)
	return string(jsonVoteResponse)
}

func handleVoteRequest(requestVote string) {
	requestVoteStructure := &request_vote{}
	json.Unmarshal([]byte(requestVote), requestVoteStructure)
	nodeInformation := readNodeInformation()
	if nodeInformation.LeaderInfo != os.Getenv("NODE_NAME") {
		currentLastLogIndex := len(nodeInformation.Log) - 1
		var currentLastLogTerm int
		if currentLastLogIndex == -1 {
			currentLastLogTerm = 0
		} else {
			currentLastLogTerm = nodeInformation.Log[currentLastLogIndex].Term
		}
		if requestVoteStructure.Term >= nodeInformation.CurrentTerm && nodeInformation.VotedFor == "" {
			nodeInformation.CurrentTerm = requestVoteStructure.Term
			if currentLastLogTerm > requestVoteStructure.LastLogTerm || (currentLastLogTerm == requestVoteStructure.LastLogTerm && currentLastLogIndex > requestVoteStructure.LastLogIndex) {
				// fmt.Println("Vote rejected!")
			} else {
				nodeInformation.VotedFor = requestVoteStructure.CandidateID
				responseString := writeVoteResponse("true", requestVoteStructure.Term)
				requestFormatStructure := &request_format{Type: "vote", Data: responseString}
				request_json_bytes, _ := json.Marshal(requestFormatStructure)
				request_json_string := string(request_json_bytes)
				go sendUDP(request_json_string, requestVoteStructure.CandidateID)
				log.Println("Voted :", request_json_string, "\tFor: ", requestVoteStructure.CandidateID)
				writeNodeInformation(nodeInformation)
			}
		}
	}
}

func getReplyString(response string) string {
	appendReply := &append_reply{Sender: os.Getenv("NODE_NAME"), Response: response}
	appendReply_json_bytes, _ := json.Marshal(appendReply)
	appendReply_json_string := string(appendReply_json_bytes)
	requestReply := &request_format{Type: "heartbeat reply", Data: appendReply_json_string}
	requestReply_json_bytes, _ := json.Marshal(requestReply)
	requestReply_json_string := string(requestReply_json_bytes)

	return requestReply_json_string
}

// if node's term is less than heartbreat, persist
func handleHeartbeat(heartbeat string) {
	//fmt.Println(heartbeat)
	heartbeatStructure := &append_entry{}
	json.Unmarshal([]byte(heartbeat), heartbeatStructure)
	nodeInformation := readNodeInformation()
	nodeInformation.LeaderInfo = heartbeatStructure.LeaderID
	nodeInformation.VotedFor = ""

	if len(heartbeatStructure.Entries) > 0 {
		// fmt.Println("Entries received!")
		if heartbeatStructure.Term < nodeInformation.CurrentTerm {
			fmt.Println("Heartbeat rejected, term was lower than current term of this node!")
			requestReply_json_string := getReplyString("false")
			go sendUDP(requestReply_json_string, heartbeatStructure.LeaderID)
		} else {
			nodeInformation.CurrentTerm = heartbeatStructure.Term
			if len(nodeInformation.Log) == 0 {
				// fmt.Println("Heartbeat accepted, first entry to log.")
				nodeInformation.Log = append(nodeInformation.Log, heartbeatStructure.Entries[0])
				requestReply_json_string := getReplyString("true")
				go sendUDP(requestReply_json_string, heartbeatStructure.LeaderID)
			} else {
				if nodeInformation.Log[heartbeatStructure.PrevLogIndex].Term != heartbeatStructure.PrevLogTerm {
					fmt.Println("Heartbeat rejected, term mismatch between follower and leader log!")
					requestReply_json_string := getReplyString("false")
					go sendUDP(requestReply_json_string, heartbeatStructure.LeaderID)
				} else {
					// fmt.Println("Log entries accepted, appending.")
					nodeInformation.Log = append(nodeInformation.Log, heartbeatStructure.Entries[0])
					requestReply_json_string := getReplyString("true")
					go sendUDP(requestReply_json_string, heartbeatStructure.LeaderID)
				}
			}
		}
	}
	if heartbeatStructure.LeaderCommit > nodeInformation.CommitIndex {
		if heartbeatStructure.LeaderCommit < heartbeatStructure.PrevLogIndex+1 {
			nodeInformation.CommitIndex = heartbeatStructure.LeaderCommit
		} else {
			nodeInformation.CommitIndex = heartbeatStructure.PrevLogIndex + 1
		}
		// fmt.Println("Commit index updated to", nodeInformation.CommitIndex)
	}
	writeNodeInformation(nodeInformation)

}

func startElection() {
	log.Println("Entered startElection!")
	nodeInformation := readNodeInformation()
	lastLogIndex := len(nodeInformation.Log) - 1
	var lastLogTerm int
	if lastLogIndex == -1 {
		lastLogTerm = 0
	} else {
		lastLogTerm = nodeInformation.Log[lastLogIndex].Term
	}
	voteStructure := &request_vote{CandidateID: os.Getenv("NODE_NAME"), Term: nodeInformation.CurrentTerm, LastLogIndex: lastLogIndex, LastLogTerm: lastLogTerm}
	vote_json_bytes, _ := json.Marshal(voteStructure)
	vote_json_string := string(vote_json_bytes)
	requestFormatStructure := &request_format{Type: "vote_request", Data: vote_json_string}
	request_json_bytes, _ := json.Marshal(requestFormatStructure)
	request_json_string := string(request_json_bytes)
	for i := 0; i < 4; i++ {
		go sendUDP(request_json_string, os.Getenv("COMPANION_"+strconv.Itoa(i)))
	}
}

func sendHeartbeats() {
	log.Println(os.Getenv("NODE_NAME") + " became a leader!")
	nodeInformation := readNodeInformation()
	nodeInformation.LeaderInfo = os.Getenv("NODE_NAME")
	writeNodeInformation(nodeInformation)
	for i := 0; i < len(nextSendingIndex); i++ {
		nextSendingIndex[i] = nodeInformation.CommitIndex + 1
	}
	go func() {
		for {
			nodeInformation := readNodeInformation()
			heartbeatStructure := &append_entry{LeaderID: os.Getenv("NODE_NAME"), Term: nodeInformation.CurrentTerm, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []log_data{}, LeaderCommit: nodeInformation.CommitIndex}
			for i := 0; i < 4; i++ {
				// if nextSendingIndex[i]-1 > len(nodeInformation.Log) {
				// 	nextSendingIndex[i] = len(nodeInformation.Log)
				// }
				previousLogIndex := nextSendingIndex[i] - 1
				var previousLogTerm int
				if previousLogIndex < 0 {
					previousLogTerm = 0
				} else {
					if previousLogIndex >= len(nodeInformation.Log) {
						previousLogIndex = len(nodeInformation.Log) - 1
					}
					previousLogTerm = nodeInformation.Log[previousLogIndex].Term
				}
				heartbeatStructure.PrevLogIndex = previousLogIndex
				heartbeatStructure.PrevLogTerm = previousLogTerm
				if len(nodeInformation.Log) > nextSendingIndex[i] {
					heartbeatStructure.Entries = []log_data{nodeInformation.Log[nextSendingIndex[i]]}
				}
				heartbeat_json_bytes, _ := json.Marshal(heartbeatStructure)
				heartbeat_json_string := string(heartbeat_json_bytes)
				requestFormatStructure := &request_format{Type: "heartbeat", Data: heartbeat_json_string}
				request_json_bytes, _ := json.Marshal(requestFormatStructure)
				request_json_string := string(request_json_bytes)
				go sendUDP(request_json_string, os.Getenv("COMPANION_"+strconv.Itoa(i)))
			}
			time.Sleep(time.Duration(nodeInformation.HeartbeatInterval) * time.Millisecond)
		}
	}()
}

func handleHeartbeatReply(heartbeatReply string) {
	heartbeatReplyStructure := &append_reply{}
	json.Unmarshal([]byte(heartbeatReply), heartbeatReplyStructure)
	nodeInformation := readNodeInformation()

	//! Figure out who sent the response, according to our indexing system
	senderID := heartbeatReplyStructure.Sender
	index := -1
	for i := 0; i < 4; i++ {
		if os.Getenv("COMPANION_"+strconv.Itoa(i)) == senderID {
			index = i
			break
		}
	}
	//! See what they sent. If it was true, update our nextSendingIndex, and increase the vote count for
	//! the number of accepted responses. Deal with the second part off this as a TODO.

	if heartbeatReplyStructure.Response == "true" {
		// fmt.Println("Received a true response from " + senderID)
		nextSendingIndex[index] += 1
	} else {
		// fmt.Println("Received a false response from " + senderID)
		nextSendingIndex[index] -= 1
	}
	uptoDateNodeCount := 0
	for i := 0; i < 4; i++ {
		if nextSendingIndex[i] == len(nodeInformation.Log) {
			uptoDateNodeCount++
		}
		if uptoDateNodeCount > 1 && nodeInformation.CommitIndex < len(nodeInformation.Log)-1 {
			nodeInformation.CommitIndex += 1
			fmt.Println("Leader committed log entry!")
		}
	}
	writeNodeInformation(nodeInformation)
}

func writeToLog(term int, key string, value string) {
	nodeInformation := readNodeInformation()
	logData := log_data{Term: term, Key: key, Value: value}
	nodeInformation.Log = append(nodeInformation.Log, logData)
	writeNodeInformation(nodeInformation)
}

func runUDPListener(c chan string) {
	address, _ := net.ResolveUDPAddr("udp", os.Getenv("NODE_NAME")+":8080")
	connection, _ := net.ListenUDP("udp", address)
	votes = 0
	for {
		p := make([]byte, 2048)
		_, _, _ = connection.ReadFromUDP(p)
		p = bytes.Trim(p, "\x00")
		requestFormatStructure := request_format{}
		_ = json.Unmarshal(p, &requestFormatStructure)
		switch requestFormatStructure.Type {
		case "heartbeat":
			go handleHeartbeat(requestFormatStructure.Data)
			c <- "heartbeat"
		case "vote":
			voteResponseStructure := &vote_response{}
			json.Unmarshal([]byte(requestFormatStructure.Data), voteResponseStructure)

			if voteResponseStructure.Response == "true" && voteResponseStructure.Term == readNodeInformation().CurrentTerm && currentState == Candidate {
				// fmt.Println("Recieved a vote!")
				votes += 1
				if votes > 2 {
					currentState = Leader
					votes = 0
					go sendHeartbeats()
					c <- "kill"
				}
			}
		case "vote_request":
			go handleVoteRequest(requestFormatStructure.Data)
		case "heartbeat reply":
			if currentState == Leader {
				handleHeartbeatReply(requestFormatStructure.Data)
			}
		default:
			controllerRequestStructure := controller_request{}
			err := json.Unmarshal(p, &controllerRequestStructure)
			if err != nil {
				fmt.Println(err)
			}
			// The controller can send one of four request types to the node:
			// 1. "CONVERT_FOLLOWER" - the node should convert to a follower
			// 2. "TIMEOUT" - The node should timeout immediately
			// 3. "SHUTDOWN" - The node should shutdown immediately
			// 4. "LEADER_INFO" - The node should send back its leader's info
			switch controllerRequestStructure.Request {
			case "CONVERT_FOLLOWER":
				fmt.Println("Entered the convert follower block!")
				currentState = Follower
			case "TIMEOUT":
				fmt.Println("Entered the timeout block, timing out!")
				c <- "restart"
			case "SHUTDOWN":
				fmt.Println("Shutting down!")
				os.Exit(0)
			case "LEADER_INFO":
				fmt.Println("Sending leader info!")
				nodeInformation := readNodeInformation()
				controllerReturnStructure := &controller_request{Sender_Name: os.Getenv("NODE_NAME"), Request: "LEADER_INFO", Term: nodeInformation.CurrentTerm, Key: "LEADER", Value: nodeInformation.LeaderInfo}
				ctrl_ret_json_bytes, _ := json.Marshal(controllerReturnStructure)
				ctrl_ret_json_str := string(ctrl_ret_json_bytes)
				// fmt.Println(ctrl_ret_json_str)
				go sendUDP(ctrl_ret_json_str, controllerRequestStructure.Sender_Name)
			case "STORE":
				nodeInformation := readNodeInformation()
				if nodeInformation.LeaderInfo == os.Getenv("NODE_NAME") {
					fmt.Println("Storing data!")
					writeToLog(nodeInformation.CurrentTerm, controllerRequestStructure.Key, controllerRequestStructure.Value)
				} else {
					fmt.Println("Returning leader data!")
					controllerReturnStructure := &controller_request{Sender_Name: os.Getenv("NODE_NAME"), Request: "LEADER_INFO", Term: nodeInformation.CurrentTerm, Key: "LEADER", Value: nodeInformation.LeaderInfo}
					ctrl_ret_json_bytes, _ := json.Marshal(controllerReturnStructure)
					ctrl_ret_json_str := string(ctrl_ret_json_bytes)
					// fmt.Println(ctrl_ret_json_str)
					go sendUDP(ctrl_ret_json_str, controllerRequestStructure.Sender_Name)
				}
			case "RETRIEVE":
				nodeInformation := readNodeInformation()
				if nodeInformation.LeaderInfo == os.Getenv("NODE_NAME") {
					controllerReturnStructure := &controller_request_retrieve{Sender_Name: os.Getenv("NODE_NAME"), Request: "RETRIEVE", Term: nodeInformation.CurrentTerm, Key: "COMMITTED_LOGS", Value: nodeInformation.Log}
					ctrl_ret_json_bytes, _ := json.Marshal(controllerReturnStructure)
					ctrl_ret_json_str := string(ctrl_ret_json_bytes)
					go sendUDP(ctrl_ret_json_str, controllerRequestStructure.Sender_Name)
				} else {
					fmt.Println("Returning leader data!")
					controllerReturnStructure := &controller_request{Sender_Name: os.Getenv("NODE_NAME"), Request: "LEADER_INFO", Term: nodeInformation.CurrentTerm, Key: "LEADER", Value: nodeInformation.LeaderInfo}
					ctrl_ret_json_bytes, _ := json.Marshal(controllerReturnStructure)
					ctrl_ret_json_str := string(ctrl_ret_json_bytes)
					go sendUDP(ctrl_ret_json_str, controllerRequestStructure.Sender_Name)
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
			nodeInformation := readNodeInformation()
			nodeInformation.VotedFor = ""
			nodeInformation.CurrentTerm += 1
			writeNodeInformation(nodeInformation)
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
	log.Println("Exited timeout!")
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
	sleepDuration := nodeInformation.TimeoutInterval + rand.Intn(100)
	var inCh = make(chan string, 1)
	var patsy = make(chan string, 1)
	log.Println("Start")
	go timeout(inCh, time.Millisecond*time.Duration(sleepDuration))
	go listener(inCh)
	<-patsy
}
