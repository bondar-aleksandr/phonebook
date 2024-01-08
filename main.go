package main

import (
	"context"
	"log"

	"github.com/bondar-aleksandr/phonebook/entities"
	"github.com/bondar-aleksandr/phonebook/phonebook"
	"github.com/bondar-aleksandr/phonebook/storage"
	"github.com/bondar-aleksandr/phonebook/storage/mysql"
)

func main() {
	ctx := context.Background()
	// create a storage

	s, err := mysql.New("root:password@tcp(localhost:3306)/phonebook?tls=skip-verify")
	if err != nil {
		log.Fatal(err)
	}
	pb := phonebook.New(s)
	person01 := entities.NewPerson("vasya", "pupkin", "")
	phone01 := entities.NewPhone(entities.Mobile, "+380951235432", true)
	phone02 := entities.NewPhone(entities.Work, "+380961112223", false)
	person01.AddPhone(phone01)
	person01.AddPhone(phone02)
	err = pb.Add(ctx, person01)
	if err != nil {
		log.Println(err)
	}
	search := storage.NewSearchData("vasya")
	person02, err := pb.Get(ctx, search)
	if err != nil {
		log.Println(err)
	}
	log.Println(person02)
	// call menu for interactive actions
}