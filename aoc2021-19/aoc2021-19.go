package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
    "math"
    "sort"
)

type Trinagle struct {
    sides [3]float64
}

type Location struct {
    x, y, z int
}
    func (l *Location) isNull() bool {
        return l.x == 0 && l.y == 0 && l.z == 0
    }

type Vector0 struct {
    x, y, z int
}

type Transform struct {
    x, y, z [3]int
}
    func (t Transform) isNull() bool {
        return t == Transform{[3]int{0, 0, 0}, [3]int{0, 0, 0}, [3]int{0, 0, 0}}
    }

type Space struct {
    beacons map[Location]int
}
    func (space *Space) init() {
        space.beacons = make(map[Location]int)
    }
    func (space *Space) savePoint(l Location) {
        if _, ok := space.beacons[l]; ok {
            space.beacons[l]++
        } else {
            space.beacons[l] = 1
        }
    }
    func (space *Space) print() {
        fmt.Printf("\r\n\r\n---\r\nFinal space (%v points):\r\n", len(space.beacons))
        for loc, count := range space.beacons {
            fmt.Printf("{%v, %v, %v}: %v\r\n", loc.x, loc.y, loc.z, count)
        }
    }

type Scanner struct {
    num int
    beaconSignals []Location
    location Location
    transform Transform
    determined bool
    tris     map[[3]int]Trinagle
}
    func (s *Scanner) print() {
        fmt.Printf("\r\n--- scanner %v ---\r\n", s.num)
        for _, beacon := range s.beaconSignals {
            fmt.Printf("%4d, %4d, %4d\r\n", beacon.x, beacon.y, beacon.z)
        }
        fmt.Print("-- triangles: \r\n")
        for beacons, triangle := range s.tris {
            fmt.Printf("%v %v\r\n", beacons, triangle)
        }
        fmt.Printf("\r\n")
    }
    func (s *Scanner) makeTriangles() {
        s.tris = make(map[[3]int]Trinagle)
        beaconsCount := len(s.beaconSignals)
        for i1 := 0; i1 < beaconsCount; i1++ {
            for i2 := 0; i2 < beaconsCount; i2++ {
                if i1 == i2 {
                    continue
                }
                for i3 := 0; i3 < beaconsCount; i3++ {
                    if i3 == i1 || i3 == i2 || i1 == i2 {
                        continue
                    }
                    triangle, trinitySorted := handleTrinity(s, i1, i2, i3)
                    if _, ok := s.tris[trinitySorted]; ok {
                        continue
                    }
                    s.tris[trinitySorted] = triangle
                }
            }
        }
    }



/*
Receives three beacons from scanner.
Creates Triangle from that points and sorts them in a certain way
*/
func handleTrinity(s *Scanner, i1, i2, i3 int) (trinagle Trinagle, trinitySorted [3]int) {
    trinitySorted = [3]int{-1, -1, -1}
    trinagle.sides[0] = vRange(s.beaconSignals[i1], s.beaconSignals[i2])
    trinagle.sides[1] = vRange(s.beaconSignals[i1], s.beaconSignals[i3])
    trinagle.sides[2] = vRange(s.beaconSignals[i2], s.beaconSignals[i3])
    switch {
    case trinagle.sides[0] >= trinagle.sides[1] && trinagle.sides[0] >= trinagle.sides[2]:
        trinitySorted[0] = i3
        if trinagle.sides[1] >= trinagle.sides[2] {
            trinitySorted[1] = i2
            trinitySorted[2] = i1
        } else {
            trinitySorted[1] = i1
            trinitySorted[2] = i2
        }
    case trinagle.sides[1] >= trinagle.sides[0] && trinagle.sides[1] >= trinagle.sides[2]:
        trinitySorted[0] = i2
        if trinagle.sides[0] >= trinagle.sides[2] {
            trinitySorted[1] = i3
            trinitySorted[2] = i1
        } else {
            trinitySorted[1] = i1
            trinitySorted[2] = i3
        }
    case trinagle.sides[2] >= trinagle.sides[0] && trinagle.sides[2] >= trinagle.sides[1]:
        trinitySorted[0] = i1
        if trinagle.sides[0] >= trinagle.sides[1] {
            trinitySorted[1] = i3
            trinitySorted[2] = i2
        } else {
            trinitySorted[1] = i2
            trinitySorted[2] = i3
        }
    }
    sort.Float64s(trinagle.sides[:])
    return
}

func applyTranform(v Vector0, t Transform) (newV Vector0) {
    newV.x = t.x[0] * v.x + t.x[1] * v.y + t.x[2] * v.z
    newV.y = t.y[0] * v.x + t.y[1] * v.y + t.y[2] * v.z
    newV.z = t.z[0] * v.x + t.z[1] * v.y + t.z[2] * v.z
    return
}

func applyTranformPoint(l Location, t Transform) (res Location) {
    v0 := applyTranform(makeVector0(Location{0, 0, 0}, l), t)
    res.x = v0.x
    res.y = v0.y
    res.z = v0.z
    return
}


func vRange(l1, l2 Location) float64 {
    return math.Pow(math.Pow(float64(l1.x - l2.x), 2) + math.Pow(float64(l1.y - l2.y), 2) + math.Pow(float64(l1.z - l2.z), 2), 0.5)
}


func findScannerTransform(s1, s2 *Scanner) {
    //fmt.Printf("\r\n\r\n---\r\nWorking with scanners %v and %v\r\n", s1.num, s2.num)

    var pointsMapFrequency map[[2]int]int

    pointsMapFrequency = make(map[[2]int]int)

    for trinity1, trinagle1 := range s1.tris {
        for trinity2, trinagle2 := range s2.tris {
            if trinagle1 == trinagle2 {
                //fmt.Printf("trinities %v and %v forms equal %v\r\n", trinity1, trinity2, trinagle1)
                for i := 0; i < 3; i++ {
                    pmfKey := [2]int{trinity1[i], trinity2[i]}
                    if _, ok := pointsMapFrequency[pmfKey]; ok {
                        pointsMapFrequency[pmfKey]++
                    } else {
                        pointsMapFrequency[pmfKey] = 1
                    }
                }
            }
        }
    }

    if len(pointsMapFrequency) == 0 {
        fmt.Printf("No common triangles found in %v and %v\r\n", s1.num, s2.num)
        return
    }

    var samePoints [][2]int
    for points, _ := range pointsMapFrequency {
        samePoints = append(samePoints, [2]int{points[0], points[1]})
        //fmt.Printf("%v -> %v\r\n", s1.beaconSignals[points[0]], s2.beaconSignals[points[1]])
    }
    //fmt.Println(samePoints)
    if len(samePoints) < 12 {
        fmt.Printf("Too few common points found in %v and %v\r\n", s1.num, s2.num)
        return
    }

    // Find rotation and facing
    // For all vectors0 try all transforms until the one will last
    currentTransforms := allTransforms[:]
    for i := 1; i < len(samePoints); i++ {
        v1 := makeVector0(s1.beaconSignals[samePoints[0][0]], s1.beaconSignals[samePoints[i][0]])
        v2 := makeVector0(s2.beaconSignals[samePoints[0][1]], s2.beaconSignals[samePoints[i][1]])
        tmpTransforms := make([]Transform, 0, len(currentTransforms))
        for t := 0; t < len(currentTransforms); t++ {
            //fmt.Printf("Comparing %v and %v\r\n", v1, applyTranform(v2, currentTransforms[t]))
            if v1 == applyTranform(v2, currentTransforms[t]) {
                //fmt.Printf("Transform %v is ok\r\n", currentTransforms[t])
                tmpTransforms = append(tmpTransforms, currentTransforms[t])
                continue
            }
        }
        currentTransforms = tmpTransforms
        if len(currentTransforms) == 1 {
            break
        }
    }
    //fmt.Printf("Transforms: %v\r\n", currentTransforms)
    if len(currentTransforms) == 0 {
        fmt.Printf("No transforms found for %v and %v\r\n", s1.num, s2.num)
        return
    }
    if len(currentTransforms) > 1 {
        fmt.Printf("More than one transforms found for %v and %v: %v\r\n", s1.num, s2.num, currentTransforms)
        return
    }

    // save the transform found
    s2.transform = currentTransforms[0]

    // find s2 shift from s1
    p1 := s1.beaconSignals[samePoints[2][0]]
    p2 := applyTranformPoint(s2.beaconSignals[samePoints[2][1]], currentTransforms[0])
    s2shift := Location{p1.x - p2.x, p1.y - p2.y, p1.z - p2.z}
    s2.location = s2shift
    //fmt.Printf("Scanner %v location is %v\r\n", s2.num, s2.location)

    // save points to the final space
    // and fix all s2 beacons
    for i := 0; i < len(s2.beaconSignals); i++ {
        s2.beaconSignals[i] = applyTranformPoint(s2.beaconSignals[i], currentTransforms[0])
        s2.beaconSignals[i].x += s2shift.x
        s2.beaconSignals[i].y += s2shift.y
        s2.beaconSignals[i].z += s2shift.z
        space.savePoint(s2.beaconSignals[i])
        //fmt.Printf("%v ", s2.beaconSignals[i])
    }
    s2.makeTriangles()
    s2.determined = true

    return
}

func makeVector0(l1, l2 Location) (res Vector0) {
    res.x = l2.x - l1.x
    res.y = l2.y - l1.y
    res.z = l2.z - l1.z
    return
}

func abs(n int) int {
    if n >= 0 {
        return n
    }
    return n * -1
}

func manhattanDistance(l1, l2 Location) (md int) {
    md = abs(l2.x - l1.x) + abs(l2.y - l1.y) + abs(l2.z - l1.z)
    return
}

func maxManhattanDistance() int {
    mdMax := 0
    for i := 0; i < len(scanners) - 1; i++ {
        for ii := i + 1; ii < len(scanners); ii++ {
            if i == ii {
                continue
            }
            md := manhattanDistance(scanners[i].location, scanners[ii].location)
            if mdMax < md {
                mdMax = md
            }
        }
    }
    return mdMax
}

/*
 * Global variables
 */
var scanners []Scanner
var space Space
var allTransforms = [24]Transform{
    Transform{[3]int{1, 0, 0}, [3]int{0, 1, 0}, [3]int{0, 0, 1}},
    Transform{[3]int{1, 0, 0}, [3]int{0, 0, 1}, [3]int{0, -1, 0}},
    Transform{[3]int{1, 0, 0}, [3]int{0, -1, 0}, [3]int{0, 0, -1}},
    Transform{[3]int{1, 0, 0}, [3]int{0, 0, -1}, [3]int{0, 1, 0}},
    Transform{[3]int{-1, 0, 0}, [3]int{0, 1, 0}, [3]int{0, 0, -1}},
    Transform{[3]int{-1, 0, 0}, [3]int{0, 0, -1}, [3]int{0, -1, 0}},
    Transform{[3]int{-1, 0, 0}, [3]int{0, -1, 0}, [3]int{0, 0, 1}},
    Transform{[3]int{-1, 0, 0}, [3]int{0, 0, 1}, [3]int{0, 1, 0}},
    Transform{[3]int{0, 1, 0}, [3]int{-1, 0, 0}, [3]int{0, 0, 1}},
    Transform{[3]int{0, 1, 0}, [3]int{0, 0, -1}, [3]int{-1, 0, 0}},
    Transform{[3]int{0, 1, 0}, [3]int{1, 0, 0}, [3]int{0, 0, -1}},
    Transform{[3]int{0, 1, 0}, [3]int{0, 0, 1}, [3]int{1, 0, 0}},
    Transform{[3]int{0, -1, 0}, [3]int{1, 0, 0}, [3]int{0, 0, 1}},
    Transform{[3]int{0, -1, 0}, [3]int{0, 0, 1}, [3]int{-1, 0, 0}},
    Transform{[3]int{0, -1, 0}, [3]int{-1, 0, 0}, [3]int{0, 0, -1}},
    Transform{[3]int{0, -1, 0}, [3]int{0, 0, -1}, [3]int{1, 0, 0}},
    Transform{[3]int{0, 0, 1}, [3]int{0, 1, 0}, [3]int{-1, 0, 0}},
    Transform{[3]int{0, 0, 1}, [3]int{-1, 0, 0}, [3]int{0, -1, 0}},
    Transform{[3]int{0, 0, 1}, [3]int{0, -1, 0}, [3]int{1, 0, 0}},
    Transform{[3]int{0, 0, 1}, [3]int{1, 0, 0}, [3]int{0, 1, 0}},
    Transform{[3]int{0, 0, -1}, [3]int{0, -1, 0}, [3]int{-1, 0, 0}},
    Transform{[3]int{0, 0, -1}, [3]int{1, 0, 0}, [3]int{0, -1, 0}},
    Transform{[3]int{0, 0, -1}, [3]int{0, 1, 0}, [3]int{1, 0, 0}},
    Transform{[3]int{0, 0, -1}, [3]int{-1, 0, 0}, [3]int{0, 1, 0}},
}

/*
 * Entry
 */
func main() {

    scanners = readInput()
    totalScanners := len(scanners)
    space.init()

    // it means, scanner 0 has no transforms
    scanners[0].transform = allTransforms[0]
    scanners[0].location = Location{}
    for i := 0; i < len(scanners[0].beaconSignals); i++ {
        space.savePoint(scanners[0].beaconSignals[i])
    }
    scanners[0].determined = true

    for i := 0; i < totalScanners; i++ {
        scanners[i].makeTriangles()
    }

    pairs := make([]bool, totalScanners * totalScanners)
    for ;; {
        changed := false
        for i := 0; i < totalScanners; i++ {
            for ii := 0; ii < totalScanners; ii++ {
                if i == ii {
                    continue
                }
                if !scanners[i].determined {
                    continue
                }
                if pairs[i * totalScanners + ii] {
                    continue
                }
                pairs[i * totalScanners + ii] = true
                pairs[ii * totalScanners + i] = true
                if scanners[ii].determined {
                    continue
                }
                findScannerTransform(&scanners[i], &scanners[ii])
                changed = true
            }
        }
        if !changed {
            break
        }
    }
    
    space.print()

    fmt.Printf("Task 2: %v", maxManhattanDistance())    
}

/*
 * Input file reading functions
 */
func readInput() (scannersRead []Scanner) {
    file := openFile("./input.txt")
    fileScanner := bufio.NewScanner(file)

    var (
        currentScanner Scanner
        location       Location
        coords         []string
    )

    for fileScanner.Scan() {
        str := fileScanner.Text()
        if str == "" {
            scannersRead = append(scannersRead, currentScanner)
            currentScanner = Scanner{}
            continue
        }
        if string(str[1]) == "-" {
            currentScanner.num = convInt(strings.Split(str, " ")[2])
            continue
        }
        coords = strings.Split(str, ",")
        location.x = convInt(coords[0])
        location.y = convInt(coords[1])
        location.z = convInt(coords[2])
        currentScanner.beaconSignals = append(currentScanner.beaconSignals, location)
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

func convInt(s string) int {
    value, err := strconv.Atoi(s)
    if err != nil {
        fmt.Println("Bad number given: %d", s)
        panic("Fatal error when reading input file")
    }
    return value
}
