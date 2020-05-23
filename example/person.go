package main

//go:generate gorec 
//go:generate gorec has_many credit_cards
type Person struct {
	PersonRecord
}

//go:generate gorec --model main.CreditCard 
type CreditCard struct {
	CreditCardRecord
}
