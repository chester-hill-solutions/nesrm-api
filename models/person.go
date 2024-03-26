package models

import (
	"fmt"

	"github.com/supabase-community/supabase-go"

	"github.com/joho/godotenv"
)

type Person struct{
  uuid string
  givenname string
  surname string
  birthdate string
  deceased string
  bio_mother_uuid string
  bio_father_uuid string
  linkedin_link string
}

func GetPersons() string {
  myEnv, _ := godotenv.Read(".env")
  client, err := supabase.NewClient(myEnv["API_URL"], myEnv["API_KEY"], nil)
  if err != nil {
    fmt.Println("cannot initalize client", err)
  }
  data, _, err := client.From("person").Select("*", "exact", false).Execute()
  return fmt.Sprintf("%T", data)
}
