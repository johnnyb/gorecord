# GoRecord

## About

This is meant to be a Go ORM similar to ActiveRecord.  
There are a million Go ORMs, but none of them quite as good as ActiveRecord.
[sqlboiler](https://github.com/volatiletech/sqlboiler) comes close, but
I don't like the reliance on separate config files, or its lack of a migrator.  
Looking for a more Go-ish solution.

Additionally, I want everything to be either (a) compiled in to the project, or (b) set via environment variable.  This makes distribution and configuration across many different systems 
much easier.

This is a brand new project, so don't expect full magic.  Also, currently this expects all models that have relationships with each other be in the same package.

For some of the thinking behind the way that this project is stuctured, see this short article on [Automated Code Generation](https://mindmatters.ai/2020/05/automated-code-generation-tools-can-solve-problems/).

## Installing

First, grab the generating command:

```
go get github.com/johnnyb/gorecord/...
```

This will get both the library and the command (if you just import the library, you won't get 
the `gorec` command).

## Building an Example Program

In this section we are going to build a complete example GoRecord app from scratch.
Create an empty directory called `testprogram` and go there.  Now initialize the module:

```
go mod init example.com/testprogram
```

Now create a `main.go` file:

```
package main
import (
	"fmt"
)
func main() {
	fmt.Println("Hello world")
}
```

Now, `go build` and `./testprogram` to make sure it works.

### Creating the database

Create the database manually. 

```
createdb -U postgres gotest
```

Set these two environment variables:
```
export DB_DRIVER=pgx
export DB_CONNECTION_STRING="user=postgres database=gotest password=REPLACE_WITH_YOUR_PWD"
```

### Creating your first migration

Now, create a subdirectory called `migrations` to store your database migrations.  Create your
first migration with the `gorec` command.  We are going to create a table called `people`:

```
gorec -action migration -directory migrations -named create_people
```

This will create a file called `migrations/TIMESTAMP_create_people`.  Open the file.  Replace the `panic` on the up migration with the following code:

```
_, err := conn.Exec("CREATE TABLE people (id serial PRIMARY KEY, name text, description text)")
return err
```

Note that your migration can actually be any Go code that you want.  However, I encourage
you to mostly keep it to just `Exec`ing SQL commands.  Note that since the `id` is `serial`
the database will be responsible for generating primary keys.  

### Connecting your code to the database and migrations

Now, in your `main.go` file, add the following to your imports list:
```
"github.com/johnnyb/gorecord/gorec"
"github.com/johnnyb/gorecord/migrator"
_ "github.com/jackc/pgx/v4/stdlib"       // Load the database driver
_ "example.com/testprogram/migrations"   // Load your migrations
```
Now add the following code to your `main()` function:
```
db, err := gorec.AutoConnect()
if err != nil {
	panic(err)
}
defer db.Close()

err = migrator.MigrateRegisteredMigrations()
if err != nil {
	panic(err)
}
```

This will connect to the database, and auto-run any migrations that you have outstanding.
Note that it keeps track of which migrations you have run, so if a migration has run it will
not run again.  The `AutoConnect` function connects to the database and sets the global
variables `gorec.GlobalConnection` which will be implicitly used throughout GoRecord.
If you need to connect to more than one database, that is supported, but is a little tricky.

Now, build a run your program: `go build` and `./testprogram`.  This will perform your migration for you!

### Create your first model

Create the directory `models`.  
In that directory, create a file called `person.go` and 
put in the following code:

```
package models

//go:generate gorec -model Person
type Person struct {
	PersonRecord
}
```

Note that the LACK of a space between `//` and `go:generate` is EXTREMELY important in Go.

Now, while in the `models` directory, run `go generate`.
This will make a file called `person.impl.go` which has all
of the details of the model.  This includes getters and setters
for all of the fields, a `PersonNew()` function, a `PersonFind()` function,
a `PersonQuerySimple()` function and a `Save()` method.  `gorec`
automagically discovers your table fields from the database itself.  
This is why you had to run migrations first.  If you add more fields, to the table,
you will need to re-run `go generate` here.

You may wonder how we went from a table named `people` to a model named `Person`.
GoRecord adopts many of the ActiveRecord naming conventions.  This means that the
model name is an uppercase singular version of the word, and the table is an underscored
plural version of the word.  GoRecord contains an inflector which will singularize/pluralize
words appropriately, though, at the current stage of development, its handling of special
cases is limited.  If you need to specify the table name, you can pass in a `-table` option
to the `gorec` command.

### Access your model from your code

In your `main.go` add the following import:
```
"example.com/testprogram/models"
```

Then, add the following code to the *end* of the `main()` function:

```
rec := models.PersonNew()
rec.SetName("Bob")
rec.SetDescription("Bob is a Great Guy")
err = rec.Save()
if(err != nil) {
	panic(err)
}
fmt.Printf("Bob's new record ID is %d\n", rec.Id())
```

If you want to get that same record again, you can add this code:
```
bobrec, err := models.PersonFind(rec.Id())
fmt.Printf("Found record %d named %s\n", bobrec.Id(), bobrec.Name());
```
If you want to find Bob by his name, you can add this code:
```
bobrecs, err := models.PersonQuerySimple("where name = $1", rec.Name())
fmt.Printf("Found %d records named Bob\n", len(bobrecs))
```
This will return a slice of `Person`s named `Bob`.

### Creating a relationship

Let's say that we want to add a one-to-many relationship.
We would do this using the `HasMany` relationship.

Let's start by creating a new migration, `create_cars`.  
From your main directory, do:
```
gorec -action migration -named create_cars -directory migrations
```
Put this in the migration:
```
_, err := conn.Exec("CREATE TABLE cars (id serial PRIMARY KEY, person_id int, make text, model text, year text)")
return err
```

Build and run your program (`./testprogram`)to perform the migration.  Now add the model.
In `models/car.go`:

```
package models

//go:generate gorec -model Car
//go:generate gorec -model Car -action BelongsTo -relationship Person
type Car struct {
	CarRecord
}
```

This will add a function `.Person()` to any car object to retrieve the related person object (also with an error result).

Note that this has an additional generator specifying the relationship between the Car and the Person, stating that a car "belongs to" that person.
In `person.go` we need to add a similar relationship line.  This has to come *after*
the first `go:generate` command in the file:

```
//go:generate gorec -model Person -action HasMany -relationship Cars
```

This will add a function `.Cars()` to any person object to retrieve the related cars objects (as well as an error field).
Additionally, you will get a `.CarsBuild()` function to create a new Car object that is
already setup with the relationship.

### Enjoy!

I hope you enjoy GoRecord.  It is still in its early stages of development, but it is already pretty fun to play with.
