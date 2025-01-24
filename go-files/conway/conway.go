package conway

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"log"
	"github.com/gorilla/websocket"
	"encoding/json"
)
type Board struct{
	mu sync.Mutex
	mboard [16][32]bool
	board_as_string string
}
var board Board = Board{
	mboard: [16][32]bool{},
	board_as_string: "",
}
var activeConnections = make(map[*websocket.Conn]bool)
var con_mu sync.Mutex

func printBoard(tempBoard *Board){
	// method is intended for debugging only
	for x := range(tempBoard.mboard){
		for y := range(tempBoard.mboard[0]){
			if tempBoard.mboard[x][y]{
				fmt.Print(1)
			} else {
				fmt.Print(0)
			}
		}
		fmt.Println("")
	}
	//fmt.Println("Str:", tempBoard.board_as_string)
	fmt.Println()
}

func changeSquare(row int, col int){
	// bit flips a single square
	board.mu.Lock()
	if(row >= 0 && col >= 0 && row < len(board.mboard) && col < len(board.mboard[0])){
		board.mboard[row][col] = !board.mboard[row][col]
	}
	board.mu.Unlock()
}

func arrToString(){
	var builder strings.Builder

	board.mu.Lock()
	for x := 0; x < len(board.mboard); x++ {
		for y := 0; y < len(board.mboard[x]); y++ {
			if board.mboard[x][y] {
				builder.WriteString("1")
			} else {
				builder.WriteString("0")
			}
		}
	}
	board.board_as_string = builder.String() 
	board.mu.Unlock()
}

func UpdateBoard(){
	board.mu.Lock()
	var tempBoard [len(board.mboard)][len(board.mboard[1])]bool 
	for x := range(board.mboard){
		for y := range(board.mboard[0]){
			amtNeighbors, err := amtNeighbors(x, y)
			if err != nil{
				fmt.Println("Neighbor Check failed")
				return
			} else if board.mboard[x][y] == true && amtNeighbors < 2{
				tempBoard[x][y] = false
			} else if board.mboard[x][y] == true && amtNeighbors > 3 {
				tempBoard[x][y] = false
			} else if board.mboard[x][y] == true{
				tempBoard[x][y] = true
			} else if amtNeighbors == 3 {
				tempBoard[x][y] = true
			}
		}
	}
	board.mboard = tempBoard
	board.mu.Unlock()
	arrToString()
	//printBoard(board.mboard)
}

func amtNeighbors(row int, col int) (int, error){
	// used for updating the board, it returns the number of active
	// neighboring cells in a chess kings radius
	if(row < 0 || col < 0 || row >= len(board.mboard) || col >= len(board.mboard[0])){
		fmt.Println("Invalid Input")
		return 0, errors.New("bad bounds") 
	}
	neigborAmt := 0
	for u := -1; u < 2; u++{
		for v := -1; v < 2; v++{
			if(u == 0 && v == 0){
				continue
			}
			if(row+u >= 0 && col+v >= 0 && row+u < len(board.mboard) && col+v < len(board.mboard[0])){
				if board.mboard[row+u][col+v]{
					neigborAmt += 1
				}
			}
		}
	}
	return neigborAmt, nil
}




type jsonData struct {
	Row int
	Column int
}

func WsHandler(conn *websocket.Conn) {

	con_mu.Lock()
	activeConnections[conn] = true
	con_mu.Unlock()

	// Listen for incoming messages from this specific connection
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Printf("Received: %s\n", message)
		if string(message) == "req-b"{
			// req-b, or request board. We respond with the current board state
			boardState := GetBoard()
			if err := conn.WriteMessage(websocket.TextMessage, []byte(boardState)); err != nil {
				log.Println("Error writing message:", err)
				break
			}
		} else {
			// if not a specific message just assume its json
			var payload jsonData
			err := json.Unmarshal([]byte(message), &payload)
			if err != nil {
				// basically just ignore this error, this is a low stakes message anyway
				log.Println("Error unmarshalling json", err)
			} else{
				changeSquare(int(payload.Row), int(payload.Column))
			}
		}
	}

	// Once the connection is closed, remove it from the active connections map
	con_mu.Lock()
	delete(activeConnections, conn)
	con_mu.Unlock()
}

func UpdateConway(){
	// this is the wrapper function for updateBoard and handles websocket updates etc
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//apply the rules to the board once
			UpdateBoard()
			message := GetBoard()
			// message each connection
			con_mu.Lock()
			for conn := range activeConnections {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					fmt.Println("Error sending conway board update:", err)
					conn.Close() // Close connection if it fails to send a message
					delete(activeConnections, conn) // Remove broken connection
				}
			}
			con_mu.Unlock()
		}
	}
}

func GetBoard() string{
	// not really necessary, this just adds a marker so that our js knows this is 
	// a board update and not some other message
	return "b"+ board.board_as_string
}
