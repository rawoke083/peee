package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/zenazn/goji"

	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"encoding/json"
	"github.com/zenazn/goji/web"
	"strconv"
)

const (
	GS_EMPTY        = 0
	GS_READY        = 2
	GS_ACTION_A_REQ = 3
	GS_KICKOFF      = 4

	GS_PLAYING = 5

	GS_DONE    = 100
	GS_DONE_P1 = 110
	GS_DONE_P2 = 120
)

const (
	PC_UP = 1

	PC_DOWN = 2

	PC_LEFT = 3

	PC_RIGHT = 4

	PC_ACTION_A = 10
	PC_ACTION_B = 12
	PC_GAME_GET = 20
)

const (
	PDEF_P1_XPOS = 10
	PDEF_P1_YPOS = 140

	PDEF_P2_XPOS = 1130
	PDEF_P2_YPOS = 140
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
	Height   int
	Width    int
}

type PeeItem struct {
	XPos int
	YPos int
	ZPos int
	Name string

	VX int
	VY int
	AX int
	AY int

	Src             string
	Visible         int
	DefaultStepSize int
	Height          int
	Width           int
}

type PeeeGame struct {
	Id          int
	State       int
	PlayerCount int
	Players     [2]PeePlayer
	PCount      int
	WorldWidth  int
	WorldHeight int
	Ball        PeeItem
	PowerUp     PeeItem
}

type PCmd struct {
	GameId    int
	Cmd       int
	Cmd2      int
	Timestamp int64
	RKey      string
}

var pgames = make([]PeeeGame, 3)

func updateWorlds() {

	for {

		for i, _ := range pgames {

			if pgames[i].State >= GS_KICKOFF && pgames[i].State <= GS_DONE {

				pgames[i].Ball.XPos = pgames[i].Ball.XPos + pgames[i].Ball.VX
				pgames[i].Ball.YPos = pgames[i].Ball.YPos + pgames[i].Ball.VY

				if pgames[i].Ball.YPos < 0 { // hitting the top wall
					pgames[i].Ball.YPos = 0
					pgames[i].Ball.VY = -pgames[i].Ball.VY
					// pgames[i].Ball.VX = -pgames[i].Ball.VX;

				} else if pgames[i].Ball.YPos > (600 - 50) { // hitting the bottom wall
					pgames[i].Ball.YPos = 600 - 50
					pgames[i].Ball.VY = -pgames[i].Ball.VY
					// pgames[i].Ball.VX = -pgames[i].Ball.VX;

				}

				// a point was scored
				if pgames[i].Ball.XPos < 0 || pgames[i].Ball.XPos > 1250 {
					if pgames[i].Ball.XPos < 5 {

						pgames[i].Players[1].Score++
					} else if pgames[i].Ball.XPos > 1250 {

						pgames[i].Players[0].Score++
					}

					pgames[i].Ball.XPos = pgames[i].WorldWidth/2 - (pgames[i].Ball.Width / 2)
					pgames[i].Ball.YPos = pgames[i].WorldHeight/2 - (pgames[i].Ball.Height / 2)
					pgames[i].Ball.VY = 0

					pgames[i].Ball.VX = -pgames[i].Ball.VX
					pgames[i].State = GS_ACTION_A_REQ

				}

				//hit left paddle

				if pgames[i].Ball.XPos < pgames[i].Players[0].XPos+pgames[i].Players[0].Width && pgames[i].Ball.XPos+50 > pgames[i].Players[0].XPos &&
					pgames[i].Ball.YPos < pgames[i].Players[0].YPos+pgames[i].Players[0].Height && pgames[i].Ball.YPos+50 > pgames[i].Players[0].YPos {

					pgames[i].Ball.VX = -pgames[i].Ball.VX
					pgames[i].Ball.VY += (pgames[i].Ball.YPos - pgames[i].Players[0].YPos)
					if pgames[i].Ball.VY > 3 {
						pgames[i].Ball.VY = 2

					} else if pgames[i].Ball.VY < -3 {
						pgames[i].Ball.VY = -2

					}

					pgames[i].Ball.XPos += 130

				}
				//hit right paddle
				if pgames[i].Ball.XPos < pgames[i].Players[1].XPos+pgames[i].Players[1].Width && pgames[i].Ball.XPos+50 > pgames[i].Players[1].XPos &&
					pgames[i].Ball.YPos < pgames[i].Players[1].YPos+pgames[i].Players[1].Height && pgames[i].Ball.YPos+50 > pgames[i].Players[1].YPos {
					// The objects are touching

					pgames[i].Ball.VX = -pgames[i].Ball.VX
					pgames[i].Ball.VY += (pgames[i].Ball.YPos - pgames[i].Players[1].YPos)
					if pgames[i].Ball.VY > 3 {
						pgames[i].Ball.VY = 2

					} else if pgames[i].Ball.VY < -3 {
						pgames[i].Ball.VY = -2

					}

					pgames[i].Ball.XPos -= 130

				}

			} //end if game state
		} //end foreach games


		time.Sleep(12000000)

	} //end for-infinite
} //end updatewords

func findGameById(gameId int) *PeeeGame {

	for i, agame := range pgames {
		if agame.Id == gameId {
			return &pgames[i]
		} //end if
	}
	return nil

} //end find game

func processCmd(cmd *PCmd) *PeeeGame {
	stepSize := 30

	game := findGameById(cmd.GameId)
	if game == nil {
		log.Println("\n\nBAD NEWS CANT FIND\n\n")
		return nil
	}

	pindex := 0

	if cmd.RKey == game.Players[1].RKey {
		pindex = 1
	}
	_ = pindex
	switch cmd.Cmd2 {
	case PC_LEFT:
		{
			game.Players[pindex].XPos = game.Players[pindex].XPos - stepSize

			if pindex == 0 && game.Players[pindex].XPos < 10 {
				game.Players[pindex].XPos = 10
			}

			if pindex == 1 && game.Players[pindex].XPos < (game.WorldWidth/2)+(game.Players[pindex].Width+100) {
				game.Players[pindex].XPos = (game.WorldWidth / 2) + (game.Players[pindex].Width + 100)
			}

			log.Println("CMD LEFR")
			break

		}
	case PC_RIGHT:
		{
			game.Players[pindex].XPos = game.Players[pindex].XPos + stepSize

			if pindex == 0 && game.Players[pindex].XPos > (game.WorldWidth/2)-(game.Players[pindex].Width+100) {
				game.Players[pindex].XPos = (game.WorldWidth / 2) - (game.Players[pindex].Width + 100)
			}

			if pindex == 1 && game.Players[pindex].XPos > 1130 {
				game.Players[pindex].XPos = 1130
			}

			log.Println("CMD RIGHT")
			break

		}
	} //end first switch

	switch cmd.Cmd {
	case PC_UP:
		{
			//log.Println("CMD UP")

			game.Players[pindex].YPos = game.Players[pindex].YPos - stepSize

			if game.Players[pindex].YPos < 1 {
				game.Players[pindex].YPos = 0
			}

			return nil
		}

	case PC_DOWN:
		{
			//log.Println("CMD DOWN")
			game.Players[pindex].YPos = game.Players[pindex].YPos + stepSize

			if game.Players[pindex].YPos > (game.WorldHeight - game.Players[pindex].Height) {
				game.Players[pindex].YPos = (game.WorldHeight - game.Players[pindex].Height)
			}

			return nil
		}

	case PC_ACTION_A:
		{
			log.Println("CMD PC_ACTION_A")
			if game.State == GS_ACTION_A_REQ {
				game.State = GS_KICKOFF
			}
			return nil
		}
	case PC_ACTION_B:
		{

			log.Println("CMD PC_ACTION_B")
			return nil
		}

	case PC_GAME_GET:
		{

			return game

		}

	} //end switch

	return nil

} //end process

func webHandlerCmd(ws *websocket.Conn) {

	for {

		//msg :=""

		//fmt.Println("DEBUG:ip=" + ws.Request().RemoteAddr)

		var reply string
		err := websocket.Message.Receive(ws, &reply)
		if err != nil {
			fmt.Println("Can't receive")
			break
		}

		rxCmd := &PCmd{}

		json.Unmarshal([]byte(reply), &rxCmd)
		if rxCmd.Cmd != 20 {
			log.Println(fmt.Sprintf("RX-CMD:GameId(%d) , CMD(%d)   ", rxCmd.GameId, rxCmd.Cmd))
		}
		/////	fmt.Println("RX:" + reply + ",ip=" + ws.Request().RemoteAddr)
		ws_resp := processCmd(rxCmd)
		if ws_resp != nil {

			str, _ := json.Marshal(ws_resp)
			websocket.Message.Send(ws, string(str))
			//log.Println(string(str))

		} //end processCmd
		//time.Sleep()
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

func RESTGameGet(c web.C, w http.ResponseWriter, r *http.Request) {

	gameidStr := c.URLParams["gameid"]
	gameId, err := strconv.Atoi(gameidStr)
	if err != nil {
		fmt.Println(err)
	}

	for i, agame := range pgames {
		if agame.Id == gameId {
			str, _ := json.Marshal(pgames[i])
			fmt.Fprintf(w, string(str))
			return
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
		fmt.Fprintf(w, err.Error())
	}

	for i := range pgames {
		if gameId == pgames[i].Id && pgames[i].State >= GS_READY {

			if pgames[i].Players[0].RKey == "" {
				pgames[i].Players[0].RKey = rkeyIdstr
				pgames[i].PCount++
				//pgames[i].State = GS_KICKOFF

			} else if pgames[i].Players[1].RKey == "" {
				pgames[i].Players[1].RKey = rkeyIdstr
				pgames[i].PCount++
			} else if pgames[i].Players[1].RKey == rkeyIdstr {
				//
			} else if pgames[i].Players[0].RKey == rkeyIdstr {

			} else {
				///http.Error(w, "Too Many Players", 429)
			}

			if pgames[i].PCount > 1 {
				pgames[i].State = GS_KICKOFF

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

	agame.WorldWidth = 1300
	agame.WorldHeight = 600

	agame.Ball.Width = 50
	agame.Ball.Height = 50

	agame.Ball.XPos = agame.WorldWidth/2 - (agame.Ball.Width / 2)
	agame.Ball.YPos = agame.WorldHeight/2 - (agame.Ball.Height / 2)

	agame.Ball.Visible = 1

	agame.Ball.VX = 4
	agame.Ball.VY = 2


	agame.Ball.Src = "http://peeepong.host/ball.png"

	agame.Players[0].XPos = PDEF_P1_XPOS
	agame.Players[0].YPos = PDEF_P1_YPOS
	agame.Players[0].RKey = ""
	agame.Players[0].Height = 160
	agame.Players[0].Width = 160

	agame.Players[1].XPos = PDEF_P2_XPOS
	agame.Players[1].YPos = PDEF_P2_YPOS
	agame.Players[1].RKey = ""

	agame.Players[1].Height = 160
	agame.Players[1].Width = 160

}

func main() {

	http.Handle("/rest/cmd", websocket.Handler(webHandlerCmd))

	goji.Post("/rest/game", RESTGameNew)
	goji.Get("/rest/game", RESTGameList)
	goji.Get("/rest/game/:gameid", RESTGameGet)

	goji.Post("/rest/gamejoin", RESTGameJoin)

	log.Println("PeeeServer")
	setupGame(&pgames[0])
	flag.Set("bind", ":8080")
	go updateWorlds()

	goji.Serve()

}
