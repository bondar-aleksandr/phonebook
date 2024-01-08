package main

import (
	"context"
	"log"
	"fmt"
	"github.com/bondar-aleksandr/phonebook/entities"
	"github.com/bondar-aleksandr/phonebook/phonebook"

	// "github.com/bondar-aleksandr/phonebook/storage"
	// "github.com/bondar-aleksandr/phonebook/storage/mysql"
	sqlcstorage "github.com/bondar-aleksandr/phonebook/storage/sqlcStorage"
)

func main() {
	ctx := context.Background()
	// create a storage

	// s, err := mysql.New("root:password@tcp(localhost:3306)/phonebook?tls=skip-verify")
	s, err := sqlcstorage.New("root:password@tcp(localhost:3306)/?tls=skip-verify&multiStatements=true")
	if err != nil {
		log.Fatal(err)
	}
	pb := phonebook.New(s)
	person01 := entities.NewPerson("vasya", "pupkin", "")
	person01.AddPhone(entities.NewPhone(entities.Mobile, "+380951235432"))
	person01.AddPhone(entities.NewPhone(entities.Work, "+380961112223"))
	person02 := entities.NewPerson("petya", "pupkin", "xz")
	// person02.AddPhone(entities.NewPhone(entities.Mobile, "+380931112233"))
	// person02.AddPhone(entities.NewPhone(entities.Home, "+380683322111"))
	err = pb.Add(ctx, person01)
	if err != nil {
		log.Println(err)
	}
	err = pb.Add(ctx, person02)
	if err != nil {
		log.Println(err)
	}
	search := entities.NewSearchData("pet")
	result, err := pb.Search(ctx, search)
	if err != nil {
		log.Println(err)
	}
	for _,v := range result {
		fmt.Println(v)
	}
	// sperson01, err := pb.Get(ctx, search)
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(person02)
	// call menu for interactive actions
}