package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

const task1steps = 10
const task2steps = 40

var (
    subs map[[2]byte]byte
    polymerBytes []byte
    cache cachedPaths
)

type cachedPaths []KnownPath
    func (cp *cachedPaths) search(chr1, chr2 byte, step int) (KnownPath, bool) {
        for _, p := range *cp {
            if p.pair == [2]byte{chr1, chr2} && p.step == step {
                return p, true
            }
        }
        return KnownPath{}, false
    }

type KnownPath struct {
    pair [2]byte
    step int
    counts Counts
}
    func (p *KnownPath) new(chr1, chr2 byte, step int) {
        p.pair = [2]byte{chr1, chr2}
        p.step = step
        p.counts = make(map[byte]int)
    }

type Counts map[byte]int
    func (counts Counts) add(chr byte) {
        if c, ok := counts[chr]; ok {
            counts[chr] = c + 1
        } else {
            counts[chr] = 1
        }
    }
    func (counts Counts) calculate() int {
        var (
            max, min int
        )
        for _, count := range counts {
            if count > max {
                max = count
            }
            if min == 0 || count < min {
                min = count
            }
        }
        return max - min
    }
    func (c Counts) merge(path KnownPath) {
        if path.step == 0 { // nil path
            return
        }
        for chr, i := range path.counts {
            if cm, ok := c[chr]; ok {
                c[chr] = cm + i
            } else {
                c[chr] = i
            }
        }
    }

func dive(chr1 byte, chr2 byte, step int) KnownPath {
    if step == 0 {
        return KnownPath{}
    }
    // search for cached known path
    if step > 3 {
        if foundPath, found := cache.search(chr1, chr2, step); found {
            return foundPath
        }
    }

    // if it wasn't found, let's dive in
    var thisPath KnownPath
    thisPath.new(chr1, chr2, step)
    if insert, ok := subs[[2]byte{chr1, chr2}]; ok {
        thisPath.counts.add(insert)
        thisPath.counts.merge(dive(chr1, insert, step - 1))
        thisPath.counts.merge(dive(insert, chr2, step - 1))
    }
    if step > 4 {
        cache = append(cache, thisPath)
    }
    return thisPath
}

func main() {

    readInput()

    var task1, task2 Counts
    task1 = make(map[byte]int)
    task2 = make(map[byte]int)

    for i := 1; i < len(polymerBytes); i++ {
        if i == 1 {
            task1.add(polymerBytes[0])
            task2.add(polymerBytes[0])
        }
        task1.add(polymerBytes[i])
        task1.merge(dive(polymerBytes[i - 1], polymerBytes[i], task1steps))

        task2.add(polymerBytes[i])
        task2.merge(dive(polymerBytes[i - 1], polymerBytes[i], task2steps))
    }

    fmt.Printf("Task 1: %v\r\n", task1.calculate())
    fmt.Printf("Task 2: %v\r\n", task2.calculate())
}

// read all input data
func readInput() {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)
    // template
    scanner.Scan()
    polymerBytes = []byte(scanner.Text())
    // skip empty string
    scanner.Scan()
    // substitutions
    subs = make(map[[2]byte]byte)
    for scanner.Scan() {
        p := strings.Split(scanner.Text(), " -> ")
        pair := [2]byte{[]byte(p[0])[0], []byte(p[0])[1]}
        subs[pair] = []byte(p[1])[0]
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
