# UBC-RMP-Scraper

Scrapes [UBC Course web page](https://courses.students.ubc.ca/cs/courseschedule?pname=subjarea), and stores it as JSON to be used for [UBC-RMP extension](https://github.com/jumbosushi/ubc-rmp)

JSON Format
```json
{
  "AANB": {
    "504": {
      "002": {
        "Doe, John": {
          "name":             "John Doe",
          "difficulty":       3.5,
          "overall":          4,
          "would_take_again": "yes",
        }
      }
    }
  }
}
```

