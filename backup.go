package main

import (
	"os"
	"io/ioutil"
)

func checkForBackup(fileName string) bool {
	//check of de file er is
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return true
	}
	return false
}

func getBackup(fileName string) []byte {
	//pls get backup file
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return dat
}

func saveBackup(data []byte, fileName string) {
	//pls save backup file
	ioutil.WriteFile(fileName, data, 0644)
}
