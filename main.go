package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
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

type Anal struct {
	Real   bool
	Cool   string
	SoDeep Deep
}

type Deep struct {
	WOOOW bool
}

type SAnal struct {
	Real   bool
	Cool   string
	SoDeep SDeep
}

type SDeep struct {
	WOOOW bool
}

type Person struct {
	ID      int
	Name    string
	Surname string
	Age     uint8
	Jopa    string
	Anal    Anal
}

type SuperPerson struct {
	ID      int    `json:"id,omitempty"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     uint8  `json:"age"`
	Anal    SAnal  `json:"anal"`
}

func NewPerson() Person {
	return Person{
		ID:      rand.Int(),
		Name:    strconv.Itoa(rand.Int()),
		Surname: strconv.Itoa(rand.Int()),
		Age:     uint8(rand.Intn(100)),
		Jopa:    "Jopa",
		Anal: Anal{
			Real: true,
			Cool: "String",
			SoDeep: Deep{
				WOOOW: true,
			},
		},
	}
}

//i - input, r - pointer to result
func Reflect(input interface{}, result interface{}) {
	r := reflect.ValueOf(result).Elem()
	el := reflect.ValueOf(input)
	numfield := el.NumField()
	fields := make(map[string]reflect.Value)
	for i := 0; i < numfield; i++ {
		field := el.Field(i)
		fields[el.Type().Field(i).Name] = field
		// fmt.Printf("%#v\n", field)
	}
	for name, field := range fields {
		fmt.Printf("Name: %v, Value: %v, CanSet: %v\n", cyan(name), cyan(field), cyan(r.FieldByName(name).CanSet()))
		if r.FieldByName(name).CanSet() {
			if field.Kind() != reflect.Struct {
				r.FieldByName(name).Set(field)
			} else {
				t2 := r.FieldByName(name)
				ReflectStruct(field, t2)
			}
		}
	}
}

func ReflectStruct(input reflect.Value, result reflect.Value) {
	numfields := input.NumField()
	fields := make(map[string]reflect.Value)
	for a := 0; a < numfields; a++ {
		field := input.Field(a)
		fields[input.Type().Field(a).Name] = field
	}
	for name, field := range fields {
		fmt.Printf("Name: %v, Value: %v, CanSet: %v\n", cyan(name), cyan(field), cyan(result.FieldByName(name).CanSet()))
		if result.FieldByName(name).CanSet() {
			if field.Kind() != reflect.Struct {
				result.FieldByName(name).Set(field)
			} else {
				t2 := result.FieldByName(name)
				ReflectStruct(field, t2)
			}
		}
	}
}

func PreetyPrint(modifier string, inputs ...interface{}) {
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		in := reflect.ValueOf(input)

		if in.Kind() == reflect.Struct {
			numfields := in.NumField()
			fmt.Printf("%v%v%v", modifier, in.Type().Name(), red("{"))
			for nm := 0; nm < numfields; nm++ {
				field := in.Field(nm)
				if field.Kind() == reflect.Struct {
					// fmt.Print("\n")
					PreetyPrint(fmt.Sprintf("%v", modifier), field.Interface())
				} else {
					fmt.Printf("%v:%v, ", in.Type().Field(nm).Name, cyan(field.Interface()))
				}
			}
			fmt.Printf(red("}"))
		} else {
			fmt.Printf("%vT: %v; V: %v\n", modifier, succes(in.Type()), cyan(i))
		}
	}
}

func main() {
	flag.Parse()
	person := NewPerson()
	PreetyPrint("", 12)
	fmt.Println("")
	PreetyPrint("", person)
	fmt.Println()
	// sperson := &SuperPerson{}
	// Reflect(person, sperson)
	// fmt.Printf(cyan(fmt.Sprintf("SPerson: %#v\n", *sperson)))
	// if isHelp {
	// 	log.Println("Help")
	// }

	// fmt.Printf("Flags:\n\tGenerate mod: %v \n", cyan(isGenerating))
	// if isGenerating {
	// 	fmt.Printf("\tNumbers size: %v\n\tOutfile: %v \n", cyan(genfileSize), cyan(inFile))
	// 	generateFile()
	// 	fmt.Println(succes("Succes"))
	// } else {
	// 	start := time.Now()
	// 	fmt.Printf("\tInput file: %v\n\tOutputfile: %v \n\tRAM size: %v \n", cyan(inFile), cyan(outFile), cyan(ramsize))
	// 	unoptimizedSort()
	// 	fmt.Println(succes("Succes"))
	// 	fmt.Printf("Time passed: %v \n", cyan(time.Since(start)))
	// }

}
