package conway

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
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

func printBoard(tempBoard *Board){
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
	fmt.Println("Str:", tempBoard.board_as_string)
	fmt.Println()
}

func HandlePost(rw http.ResponseWriter, req *http.Request){
	// This functions catches the post request and parses it before passing it along

	// This gets the raw data in the post resquest
	reqBody, err := io.ReadAll(req.Body) 
	if err != nil {
		fmt.Println("Failed to read post request: ", err)
		return
	}
	// Decode the data from the raw ascii

	decodedParams, err := url.ParseQuery(string(reqBody))
	if err != nil {
		fmt.Println("Failed to parse a post request", err)
		return
	}
	// Convert the data to integers
	col, err := strconv.ParseInt(decodedParams.Get("col"),10,0)
	if err != nil {
		fmt.Println("Failed to convert a post req to int", err)
		return
	}
	row, err := strconv.ParseInt(decodedParams.Get("row"),10,0)
	if err != nil {
		fmt.Println("Failed to convert a post req to int", err)
		return
	}

	board.mu.Lock()
	if(row >= 0 && col >= 0 && int(row) < len(board.mboard) && int(col) < len(board.mboard[0])){
		board.mboard[row][col] = !board.mboard[row][col]
	}
	board.mu.Unlock()
}

func HandleGet(rw http.ResponseWriter, req *http.Request){
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	arrToString()
	board.mu.Lock()
	rw.Write([]byte(board.board_as_string))
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
