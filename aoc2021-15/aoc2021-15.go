package main

import (
    "fmt"
    "os"
    "bufio"
)

type Cavern struct {
    x, y int
    risks []byte
}
    func (cavern *Cavern) getRiskAt(x, y int) (byte, bool) {
        if x < 0 || x >= cavern.x || y < 0 || y >= cavern.y {
            return 0, false
        }
        return cavern.risks[y * cavern.x + x], true
    }
    func (cavern *Cavern) setRiskAt(x, y int, r byte) {
        cavern.risks[y * cavern.x + x] = r
    }
    func (cavern *Cavern) makeLarger(scale int) Cavern {
        var newCavern Cavern
        newCavern.x = cavern.x * scale
        newCavern.y = cavern.y * scale
        newCavern.risks = make([]byte, newCavern.x * newCavern.y)
        for y := 0; y < cavern.y; y++ {
            for x := 0; x < cavern.x; x++ {
                r, _ := cavern.getRiskAt(x, y)
                for yS := 0; yS < scale; yS++ {
                    for xS := 0; xS < scale; xS++ {
                        nx := x + xS * cavern.x
                        ny := y + yS * cavern.y
                        nr := byte((int(r) + yS + xS - 1) % 9 + 1)
                        newCavern.setRiskAt(nx, ny, nr)
                    }
                }
            }
        }
        return newCavern
    }

type Dijkstra struct {
    cavern    *Cavern
    distances []int
    step      int
}
    func (dij *Dijkstra) getDistanceAt(x, y int) (int, bool) {
        if x < 0 || x >= dij.cavern.x || y < 0 || y >= dij.cavern.y {
            return 0, false
        }
        return dij.distances[y * dij.cavern.x + x], true
    }
    func (dij *Dijkstra) setDistanceAt(x, y, distance int) bool {
        if x < 0 || x >= dij.cavern.x || y < 0 || y >= dij.cavern.y {
            return false
        }
        dij.distances[y * dij.cavern.x + x] = distance
        return true
    }
    func (dij *Dijkstra) findPath() {
        dij.distances = make([]int, dij.cavern.x * dij.cavern.y)
        for step := 1; true; step++ {
            if !dij.proceedStep() {
                fmt.Printf("No changes at step %v\r\n", step)
                break
            }
        }
    }
    func (dij *Dijkstra) proceedStep() bool {
        somethingChanged := false
        for y := 0; y < dij.cavern.y; y++ {
            for x := 0; x < dij.cavern.x; x++ {
                if x == 0 && y == 0 {
                    continue
                }
                pointRisk, _ := dij.cavern.getRiskAt(x, y)
                pointDistance, _ := dij.getDistanceAt(x, y)
                newPointDistance := pointDistance
                for _, np := range getNearPoints(x, y) {
                    npDistance, ok := dij.getDistanceAt(np[0], np[1])
                    if !ok || (npDistance == 0 && np != [2]int{0, 0}) {
                        continue
                    }
                    if (npDistance + int(pointRisk)) < pointDistance || newPointDistance == 0 {
                        newPointDistance = npDistance + int(pointRisk)
                    }
                }
                if newPointDistance != pointDistance {
                    dij.setDistanceAt(x, y, newPointDistance)
                    somethingChanged = true
                }
            }
        }
        return somethingChanged
    }
    func (dij *Dijkstra) print() {
        for y := 0; y < dij.cavern.y; y++ {
            for x := 0; x < dij.cavern.x; x++ {
                risk, _ := dij.cavern.getRiskAt(x, y)
                fmt.Printf("%2d ", risk)
            }
            fmt.Printf(" | ")
            for x := 0; x < dij.cavern.x; x++ {
                dist, _ := dij.getDistanceAt(x, y)
                fmt.Printf("%2d ", dist)
            }
            fmt.Printf("\r\n")
        }
    }
    func (dij *Dijkstra) printLast() {
        d, _ := dij.getDistanceAt(dij.cavern.x - 1, dij.cavern.y - 1)
        fmt.Printf("\r\nLast distance is %v\r\n", d)
    }

func getNearPoints(x, y int) [4][2]int {
    p := [4][2]int{
            {x - 1, y},
            {x, y - 1},
            {x, y + 1},
            {x + 1, y},
        }
    return p
}

func getNearPoints2(x, y int) [8][2]int {
    p := [8][2]int{
            {x - 1, y - 1},
            {x - 1, y},
            {x, y - 1},
            {x + 1, y - 1},
            {x - 1, y + 1},
            {x, y + 1},
            {x + 1, y},
            {x + 1, y + 1},
        }
    return p
}

func main() {

    var (
       cavern  Cavern
       cavern2 Cavern
       dij     Dijkstra
    )

    readCavern(&cavern)

    dij.cavern = &cavern
    dij.findPath()
    dij.print()
    dij.printLast()

    cavern2 = cavern.makeLarger(5)

    dij.cavern = &cavern2
    dij.findPath()
    dij.printLast()

    //fmt.Println(dij)
}

// read all input data
func readCavern(cavern *Cavern) {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        row := scanner.Bytes()
        if cavern.x == 0 {
            cavern.x = len(row)
        } else if cavern.x != len(row) {
            panic("Bad row length found")
        }
        for i := 0; i < len(row); i++ {
            cavern.risks = append(cavern.risks, row[i] - 48)    // cheat byte to digit transform
        }
        cavern.y++
    }
    return
}

// file opening routine
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}
