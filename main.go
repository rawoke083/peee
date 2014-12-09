package main

import
(
"fmt"
 "code.google.com/p/go.net/websocket"
    "net/http"
    "log"
)

func webHandlerCmd(ws *websocket.Conn) {

	for {
		msg := "\nwebHandlerCmd:"
		
		fmt.Println("TX: " + msg+",ip="+ws.Request().RemoteAddr)
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			fmt.Println("Can't send")
			break
		}

		var reply string
		err = websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Can't receive")
			break
		}
		fmt.Println("RX:" + reply + ",ip="+ws.Request().RemoteAddr)
	}  
}
func webHandler(ws *websocket.Conn) {
   n:= 0;
	for {
		msg := "Hello  " + string(n+48)
		fmt.Println("Sending to client: " + msg+""+ ws.Request().RemoteAddr )
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			fmt.Println("Can't send")
			break
		}

		var reply string
		err = websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Can't receive")
			break
		}
		fmt.Println("Received back from client: " + reply)
	}  
}



func main(){
	//static file handler.
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
 	http.Handle("/static", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.Handle("/echo", websocket.Handler(webHandler))
	http.Handle("/cmd", websocket.Handler(webHandlerCmd))
	log.Println("PeeeServer")

    err := http.ListenAndServe(":8080", nil)
    if nil != err {
                panic(err)
                        }
	fmt.Println("Pee")
}


