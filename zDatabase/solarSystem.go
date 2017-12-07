package zDatabase

import (
	"eveKillmailCrawler/killmail"
)

//faction Id:
//Amarr = 500003
//Minmatar = 500002

const Amarr = 500003
const Minmatar = 500002

type SolarSystem struct {
	Id          int           `json:"Id"`
	KillmailIDs []int         `json:"KillmailIDs"`
	Amarr       map[int]*Ship `json:"Amarr"`    //WHERE INT IS TYPE Id
	Minmatar    map[int]*Ship `json:"Minmatar"` //WHERE INT IS TYPE Id
	Neutral     map[int]*Ship `json:"Neutral"`  //WHERE INT IS TYPE Id
}

func (s *SolarSystem) AddVictim(killmail *killmail.ZKillmail) {
	//loop door alle zKillmail Id's en check of die al is geweest is zo jah quit functie
	for _, killmailID := range s.KillmailIDs {
		if killmailID == killmail.KillmailID {
			return
		}
	}

	//voeg de zKillmail Id toe aan de lijst
	s.KillmailIDs = append(s.KillmailIDs, killmail.KillmailID)

	switch killmail.Victim.FactionID {
	case Amarr:
		s.addVictimToFaction(s.Amarr, killmail)
	case Minmatar:
		s.addVictimToFaction(s.Minmatar, killmail)
	default:
		s.addVictimToFaction(s.Neutral, killmail)
	}
}

func (s *SolarSystem) addVictimToFaction(faction map[int]*Ship, killmail *killmail.ZKillmail) {
	//check of de struct al bestaat, van het schip
	checkAndCreateShip(killmail.Victim.ShipTypeID, faction)

	//haal de struct op van het schip
	var victimShip *Ship = faction[killmail.Victim.ShipTypeID]

	// :( its ded
	victimShip.Loss++

	//loop door alle items heen en voeg deze toe
	//een map op te tracken welke items als geweest zijn
	addedItems := make(map[int]bool)

	for _, item := range killmail.Victim.Items {
		//check of de Item er al is in geweest
		if addedItems[item.ItemTypeID] {
			continue
		}
		addedItems[item.ItemTypeID] = true

		//hoeveel is er gedropt
		quantity := 0

		if item.QuantityDestroyed != nil {
			quantity += *item.QuantityDestroyed
		}
		if item.QuantityDropped != nil {
			quantity += *item.QuantityDropped
		}

		switch item.Flag {
		case 5:
			victimShip.addToCargo(item.ItemTypeID, quantity)
		default:
			victimShip.addToFitted(item.ItemTypeID, quantity)
		}
	}
}

func (s *SolarSystem) AddAttacker(killmail *killmail.ZKillmail) {
	for _, attacker := range killmail.Attackers {

		//als het een npc is dan skip
		if attacker.CharacterID == nil {
			continue
		}

		//faction Id staat niet altijd in json, dus is deze op zijn default value(0 in dit geval)
		switch attacker.FactionID {
		case Amarr:
			s.addAttackerToFaction(s.Amarr, attacker.ShipTypeID)
		case Minmatar:
			s.addAttackerToFaction(s.Minmatar, attacker.ShipTypeID)
		default:
			s.addAttackerToFaction(s.Neutral, attacker.ShipTypeID)
		}
	}
}

func (s *SolarSystem) addAttackerToFaction(faction map[int]*Ship, shipID int) {
	checkAndCreateShip(shipID, faction)
	faction[shipID].Killer++
}

func checkAndCreateShip(typeID int, ships map[int]*Ship) {
	if ships[typeID] == nil {
		ships[typeID] = &Ship{
			Killer: 0,
			Id:     typeID,
			Loss:   0,
			Solo:   0,
			Fitted: make(map[int]*Item),
			Cargo:  make(map[int]*Item),
		}
	}
}
