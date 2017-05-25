package main

import (
	"bytes"
	"dockzen-agent/types/dockzenl"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Request struct {
	Num  int
	Resp chan Response
}

type Response struct {
	Num      int
	WorkerID int
}

var clientWS *websocket.Conn
var clientNotifyWS *websocket.Conn

var dockzenLauncherClient net.Conn
var recvDockzenAgentCh chan string
var sendDockzenNotifyCh chan string

const (
	DockzenLauncherSocket  string = "dockzen_launcher.sock"
	DockzenNotifySocket    string = "dockzen_agent_notify.sock"
	DefaultSocketPath      string = "/var/run/"
	DockzenAgentPort       string = "8080"
	DockzenAgentNotifyPort string = "8082"
	NumOfServer            int    = 4
)

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Println(err)
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func handleConn(client net.Conn, notifier chan string) {
	data := make([]byte, 0)
	for {
		dataBuf := make([]byte, 1024)
		nr, err := client.Read(dataBuf)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("nr size [%d]\n", nr)
		if nr == 0 {
			break
		}

		dataBuf = dataBuf[:nr]
		data = append(data, dataBuf...)
	}

	fmt.Println(string(data))
	notifier <- string(data)
	fmt.Println("end handleConn")
}

func dockzenAgentWSHandler(ws *websocket.Conn) {
	fmt.Println("dockzen-agent handler")

	clientWS = ws

	for {
		msg := make([]byte, 1024)

		n, err := ws.Read(msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("size : %d\n", n)
		recvDockzenAgentCh <- string(msg)
	}
}

func dockzenNotifyWSHandler(ws *websocket.Conn) {
	fmt.Println("dockzen-agent-notify handler")

	clientNotifyWS = ws

	for {
		msg := make([]byte, 1024)

		n, err := ws.Read(msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("size : %d\n", n)
		sendDockzenNotifyCh <- string(msg)
	}
}

func receiver(recvFromLauncherAPICh, recvFromLauncherNotifyCh, recvFromAgentCh, sendToAgentCh, sendToLauncherCh chan string) {

	for {
		select {
		case msg1 := <-recvFromLauncherAPICh:
			fmt.Println("recvFromLauncherAPICh : " + msg1)

		case msg2 := <-recvFromLauncherNotifyCh:
			fmt.Println("recvFromLauncherNotifyCh : " + msg2)

		case msg3 := <-recvFromAgentCh:
			fmt.Println("recvFromAgentCh : " + msg3)
			trimmedMsg := bytes.Trim([]byte(msg3), "\x00")

			rcv := dockzenl.Cmd{}
			json.Unmarshal([]byte(trimmedMsg), &rcv)
			fmt.Println(rcv.Cmd)

			switch rcv.Cmd {
			case "GetContainersInfo":
				fmt.Println("Call GetContainersInfo")
				containersInfo, err := getContainersInfo(sendToLauncherCh)
				if err != nil {
					fmt.Println(err)
				}

				sendToAgentCh <- string(containersInfo)
			case "UpdateImage":
				fmt.Println("Call UpdateImage")

				updateReturn, err := updateImage(sendToLauncherCh, trimmedMsg)
				if err != nil {
					fmt.Println(err)
				}

				sendToAgentCh <- string(updateReturn)
			default:
				fmt.Println("Not Support the Command")
			}
		}
	}
}

func sender(sendToAgentCh, sendToLauncherCh, sendtToAgentNotifyCh chan string) {

	for {
		select {
		case msg1 := <-sendToAgentCh:
			fmt.Println("sendToAgentCh")
			n, err := clientWS.Write([]byte(msg1))
			if err != nil {
				fmt.Printf("error : %s", err)
			} else {
				fmt.Printf("Sender[%d]: %s\n", n, msg1)
			}

		case msg2 := <-sendToLauncherCh:
			fmt.Println("sendToLauncherCh")
			n, err := dockzenLauncherClient.Write([]byte(msg2))
			if err != nil {
				fmt.Printf("error : %s", err)
			} else {
				fmt.Printf("Sender[%d]: %s\n", n, msg2)
			}

		case msg3 := <-sendtToAgentNotifyCh:
			fmt.Println("sendtToAgentNotifyCh")
			n, err := clientNotifyWS.Write([]byte(msg3))
			if err != nil {
				fmt.Printf("error : %s", err)
			} else {
				fmt.Printf("Sent[%d]: %s\n", n, msg3)
			}
		}
	}
}

func main() {
	fmt.Println("Start Event Initialize...")

	var dockzenAgentNotiServer net.Listener
	var err error

	/* Create Receive channel
	dockzenAgentApiCh is from dockzen-launcher for api responses
	dockzenAgentNotifyCh is from dockzen-launcher for notify
	dockzenAgentCh is from dockzen-agentconn for cmd
	*/
	recvDockzenAgentApiCh := make(chan string)
	recvDockzenAgentNotifyCh := make(chan string)
	recvDockzenAgentCh = make(chan string)

	sendDockzenAgentCh := make(chan string)
	sendDockzenLauncherCh := make(chan string)
	sendDockzenNotifyCh = make(chan string)

	var wg sync.WaitGroup
	wg.Add(NumOfServer)
	/* DockzenAgent Server */

	go func() {
		fmt.Println("Start Dockzen-Agent Server")
		http.Handle("/wsAgent", websocket.Handler(dockzenAgentWSHandler))
		wg.Done()
		err = http.ListenAndServe(":"+DockzenAgentPort, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	go func() {
		fmt.Println("Start Dockzen-Agent-Notify Server")
		http.Handle("/wsNotify", websocket.Handler(dockzenNotifyWSHandler))
		wg.Done()
		err = http.ListenAndServe(":"+DockzenAgentNotifyPort, nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()

	/* DockzenAgentNotify Server */
	go func() {
		listenAddress := DefaultSocketPath + DockzenNotifySocket
		dockzenAgentNotiServer, err = net.Listen("unix", listenAddress)
		if err != nil {
			fmt.Println("Could not start dockzenAgentNotiServer : ", err)
			return
		}

		defer dockzenAgentNotiServer.Close()
		wg.Done()

		server := clientConns(dockzenAgentNotiServer)
		for {
			go handleConn(<-server, recvDockzenAgentNotifyCh)
		}

	}()

	dockzenLauncherClient, err = net.Dial("unix", DefaultSocketPath+DockzenLauncherSocket)
	if err != nil {
		fmt.Println(err)
		return
	}

	wg.Done()

	defer dockzenLauncherClient.Close()

	fmt.Println("Setup exception case")
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func(client net.Conn, linstener net.Listener, c chan os.Signal) {
		sig := <-c

		close(recvDockzenAgentApiCh)
		close(recvDockzenAgentNotifyCh)
		close(recvDockzenAgentCh)

		close(sendDockzenAgentCh)
		close(sendDockzenLauncherCh)
		close(sendDockzenNotifyCh)

		dockzenAgentNotiServer.Close()
		dockzenLauncherClient.Close()

		fmt.Printf("Caught signal %s: shutting down.", sig)
		os.Exit(0)
	}(dockzenLauncherClient, dockzenAgentNotiServer, sigc)

	wg.Wait()

	fmt.Println("Initialized")

	go receiver(recvDockzenAgentApiCh, recvDockzenAgentNotifyCh, recvDockzenAgentCh, sendDockzenAgentCh, sendDockzenLauncherCh)

	go sender(sendDockzenAgentCh, sendDockzenLauncherCh, sendDockzenNotifyCh)

	for {
		time.Sleep(time.Second)
	}

}
