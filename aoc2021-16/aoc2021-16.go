package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

type Bit byte   // :)

var hexMap = map[string][4]byte{
    "0" : {0, 0, 0, 0},
    "1" : {0, 0, 0, 1},
    "2" : {0, 0, 1, 0},
    "3" : {0, 0, 1, 1},
    "4" : {0, 1, 0, 0},
    "5" : {0, 1, 0, 1},
    "6" : {0, 1, 1, 0},
    "7" : {0, 1, 1, 1},
    "8" : {1, 0, 0, 0},
    "9" : {1, 0, 0, 1},
    "A" : {1, 0, 1, 0},
    "B" : {1, 0, 1, 1},
    "C" : {1, 1, 0, 0},
    "D" : {1, 1, 0, 1},
    "E" : {1, 1, 1, 0},
    "F" : {1, 1, 1, 1},
}

type Source struct {
    reader *strings.Reader
    pos int
}
    func (s *Source) loadString(str string) {
        var (
            binaryStr []byte
            builder strings.Builder
        )
        r := strings.NewReader(str)
        len := r.Len();
        for i := 0; i < len; i++ {
            byteRead, _ := r.ReadByte()
            bytesToWrite := hexMap[string([]byte{byteRead})]
            for _, byte := range bytesToWrite {
                builder.WriteByte(byte)
            }
        }
        binaryStr = []byte(builder.String())
        // fmt.Println(binaryStr)
        s.reader = strings.NewReader(string(binaryStr))
        s.pos = 0
    }
    func (s *Source) nextBit() Bit {
        byteRead, err := s.reader.ReadByte()
        if err != nil {
            panic("Can't read another byte")
        }
        s.pos++
        return Bit(byteRead)
    }
    func (s *Source) nextBits(amount int) (bytesRead []Bit) {
        if s.reader.Len() < amount {
            panic("Can't read so much, end of the source")
        }
        bytesRead = make([]Bit, amount)
        for i := 0; i < amount; i++ {
            bytesRead[i] = s.nextBit()
        }
        return
    }
    func (source *Source) readPacket() Packet {
        var newPacket Packet
        // read version
        newPacket.version = dec(source.nextBits(3))
        // read type
        newPacket.typeId = dec(source.nextBits(3))
        // call something according to the type
        if newPacket.typeId == 4 {
            newPacket.readLiteral(source)
        } else {
            newPacket.readOperator(source)
        }
        return newPacket
    }

type Message struct {
    rootPacket Packet
}
    func (message *Message) readFrom(source *Source) {
        message.rootPacket = source.readPacket()
    }
    func (message *Message) task1() int {
        return message.rootPacket.versionSum()
    }
    func (message *Message) task2() int {
        return message.rootPacket.value()
    }

type Packet struct {
    version, typeId int
    subPackets []Packet
    literalValue int
}
    func (p *Packet) readOperator(source *Source) {
        // read length type
        lengthType := source.nextBit()
        // there're two ways
        if lengthType == 0 {
            // total length
            totalLength := dec(source.nextBits(15))
            pos0 := source.pos
            for source.pos < pos0 + totalLength {
                packet := source.readPacket()
                p.subPackets = append(p.subPackets, packet)
            }
            if source.pos != pos0 + totalLength {
                panic("Operator subpacket's length doesn't match supposed length")
            }
        } else if lengthType == 1 {
            // number of subpackets
            numberOfSubpackets := dec(source.nextBits(11))
            for i := 0; i < numberOfSubpackets; i++ {
                p.subPackets = append(p.subPackets, source.readPacket())
            }
        } else {
            panic("Unknown lengthType")
        }
    }
    func (p *Packet) readLiteral(source *Source) {
        var literalBits []Bit
        for groupIsNotLast := 1; groupIsNotLast == 1; {
            groupIsNotLast = int(source.nextBit())
            literalBits = append(literalBits, source.nextBits(4)...)
        }
        p.literalValue = dec(literalBits)
    }
    func (p *Packet) versionSum() int {
        var sum int
        sum = p.version
        for _, packet := range p.subPackets {
            sum += packet.versionSum()
        }
        return sum
    }
    func (p * Packet) subPacketsValues() []int {
        var vals []int
        vals = make([]int, len(p.subPackets))
        for i := 0; i < len(p.subPackets); i++ {
            vals[i] = p.subPackets[i].value()
        }
        return vals
    }
    func (p * Packet) value() int {
        switch p.typeId {
            case 0:
                return oSum(p.subPacketsValues())
            case 1:
                return oProduct(p.subPacketsValues())
            case 2:
                return oMin(p.subPacketsValues())
            case 3:
                return oMax(p.subPacketsValues())
            case 5:
                return oGreaterThan(p.subPacketsValues())
            case 6:
                return oLessThan(p.subPacketsValues())
            case 7:
                return oEqual(p.subPacketsValues())
        }
        return p.literalValue
    }

func dec(bits []Bit) int {
    var length, decSum int
    length = len(bits)
    for i := 0; i < length; i++ {
        decSum += int(bits[i]) << (length - i - 1)
    }
    return decSum
}

func oSum(ls []int) int {
    sum := 0
    for i := 0; i < len(ls); i++ {
        sum += ls[i]
    }
    return sum
}
func oProduct(ls []int) int {
    res := ls[0]
    for i := 1; i < len(ls); i++ {
        res = res * ls[i]
    }
    return res
}
func oMin(ls []int) int {
    min := ls[0]
    for i := 1; i < len(ls); i++ {
        if min > ls[i] {
            min = ls[i]
        }
    }
    return min
}
func oMax(ls []int) int {
    max := ls[0]
    for i := 1; i < len(ls); i++ {
        if max < ls[i] {
            max = ls[i]
        }
    }
    return max
}
func oGreaterThan(ls []int) int {
    if ls[0] > ls[1] {
        return 1
    }
    return 0
}
func oLessThan(ls []int) int {
    if ls[0] < ls[1] {
        return 1
    }
    return 0    
}
func oEqual(ls []int) int {
    if ls[0] == ls[1] {
        return 1
    }
    return 0    
}

func main() {
    var source Source
    source.loadString(readInput())

    var message Message
    message.readFrom(&source)

    fmt.Println(message)

    task1 := message.task1()
    fmt.Printf("Task1: %d\r\n", task1)

    task2 := message.task2()
    fmt.Printf("Task2: %d\r\n", task2)
}

// read all input data
func readInput() string {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)
    scanner.Scan()
    return scanner.Text()
}

// file opening routine
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}
