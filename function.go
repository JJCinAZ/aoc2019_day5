// function.go
package function

import (
	"fmt"
	"strconv"
	"strings"
)

func parseInput(input string) []int {
	a := strings.Split(input, ",")
	pgm := make([]int, len(a))
	for i := range a {
		pgm[i], _ = strconv.Atoi(a[i])
	}
	return pgm
}

type OpCode struct {
	op        int
	parmModes [3] int
}

func decodeOp(op int) OpCode {
	result := OpCode{}
	result.parmModes[2] = op / 10000
	op = op % 10000
	result.parmModes[1] = op / 1000
	op = op % 1000
	result.parmModes[0] = op / 100
	result.op = op % 100
	return result
}

func execPgm(pgm []int, inputBuffer InputBuffer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	ip := 0
PGMLOOP:
	for {
		opcode := decodeOp(pgm[ip])
		switch opcode.op {
		case 99:
			break PGMLOOP
		case 1: // Addition
			v1, v2 := getParamsValues(opcode, pgm, ip)
			op3 := pgm[ip+3]
			pgm[op3] = v1 + v2
			ip += 4
		case 2: // Multiplication
			v1, v2 := getParamsValues(opcode, pgm, ip)
			op3 := pgm[ip+3]
			pgm[op3] = v1 * v2
			ip += 4
		case 3: // Input
			pgm[pgm[ip+1]] = inputBuffer.get()
			fmt.Printf("INPUT:%d\n", pgm[pgm[ip+1]])
			ip += 2
		case 4: // Output
			fmt.Printf("OUTPUT:%d\n", pgm[pgm[ip+1]])
			ip += 2
		case 5: // Jump-if-true
			v1, v2 := getParamsValues(opcode, pgm, ip)
			if v1 != 0 {
				ip = v2
			} else {
				ip += 3
			}
		case 6: // Jump-if-false
			v1, v2 := getParamsValues(opcode, pgm, ip)
			if v1 == 0 {
				ip = v2
			} else {
				ip += 3
			}
		case 7: // Less-than
			v1, v2 := getParamsValues(opcode, pgm, ip)
			op3 := pgm[ip+3]
			if v1 < v2 {
				pgm[op3] = 1
			} else {
				pgm[op3] = 0
			}
			ip += 4
		case 8: // Equals
			v1, v2 := getParamsValues(opcode, pgm, ip)
			op3 := pgm[ip+3]
			if v1 == v2 {
				pgm[op3] = 1
			} else {
				pgm[op3] = 0
			}
			ip += 4
		default:
			panic(fmt.Errorf("illegal opcode at offset %d", ip))
		}
	}
	return nil
}

func getParamsValues(opcode OpCode, pgm []int, ip int) (int, int) {
	v1 := pgm[ip+1]
	if opcode.parmModes[0] == 0 {
		v1 = pgm[v1]
	}
	v2 := pgm[ip+2]
	if opcode.parmModes[1] == 0 {
		v2 = pgm[v2]
	}
	return v1, v2
}

type InputBuffer struct {
	buff []int
	position int
}

func (buffer *InputBuffer) push(x int) {
	buffer.buff = append(buffer.buff, x)
}

func (buffer *InputBuffer) get() int {
	if buffer.position >= len(buffer.buff) {
		panic("EOF")
	}
	x := buffer.buff[buffer.position]
	buffer.position++
	return x
}