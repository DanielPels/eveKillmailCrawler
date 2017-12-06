package zDatabase

import (
	"encoding/json"
	"eveKillmailCrawler/killmail"
)

type ZDatabase struct {
	SolarSystems map[int]*SolarSystem `json:"SolarSystems"`
}

func (d *ZDatabase) AddSolarSystem(id int) {
	if d.SolarSystems[id] == nil {
		d.SolarSystems[id] = &SolarSystem{
			Id:          id,
			KillmailIDs: make([]int, 0),
			Minmatar:    make(map[int]*Ship),
			Amarr:       make(map[int]*Ship),
			Neutral:     make(map[int]*Ship),
		}
	}
}

func (d ZDatabase) AddKillmails(solarSystemID int, kill killmail.ZKillmail) {
	d.AddSolarSystem(solarSystemID)
	d.SolarSystems[solarSystemID].AddVictim(&kill)
	d.SolarSystems[solarSystemID].AddAttacker(&kill)
}

func (d ZDatabase) ExportToJson() []byte {
	b, _ := json.Marshal(d)
	return b
}

func (d *ZDatabase) ImportFromJson(b []byte) {
	json.Unmarshal(b, d)
}

func New() *ZDatabase {
	return &ZDatabase{SolarSystems: make(map[int]*SolarSystem)}
}
