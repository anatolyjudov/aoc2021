package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

var cavesMap map[string][]string

// Task 1 pathes
var pathesFound1 [][]string

// Task 2 pathes
var pathesFound2 [][]string

func readCaves() {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        way := scanner.Text()
        cvs := strings.Split(way, "-")       
        cavesMap[cvs[0]] = append(cavesMap[cvs[0]], cvs[1])
        cavesMap[cvs[1]] = append(cavesMap[cvs[1]], cvs[0])
    }

}

/*
 * Task 2
 */
func explore2(path []string, smallRevisited bool, cave string) {
    path = append(path, cave)
    if cave == "end" {
        pathesFound2 = append(pathesFound2, path)
        //fmt.Printf("%v -> found path\r\n", path)
        return
    }
    for _, nextCave := range cavesMap[cave] {
        if nextCave == "start" {
            continue
        }
        if isSmallCave(nextCave) && inArray(path, nextCave) {
            if smallRevisited {
                continue
            }
            explore2(path, true, nextCave)
            continue
        }
        explore2(path, smallRevisited, nextCave)
    }
    //fmt.Printf("%v -> no way there\r\n", path)
}

/*
 * Task 1
 */
func explore1(path []string, cave string) {
    path = append(path, cave)
    if cave == "end" {
        pathesFound1 = append(pathesFound1, path)
        //fmt.Printf("%v -> found path\r\n", path)
        return
    }
    for _, nextCave := range cavesMap[cave] {
        if isSmallCave(nextCave) && inArray(path, nextCave) {
            continue
        }
        explore1(path, nextCave)
    }
    //fmt.Printf("%v -> no way there\r\n", path)
}

func isSmallCave(cave string) bool {
    return []byte(cave)[0] > 96
}

func inArray(s []string, e string) bool {
    for i := 0; i < len(s); i++ {
        if e == s[i] {
            return true
        }
    }
    return false
}

//
func main() {
    cavesMap = make(map[string][]string)
    readCaves()
    fmt.Println(cavesMap)

    explore1([]string{}, "start")
    explore2([]string{}, false, "start")

    fmt.Printf("Task 1 total pathes found: %v", len(pathesFound1))
    fmt.Printf("Task 2 total pathes found: %v", len(pathesFound2))
}

// file opening routine
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}