# sqlformatter [![Go Report Card](https://goreportcard.com/badge/github.com/brettcodling/sqlformatter)](https://goreportcard.com/report/github.com/brettcodling/sqlformatter)
Golang port of https://github.com/jdorn/sql-formatter

## Usage

```go
package main

import (
  "fmt"

  "github.com/brettcodling/sqlformatter"
)

func main() {
  fmt.Println(sqlformatter.Format("SELECT * FROM test.test"))
}
```
