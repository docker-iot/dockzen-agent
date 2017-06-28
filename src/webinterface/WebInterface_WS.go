package webinterface

import (
	"log"
	"golang.org/x/net/websocket"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"encoding/json"
	"syscall"
	"time"
	set "services"
	dockzen_h "include"
)

var wss_prefix = "ws://"

type Command struct {
	Cmd string `json:"cmd"`
}

var chSignal chan os.Signal
var done chan bool

func WI_init(){

	log.Printf("[%s] Web connection start !!!\n", __FILE__)

	for {

		go ws_mainLoop()

		<-done
		time.Sleep(time.Second)
	}

}

func ws_mainLoop() (err error) {

	go func() {
		<-chSignal
		done <- true
		return
	}()

	var server_url = set.GetServerURL("")
	if server_url == "" {
		log.Printf("[%s] Server URL error !!!", __FILE__ )
		return nil
	}

	var wss_server_url = wss_prefix + server_url
	ws, err := wsProxyDial(wss_server_url, "tcp", wss_server_url)
	if err != nil {
		log.Printf("[%s] wsProxyDial : ",__FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return err
	}

	defer ws.Close()

	/* connect test2 : message driven
	 */
	messages := make(chan string)
	go wsReceive(wss_server_url, ws, messages)

	name, _ := GetHardwareAddress()

	err = wsReqeustConnection(ws, name)

	for {
		msg := <-messages

		rcv := Command{}
		json.Unmarshal([]byte(msg), &rcv)
		log.Printf(rcv.Cmd)

		switch rcv.Cmd {
		case "connected":
			log.Printf("[%s] connected succefully~~", __FILE__)
		case "GetContainersInfo":
			send_info, ret := WS_GetContainerLists()
			if ret == dockzen_h.DOCKZEN_ERROR_NONE {
				websocket.JSON.Send(ws, send_info)
			}
		case "UpdateImage":
			log.Printf("[%s] command <UpdateImage>", __FILE__)
			update_msg, r := ParseUpdateParam(msg)
			if r == nil {
				send_update, ret := WS_UpdateImage(update_msg)
				if ret == dockzen_h.DOCKZEN_ERROR_NONE {
					websocket.JSON.Send(ws, send_update)
				}
			} else {
				log.Printf("[%s] UpdateImage message null !!!")
			}

		default:
			log.Printf("[%s] add command of {%s}", __FILE__, rcv.Cmd)
		}

	}
}

func wsReceive(wss_server_url string, ws *websocket.Conn, chan_msg chan string) (err error) {

	var read_buf string

	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		for {
			log.Printf("[%s] panic recovery !!!", __FILE__)
			ws, err = wsProxyDial(wss_server_url, "tcp", wss_server_url)

			if err != nil {
				log.Printf("[%s] wsProxyDial : %s ", __FILE__, err)
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}()

	for {
		err = websocket.Message.Receive(ws, &read_buf)
		if err != nil {
			log.Printf("[%s] wsReceive : %s", __FILE__, err)
			syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
			break
		}
		log.Printf("[%s] received: %s", __FILE__, read_buf)
		chan_msg <- read_buf
	}

	return err
}

func wsReqeustConnection(ws *websocket.Conn, name string) (err error) {
	send := ConnectReq{}
	send.Cmd = "request"
	send.Name = name

	websocket.JSON.Send(ws, send)

	return nil
}


func wsProxyDial(url_, protocol, origin string) (ws *websocket.Conn, err error) {

	log.Printf("[%s] http_proxy {%s}\n", __FILE__, os.Getenv("HTTP_PROXY"))

	// comment out in case of testing without proxy
	if strings.Contains(url_, "10.113.") {
		return websocket.Dial(url_, protocol, origin)
	}

	if os.Getenv("HTTP_PROXY") == "" {
		return websocket.Dial(url_, protocol, origin)
	}

	purl, err := url.Parse(os.Getenv("HTTP_PROXY"))
	if err != nil {
		log.Printf("[%s] Parse : ", __FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Printf("====================================")
	log.Printf("    websocket.NewConfig")
	log.Printf("====================================")
	config, err := websocket.NewConfig(url_, origin)
	if err != nil {
		log.Printf("[%s] NewConfig : ", __FILE__,  err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	if protocol != "" {
		config.Protocol = []string{protocol}
	}

	log.Printf("====================================")
	log.Printf("    HttpConnect")
	log.Printf("====================================")
	client, err := wsHttpConnect(purl.Host, url_)
	if err != nil {
		log.Printf("[%s] HttpConnect : ", __FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Printf("====================================")
	log.Printf("    websocket.NewClient")
	log.Printf("====================================")
	return websocket.NewClient(config, client)
}

func wsHttpConnect(proxy, url_ string) (io.ReadWriteCloser, error) {
	log.Printf("[%s] proxy =", __FILE__, proxy)
	proxy_tcp_conn, err := net.Dial("tcp", proxy)
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] proxy_tcp_conn =", __FILE__, proxy_tcp_conn)
	log.Printf("[%s] url_ =", __FILE__, url_)

	turl, err := url.Parse(url_)
	if err != nil {
		log.Printf("[%s] Parse : ", __FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}

	log.Printf("[%s] proxy turl.Host =", __FILE__, string(turl.Host))

	req := http.Request{
		Method: "CONNECT",
		URL:    &url.URL{},
		Host:   turl.Host,
	}

	proxy_http_conn := httputil.NewProxyClientConn(proxy_tcp_conn, nil)
	//cc := http.NewClientConn(proxy_tcp_conn, nil)

	log.Printf("[%s] proxy_http_conn =", __FILE__, proxy_http_conn)

	resp, err := proxy_http_conn.Do(&req)
	if err != nil && err != httputil.ErrPersistEOF {
		log.Printf("[%s] ErrPersistEOF : ", __FILE__, err)
		syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
		return nil, err
	}
	log.Printf("[%s] proxy_http_conn<resp> =", __FILE__, (resp))

	rwc, _ := proxy_http_conn.Hijack()

	return rwc, nil

}
