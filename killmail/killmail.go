package killmail

import "time"

type ZKillmail struct {
	KillmailID   int       `json:"killmail_id"`
	KillmailTime time.Time `json:"killmail_time"`
	Victim struct {
		DamageTaken   int `json:"damage_taken"`
		ShipTypeID    int `json:"ship_type_id"`
		CharacterID   int `json:"character_id"`
		CorporationID int `json:"corporation_id"`
		FactionID     int `json:"faction_id"`
		Items []struct {
			ItemTypeID        int  `json:"item_type_id"`
			Singleton         int  `json:"singleton"`
			Flag              int  `json:"flag"`
			QuantityDestroyed *int `json:"quantity_destroyed,omitempty"`
			QuantityDropped   *int `json:"quantity_dropped,omitempty"`
		} `json:"items"`
		Position struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
			Z float64 `json:"z"`
		} `json:"position"`
	} `json:"victim"`
	Attackers []struct {
		SecurityStatus float64 `json:"security_status"`
		FinalBlow      bool    `json:"final_blow"`
		DamageDone     int     `json:"damage_done"`
		CharacterID    *int    `json:"character_id,omitempty"`
		CorporationID  *int    `json:"corporation_id,omitempty"`
		AllianceID     *int    `json:"alliance_id,omitempty"`
		FactionID      int     `json:"faction_id"`
		ShipTypeID     int     `json:"ship_type_id"`
		WeaponTypeID   *int    `json:"weapon_type_id,omitempty"`
	} `json:"attackers"`
	SolarSystemID int `json:"solar_system_id"`
	Zkb struct {
		LocationID  int     `json:"locationID"`
		Hash        string  `json:"hash"`
		FittedValue float64 `json:"fittedValue"`
		TotalValue  float64 `json:"totalValue"`
		Points      int     `json:"points"`
		Npc         bool    `json:"npc"`
		Solo        bool    `json:"solo"`
		Awox        bool    `json:"awox"`
	} `json:"zkb"`
}