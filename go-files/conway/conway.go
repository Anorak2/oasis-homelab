package conway

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)
var mboard [16][32]bool
var board_as_string string

func printBoard(board [16][32]bool){
	for x := range(board){
		for y := range(board[0]){
			if board[x][y]{
				fmt.Print(1)
			} else {
				fmt.Print(0)
			}
		}
		fmt.Println("")
	}
	fmt.Println("Str:", board_as_string)
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
	if(row >= 0 && col >= 0 && int(row) < len(mboard) && int(col) < len(mboard[0])){
		mboard[row][col] = !mboard[row][col]
		//printBoard(mboard)
		//updateBoard()
		//printBoard(mboard)
	}
}

func HandleGet(rw http.ResponseWriter, req *http.Request){
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	arrToString()
	rw.Write([]byte(board_as_string))
}

func arrToString(){
	var builder strings.Builder

	for x := 0; x < len(mboard); x++ {
		for y := 0; y < len(mboard[x]); y++ {
			if mboard[x][y] {
				builder.WriteString("1")
			} else {
				builder.WriteString("0")
			}
		}
	}
	board_as_string = builder.String() 
}

func UpdateBoard(){
	var tempBoard [len(mboard)][len(mboard[1])]bool 
	for x := range(mboard){
		for y := range(mboard){
			amtNeighbors, err := amtNeighbors(x, y)
			if err != nil{
				return
			} else if mboard[x][y] == true && amtNeighbors < 2{
				tempBoard[x][y] = false
			} else if mboard[x][y] == true && amtNeighbors > 3 {
				tempBoard[x][y] = false
			} else if mboard[x][y] == true{
				tempBoard[x][y] = true
			} else if amtNeighbors == 3 {
				tempBoard[x][y] = true
			}
		}
	}
	mboard = tempBoard
	arrToString()
	//printBoard(mboard)
}

func amtNeighbors(row int, col int) (int, error){
	
	if(row < 0 || col < 0 || row >= len(mboard) || col >= len(mboard[0])){
		fmt.Println("Invalid Input")
		return 0, errors.New("bad bounds") 
	}
	neigborAmt := 0
	for u := -1; u < 2; u++{
		for v := -1; v < 2; v++{
			if(u == 0 && v == 0){
				continue
			}
			if(row+u >= 0 && col+v >= 0 && row+u < len(mboard) && col+v < len(mboard[0])){
				if mboard[row+u][col+v]{
					neigborAmt += 1
				}
			}
		}
	}
	return neigborAmt, nil
}
