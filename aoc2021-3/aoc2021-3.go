package main

import (
    "fmt"
    "os"
    "bufio"
)

// global: length of each measure
var lengthBits int;

//
func main() {
    
    measures := readMeasures()

    task1 := rateGammaEpsilon(measures)
    fmt.Printf("Task 1: %d\r\n", task1)

    task2  := rateOxygen(measures) * rateCO2(measures)
    fmt.Printf("Task 2: %d\r\n", task2)

}

/*
 * Task 1
 * Gamma & Epsilon rates
 */
func rateGammaEpsilon(measures []int) int {
    var gamma, epsilon int

    for bitNum := 0; bitNum < lengthBits; bitNum++ {
        mostBit := checkMostBit(&measures, bitNum)
        if (mostBit == 1) {
            gamma += 1 << (lengthBits - bitNum - 1)
        } else {
            epsilon += 1 << (lengthBits - bitNum - 1)
        }
    }    

    return gamma * epsilon
}

func checkMostBit(measures *[]int, bitNum int) int {
    sum := 0
    for _, m := range *measures {
        if checkBit(m, bitNum) == 1 {
            sum++
        } else {
            sum--
        }
    }
    if (sum >= 0) {
        return 1
    }
    return 0
}

/*
 * Task 2
 */

// Oxygen generator rate
func rateOxygen(measures []int) int {
    var oxygen int

    var lookIn = &measures;

    for bitNum := 0; bitNum < lengthBits; bitNum++ {
        ones, zeroes := filterByBit(lookIn, bitNum)
        if len(ones) >= len(zeroes) {
            oxygen += 1 << (lengthBits - bitNum - 1)
            lookIn = &ones
        } else {
            lookIn = &zeroes
        }
        if len(*lookIn) == 1 {
            // this would be the same as final result if we'd continue to go through the rest of iterations
            last := *lookIn
            return last[0]
        }
    }
    return oxygen
}

// CO2 generator rate
func rateCO2(measures []int) int {
    var co2 int

    var lookIn = &measures;

    for bitNum := 0; bitNum < lengthBits; bitNum++ {
        ones, zeroes := filterByBit(lookIn, bitNum)
        if len(zeroes) > len(ones) {
            co2 += 1 << (lengthBits - bitNum - 1)
            lookIn = &ones
        } else {
            lookIn = &zeroes
        }
        if len(*lookIn) == 1 {
            // this would be the same as final result if we'd continue to go through the rest of iterations
            last := *lookIn
            return last[0]
        }
    }
    return co2
}

func filterByBit(measures *[]int, bitNum int) ([]int, []int) {
    var ones, zeroes []int
    for _, m := range *measures {
        if checkBit(m, bitNum) == 1 {
            ones = append(ones, m)
        } else {
            zeroes = append(zeroes, m)
        }
    }
    return ones, zeroes
}

/*
 * Common functions
 */

func checkBit(number int, bitNum int) int {
    mask := 1 << (lengthBits - bitNum - 1)
    if number & mask == 0 {
        return 0
    }
    return 1
}


// read all measures
func readMeasures() []int {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    var measures []int

    for scanner.Scan() {
        m, err := readBinary(scanner.Bytes())
        if m == -1 {
            panic(err)
        }
        measures = append(measures, m)
    }

    return measures
}

// read one binary
func readBinary(data []uint8) (int, error) {
    var res int
    var len = len(data)
    if lengthBits == 0 {
        lengthBits = len
    } else if len != lengthBits {
        return -1, fmt.Errorf("Incorrect length of measure data %d", data)
    }
    for i := 0; i < len; i++ {
        switch data[i] {
        case 49:
            res += 1 << (len - i - 1)
        case 48:
            break
        default:
            return -1, fmt.Errorf("Incorrect measure data %d", data)
        }
    }
    return res, nil
}

// file opening routine
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}