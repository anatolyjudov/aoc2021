package main

import (
    "fmt"
    "os"
    "bufio"
)

type Pattern [9]bool
    func (p *Pattern) toInt() (res int) {
        for i := 0; i < 9; i++ {
            if p[i] {
                res += (1 << i)
            }
        }
        return
    }
    func (p *Pattern) fromInt(n int) {
        if n >= 512 {
            panic("Too big number to convert to Pattern")
        }
        for i := 8; i >= 0; i-- {
            if n >= (1 << i) {
                n = n - (1 << i)
                p[i] = true
            } else {
                p[i] = false
            }
        }
    }
    func (p *Pattern) print() {
        for i := 8; i >= 0; i-- {
            if p[i] {
                fmt.Print("#")
            } else {
                fmt.Print(".")
            }            
        }
        fmt.Print("\r\n")
    }

type Enchancer struct {
    mapData string
    enchMap [512]bool
}
    func (e *Enchancer) set(pos int, val byte) {
        if string(val) == "." {
            e.enchMap[pos] = false
        } else if string(val) == "#" {
            e.enchMap[pos] = true
        } else {
            panic("Bad value given for enchMap")
        }        
    }
    func (e *Enchancer) get(pattern Pattern) (res bool) {

        return e.enchMap[pattern.toInt()]
    }
    func (e *Enchancer) getByInt(val int) (res bool) {
        return e.enchMap[val]
    }
    func (e *Enchancer) print() {
        fmt.Print("Image enchancement algorithm:\r\n")
        for i := 0; i < 512; i++ {
            if e.enchMap[i] {
                fmt.Print("#")
            } else {
                fmt.Print(".")
            }            
        }
        fmt.Print("\r\n")
    }

type Image struct {
    sizeX, sizeY int
    pixels []bool
    outerPixel bool
}
    func (i *Image) readStr(x, y int, data string) {
        i.sizeX = x
        i.sizeY = y
        i.pixels = make([]bool, x * y)
        for y := 0; y < i.sizeY; y++ {
            for x := 0; x < i.sizeX; x++ {
                switch string(data[y * i.sizeX + x]) {
                case "#":
                    i.pixels[y * i.sizeX + x] = true
                case ".":
                    i.pixels[y * i.sizeX + x] = false
                default:
                    panic("Unknown symbol found when reading initial image")
                }
            }
        }
    }
    func (i *Image) print() {
        fmt.Printf("\r\nImage %d x %d\r\n", i.sizeX, i.sizeY)
        for y := 0; y < i.sizeY; y++ {
            for x := 0; x < i.sizeX; x++ {
                if i.pixels[y * i.sizeX + x] {
                    fmt.Print("#")
                } else {
                    fmt.Print(".")
                }
            }
            fmt.Print("\r\n")
        }
        if i.outerPixel {
            fmt.Print("Outer pixels are #\r\n")
        } else {
            fmt.Print("Outer pixels are .\r\n")
        }
        fmt.Print("\r\n")
    }
    func (i *Image) getPattern(xP, yP int) (p Pattern) {
        for y := -1; y <= 1; y++ {
            for x := -1; x <= 1; x++ {
                if (xP + x < 0) || (xP + x >= i.sizeX) || (yP + y < 0) || (yP + y >= i.sizeY) {
                    p[8 - (y + 1) * 3 - x - 1] = i.outerPixel
                } else {
                    p[8 - (y + 1) * 3 - x - 1] = i.pixels[(yP + y) * i.sizeX + (xP + x)]
                }
            }
        }
        return
    }
    func (i *Image) grow(factor int) {
        newSizeX := i.sizeX + factor * 2
        newSizeY := i.sizeY + factor * 2
        newPixels := make([]bool, newSizeX * newSizeY)

        // fill outer pixels
        for ii := 0; ii < len(newPixels); ii++ {
            newPixels[ii] = i.outerPixel
        }

        // copying old image
        for y := 0; y < i.sizeY; y++ {
            orig0 := y * i.sizeX
            orig1 := y * i.sizeX + i.sizeX
            dest0 := (y + factor) * newSizeX + factor
            dest1 := (y + factor) * newSizeX + i.sizeX + factor
            copy(newPixels[dest0:dest1], i.pixels[orig0:orig1])
        }

        // remembering new state
        i.sizeX = newSizeX
        i.sizeY = newSizeY
        i.pixels = newPixels
    }
    func (i *Image) enchance(enchancer *Enchancer) {
        i.grow(1)
        newPixels := make([]bool, i.sizeX * i.sizeY)
        for y := 0; y < i.sizeY; y++ {
            for x := 0; x < i.sizeX; x++ {
                newPixels[y * i.sizeX + x] = enchancer.get(i.getPattern(x, y))
            }
        }
        i.pixels = newPixels
        if i.outerPixel {
            if !enchancer.getByInt(511) {
                i.outerPixel = false
            }
        } else {
            if enchancer.getByInt(0) {
                i.outerPixel = true
            }
        }
    }
    func (i *Image) count() (count int) {
        for y := 0; y < i.sizeY; y++ {
            for x := 0; x < i.sizeX; x++ {
                if i.pixels[y * i.sizeX + x] {
                    count++
                }
            }
        }
        return
    }

/*
 * Entry
 */
func main() {
    enchancer, im := readInput()

    enchancer.print()
    im.print()

    for t := 0; t < 50; t++ {
        im.enchance(&enchancer)
    }
    fmt.Println(im.count())
}

/*
 * Input file reading functions
 */
func readInput() (enchancer Enchancer, im Image) {
    file := openFile("./input.txt")
    fileScanner := bufio.NewScanner(file)

    // read enchacement algorithm
    p := 0
    for fileScanner.Scan() {
        str := fileScanner.Text()
        if str == "" {
            break
        }
        for i := 0; i < len(str); i++ {
            enchancer.set(p, str[i])
            p++
        }
    }

    // read initial image
    var (
        xSize, ySize int
        data         string
    )
    for fileScanner.Scan() {
        str := fileScanner.Text()
        if str == "" {
            break
        }
        if xSize == 0 {
            xSize = len(str)
        } else if xSize != len(str) {
            panic("Different sized string found when reading initial image")
        }
        ySize++
        data += str
    }
    im.readStr(xSize, ySize, data)

    return
}
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}