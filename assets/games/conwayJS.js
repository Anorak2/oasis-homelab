const gridWidth = 32; //32
const gridHeight = 16; //16
function canvasClick() {
	const xPos = event.clientX;
	const yPos = event.clientY;
	const width = window.innerWidth*.75;
	const height = window.innerHeight*.75;
	const leftPad = .125 * window.innerWidth;
	const topPad = 50;
	const column = Math.floor((xPos-leftPad)/(width/gridWidth));
	const row = Math.floor((yPos-topPad)/(height/gridHeight));
	let params = new URLSearchParams();
	params.append('col', column);
	params.append('row', row);
	fetch("http://localhost:8080/games/conway/post",
		{
			method: "POST",
			body: params.toString(), 
			headers: 
			{
				"Content-Type": "application/x-www-form-urlencoded"
			}
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

function drawBoxes(input){
	const width = window.innerWidth*.75;
	const height = window.innerHeight*.75;
	const xStep = width/gridWidth; 
	const yStep = height/gridHeight;
	ctx.fillStyle = "white";
	var count = 0;
	for(let row = 0; row < gridHeight; row++){
		for(let col = 0; col < gridWidth; col++){
			if(input[count]==1){
				ctx.fillRect(xStep*col,yStep*row,xStep,yStep);
			}
			count += 1;
		}
	}
}

function updateFullCanvas(){
	sizeCanvas();
	drawGrid();
	fetchBoard();
}
function simpleUpdate(){
	ctx.clearRect(0,0,c.width, c.height);
	drawGrid();
	fetchBoard();
}

function fetchBoard(){
	fetch("http://localhost:8080/games/conway/get",
		{
			method: "GET",
		}).then((response) => response.text())
		.then(data =>{
			drawBoxes(data);
		})
	.catch(error => console.error("Error:", error));
	
}



var c = document.getElementById("mainCanvas");
var ctx = c.getContext("2d");
updateFullCanvas();
setInterval(simpleUpdate, 500);
c.onclick = canvasClick;
window.onresize = updateFullCanvas;
