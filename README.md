# UBC-RMP-Data :card_file_box:

![ubc-rmp-data demo](https://user-images.githubusercontent.com/9669739/52615922-b411c880-2e4b-11e9-9d80-fc00f31b0b3e.gif)

Scrapes [UBC Course Schedule page](https://courses.students.ubc.ca/cs/courseschedule?pname=subjarea), and stores it as JSON to be used for [UBC-RMP Chrome extension](https://github.com/jumbosushi/ubc-rmp)

## Installation

```
go get github.com/jumbosushi/ubc-rmp-data
```

## Development

```
git clone git@github.com:jumbosushi/ubc-rmp-data.git
cd ubc-rmp-data
make
# Start scraper
./bin/ubc-rmp-data
```

## Data Format

Two JSON files are available under `/data`

## ubcCourseInfo.json

```json
#
{
    "AANB": {
        "504": {
            "002": 1234
            }
        }
    }
}
```

## ubcInstrInfo.json

```json
{
    "1234": {
        "ubcid":            1234,
        "name":             "Doe, John",
        "difficulty":       3.5,
        "overall":          4,
        "would_take_again": "yes",
    }
}
```
