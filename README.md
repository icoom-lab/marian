Clean structs conditionally based on tags on the fields.


**NOTE**: This package is tested only on Go 1.18+

Example

```go
package main

import (
	"fmt"
	"github.com/icoom-lab/marian"
)

type Account struct {
	Id   int    `json:"id,omitempty" role:"admin"`
	Name string `json:"name,omitempty" role:"admin,normal"`
}

func main() {
	account := Account{
		Id:   1,
		Name: "Jhon",
	}

	fmt.Println(account)
	// {id:1, Name:"Jhon"}

	marian.CleanStruct("role", "admin", &account)
	fmt.Println(account)
	// {id:1, Name:"Jhon"}

	marian.CleanStruct("role", "normal", &account)
	fmt.Println(account)
	// {Name:"Jhon"}
}
```

