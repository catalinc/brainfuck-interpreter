/*
 * Brainfuck interpreter.
 * See http://en.wikipedia.org/wiki/Brainfuck for more details.
 */
package main

import (
	"fmt"
	"os"
	"flag"
	"io/ioutil"
)

type Interpreter struct {
	program []byte // program data
	mem     []byte // working memory
	ip      int    // instruction pointer
	dp      int    // data pointer
}

func NewInterpreter(program []byte, memSize int) *Interpreter {
	return &Interpreter{program, make([]byte, memSize), 0, 0}
}

func (self *Interpreter) Run() {
	for self.ip < len(self.program) {
		switch self.program[self.ip] {
		case '<':
			self.dp--
			if self.dp < 0 {
				self.Fatal("execution error")
			}
		case '>':
			self.dp++
			if self.dp >= len(self.mem) {
				self.Fatal("out of memory")
			}
		case '+':
			self.mem[self.dp]++
		case '-':
			self.mem[self.dp]--
		case '.':
			fmt.Printf("%c", self.mem[self.dp])
		case ',':
			fmt.Scanf("%c", &self.mem[self.dp])
		case '[':
			if self.mem[self.dp] == 0 {
				open := 0
			L1:
				for {
					switch self.program[self.ip] {
					case '[':
						open++
					case ']':
						open--
						if open == 0 {
							break L1
						}
					}
					self.ip++
					if self.ip >= len(self.program) {
						self.Fatal("execution error")
					}
				}
			}
		case ']':
			if self.mem[self.dp] != 0 {
				closed := 0
			L2:
				for {
					switch self.program[self.ip] {
					case '[':
						closed--
						if closed == 0 {
							break L2
						}
					case ']':
						closed++
					}
					self.ip--
					if self.ip < 0 {
						self.Fatal("execution error")
					}
				}
			}
		}
		self.ip++
	}
}

func (self *Interpreter) Fatal(s string) {
	fmt.Printf("%s (ip=%d, dp=%d)\n", s, self.ip, self.dp)
	os.Exit(1)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options] file\nwhere options are:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

var memorySize = flag.Int("mem", 30000, "memory size")

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		usage()
	}

	program, err := ioutil.ReadFile(flag.Arg(0))
	if err == nil {
		NewInterpreter(program, *memorySize).Run()
	} else {
		fmt.Printf("cannot read %s (%v)\n", os.Args[1], err)
		os.Exit(1)
	}
}
