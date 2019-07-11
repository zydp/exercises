package main

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

const (
	TEST_FILE_NAME = "pefinder.txt"
)

func main() {

	//lineByLine(TEST_FILE_NAME)
	//wordByWord(TEST_FILE_NAME)
	charByChar(TEST_FILE_NAME)
	ReadingFromRandom() /*for linux*/
	WriteingToFile()
	Copy(TEST_FILE_NAME, "dest.txt")
	CopyUsingIoutil(TEST_FILE_NAME, "dest.txt")
	AppendToFile(TEST_FILE_NAME)
}

/*逐行读取*/
func lineByLine(file string) error {
	var err error
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			break
		}
		fmt.Print(line)
	}
	return nil
}

/*逐个单词读取*/
func wordByWord(file string) error {
	var err error
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			return err
		}
		r := regexp.MustCompile("[^\\s]+")
		words := r.FindAllString(line, -1)
		for i := 0; i < len(words); i++ {
			fmt.Println(words[i])
		}
	}
	return nil
}

/*逐个字节读取*/
func charByChar(file string) error {
	var err error
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error reading file %s", err)
			return err
		}
		for _, x := range line {
			fmt.Println(string(x))
		}
	}
	return nil
}

/*读取指定大小*/
func readSize(f *os.File, size int) []byte {
	buffer := make([]byte, size)
	n, err := f.Read(buffer)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return buffer[0:n]
}

func ReadingCsvFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Cannot open", file)
		fmt.Println(err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	allRecords, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	xP := []float64{}
	yP := []float64{}
	for _, rec := range allRecords {
		x, _ := strconv.ParseFloat(rec[0], 64)
		y, _ := strconv.ParseFloat(rec[1], 64)
		xP = append(xP, x)
		yP = append(yP, y)
	}
	/*--*/
}

func WriteingToFile() {
	s := []byte("Data to write\n")
	/*Fprintf*/
	f1, err := os.Create("f1.txt")
	if err != nil {
		fmt.Println("Cannot create file", err)
		return
	}
	defer f1.Close()
	fmt.Fprintf(f1, string(s))

	/*WriteString*/
	f2, err := os.Create("f2.txt")
	if err != nil {
		fmt.Println("Cannot create file", err)
		return
	}
	defer f2.Close()
	n, err := f2.WriteString(string(s))
	fmt.Printf("wrote %d bytes\n", n)

	/*bufio.NewWriter*/
	f3, err := os.Create("f3.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	w := bufio.NewWriter(f3)
	n, err = w.WriteString(string(s))
	fmt.Printf("wrote %d bytesn", n)
	w.Flush()

	/*ioutil.WriteFile*/
	f4 := "f4.txt"
	err = ioutil.WriteFile(f4, s, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	/*io.WriteString*/
	f5, err := os.Create("f5.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	n, err = io.WriteString(f5, string(s))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("wrote %d bytes\n", n)
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}
	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CopyUsingIoutil(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		fmt.Println("Error creating the new file", dst)
		fmt.Println(err)
		os.Exit(1)
	}
	return err
}

func AppendToFile(filename string) (err error){
	f, err := os.OpenFile(filename,	os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer f.Close()
	fmt.Fprintf(f, "%s\n", "ABCDEFG")
	return err
}


func ReadingFromRandom() {
	f, err := os.Open("/dev/random")
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	var seed int64
	binary.Read(f, binary.LittleEndian, &seed)
	fmt.Println("Seed:", seed)
}
