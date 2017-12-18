package main

import (
	"net/http"
	"eveKillmailCrawler/staticData"
	"eveKillmailCrawler/market"
	"html/template"
	"eveKillmailCrawler/zDatabase"
)

func NewWebServer() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/amarr", handleAmarr)
	http.HandleFunc("/minmatar", handleMinmatar)
	http.ListenAndServe(":8080", nil)
}

type tableElement struct {
	Name      string  `json:"Name"`
	Total     int     `json:"Total"`
	PriceJita float64 `json:"PriceJita"`
	PriceSell float64 `json:"PriceSell"`
}

type combinedData struct {
	LostShips   []tableElement
	LostModules []tableElement
	LostCharges []tableElement
	KillerShips []tableElement
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	t := template.New("Test")
	t.Parse(html)
	t.Execute(w, getCombinedData(0))
}

func handleAmarr(w http.ResponseWriter, r *http.Request) {
	t := template.New("Test")
	t.Parse(html)
	t.Execute(w, getCombinedData(zDatabase.Amarr))
}

func handleMinmatar(w http.ResponseWriter, r *http.Request) {
	t := template.New("Test")
	t.Parse(html)
	t.Execute(w, getCombinedData(zDatabase.Minmatar))
}

func getCombinedData(f int) combinedData {

	lostShipsData := database.GetMostLostShipSorted(f)
	lostShips := make([]tableElement, 0)

	for i := 0; i <= 30; i++ {
		if staticData.GetCategoryIDFromTypeID(lostShipsData[i].Id) == 6 {
			lostShips = append(lostShips, tableElement{
				Name:      staticData.GetTypeIDName(lostShipsData[i].Id),
				Total:     lostShipsData[i].Count,
				PriceJita: market.GetPriceOfTypeID(lostShipsData[i].Id),
				PriceSell: market.GetPriceOfTypeID(lostShipsData[i].Id) * 1.1,
			})
		}
	}

	lostModulesData := database.GetMostFittedItemSorted(f)
	lostModules := make([]tableElement, 0)

	for i := 0; i <= 100; i++ {
		if staticData.GetCategoryIDFromTypeID(lostModulesData[i].Id) == 7 {
			lostModules = append(lostModules, tableElement{
				Name:      staticData.GetTypeIDName(lostModulesData[i].Id),
				Total:     lostModulesData[i].Count,
				PriceJita: market.GetPriceOfTypeID(lostModulesData[i].Id),
				PriceSell: market.GetPriceOfTypeID(lostModulesData[i].Id) * 1.1,
			})
		}
	}

	lostChargesData := database.GetMostCargoItemSorted(f)
	lostCharges := make([]tableElement, 0)

	for i := 0; i <= 30; i++ {
		if staticData.GetCategoryIDFromTypeID(lostChargesData[i].Id) == 8 {
			lostCharges = append(lostCharges, tableElement{
				Name:      staticData.GetTypeIDName(lostChargesData[i].Id),
				Total:     lostChargesData[i].Count,
				PriceJita: market.GetPriceOfTypeID(lostChargesData[i].Id),
				PriceSell: market.GetPriceOfTypeID(lostChargesData[i].Id) * 1.1,
			})
		}
	}

	killerShipsData := database.GetMostKillerShipSorted(f)
	killerShips := make([]tableElement, 0)

	for i := 0; i <= 30; i++ {
		if staticData.GetCategoryIDFromTypeID(killerShipsData[i].Id) == 6 {
			killerShips = append(killerShips, tableElement{
				Name:      staticData.GetTypeIDName(killerShipsData[i].Id),
				Total:     killerShipsData[i].Count,
				PriceJita: market.GetPriceOfTypeID(killerShipsData[i].Id),
				PriceSell: market.GetPriceOfTypeID(killerShipsData[i].Id) * 1.1,
			})
		}
	}

	return combinedData{
		KillerShips: killerShips,
		LostCharges: lostCharges,
		LostModules: lostModules,
		LostShips:   lostShips,
	}
}

var html = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Eve killmail crawler</title>
    <style>
        table {
            font-family: arial, sans-serif;
            border-collapse: collapse;
            width: 100%;
        }

        td, th {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }

        tr:nth-child(even) {
            background-color: #dddddd;
        }
    </style>
</head>
<body>
<h1>30 most killed ships</h1>
<table>
    <tr>
        <th>Ship</th>
        <th>Total</th>
        <th>Jita price</th>
        <th>Sell price</th>
    </tr>
{{range .LostShips}}
    <tr>
        <td>{{.Name}}</td>
        <td>{{.Total}}</td>
        <td>{{printf "%.2f" .PriceJita}}</td>
        <td>{{printf "%.2f" .PriceSell}}</td>
    </tr>
{{end}}
</table>
<hr>

<h1>100 most lost modules</h1>
<table>
    <tr>
        <th>Item</th>
        <th>Total</th>
        <th>Jita price</th>
        <th>Sell price</th>
    </tr>
{{range .LostModules}}
    <tr>
        <td>{{.Name}}</td>
        <td>{{.Total}}</td>
        <td>{{printf "%.2f" .PriceJita}}</td>
        <td>{{printf "%.2f" .PriceSell}}</td>
    </tr>
{{end}}
</table>
<hr>

<h1>30 most lost charges</h1>
<table>
    <tr>
        <th>Item</th>
        <th>Total</th>
        <th>Jita price</th>
        <th>Sell price</th>
    </tr>
{{range .LostCharges}}
    <tr>
        <td>{{.Name}}</td>
        <td>{{.Total}}</td>
        <td>{{printf "%.2f" .PriceJita}}</td>
        <td>{{printf "%.2f" .PriceSell}}</td>
    </tr>
{{end}}
</table>
<hr>

<h1>30 highest KILLING ships</h1>
<table>
    <tr>
        <th>Item</th>
        <th>Total</th>
        <th>Jita price</th>
        <th>Sell price</th>
    </tr>
{{range .KillerShips}}
    <tr>
        <td>{{.Name}}</td>
        <td>{{.Total}}</td>
        <td>{{printf "%.2f" .PriceJita}}</td>
        <td>{{printf "%.2f" .PriceSell}}</td>
    </tr>
{{end}}
</table>
<hr>
</body>
</html>`
