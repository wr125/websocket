package main

import (
	"fmt"
)

type Status byte

type Board []byte

type Turn struct {
	Board  Board
	Status Status
}

//Competitor is an instance of Player 1 and Player 2
type Competitor struct {
	Name   string          //Name of Players
	Tactic func(Board) int //Returns updated state of game
	Rating int16           //Is the players rating
}

//Game receives 3 structs of type []Turn, Player 1 and Player 2
type Game struct {
	OPlayer Competitor
	XPlayer Competitor
	Turns   []Turn
}

const (
	Inplay Status = 1 << iota
	Owin          //2
	Xwin          //4
	Draw          //8
)

func (s Status) String() string {
	/* This string function ensures that Status
	    can be used where a string in require,such as
	    in fmt.  Here we assign the strings associated
		with each status

		Note, we could do the same with O and X*/

	if s == Inplay {
		return "Playing"
	}
	if s == Owin {
		return "O Wins"
	}
	if s == Xwin {
		return "X Wins"
	}
	if s == Draw {
		return "A Draw"
	}

	return "Unknown"
}

func main() {

	board := []byte("O        ") //Struct Turn which is a slice of Turn struct
	state := Inplay

	o := Competitor{"random_only", //Player1 returns a board state
		func(b Board) int { return 0 },
		1000,
	}

	x := Competitor{
		"center_first", //Player2 returns a board state
		func(b Board) int { return 0 },
		1000,
	}
	t1 := Turn{ //variable t1 of turn which returns board and state
		board,
		state,
	}
	t2 := Turn{[]byte("OX       "), Inplay}
	t3 := Turn{[]byte("OX O     "), Inplay}
	t4 := Turn{[]byte("OX OX    "), Inplay}
	t5 := Turn{[]byte("OX OX O  "), Owin}

	var turns []Turn
	turns = append(turns, t1)
	turns = append(turns, t2)
	turns = append(turns, t3)
	turns = append(turns, t4)
	turns = append(turns, t5)

	g := Game{o, x, turns} //Return both Players and their turns
	fmt.Printf("Board is %s \n", string(board))
	fmt.Printf("player O is called %s \n", o.Name)
	fmt.Printf("player X is called %s \n\n", x.Name)
	fmt.Printf("Game %#v \n\n\n", g)
	fmt.Printf("Noughts competitor is %s Rating is %d\n", g.OPlayer.Name, g.OPlayer.Rating)
	fmt.Printf("Crosses competitor is %s Rating is %d\n\n", g.XPlayer.Name, g.XPlayer.Rating)
	fmt.Printf("Turns\n")
	fmt.Printf(" %s \t %s\n", string(g.Turns[0].Board), g.Turns[0].Status)
	fmt.Printf(" %s \t %s\n", string(g.Turns[1].Board), g.Turns[1].Status)
	fmt.Printf(" %s \t %s\n", string(g.Turns[2].Board), g.Turns[2].Status)
	fmt.Printf(" %s \t %s\n", string(g.Turns[3].Board), g.Turns[3].Status)
	fmt.Printf(" %s \t %s\n", string(g.Turns[4].Board), g.Turns[4].Status)
}
