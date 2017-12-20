# eveKillmailCrawler
Fun project to learn GO.

The idea is to crawl Zkillboard to find out what people fly/fit on their internet spaceship(EVE Online). This data can be filterd to see what FW pilots fly/fit(amarr, minmatar only).

This tool can be usefull to find out what to buy and seed in lowsec market areas

[Tool to make eve TypeIDs easier for GO](https://github.com/DanielPels/EveTypeIDFix)

How it works:
* Give the crawler solar system ID's to crawl
* Crawls zkillboards api
	* Every 30 sec makes a http call to https://zkillboard.com/api/system/
* Saves data from killmails in a simple struct(sort of a database)
* Makes a simple text file backup of gatherd killmails every 60 seconds
* Gets market prices from [esi eve online](https://esi.tech.ccp.is/ui/)
* Serve http server on :8080
	* Html file shows table with item name, total count, jita price, jita+10% price
	* **Paths**:
	* / gives all combined data(so amarr, neutral, minmatar)
	* /amarr gives all combined data for amarr FW pilots only!
	* /minmatar gives all combined data for minmatar FW pilots only!


Thats it! easy and usefull. Was a great starting project.

If intrested fun comments are in the code! ;)