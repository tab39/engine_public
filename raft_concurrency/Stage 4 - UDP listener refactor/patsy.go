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
			// TODO: Refactor handleHeartbeat
			go handleHeartbeat(requestFormatStructure.Data)
			c <- "heartbeat"
		case "vote":
			if requestFormatStructure.Data == "true" && currentState == Candidate {
				fmt.Println("Recieved a vote!")
				votes += 1
				if votes > 1 {
					currentState = Leader
					votes = 0
					c <- "kill"
					go sendHeartbeats()
				}
			}
		case "vote_request":
			//TODO: Refactor handleVoteRequest
			go handleVoteRequest(requestFormatStructure.Data)
		case "Controller":
			controllerRequestStructure := controller_request{}
			_ = json.Unmarshal(p, &controllerRequestStructure)
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
				go sendUDP(ctrl_ret_json_str, controllerRequestStructure.Sender_Name)
			case "STORE":
				nodeInformation := readNodeInformation()
				if nodeInformation.LeaderInfo == os.Getenv("NODE_NAME") {
					fmt.Println("Storing data!")
					go writeToLog(nodeInformation.CurrentTerm, controllerRequestStructure.Key, controllerRequestStructure.Value)
				} else {
					fmt.Println("Returning leader data!")
					controllerReturnStructure := &controller_request{Sender_Name: os.Getenv("NODE_NAME"), Request: "LEADER_INFO", Term: nodeInformation.CurrentTerm, Key: "LEADER", Value: nodeInformation.LeaderInfo}
					ctrl_ret_json_bytes, _ := json.Marshal(controllerReturnStructure)
					ctrl_ret_json_str := string(ctrl_ret_json_bytes)
					fmt.Println(ctrl_ret_json_str)
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

// Copy of original:

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
			go handleHeartbeat(heartbeatStructure)
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
					go sendUDP(ctrl_ret_json_str, controllerRequestStructure.Sender_Name)
				case "STORE":
					nodeInformation := readNodeInformation()
					if nodeInformation.LeaderInfo == os.Getenv("NODE_NAME") {
						fmt.Println("Storing data!")
						go writeToLog(nodeInformation.CurrentTerm, controllerRequestStructure.Key, controllerRequestStructure.Value)
					} else {
						fmt.Println("Returning leader data!")
						controllerReturnStructure := &controller_request{Sender_Name: os.Getenv("NODE_NAME"), Request: "LEADER_INFO", Term: nodeInformation.CurrentTerm, Key: "LEADER", Value: nodeInformation.LeaderInfo}
						ctrl_ret_json_bytes, _ := json.Marshal(controllerReturnStructure)
						ctrl_ret_json_str := string(ctrl_ret_json_bytes)
						fmt.Println(ctrl_ret_json_str)
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
			} else {
				if string(p) != "true" {
					go handleVoteRequest(voteRequestStructure)
				}
			}
		}
	}
}