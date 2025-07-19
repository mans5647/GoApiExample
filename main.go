package main

import (
	"first/types"
	"net/http"
	"strconv"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/gin-gonic/gin"
)




const rows int = 11;
const columns int = 22;

var maze = [rows][columns] int {
	{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
	{1,1,1,1,1,1,1,1,1,1,0,0,0,0,1,1,1,0,0,0,0,0},
	{0,0,0,0,0,0,0,0,0,1,0,0,0,0,1,0,1,0,1,1,0,0},
	{0,0,1,1,1,1,1,1,0,1,1,1,1,0,1,0,1,0,1,0,0,0},
	{0,0,1,0,0,0,0,1,1,1,0,0,1,0,1,0,1,1,1,0,1,0},
	{0,0,1,0,1,0,0,0,0,0,0,0,1,1,1,0,0,0,1,0,1,0},
	{0,0,1,1,1,1,1,0,0,0,1,0,0,1,0,0,0,0,1,0,1,0},
	{0,0,0,0,0,0,0,0,0,0,1,1,1,1,1,0,0,0,1,1,1,0},
	{0,0,1,1,1,1,0,1,1,1,1,0,0,0,1,1,0,0,0,0,0,0},
	{0,0,1,0,0,1,1,1,0,0,0,0,0,0,0,1,1,1,1,1,1,999},
	{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
}

const EXIT_VALUE = 999
const API_BASE string = "maze/"
const API_TYPE_MOVE string = "move"
const API_TYPE_REQUIRE_LOC = "/myloc"
const QUERY_WHERE string = "where"

const API_MAZE_MOVE = API_BASE + API_TYPE_MOVE
const API_MAZE_LOOK = API_BASE + API_TYPE_REQUIRE_LOC
var pos types.Position;


type Direction int

const (
	None = -1
	Left Direction = 1
	Right Direction = 4
	Up Direction = 8
	Bottom Direction = 16
)


type LocationInfo struct
{
	WallCount int
	DirectionsAvailable int
}

type WalkResult struct
{
	Walked bool
	Where int
	Exit bool
	CurrentLocation types.Position
}

func CorrectLocation() {

	if (pos.Column >= (columns - 1) ) { pos.Column = columns - 2 }
	if (pos.Row >= (rows - 1) ) {	pos.Row = rows - 2	}
	if (pos.Row < 0) { pos.Row = 1 }
	if (pos.Column < 0) { pos.Column = 1 }


}

func AvailableTo(dir Direction) bool {
	var value int = 1;
	switch (dir) {
	case Right:
		value = maze[pos.Row][pos.Column + 1];
	case Left:
		value = maze[pos.Row][pos.Column - 1];
	case Up:
		value = maze[pos.Row - 1][pos.Column];
	case Bottom:
		value = maze[pos.Row + 1][pos.Column];
	}

	return value == 1 || value == EXIT_VALUE
}

func Go(direction int) bool {
	var c bool;

	CorrectLocation()
	
	switch (direction) {
	case 0: // right

		c = AvailableTo(Right);

		if (c) {

			pos.Column += 1;

		}
	case 1: // left


	c = AvailableTo(Left);

	if (c) {
		pos.Column -= 1;
	}

	case 2: // up


	c = AvailableTo(Up);

		if (c) {
			pos.Row -= 1;
		}

	case 3: // bottom


		c = AvailableTo(Bottom);

		if (c) {
			pos.Row += 1
		}
	}

	


	return c;
}


func Move (c * gin.Context) {

	
	value := c.Query(QUERY_WHERE)

	var walkResult WalkResult;
	walkResult.Walked = false;
	walkResult.Exit = false
	walkResult.Where = -1
	if (!stringUtils.IsEmpty(value)) {

		num , err := strconv.Atoi(value)
		if (err == nil) {
		result := Go(num)

		if (result) {
			walkResult.Walked = true;
			walkResult.Where = num
		}
		
	}
	}


	var response_code = http.StatusOK;
	if (!walkResult.Walked) {
		response_code = http.StatusBadRequest
	}

	if (maze[pos.Row][pos.Column] == EXIT_VALUE) { walkResult.Exit = true }
	walkResult.CurrentLocation = pos
	c.JSON(response_code, walkResult)
}

func LookHandler(c * gin.Context) {

	info := LookAround()

	c.JSON(200, info)
}

func LookAround() * LocationInfo { 
	
	info := LocationInfo{}

	info.DirectionsAvailable = 0
	info.WallCount = 0

	var offsets = [4][2] int {

		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	}


	for i := 0; i < len(offsets); i++ {

		var p = offsets[i]

		dy := p[0]
		dx := p[1]

		ny := dy + pos.Row
		nx := dx + pos.Column

		if (0 <= ny && ny < len(maze) && 0 <= nx && ny < len(maze[0])) {

			if (maze[ny][nx] == 0) { info.WallCount++ }

		}
	}


	fwd := AvailableTo(Right)
	bwd := AvailableTo(Left)
	up := AvailableTo(Up)
	bot := AvailableTo(Bottom)

	if (fwd) {
		info.DirectionsAvailable |= int(Right)
	}

	if (bwd) {
		info.DirectionsAvailable |= int(Left)
	}

	if (up) {
		info.DirectionsAvailable |= int(Up)
	}

	if (bot) {
		info.DirectionsAvailable |= int(Bottom)
	}


	return &info

}


func main() {

	pos.Row = 1
	pos.Column = 1
	

	router := gin.Default();

	router.GET(API_MAZE_MOVE, Move);
	router.GET(API_MAZE_LOOK, LookHandler)
	router.Run("127.0.0.1:8080")
}