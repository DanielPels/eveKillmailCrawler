package staticData

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type CategoryIDs struct {
	CategoryID int    `json:"categoryID"`
	Name       string `json:"name"`
}

type TypeIDs struct {
	TypeID  int    `json:"typeID"`
	Name    string `json:"name"`
	GroupID int    `json:"groupID"`
}

type GroupIDs struct {
	GroupID    int    `json:"groupID"`
	Name       string `json:"name"`
	CategoryID int    `json:"categoryID"`
}

var typeIDs map[int]TypeIDs
var groupIDs map[int]GroupIDs
var categoryIDs map[int]CategoryIDs

func Init(typeIDPath string, groupIDPath string, categoryIDPath string) {
	if checkIfFileExists(typeIDPath) && checkIfFileExists(groupIDPath) && checkIfFileExists(categoryIDPath) {

		fmt.Println("Starting loading static data")

		typeIDs = make(map[int]TypeIDs)
		groupIDs = make(map[int]GroupIDs)
		categoryIDs = make(map[int]CategoryIDs)

		typeIDSlice := make([]TypeIDs, 0)
		json.Unmarshal(getFile(typeIDPath), &typeIDSlice)
		for _, value := range typeIDSlice {
			typeIDs[value.TypeID] = value
		}

		groupIDSlice := make([]GroupIDs, 0)
		json.Unmarshal(getFile(groupIDPath), &groupIDSlice)
		for _, value := range groupIDSlice {
			groupIDs[value.GroupID] = value
		}

		categoryIDSlice := make([]CategoryIDs, 0)
		json.Unmarshal(getFile(categoryIDPath), &categoryIDSlice)
		for _, value := range categoryIDSlice {
			categoryIDs[value.CategoryID] = value
		}

		fmt.Println("Finished loading static data")
	}
}

func checkIfFileExists(fileName string) bool {
	//check of de file er is
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return true
	}
	return false
}

func getFile(fileName string) []byte {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return dat
}

//region Names
func GetTypeIDName(id int) string {
	return typeIDs[id].Name
}

func GetTypeIDGroupName(id int) string {
	return groupIDs[typeIDs[id].GroupID].Name
}

func GetTypeIDCategoryName(id int) string {
	return categoryIDs[groupIDs[typeIDs[id].GroupID].CategoryID].Name
}

//endregion

func GetGroupIDFromTypeID(id int) int {
	return groupIDs[typeIDs[id].GroupID].GroupID
}

func GetCategoryIDFromTypeID(id int) int {
	return categoryIDs[groupIDs[typeIDs[id].GroupID].CategoryID].CategoryID
}
