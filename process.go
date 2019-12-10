package main

import (
	"errors"
	"fmt"
	"strings"
	"os"
	"encoding/hex"
	"io/ioutil"
	"crypto/md5"
)

type Format struct {
	txt  string
	res string
}

func ReadDir(dirName string) ([]string, error) {
	var files []string
	f, err := ioutil.ReadDir(dirName)
	if err != nil {
	 return files, err
	}
	for _, file := range f {
		files = append(files, file.Name())
	}
	return files, nil
 }

func getHashHelper(filePath string)(chan string){
	r := make(chan string)
	var MD5String string
	go func() {
		defer close(r)
		file, err := ioutil.ReadFile(filePath)
		if err != nil {exitWithErr("error hash file")}
		sum := md5.Sum(file)
		MD5String = hex.EncodeToString(sum[:])
		r <- MD5String
	}()
	return r
}

func writeFileHelper(filePath string, data []byte)(chan string){
	r := make(chan string)
	go func() {
		defer close(r)
		err := ioutil.WriteFile(filePath, data, 0644)
		if err != nil {exitWithErr("error write file")}
		r <- "ok"
	}()
	return r
}

func separateInputsFiles(inputFiles []string)([]string) {
	s:= []string{}
	for _, fileName := range inputFiles {
		s= append(s, strings.Split(fileName,".")[0])
	} 
	return s
}

func exitWithErr(message string)(){
	fmt.Println(errors.New(message))
	os.Exit(1)
}

func getHash(dirInput string, inputFiles []string, format string)  []string {
	hashes:= []string{}
	var chans = []chan string{}
	for i := 0; i < len(inputFiles); i++ {
		chans = append(chans, getHashHelper(dirInput +"/"+inputFiles[i]+format))
	}
	for i := range chans {
		hashes= append(hashes,<-chans[i])
	}
	return hashes
}

func writeFile(dirOutput string, inputFiles []string, format string, hashes []string)  string {
	var chans = []chan string{}
	for i := 0; i < len(hashes); i++ {
		chans = append(chans, writeFileHelper(dirOutput +"/"+inputFiles[i]+format, []byte(hashes[i])))
	}
	for i := range chans {
		if <-chans[i]!= "ok" {exitWithErr("error write file")}
	}
	return "ok"
}


func main() {
	if len(os.Args[0:]) <= 2 {exitWithErr("must be addition args")}
	format := Format{
		txt:  ".txt",
		res: ".res",
	}
	dirInput, dirOutput := os.Args[1],os.Args[2]
	files, err := ReadDir(dirInput)
	if err != nil {exitWithErr("error read dir")}
	inputFiles:=	separateInputsFiles(files)
	hashes:= getHash(dirInput,inputFiles, format.txt)
	fmt.Printf("Total number of processed files: %d",len(hashes))
	writeFile(dirOutput,inputFiles, format.res, hashes)
}
