package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/zenazn/goji"

	"flag"
	"fmt"
	"log"
	"net/http"
	"math/rand"
	"time"

	"encoding/json"
	"strconv"
)

const (
	GS_EMPTY   = 0
	GS_READY   = 2
	GS_PLAYING = 4
	GS_DONE_P1 = 10
	GS_DONE_P2 = 12
)

const (
	PDEF_P1_XPOS = 10
	PDEF_P1_YPOS = 10

	PDEF_P2_XPOS = 10
	PDEF_P2_YPOS = 50
)

type PeePlayer struct {
	XPos     int
	YPos     int
	ZPos     int
	Name     string
	Points   int
	Damage   int
	PowerUp  int
	EpocTime int64
	Ip       string
	Score    int
	RKey     string
}

type PeeeGame struct {
	Id          int
	State       int
	PlayerCount int
	Players     [2]PeePlayer
	PCount      int
}

type PCmd struct {
	GameId    int
	Cmd       int
	Timestamp int64
}

var pgames = make([]PeeeGame, 3)

func webHandlerCmd(ws *websocket.Conn) {

	for {

		//msg :=""

		fmt.Println("DEBUG:ip=" + ws.Request().RemoteAddr)

		var reply string
		err := websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Can't receive")
			break
		}

		rxCmd := &PCmd{}

		json.Unmarshal([]byte(reply), &rxCmd)
		log.Println(fmt.Sprintf("RX-CMD:GameId(%d) , CMD(%d)   ", rxCmd.GameId, rxCmd.Cmd))

		/*


			err = websocket.Message.Send(ws, msg)
			if err != nil {
				fmt.Println("Can't send")
				break
			}
		*/

		fmt.Println("RX:" + reply + ",ip=" + ws.Request().RemoteAddr)
	}
}
func webHandler(ws *websocket.Conn) {
	n := 0
	for {
		msg := "Hello  " + string(n+48)
		fmt.Println("Sending to client: " + msg + "" + ws.Request().RemoteAddr)
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

	str, _ := json.Marshal(pgames)
	fmt.Fprintf(w, string(str))

}

func RESTGameNew(w http.ResponseWriter, r *http.Request) {

	for i, agame := range pgames {
		if agame.State == GS_EMPTY {
			setupGame(&pgames[i])

			log.Println("AFTER", pgames[i].Id, "xxxx", i)
			str, _ := json.Marshal(pgames[i])
			fmt.Fprintf(w, string(str))
			break
		} //end if
	}
}

func RESTGameJoin(w http.ResponseWriter, r *http.Request) {
	gameIdstr := r.URL.Query()["gameid"][0]
	rkeyIdstr := r.URL.Query()["rkey"][0]

	log.Println("RESTGAMEJOIN:gi=" + gameIdstr + ",rkey:" + rkeyIdstr)

	gameId, err := strconv.Atoi(gameIdstr)
	if err != nil {
		fmt.Println(err)
	}

	for i, _ := range pgames {
		if gameId == pgames[i].Id && pgames[i].State == GS_READY {

			if pgames[i].Players[0].RKey == "" {
				pgames[i].Players[0].RKey = rkeyIdstr
				pgames[i].PCount++
			} else if pgames[i].Players[1].RKey == "" {
				pgames[i].Players[1].RKey = rkeyIdstr
				pgames[i].PCount++
			} else if pgames[i].Players[1].RKey == rkeyIdstr {
				//
			} else if pgames[i].Players[0].RKey == rkeyIdstr {

			} else {
				http.Error(w, "Too Many Players", 429)
			}

			str, _ := json.Marshal(pgames[i])
			fmt.Fprintf(w, string(str))
			fmt.Println("%#v", pgames[i])
			return
		} //end if
	}

}

func setupGame(agame *PeeeGame) {
	rand.Seed(time.Now().UTC().UnixNano())
	agame.Id = rand.Intn(1000000) + 1
	agame.State = GS_READY
	agame.Players[0].XPos = PDEF_P1_XPOS
	agame.Players[0].YPos = PDEF_P1_YPOS
	agame.Players[0].RKey = ""

	agame.Players[1].XPos = PDEF_P2_XPOS
	agame.Players[1].YPos = PDEF_P2_YPOS
	agame.Players[1].RKey = ""
}

func main() {

	http.Handle("/echo", websocket.Handler(webHandler))
	http.Handle("/rest/cmd", websocket.Handler(webHandlerCmd))

	goji.Post("/rest/game", RESTGameNew)
	goji.Get("/rest/game", RESTGameList)
	goji.Post("/rest/gamejoin", RESTGameJoin)
	goji.Post("/rest/game", RESTGameList)

	log.Println("PeeeServer")

	flag.Set("bind", ":8080")

	goji.Serve()

}
