package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"
)

//Нужно отсортировать большой файл с ограничением на RAM
var isGenerating bool
var genfileSize int

var inFile string
var outFile string
var ramsize int

var cyan func(a ...interface{}) string
var succes func(a ...interface{}) string
var red func(a ...interface{}) string

func init() {
	flag.IntVar(&genfileSize, "gsize", 100, "En: How much numbers will be in file \nru: Как много цифр будет в файле")
	flag.IntVar(&ramsize, "rsize", 4096, "en: RAM Size \nru: Размер оперативной памяти")
	flag.BoolVar(&isGenerating, "gen", false, "en: Is Generating file with numbers\nru: Флаг для генерации файла с тестовыми цифрами")
	flag.StringVar(&inFile, "ifile", "a.txt", "en: Path to input file\nru: Путь к входному файлу")
	flag.StringVar(&outFile, "ofile", "b.txt", "en: Path to output file\nru: Путь к выходному файлу")
	cyan = color.New(color.FgCyan).SprintFunc()
	succes = color.New(color.FgGreen).SprintFunc()
	red = color.New(color.FgRed).SprintFunc()
}

func generateFile() {
	file, err := os.OpenFile(inFile, os.O_RDWR+os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	rand.Seed(time.Now().Unix())
	for i := 0; i < genfileSize; i++ {
		num := rand.Int()
		s := strconv.Itoa(num)
		file.WriteString(s + "\n")
	}
}

func unoptimizedSort() {
	file, err := os.OpenFile(inFile, os.O_RDWR+os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// The bufio.ScanLines is used as an
	// input to the method bufio.Scanner.Split()
	// and then the scanning forwards to each
	// new line using the bufio.Scanner.Scan()
	// method.
	scanner.Split(bufio.ScanLines)
	var nums []int

	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Printf("Error: %v not a number", red(scanner.Text()))
			continue
		}
		nums = append(nums, num)
	}
	sort.Ints(nums)
	ofile, err := os.OpenFile(outFile, os.O_RDWR+os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer ofile.Close()
	for _, num := range nums {
		ofile.WriteString(strconv.Itoa(num) + "\n")
	}
}

func main() {
	flag.Parse()
	// if isHelp {
	// 	log.Println("Help")
	// }

	fmt.Printf("Flags:\n\tGenerate mod: %v \n", cyan(isGenerating))
	if isGenerating {
		fmt.Printf("\tNumbers size: %v\n\tOutfile: %v \n", cyan(genfileSize), cyan(inFile))
		generateFile()
		fmt.Println(succes("Succes"))
	} else {
		start := time.Now()
		fmt.Printf("\tInput file: %v\n\tOutputfile: %v \n\tRAM size: %v \n", cyan(inFile), cyan(outFile), cyan(ramsize))
		unoptimizedSort()
		fmt.Println(succes("Succes"))
		fmt.Printf("Time passed: %v \n", cyan(time.Since(start)))
	}

}
