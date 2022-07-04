package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type operator interface {
	calc(a, b int) int
	string() string
}

type add struct {}
func(a add) calc(x, y int) int {
	return x + y
}
func(a add) string() string {
	return " + "
}

type subtract struct {}
func(s subtract) calc(x, y int) int {
	return x - y
}
func(s subtract) string() string {
	return " - "
}

type multiply struct {}
func(m multiply) calc(x, y int) int {
	return x * y
}
func(m multiply) string() string {
	return " * "
}

type divide struct {}
func(d divide) calc(x, y int) int {
	return x / y
}
func(d divide) string() string {
	return " // "
}

type equation struct {
	left int 
	operator operator
	right int
	equals int
}
func (e equation) string() string {
	return strconv.Itoa(e.left) + e.operator.string() + strconv.Itoa(e.right) + " = "
}

type equationIterator struct {
	rand *rand.Rand
	n int
}
func(eg *equationIterator) next() equation {
	left, right := eg.rand.Intn(15) + 5, eg.rand.Intn(15) + 5
	operations := []operator{
		add{},
		subtract{},
		multiply{},
		divide{},
	}
	op := operations[eg.rand.Intn(len(operations))]
	equals := op.calc(left, right)

	return equation{
		left: left,
		right: right,
		operator: op,
		equals: equals,
	}
}
func (eg *equationIterator) done() bool {
	if eg.n == 0 {
		return false
	}
	prev := eg.n
	eg.n--
	return prev == 1
}
func newIterator() equationIterator {
	src := rand.NewSource(time.Now().UnixMilli())
	return equationIterator{
		rand: rand.New(src),
	}
}

type counts struct {
	correct int
	incorrect int
}
func play(current chan counts) {
	reader := bufio.NewReader(os.Stdin)
	iterator := newIterator() 
	for !iterator.done() {
		eq := iterator.next()
		fmt.Print(eq.string())
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		answer, err := strconv.Atoi(text)
		if err != nil || answer != int(eq.equals) {
			fmt.Println("wrong:", eq.string() + strconv.Itoa(eq.equals))
			current <- counts{0, 1}
		} else {
			current <- counts{1, 0} 
		}
	}
}

func main() {
	current := make(chan counts)
	total := counts{}
	
	go play(current)
	timer := time.NewTimer(time.Second * 30)
	loop: 
		for {
			select {
			case n := <-current:
				if n.correct != 0 {
					total.correct += n.correct
				}
				if n.incorrect != 0 {
					total.incorrect += n.incorrect
				}
			case <-timer.C:
				timer.Stop()
				close(current)
			  break loop
			}
		}
	fmt.Println("\n-----fin-----")
	fmt.Println("  correct:", total.correct)
	fmt.Println("incorrect:", total.incorrect)
}