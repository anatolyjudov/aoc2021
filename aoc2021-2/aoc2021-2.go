package main

import (
    "fmt"
    "os"
    "strings"
    "strconv"
    "bufio"
)

// Submarine's command data structure and reader
type Command struct {
    action string
    value int
}
    func (command *Command) read(order string) (bool, error) {
        orderParts := strings.Split(order, " ")
        value, err := strconv.Atoi(orderParts[1])
        if err != nil {
            return false, err
        }
        if orderParts[0] != "forward" && orderParts[0] != "up" && orderParts[0] != "down" {
            return false, fmt.Errorf("Unknown action %d", orderParts[0])
        }
        command.action = orderParts[0]
        command.value = value
        return true, nil
    }

// Transition data structure
// The final result calculator is also here
type Transition struct {
    horizontal, vertical int
}
    func (t *Transition) deeper(change int) {
        t.vertical += change;
    }
    func (t *Transition) further(change int) {
        t.horizontal += change;
    }
    func (t Transition) calculate() int {
        return t.horizontal * t.vertical
    }

// An interface for all kinds of submarines
type SubmarineInterface interface {
    do(c Command)
    getTransition() Transition
}

// Simple submarine (first part of the day's task)
type SimpleSubmarine struct {
    transition Transition
}
    func (s *SimpleSubmarine) do(c Command) {
        switch c.action {
        case "forward":
            s.transition.further(c.value);
        case "up":
            s.transition.deeper(-1*c.value)
        case "down":
            s.transition.deeper(c.value)
        }
    }
    func (s *SimpleSubmarine) getTransition() Transition {
        return s.transition;
    }

// Aimed submarine (second part of the day's task)
type AimedSubmarine struct {
    aim int
    transition Transition
}
    func (s *AimedSubmarine) do(c Command) {
        switch c.action {
        case "forward":
            s.transition.further(c.value);
            s.transition.deeper(s.aim * c.value)
        case "up":
            s.aim -= c.value
        case "down":
            s.aim += c.value
        }
    }
    func (s *AimedSubmarine) getTransition() Transition {
        return s.transition;
    }

// File opening routine
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}

//
func main() {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    var (
        command Command
        submarines [2]SubmarineInterface
    )

    submarines[0] = &SimpleSubmarine{}
    submarines[1] = &AimedSubmarine{}

    for scanner.Scan() {
        ok, err := command.read(scanner.Text());
        if !ok {
            panic(err);
        }
        for _, submarine := range submarines {
            submarine.do(command);
        }
    }
    for _, submarine := range submarines {
        fmt.Println(submarine.getTransition().calculate())
    }
}