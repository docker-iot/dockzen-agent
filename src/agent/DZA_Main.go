package main

import (
	"agent/types/dockzenl"
	"services"
	"webinterface"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

const (
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
					//containersInfo, err := services.DZA_Mon_GetContainersInfo()
					//if err != nil {
					//	respQueue <- Response{req.Num, "Error", nil}
					//} else {
					//	respQueue <- Response{req.Num, req.Cmd, containersInfo}
					//}

				case "UpdateImage":
					updateReturn, err := services.DZA_Update_Do(req.HttpReq)
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

	//reqQueue := make(chan Request, maxQueue)
	//defer close(reqQueue)

	//respQueue = make(chan Response, maxQueue)
	//defer close(respQueue)

	//dispatcher := NewDispatcher(reqQueue)
	//dispatcher.run()

	//dockzenAgentNotifyCh := make(chan string)
	//var dockzenAgentNotiServer net.Listener
	//var err error

	/* DockzenAgentNotify Server */
	//go func() {
//		listenAddress := DockzenNotifySocket
//		dockzenAgentNotiServer, err = net.Listen("unix", listenAddress)
//		if err != nil {
//			fmt.Println("Could not start dockzenAgentNotiServer : ", err)
//			return
//		}

//		fmt.Println("Start Dockzen-Agent-Notify Server")
//		defer dockzenAgentNotiServer.Close()

//		server := clientConns(dockzenAgentNotiServer)
//		for {
//			go handleConn(<-server, dockzenAgentNotifyCh)
//		}

//	}()

	//Test for notify
//	go func() {
//		for {
//			select {
//			case msg1 := <-dockzenAgentNotifyCh:
//				fmt.Println("Message1 :" + msg1)
//			}
//		}
//	}()

	log.Printf("WI init function !!!")
	webinterface.WI_init()

	//listenAddress := services.ContainerServiceSocket
	//router := mux.NewRouter()
//	setupApi(router, reqQueue, respQueue)

//	listener, err := net.Listen("unix", listenAddress)

//	if err != nil {
//		log.Fatalf("Could not listen on %s: %v", listenAddress, err)
//		return
//	}

//	defer listener.Close()
//	sigc := make(chan os.Signal, 1)
//	signal.Notify(sigc,
//		syscall.SIGHUP,
//		syscall.SIGINT,
//		syscall.SIGTERM,
//		syscall.SIGQUIT)

//	go func(listener net.Listener, c chan os.Signal) {
//		sig := <-c
//		listener.Close()
//		dockzenAgentNotiServer.Close()
//		log.Printf("Caught signal %s: shutting down.", sig)
//		os.Exit(0)
//	}(listener, sigc)

//	log.Printf("Starting HTTP server on %s\n", listenAddress)
//	if err = http.Serve(listener, router); err != nil {
//		log.Fatalf("Could not start HTTP server: %v", err)
//	}
}
