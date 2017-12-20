package main

import (
	"fmt"
	"time"
	"os"
	"github.com/DanielPels/eveKillmailCrawler/crawler"
	"github.com/DanielPels/eveKillmailCrawler/staticData"
	"github.com/DanielPels/eveKillmailCrawler/market"
	"github.com/DanielPels/eveKillmailCrawler/zDatabase"
)

//dits een persoonlijk project om lol mee te hebben ;) en GO te leren

//Levens lessen in GO:
//fuck packages, het werkt zo: de directory naam is zo wat de package naam.
//Als je een package hebt is het soort van een object, tis meer een bundel van code als een object(aka package)
//de grap is alles met HOOFDLETTER wordt ge export / exposed dus dat kan je van buiten aanroepen
//alles wat dus een kleine letter heeft zal NIET ge export / exposed worden
//story of my life 8)

//Letop!
//maps zijn altijd pointers/references http://goinbigdata.com/golang-pass-by-pointer-vs-pass-by-value/

//todo lijst:
//maak functies zo dat de killDatabase gevult kan worden - DONE
//maak functie die de killDatabase backuped naar file - DONE
//maak een crawler voor zkillboard die de killDatabase vult - DONE
//maak een functie die de naam van een typeID ophaalt - DONE
//maak een object dat market data ophaalt en je de prijs kan opvragen afhanelijk van typeID - DONE
//maak functies die nuttige info uit de database halen - DONE
//maak een http web server die de data mooi kan representeren - DONE
//	er moet uitgerekend worden voor hoeveel de item ongeveer verkocht moet worden - DONE
//json bestand komen van welke system gecrawled moeten worden
//via http post een systeem kunnen toevoegen
//maak een market order tracker
//	moet items die sold zijn kunnen tracken
//	moet checken welke items in gekocht moeten worden
//	moet een mail kunnen uitsturen met een lijst van items die gekocht moeten worden
//maak een functie die zoekt welke item het meeste profijt geeft

var backupFileName = "data.json"
var backupTicker *time.Ticker
var database *zDatabase.ZDatabase

func main() {
	//testCodeMapsEnPointers()
	//testCodeTime()

	//load static data
	staticData.Init("staticData/typeIDs.json", "staticData/groupIDs.json", "staticData/categoryIDs.json")
	//init market
	market.Init()
	//create new database
	database = zDatabase.New()
	//check if we had a backup
	if checkForBackup(backupFileName) {
		//load if we had
		database.ImportFromJson(getBackup(backupFileName))
	}

	//make new crawler and add system
	go func() {
		crawler := zCrawler.New("2017-01-01T00:00:00Z", database)
		crawler.AddSystem(30002960)
		crawler.AddSystem(30002959)
		crawler.AddSystem(30002958)
		crawler.AddSystem(30002961)
		crawler.AddSystem(30002957)
		crawler.AddSystem(30002979)
		crawler.AddSystem(30002980)
		crawler.AddSystem(30002978)
		crawler.AddSystem(30002981)
		crawler.AddSystem(30002977)
		crawler.AddSystem(30002976)
		crawler.AddSystem(30003088)
		crawler.AddSystem(30002962)
		crawler.AddSystem(30002538)
		crawler.AddSystem(30003063)
		crawler.AddSystem(30003069)
		crawler.AddSystem(30003068)
		crawler.AddSystem(30002975)
		crawler.AddSystem(30002541)
		crawler.AddSystem(30002542)
		crawler.Start()
	}()

	//begin backup timer
	backupTicker = time.NewTicker(time.Second * time.Duration(60))
	go func() {
		for range backupTicker.C {
			//make backup
			saveBackup(database.ExportToJson(), backupFileName)
		}
	}()

	//start webserver
	NewWebServer()
}

func testCodeMapsEnPointers() {
	//zie links voor meer info
	//https://golang.org/doc/faq#Pointers
	//http://piotrzurek.net/2013/09/20/pointers-in-go.html
	//https://gist.github.com/josephspurrier/7686b139f29601c3b370

	fmt.Println("")
	fmt.Println("BEGIN POINTERS EN MAPS")
	//maak een map aan waar van de keys ints zijn en de value pointers naar de struct ship
	pointerMaps := make(map[int]*zDatabase.Ship)
	//als je een value ophaalt zal die nil zijn omdat de struct nog niet bestaat en niet toegevoegt is tot de map
	s := pointerMaps[0]
	if s == nil {
		//maak een struct aan en stop deze in de map
		pointerMaps[0] = &zDatabase.Ship{
			Killer: 0,
			Id:     0,
			Loss:   0,
			Solo:   0,
			Fitted: make(map[int]*zDatabase.Item),
			Cargo:  make(map[int]*zDatabase.Item),
		}
	}
	//je kan nu de struct wel ophalen omdat die al bestaat(je krijgt dus de pointer naar de struct ship)
	k := pointerMaps[0]
	//je kan nu de variable editen van de struct
	k.Solo++
	k.Killer++
	k.Killer++
	//print shizzle
	fmt.Println(&k)
	fmt.Printf("%+v\n", k)

	//als we een nieuwe struct aan maken zal deze in een ander stuk geheugen geplaatst worden
	pointerMaps[1] = &zDatabase.Ship{
		Killer: 0,
		Id:     0,
		Loss:   0,
		Solo:   0,
		Fitted: make(map[int]*zDatabase.Item),
		Cargo:  make(map[int]*zDatabase.Item),
	}

	//haal de pointer naar ship op
	a := pointerMaps[1]
	//print uit dat het daadwerkelijk een ander address heeft
	fmt.Println(&a)
	fmt.Printf("%+v\n", a)

	//wat ook kan is een pointer address pakken en deze in een andere variable stoppen
	//hier stop je memory address van K in Q(met het & teken haal je de memory address op)
	q := &k
	//nu is het mogelijk om met een * de pointer te dereference
	fmt.Println(*q)
	fmt.Println("END POINTERS EN MAPS")
	fmt.Println("")
}

func testCodeTime() {
	//links:
	//https://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format

	fmt.Println("")
	fmt.Println("BEGIN TIME")
	//time parse geeft altijd 2 variable
	//eerste is de tijd
	//tweede is error
	//als error nil is dan is alles goed gegaan zo niet dan rekt
	t1, e := time.Parse(time.RFC3339, "2017-12-01T14:50:46Z")
	if e != nil {
		os.Exit(1)
	}

	//met de underscore zeg je dat de variable weg gegooit kan worden(de error dus)
	t2, _ := time.Parse(time.RFC3339, "2017-12-01T14:44:18Z")

	//je kan tijd van elk kaar afhalen om zo te zien of een datum voorbij is
	fmt.Println(t2.Sub(t1).Seconds())

	//je kan time ook formatten zie hier onder
	fmt.Println("formatted: " + t1.Format("2006010215") + "00")

	fmt.Println("END OF TIME")
	fmt.Println("")
}
