# GoRecord

## About

This is meant to be a Go ORM similar to ActiveRecord.  
There are a million Go ORMs, but none of them quite as good as ActiveRecord.
`sqlboiler` comes close ( https://github.com/volatiletech/sqlboiler ), but
I don't like the reliance on config files.  Looking for a more Go-ish solution.

This is a brand new project, so don't expect full magic.  Also, currently this expects all models that have relationships with each other be in the same package.

## Usage

First, grab the generating command:

```
go get github.com/johnnyb/gorecord/...
```

This will get both the library and the command (if you just import the library, you won't get the `gorec` command).
Then, in your code:

```
import "github.com/johnnyb/gorecord/gorec"
```

This gives you access to the `gorec` module for the basics.
In your own code, you can use GoRecord anywhere, but I suggest creating a module called `models` for your database code.

Now, let's say you wanted to have a `Person` object that was tied to a `people` table linked to a `credit_cards` table.
First, create a table in your database called `people` and a table called `credit_cards`:
```
create table people (id int serial primary key, name text);
create_table credit_cards (id int serial primary key, person_id int, info text);
```

Next, create a file called `person.go` that looks like the following:

```
package models
import (
	"github.com/johnnyb/gorecord/gorec"
)

//go:generate gorec -model Person 
//go:generate gorec -model Person -action HasMany -relationship CreditCards
type Person struct {
	PersonRecord  // This struct will be auto-generated with backing fields for the database
	// You can add non-database fields here
}

// GoRecord will define many functions on Person already, 
// but feel free to define your own here!
```

Then, in `credit_cards.go` do:
```
package models
import (
	"github.com/johnnyb/gorecord/gorec"
)

//go:generate gorec -model CreditCard
type CreditCard struct {
	CreditCardRecord  // This struct will be auto-generated with backing fields for the database
	// You can add non-database fields here
}

// GoRecord will define many functions on CreditCard already, 
// but feel free to define your own here!
```

Now, in your main code, you can do:

```
package main
import (
	_ "github.com/jackc/pgx/v4/stdlib" // Make sure your DB driver is available
	"fmt"
	"path/to/your/models"
	"github.com/johnnyb/gorecord/gorec"
)

func main() {
	// Setup the global connection
	db, err := gorec.AutoConnect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	person := models.PersonNew()
	person.SetName("Fred")
	err = person.Save()
	if err != nil {
		panic(err)
	}

	cc := person.CreditCardsBuild()
	cc.info = "My CC Info"
	err = cc.Save()
	if err != nil {
		panic(err)
	}

	for _, cc := range person.CreditCards() {
		// Do stuff with credit cards
	}

	personId := p.Id()

	personAgain := models.PersonFind(personId)
	personAgain.SetName("Bob")
	personAgain.Save()

	fmt.Printf("These are the IDs of the people named Bob:\n");
	people := models.PersonQuerySimple("where name = $1", "Bob")
	for p := range people {
		fmt.Printf("%d", p.Id())
	}
}
```
