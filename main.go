package main

import
(
"fmt"
 "code.google.com/p/go.net/websocket"
    "net/http"
    "log"
    "flag"
    "github.com/zenazn/goji"
 /// gojiweb   "github.com/zenazn/goji/web"
	"math/rand"
	"time"
		
 "encoding/json"
)

type PeePlayer struct {
  XPos int  
  YPos int
  ZPos int
  Name string
  Points int
  Damage int
  PowerUp int
  EpocTime int64
  Ip string
  Score int
  
}


type PeeGame  struct {
	Id int
	GState int
	PlayerCount int
	Players []PeePlayer
	
	
	
	
}

type PCmd struct {
	GameId int
	Cmd int
	Timestamp int64
}


var pgame = PeeGame{}

func webHandlerCmd(ws *websocket.Conn) {

	for {
		
		//msg :=""
		
		fmt.Println("DEBUG:ip="+ws.Request().RemoteAddr)
		

		var reply string
		err := websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Can't receive")
			break
		}
		
		rxCmd  := &PCmd{}
	
		json.Unmarshal([]byte(reply), &rxCmd)
		log.Println(fmt.Sprintf("RX-CMD:GameId(%d) , CMD(%d)   ",rxCmd.GameId,rxCmd.Cmd))
		
		/*
		
		
		err = websocket.Message.Send(ws, msg)
		if err != nil {
			fmt.Println("Can't send")
			break
		}
		*/
		
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

func RESTGameList(w http.ResponseWriter, r *http.Request) {
   
    str, _ := json.Marshal(pgame)
	fmt.Fprintf(w, string(str)) 
    
}


func setupGame(){
	rand.Seed( time.Now().UTC().UnixNano())
	pgame.Id = rand.Intn(1000000)+1

}


func main(){
	
	setupGame()
	
	http.Handle("/echo", websocket.Handler(webHandler))
	http.Handle("/rest/cmd", websocket.Handler(webHandlerCmd))
	
	
	goji.Get("/rest/game", RESTGameList)
	goji.Post("/rest/game", RESTGameList)
	
	
	log.Println("PeeeServer")

	flag.Set("bind", ":8080")

    goji.Serve()

}


