const gridWidth = 32; //32
const gridHeight = 16; //16
let ws;
var board = ""

function canvasClick() {
	const xPos = event.clientX;
	const yPos = event.clientY;
	const width = window.innerWidth*.75;
	const height = window.innerHeight*.75;
	const leftPad = .125 * window.innerWidth;
	const topPad = 50;
	const columndata = Math.floor((xPos-leftPad)/(width/gridWidth));
	const rowdata = Math.floor((yPos-topPad)/(height/gridHeight));
	const data = {row: rowdata, column: columndata}

    waitForSocketConnection(ws, function(){
		ws.send(JSON.stringify(data));
    });
}

function sizeCanvas(){
	// This is a function to set up the basic canvas
	// It gets run every time the screen gets resized 
	const width = window.innerWidth;
	const height = window.innerHeight;
	c.height = height * 0.75;
	c.width = width * .75;
	c.style.position = "absolute";
	c.style.top = "50px";
	c.style.left = width * .125+"px";
}
function drawGrid(){
	// This function draws all of the lines on the screen to form a grid
	const width = window.innerWidth*.75;
	const height = window.innerHeight*.75;
	xStep = width/gridWidth; 
	yStep = height/gridHeight;
	for (let index = 0; index < gridWidth; index++){
		ctx.beginPath();
		ctx.moveTo(xStep*index, 0);
		ctx.lineTo(xStep*index, height);
		ctx.strokeStyle = "gray";
		ctx.lineWidth = 1;
		ctx.stroke();
	}
	for (let index = 0; index < gridHeight; index++){
		ctx.beginPath();
		ctx.moveTo(0, yStep*index);
		ctx.lineTo(width, yStep*index);
		ctx.strokeStyle = "white";
		ctx.lineWidth = 0.5;
		ctx.stroke();
	}
}

function drawBoxes(){
	const width = window.innerWidth*.75;
	const height = window.innerHeight*.75;
	const xStep = width/gridWidth; 
	const yStep = height/gridHeight;
	ctx.fillStyle = "white";
	var count = 0;
	for(let row = 0; row < gridHeight; row++){
		for(let col = 0; col < gridWidth; col++){
			if(board[count]==1){
				ctx.fillRect(xStep*col,yStep*row,xStep,yStep);
			}
			count += 1;
		}
	}
}
function drawSquareFlip(row, col){
	console.log(row + " " + col);
	//console.log(board);
	var s = (row*gridWidth)+col;
	//console.log(s);
	if(board[s] === "1"){
		board = board.substring(0,s) + "0" + board.substring(s+1);
	}else {
		board = board.substring(0,s) + "1" + board.substring(s+1);
	}
	console.log(board)

	ctx.clearRect(0,0,c.width, c.height);
	drawGrid();
	drawBoxes();
}

function updateFullCanvas(){
	sizeCanvas();
	drawGrid();
	requestBoardUpdate();
}


function connect() {
	ws = new WebSocket("ws://localhost:8080/games/ws/conway");

	ws.onopen = function() {
		console.log("Connected to WebSocket server");
	};

	ws.onmessage = function(event) {
		response = event.data;
		if(response[0] == "b"){
			ctx.clearRect(0,0,c.width, c.height);
			drawGrid();
			board = response.substring(1);
			drawBoxes();
		} else if(response[0] == "s"){
			var square = JSON.parse(response.substring(1));
			drawSquareFlip(square.Row, square.Column);
		}
	};

	ws.onclose = function() {
		console.log("WebSocket connection closed, retrying...");
		setTimeout(connect, 1000); // Reconnect after 1 second
	};

	ws.onerror = function(error) {
		console.error("WebSocket error:", error);
	};
}

function requestBoardUpdate(){
    waitForSocketConnection(ws, function(){
        console.log("requested board update");
		ws.send("req-b");
    });
}

function waitForSocketConnection(socket, callback){
    setTimeout(
        function () {
            if (socket.readyState === 1) {
                if (callback != null){
                    callback();
                }
            } else {
                waitForSocketConnection(socket, callback);
            }

        }, 15); // time to wait in miliseconds
}
connect();

var c = document.getElementById("mainCanvas");
var ctx = c.getContext("2d");
updateFullCanvas();
c.onclick = canvasClick;
window.onresize = updateFullCanvas;
