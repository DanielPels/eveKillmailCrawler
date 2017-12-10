package zCrawler

import (
	"time"
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
	"eveKillmailCrawler/zDatabase"
	"eveKillmailCrawler/killmail"
)

/*
crawler moet wat doen?
	per systeem alle kills ophalen en deze in de zdatabase proppen

de crawler moet stoppen als een X datum berijkt is(willen niet killmails uit het jaar nil hebben)

we hebben dus nodig:
	end datum
	solar system ID
	Lijst van systems die we crawlen
	een loop om elke X seconden een http call te maken naar de Zkillboard api
	soort van queue(wat een stom woord) voor welk solar system gerequest moet worden
*/

const baseURL string = "https://zkillboard.com/api/system/"
const apiRequestDelay int = 30

type systemInQueue struct {
	Id       int
	LastKill *time.Time
}

type ZCrawler struct {
	queue      []*systemInQueue
	ticker     *time.Ticker
	finishTime time.Time
	database   *zDatabase.ZDatabase
}

func (z *ZCrawler) AddSystem(id int) {
	for _, system := range z.queue {
		if system.Id == id {
			return
		}
	}
	z.queue = append(z.queue, &systemInQueue{Id: id, LastKill: nil})
	fmt.Println("Added " + strconv.Itoa(id) + " system to queue")
}

func (z *ZCrawler) Start() {
	fmt.Println("Starting crawler")
	if z.ticker != nil {
		fmt.Println("Crawler already active!")
		return
	}
	//dit was dom je kan dus niet time.Second keer Int doen
	//het moet een duration zijn
	z.ticker = time.NewTicker(time.Second * time.Duration(apiRequestDelay))
	go func() {
		z.nextInQueue()
		for range z.ticker.C {
			fmt.Println("")
			fmt.Println("Tick!")
			z.nextInQueue()
		}
	}()
}

func (z *ZCrawler) Stop() {
	fmt.Println("Stopping crawler")
	if z.ticker != nil {
		z.ticker.Stop()
		z.ticker = nil
	}
}

func (z *ZCrawler) nextInQueue() {
	//loop door alle systems heen die in de queue zitten
	for _, system := range z.queue {
		//check of er wel een tijd is zo niet dan gewoon de kills ophalen
		if system.LastKill != nil {
			//check of tijd voorbij is
			if hasTimePast(*system.LastKill, z.finishTime) {
				//zo jah skip dit system
				continue
			}
		}
		//idk? als goed is gewoon lekkere nieuw go routine
		go func() {
			if system.LastKill != nil {
				z.getKillmails(system.Id, system.LastKill.Format("2006010215")+"00")
			} else {
				z.getKillmails(system.Id, "")
			}
		}()
	}
}

func (z *ZCrawler) getKillmails(Id int, time string) {

	requestURL := baseURL + strconv.Itoa(Id) + "/"

	if time != "" {
		requestURL += "endTime/" + time + "/"
	}

	fmt.Println("Requesting URL: " + requestURL)

	//haal de kills op
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("ERROR HTTP GET")
		return
	}
	defer resp.Body.Close()

	var killMails *[]killmail.ZKillmail
	json.NewDecoder(resp.Body).Decode(&killMails)
	z.proccessKillmails(Id, *killMails)
}

func (z *ZCrawler) proccessKillmails(Id int, kills []killmail.ZKillmail) {
	for _, kill := range kills {
		z.database.AddKillmails(Id, kill)
	}

	if len(kills) > 0 {
		for _, system := range z.queue {
			if system.Id == Id {
				system.LastKill = &kills[len(kills)-1].KillmailTime
			}
		}
	}

	fmt.Println("Got kills for system: " + strconv.Itoa(Id))
	fmt.Println("Total of " + strconv.Itoa(len(kills)) + " added")
	fmt.Println("Last kill at: " + kills[len(kills)-1].KillmailTime.String())
}

//is T1 voorbij T2 zo jah true, zo niet false
func hasTimePast(t1 time.Time, t2 time.Time) bool {
	if t1.Sub(t2).Hours() > 0 {
		return false
	}
	return true
}

func New(finishTime string, db *zDatabase.ZDatabase) *ZCrawler {
	time, _ := time.Parse(time.RFC3339, finishTime)
	return &ZCrawler{
		queue:      make([]*systemInQueue, 0),
		finishTime: time,
		ticker:     nil,
		database:   db,
	}
}
