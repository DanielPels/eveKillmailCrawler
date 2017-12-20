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





What should the crawler do:
	Per system, get all kills, put into database
	Once X time has been hit stop crawling that system(too old data not in the meta)

What we need:
	End data,
	Solarsystem ID
	List of systems to crawl,
	Loop that ticks every X seconds to crawl zkillboard,
	A queue(word that looks funny) for every solarsystem to request
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

//adds a system to the queue list
func (z *ZCrawler) AddSystem(id int) {
	for _, system := range z.queue {
		if system.Id == id {
			return
		}
	}
	z.queue = append(z.queue, &systemInQueue{Id: id, LastKill: nil})
	fmt.Println("Added " + strconv.Itoa(id) + " system to queue")
}

//starts the crawler tick
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

//stop the crawler
func (z *ZCrawler) Stop() {
	fmt.Println("Stopping crawler")
	if z.ticker != nil {
		z.ticker.Stop()
		z.ticker = nil
	}
}

func (z *ZCrawler) nextInQueue() {
	//loop trough all systems in the queue
	for _, system := range z.queue {
		//check if we have a last kill
		if system.LastKill != nil {
			//if we do then check if we have passed this time
			if hasTimePast(*system.LastKill, z.finishTime) {
				//true == past time continue to next for loop
				continue
			}
		}
		//create new go routine to request systems
		go func() {
			if system.LastKill != nil {
				z.getKillmails(system.Id, system.LastKill.Format("2006010215")+"00")
			} else {
				z.getKillmails(system.Id, "")
			}
		}()
		//break the loop to only request 1 system per tick
		break
	}
}

//request killmails from zkillboard api
func (z *ZCrawler) getKillmails(Id int, time string) {

	requestURL := baseURL + strconv.Itoa(Id) + "/"

	//if we have a time add time parameters in URL
	if time != "" {
		requestURL += "endTime/" + time + "/"
	}

	fmt.Println("Requesting URL: " + requestURL)

	//request the kills
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("ERROR HTTP GET")
		return
	}
	defer resp.Body.Close()

	//decode killmails from json
	var killMails *[]killmail.ZKillmail
	json.NewDecoder(resp.Body).Decode(&killMails)
	z.proccessKillmails(Id, *killMails)
}

//proccess the killmails
func (z *ZCrawler) proccessKillmails(systemID int, kills []killmail.ZKillmail) {
	//loop trough kills and add to database
	for _, kill := range kills {
		z.database.AddKillmails(systemID, kill)
	}

	//if we got a kill
	//get system and add kill time stamp as last killmailTime
	if len(kills) > 0 {
		for _, system := range z.queue {
			if system.Id == systemID {
				system.LastKill = &kills[len(kills)-1].KillmailTime
			}
		}
	}

	//print ;)
	fmt.Println("Got kills for system: " + strconv.Itoa(systemID))
	fmt.Println("Total of " + strconv.Itoa(len(kills)) + " added")
	fmt.Println("Last kill at: " + kills[len(kills)-1].KillmailTime.String())
}

//is T1 voorbij T2 zo jah true, zo niet false
//time diff
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
