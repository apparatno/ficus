# Ficus

Overvåker utleggsmappen på Google Drive
og sier ifra på Slack.

## Kom i gang

Du trenger:

* en serviceaccount credentials JSON fil ved navn `credentials.json`
* en app token fra Slack (med tilhørende app integrasjon satt opp)
* en Google Team drive ID
* en ID på mappen som skal sjekkes

IDer kan finnes vha URLer på [Google Drive](https://drive.google.com).

## Kjøre

Programmet kan kjøres som en frittstående kommando
eller med cron.

### Argumenter

Programmet har disse argumentene:

```bash
$> ficus -h
Usage of ficus:
  -db string
    	path to database JSON file. Defaults to ./db.json (default "db.json")
  -driveid string
    	ID of the Google Drive to use
  -no-slack
    	don't send Slack messages
  -root string
    	ID of the folder to scan
```

For å autentisere mot Slack må du i tillegg eksportere Slack token
som `FICUS_SLACK_TOKEN`.
Det er mulig å skippe Slack ved å bruke `-no-slack`
(fortrinsvis for testformål).

### Databasen

Programmet vil lagre mapper den har funnet i en JSON fil på disk.
Man kan selv velge hvor filen skal lagres (med `-db` argumentet)
eller bare lagre i current dir.

JSON filen holder styr på når det sist ble sett filer i hver undermappe til rot-mappen
slik at man kun blir informert om nye endringer.
Finnes filen ikke vil programmet informere om alle filer den kommer over
og lagre dette i databasen.
