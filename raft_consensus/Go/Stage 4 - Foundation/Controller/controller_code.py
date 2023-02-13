from concurrent.futures import thread
import json
from nis import match
import socket
import traceback
import time
import threading

# Initialize
sender = "Controller"
target = "node1"
port = 8080
sentPayloadList = []
requestList = []
leaderInfo = "node1"

def generate_payload(request,key=None,value=None):
    # Read Message Template
    requestList.append(request)
    payload = json.load(open("Message.json"))
    payload['sender_name'] = sender
    payload['request'] = request
    payload['key'] = key
    payload['value'] = value
    print(f"Request Created : {payload}")
    sentPayloadList.append(payload)
    return payload

def listenerProcess():
    global leaderInfo
    print("Entered listener process!")
    skt = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
    skt.bind((sender, port))
    # UDPListenerSocket = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
    # UDPListenerSocket.bind(("controller",port))
    while(True):
        bytesAddressPair = skt.recvfrom(1024)
        message = bytesAddressPair[0].decode('utf-8')
        jsonMessage = json.loads(message)
        print(message)
        if(jsonMessage['request']=="LEADER_INFO"):
            leaderInfo = jsonMessage['value']
            print("Updated leader name to ",leaderInfo)
            if(requestList[-1]=="STORE" or requestList[-1]=="RETRIEVE"):
                payload = sentPayloadList[-1]
                sendUDP(jsonMessage['value'],payload)
            else:
                print("Received leader info!\n",json.dumps(jsonMessage))
        elif(jsonMessage['request']=="RETRIEVE"):
            print("Received logs!\n",jsonMessage['value'])
        
def sendUDP(node,payload):
    skt = socket.socket(family=socket.AF_INET, type=socket.SOCK_DGRAM)
    skt.bind((sender, 8081))
    # Send Message
    try:
        # Encoding and sending the message
        skt.sendto(json.dumps(payload).encode('utf-8'), (node, port))
    except:
        #  socket.gaierror: [Errno -3] would be thrown if target IP container does not exist or exits, write your listener
        print(f"ERROR WHILE SENDING REQUEST ACROSS : {traceback.format_exc()}")
    skt.close()

def main():
    # Wait following seconds below sending the controller request
    listenerThread = threading.Thread(target=listenerProcess)
    listenerThread.start()
    time.sleep(5)
    payload = generate_payload("STORE","1","Mihir")
    sendUDP("node1",payload)
    time.sleep(5)
    payload = generate_payload("STORE","2","Tarun")
    sendUDP(leaderInfo,payload)
    time.sleep(10)
    payload = generate_payload("SHUTDOWN")
    sendUDP(leaderInfo,payload)
    time.sleep(20)
    deadNode = int(leaderInfo.split("node")[1])
    if(deadNode >= 1 and deadNode<5):
        deadNode+=1
    else:
        deadNode-=1
    nodeData = "node"+str(deadNode)
    payload = generate_payload("STORE","3","Arctic Monkey")
    sendUDP(nodeData,payload)
    time.sleep(10)
    payload = generate_payload("RETRIEVE")
    sendUDP(nodeData,payload)
    time.sleep(10)

if(__name__=="__main__"):
    main()