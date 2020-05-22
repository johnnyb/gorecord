# GoRecord

## About

This is meant to be a Go ORM similar to ActiveRecord.  
There are a million Go ORMs, but none of them quite as good as ActiveRecord.
`sqlboiler` comes close ( https://github.com/volatiletech/sqlboiler ), but
I don't like the reliance on config files.  Looking for a more Go-ish solution.

This is a brand new project, so don't expect full magic.

## Usage

In your code:

```
import "github.com/johnnyb/gorecord/gorec"
```

This gives you access to the `gorec` module for the basics.
In your own code, you can put it anywhere, but I suggest creating a module called `models` for your code.

Then, let's say you wanted to have a `Person` object that was tied to a `people` table.
Create a file called `person.go` that looks like the following:

```
package models
import (
	"github.com/johnnyb/gorecord/gorec"
)

//go:generate gorec --model Person --table people
type Person struct {
	PersonRecord  // This struct will be auto-generated
	// You can add non-database fields here
}

// GoRecord will define many functions on Person already, 
// but feel free to define your own here!
```

Now, in your main code, you can do:

```
package main
import (
	"fmt"
	"database/sql"
	"path/to/your/models"
	"github.com/johnnyb/gorecord/gorec"
)

func main() {
	db, err := sql.Open("SqlDriverNameHere", "SqlConnectionStringhere")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	gorec.Connection = db

	person := models.Person{}
	person.SetName("Fred")
	err = person.Save()
	if err != nil {
		panic(err)
	}

	personId := p.Id()

	personAgain := models.FindPerson(personId)
	personAgain.SetName("Bob")
	personAgain.Save()

	fmt.Printf("These are the IDs of the people named Bob:\n");
	people := models.QueryPerson("name = $1", "Bob")
	for p := range people {
		fmt.Printf("%d", p.Id()
	}
}
```
