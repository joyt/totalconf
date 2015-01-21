# totalconf
Total configuration solution for Golang built on top of [rakyll/globalconf] (github.com/rakyll/globalconf). Define config variables using flag package's syntax and easily reuse them. You can redefine multiple flags/config variables of the same name.

# Usage
~~~ go
import "github.com/joyt/totalconf"
~~~
Define config vars in your packages:
```go
var name = totalconf.String("name", "default value", "usage")
```
and call parse with options in your `main()`:
```go
opts := new(totalconf.Options)
opts.Filename = "myapp.conf"
totalconf.Parse(opts)
```
See [rakyll/globalconf] (github.com/rakyll/globalconf) for detailed use of options.

Unlike the [flag package](http://golang.org/pkg/flag/) you can redefine flags/vars with the same name/label and reuse them across packages so that they don't have to be passed around explicitly. You can even define them with different types and they will resolve properly as long as the input is parseable for both types. Note that if you have different default values set for vars with the same name/label they default to their respective default values if not set. See example below.

# Example
In the main file:
~~~ go
package main

import (
  "github.com/joyt/totalconf"
  "server" // local package
  "fmt"
)

var (
  port = totalconf.Int("port", 8080, "Server port")
  name1 = totalconf.String("name", "Bertha", "Name 1")
  name2 = totalconf.String("name", "Jenna", "Name 2")
)

func main() {
  totalconf.Parse(nil)
  if *port < 8000 {
    fmt.Println("Port is too low")
  } else {
    fmt.Println("Starting", *name1, *name2, "on", *port)
    server.Start()
  }
}

~~~
In another package:
~~~ go
package server

import (
  "fmt"
  "http"
  "log"
  
  "github.com/joyt/totalconf"
)

var port = totalconf.String("port", "8000", "Server port")

func Start() {
  fmt.Println("Server listening on", port)
  log.Println(http.ListenAndServe(":" + port, nil))
}
~~~

```
> go run program.go
> Starting Bertha Jenna on 8080
> Server listening on 8000

> go run program.go -port=9000
> Port is too low

> go run program.go -port=7000
> Starting Bertha Jenna on 7000
> Server listening on 7000

> go run program.go -name=Freddie
> Starting Freddie Freddie on 8080
> Server listening on 8000
```
