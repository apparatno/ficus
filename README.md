# Ficus

Overvåker utleggsmappen på Google Drive
og sier ifra på Slack.

## Kom i gang

Du trenger:

* en serviceaccount credentials JSON fil ved navn `credentials.json`
* en app token fra Slack (med tilhørende app integrasjon)
* en Google drive ID på mappene som skal overvåkes deres foreldermappe
* en JSON "database" (se under)

### Databasen

Databasen er en JSON fil med følgende format:

```json
[
  {
    "id": "google-drive-mappe-id",
    "user": "navn på mappens eier"
  }
]
```

## Kjør!

Exporter slack tokenet som `FICUS_SLACK_TOKEN`
og IDen på foreldermappen som `FICUS_DRIVE_ID`
og start appen.

Den vil gjøre en første sjekk og deretter kjøre hvert 10. minutt.
