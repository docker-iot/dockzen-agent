package main

import (
	"bytes"
	"dockzen-agent/api/types"
	"dockzen-agent/types/dockzenl"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	DockerLauncherSocket string = "/var/run/dockzen_launcher.sock"
	DockzenNotifySocket  string = "/var/run/dockzen_agent_notify.sock"
	maxQueue             int    = 1
)

var respQueue chan Response

type Dispatcher struct {
	workerPool chan chan Request
	req        chan Request
}

type Worker struct {
	req        chan Request
	workerPool chan chan Request
	quitChan   chan bool
}

type Response struct {
	Num  int
	Cmd  string
	Body []byte
}

type Request struct {
	Num     int
	Cmd     string
	HttpReq *http.Request
}

func readData(client net.Conn) ([]byte, error) {

	data := make([]byte, 0)

	for {
		dataBuf := make([]byte, 1024)
		nr, err := client.Read(dataBuf)
		if err != nil {
			break
		}

		log.Printf("nr size [%d]", nr)
		if nr == 0 {
			break
		}

		dataBuf = dataBuf[:nr]
		data = append(data, dataBuf...)
	}

	log.Printf("receive data[%s]\n", string(data))
	//delete null character
	withoutNull := bytes.Trim(data, "\x00")

	rcv := dockzenl.Cmd{}
	err := json.Unmarshal([]byte(withoutNull), &rcv)
	log.Printf("rcv.Cmd = %s", rcv.Cmd)

	if rcv.Cmd == "GetContainersInfo" {
		log.Printf("Success\n")
		return withoutNull, nil
	} else if rcv.Cmd == "UpdateImage" {
		log.Printf("Success\n")
		return withoutNull, nil
	} else {
		log.Printf("error commnad[%s]\n", err)
	}

	return nil, errors.New("Error Cmd from Dockerl")
}

func writeData(client net.Conn, cmd string, m map[string]string) error {
	var send_str []byte
	var err error

	if cmd == "GetContainersInfo" {
		send := dockzenl.Cmd{}
		send.Cmd = "GetContainersInfo"
		send_str, err = json.Marshal(send)
	} else if cmd == "UpdateImage" {

		send := dockzenl.UpdateImageParameters{}
		send.Cmd = "UpdateImage"

		send.Param = dockzenl.UpdateParam{
			ContainerName: m["ContainerName"],
			ImageName:     m["ImageName"],
		}
		send_str, err = json.Marshal(send)

	} else {
		return errors.New("Invalid Cmd")
	}

	log.Printf(string(send_str))
	length := len(send_str)

	message := make([]byte, 0, length)
	message = append(message, send_str...)

	_, err = client.Write([]byte(message))
	if err != nil {
		log.Printf("error: %v\n", err)
		return err
	}

	log.Printf("sent: %s\n", message)
	err = client.(*net.UnixConn).CloseWrite()
	if err != nil {
		log.Printf("error: %v\n", err)
		return err

	}

	return nil
}

func getDockerLauncherInfo_Stub() dockzenl.GetContainersInfoReturn {
	send := dockzenl.GetContainersInfoReturn{
		Containers: []dockzenl.Container{
			{
				ContainerName:   "aaaa",
				ImageName:       "tizen1",
				ContainerStatus: "created",
			},
			{
				ContainerName:   "bbbb",
				ImageName:       "tizen2",
				ContainerStatus: "exited",
			},
		},
	}

	return send
}

func updateImage_Stub() dockzenl.UpdateImageReturn {
	send := dockzenl.UpdateImageReturn{
		State: dockzenl.DeviceState{
			CurrentState: "Updating",
		},
	}

	return send
}

func getContainersInfo() ([]byte, error) {
	log.Printf("GetContainersInfo")
	/*
		stub := getDockerLauncherInfo_Stub()
		var send_stub []byte

		send_stub, _ = json.Marshal(stub)
		log.Printf(string(send_stub))

		return send_stub, nil
	*/
	client, err := net.Dial("unix", DockerLauncherSocket)
	if err != nil {
		log.Fatal("Dial error", err)
		return nil, err
	}

	defer client.Close()

	// Send Command to dockerl
	err = writeData(client, "GetContainersInfo", nil)
	if err != nil {
		return nil, err
	}

	// Receive Command from dockerl
	data := make([]byte, 0)
	data, err = readData(client)
	if err != nil {
		return nil, err
	}

	log.Printf("end\n")
	return data, nil
}

func updateImageRequest(request *http.Request) ([]byte, error) {
	log.Printf("UpdateImageRequest")
	/*
		stub := updateImage_Stub()
		var send_stub []byte

		send_stub, _ = json.Marshal(stub)
		log.Printf(string(send_stub))

		return send_stub, nil
	*/
	client, err := net.Dial("unix", DockerLauncherSocket)
	if err != nil {
		log.Fatal("Dial error", err)
		return nil, err
	}

	defer client.Close()

	// Send Command to dockerl
	ImageName, ContainerName, err := parseUpdateImageParam(request)
	if err != nil {
		return nil, errors.New("Invalid Parameter")
	}
	log.Printf("ImageName : %s\n", ImageName)
	log.Printf("ContainerName : %s\n", ContainerName)
	m := make(map[string]string)
	m["ImageName"] = ImageName
	m["ContainerName"] = ContainerName

	err = writeData(client, "UpdateImage", m)
	if err != nil {
		return nil, err
	}

	// Receive Command from dockerl
	data := make([]byte, 0)
	data, err = readData(client)
	if err != nil {
		return nil, err
	}

	log.Printf("end\n")
	return data, nil
}

func parseUpdateImageParam(request *http.Request) (ImageName, ContainerName string, err error) {

	var body types.UpdateImageParams

	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&body)

	log.Printf("body.ImageName = %s\n", body.ImageName)
	log.Printf("body.ContainerName = %s\n", body.ContainerName)

	ImageName = body.ImageName
	ContainerName = body.ContainerName

	return ImageName, ContainerName, err
}

func apiGetHandler(w http.ResponseWriter, r *http.Request, reqs chan Request, resps chan Response) {
	vars := mux.Vars(r)
	Cmd := vars["Cmd"]
	log.Printf("Cmd: [%s]", Cmd)

	// num is always 1, because, request will be handled the earier one is finished
	req := Request{Cmd: Cmd, Num: 1, HttpReq: r}
	reqs <- req

	currentReqNum := 1
	var respData Response
	for {
		respData = <-respQueue
		if currentReqNum == respData.Num {
			fmt.Printf("done: [%d]\n", currentReqNum)
			break
		}
	}

	// Make resps
	w.Header().Set("Content-Type", "application/json")
	if respData.Cmd == "GetContainersInfo" {
		w.WriteHeader(http.StatusOK)
		w.Write(respData.Body)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println("Complete RequestHandler")
	return
}

func apiPostHandler(w http.ResponseWriter, r *http.Request, reqs chan Request, resps chan Response) {
	vars := mux.Vars(r)
	Cmd := vars["Cmd"]
	log.Printf("Cmd: [%s]", Cmd)

	// num is always 1, because, request will be handled the earier one is finished
	req := Request{Cmd: Cmd, Num: 1, HttpReq: r}

	reqs <- req

	currentReqNum := 1
	var respData Response
	for {
		respData = <-respQueue
		if currentReqNum == respData.Num {
			fmt.Printf("done: [%d]\n", currentReqNum)
			break
		}
	}

	// Make resps
	w.Header().Set("Content-Type", "application/json")
	if respData.Cmd == "UpdateImage" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Println("Complete RequestHandler")
	return
}

func NewWorker(workerPool chan chan Request) Worker {
	return Worker{
		req:        make(chan Request),
		workerPool: workerPool,
		quitChan:   make(chan bool),
	}
}

func (w Worker) start() {
	go func() {
		for {
			w.workerPool <- w.req

			select {
			case req := <-w.req:
				switch req.Cmd {
				case "GetContainersInfo":
					containersInfo, err := getContainersInfo()
					if err != nil {
						respQueue <- Response{req.Num, "Error", nil}
					} else {
						respQueue <- Response{req.Num, req.Cmd, containersInfo}
					}

				case "UpdateImage":
					updateReturn, err := updateImageRequest(req.HttpReq)
					if err != nil {
						log.Printf("Error [%s]", err)
						respQueue <- Response{req.Num, "Error", nil}
					} else {
						respQueue <- Response{req.Num, req.Cmd, updateReturn}
					}

				}

				log.Printf("Completed\n")
			case <-w.quitChan:
				log.Printf("worker stopping\n")
				return
			}
		}
	}()
}

func (w Worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}

func NewDispatcher(req chan Request) *Dispatcher {
	workerpool := make(chan chan Request, 1)

	return &Dispatcher{
		req:        req,
		workerPool: workerpool,
	}
}

func (d *Dispatcher) run() {
	woker := NewWorker(d.workerPool)
	woker.start()

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case req := <-d.req:
			go func() {
				log.Printf("fetching workerRequest for : %d\n", req.Num)
				workerRequest := <-d.workerPool
				log.Printf("adding [%d] to workerRequest\n", req.Num)
				workerRequest <- req
			}()
		}
	}
}

func setupApi(r *mux.Router, req chan Request, resp chan Response) {

	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/get/{Cmd}", func(w http.ResponseWriter, r *http.Request) {
		apiGetHandler(w, r, req, resp)
	}).Methods("GET")
	s.HandleFunc("/post/{Cmd}", func(w http.ResponseWriter, r *http.Request) {
		apiPostHandler(w, r, req, resp)
	}).Methods("POST")
}

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

	notifier <- string(data)
	fmt.Println("end handleConn")
}

func main() {
	log.Printf("Container-Service Agent starting")

	reqQueue := make(chan Request, maxQueue)
	defer close(reqQueue)

	respQueue = make(chan Response, maxQueue)
	defer close(respQueue)

	dispatcher := NewDispatcher(reqQueue)
	dispatcher.run()

	dockzenAgentNotifyCh := make(chan string)
	var dockzenAgentNotiServer net.Listener
	var err error

	/* DockzenAgentNotify Server */
	go func() {
		listenAddress := DockzenNotifySocket
		dockzenAgentNotiServer, err = net.Listen("unix", listenAddress)
		if err != nil {
			fmt.Println("Could not start dockzenAgentNotiServer : ", err)
			return
		}

		fmt.Println("Start Dockzen-Agent-Notify Server")
		defer dockzenAgentNotiServer.Close()

		server := clientConns(dockzenAgentNotiServer)
		for {
			go handleConn(<-server, dockzenAgentNotifyCh)
		}

	}()

	//Test for notify
	go func() {
		for {
			select {
			case msg1 := <-dockzenAgentNotifyCh:
				fmt.Println("Message1 :" + msg1)
			}
		}
	}()

	listenAddress := types.ContainerServiceSocket
	router := mux.NewRouter()
	setupApi(router, reqQueue, respQueue)

	listener, err := net.Listen("unix", listenAddress)

	if err != nil {
		log.Fatalf("Could not listen on %s: %v", listenAddress, err)
		return
	}

	defer listener.Close()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func(listener net.Listener, c chan os.Signal) {
		sig := <-c
		listener.Close()
		dockzenAgentNotiServer.Close()
		log.Printf("Caught signal %s: shutting down.", sig)
		os.Exit(0)
	}(listener, sigc)

	log.Printf("Starting HTTP server on %s\n", listenAddress)
	if err = http.Serve(listener, router); err != nil {
		log.Fatalf("Could not start HTTP server: %v", err)
	}
}
