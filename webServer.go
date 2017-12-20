package main

import (
	"net/http"
	"eveKillmailCrawler/staticData"
	"eveKillmailCrawler/market"
	"html/template"
	"eveKillmailCrawler/zDatabase"
)

const priceMultiplier = 1.1

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
	return combinedData{
		KillerShips: generateTableElements(database.GetMostKillerShipSorted(f), 30, 6),
		LostCharges: generateTableElements(database.GetMostCargoItemSorted(f), 30, 8),
		LostModules: generateTableElements(database.GetMostFittedItemSorted(f), 100, 7),
		LostShips:   generateTableElements(database.GetMostLostShipSorted(f), 30, 6),
	}
}

func generateTableElements(items []*zDatabase.Item, total int, categoryId int) []tableElement {
	data := make([]tableElement, 0)

	for i := 0; i <= total; i++ {
		if staticData.GetCategoryIDFromTypeID(items[i].Id) == categoryId {
			data = append(data, tableElement{
				Name:      staticData.GetTypeIDName(items[i].Id),
				Total:     items[i].Count,
				PriceJita: market.GetPriceOfTypeID(items[i].Id),
				PriceSell: market.GetPriceOfTypeID(items[i].Id) * priceMultiplier,
			})
		}
	}

	return data
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
