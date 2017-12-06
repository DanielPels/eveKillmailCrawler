package zDatabase

type Ship struct {
	Id     int           `json:"Id"`
	Solo   int           `json:"Solo"`
	Killer int           `json:"Killer"`
	Loss   int           `json:"Loss"`
	Fitted map[int]*Item `json:"Fitted"` //WHERE INT IS TYPE Id
	Cargo  map[int]*Item `json:"Cargo"`  //WHERE INT IS TYPE Id
}

func (s *Ship) addToCargo(item int, quantity int) {
	//check of Item wel bestaat
	checkAndCreateItem(item, s.Cargo)
	s.Cargo[item].Count++
	s.Cargo[item].Quantity += quantity
}

func (s *Ship) addToFitted(item int, quantity int) {
	//check of Item wel bestaat
	checkAndCreateItem(item, s.Fitted)
	s.Fitted[item].Count++
	s.Fitted[item].Quantity += quantity
}

func checkAndCreateItem(typeID int, items map[int]*Item) {
	if items[typeID] == nil {
		items[typeID] = &Item{
			Id:       typeID,
			Quantity: 0,
			Count:    0,
		}
	}
}
