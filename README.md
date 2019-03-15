# UBC-RMP-Data :card_file_box:

![ubc-rmp-data demo](https://user-images.githubusercontent.com/9669739/52615922-b411c880-2e4b-11e9-9d80-fc00f31b0b3e.gif)

Scrapes [UBC Course Schedule page](https://courses.students.ubc.ca/cs/courseschedule?pname=subjarea), & Rate My Prof and stores it as JSON to be used for [UBC-RMP Chrome extension](https://github.com/jumbosushi/ubc-rmp)

Includes data for over 8200 sections with more than 3100 ratings per term

## Installation

```
go get github.com/jumbosushi/ubc-rmp-data
```

## Development

```
git clone git@github.com:jumbosushi/ubc-rmp-data.git
cd ubc-rmp-data
make
./bin/ubc-rmp-data
```

## Data Format

Two JSON files are available under `/data` per term per campus (ex. `2019_S_UBC_courseToInstrID.json`)

## courseToInstrID.json


```json
{
  "APBI": {
    "361": {},      // No lecture section
    "398": {
        "001": [    // Two instructors
            1304835,
            1310500
        ],
        "002": [
            1304835,
            1310500
        ]
    },
}
```

## instrIDToRating.json

```json
{
    "1324945": {
        "difficulty":        2.6,
        "name":              "ICHIKAWA, JONATHAN",
        "overall":           4.3,
        "rmpid":             1676955,
        "ubcid":             1324945,
        "would_take_again":  "85%"
    }
}
```
