package main

import (
	"bufio"
	"fmt"
	"os"

	"dasa.cc/tis100mc"
)

var n0 = `# N0
L: ADD 1
MOV ACC, RIGHT
SAV
SWP
SUB 7
RET:
JEZ RET
SWP
JMP L`

var n0two = `#
# MOV 1, NIL
L: ADD 1
MOV ACC, RIGHT
SAV
SWP
SUB 7
JEZ RET
SWP
JMP L
RET: JRO 0`

var n1 = `# N1
ADD LEFT`

func main() {
	sys := tis100mc.NewSystem()
	sys.Node(0).LoadProgram([]byte(n0two))
	sys.Node(1).LoadProgram([]byte(n1))

	prompt := true

LOOP:
	for {
		if prompt {
			fmt.Printf("> ")
			r := bufio.NewReader(os.Stdin)
			l, err := r.ReadString('\n')
			if err != nil {
				panic(err)
			}
			switch l {
			case "q\n":
				break LOOP
			case "r\n":
				prompt = false
			case "p\n":
				fmt.Println("Node 0 ACC/BAK:", sys.Node(0).ACC(), sys.Node(0).BAK())
				fmt.Println("Node 1 ACC/BAK:", sys.Node(1).ACC(), sys.Node(1).BAK())
				fmt.Println("---")
				continue LOOP
			}
		}

		fmt.Println("---")
		i, op := sys.Node(0).Operation()
		fmt.Println("Node 0:", i, op)
		i, op = sys.Node(1).Operation()
		fmt.Println("Node 1:", i, op)
		sys.Step()
		// fmt.Printf("---\n0: %+v\n1: %+v\n---\n", sys.nodes[0], sys.nodes[1])
		fmt.Println("Node 0 ACC/BAK:", sys.Node(0).ACC(), sys.Node(0).BAK())
		fmt.Println("Node 1 ACC/BAK:", sys.Node(1).ACC(), sys.Node(1).BAK())
		fmt.Println("---")
	}
}
