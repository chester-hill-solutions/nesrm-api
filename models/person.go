package models

import (
	"context"
	"fmt"
	"log"
	"time"

"github.com/chester-hill-solutions/nesrm_api/pgConnector"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Person struct{
  UUID string
  Created_at time.Time
  Givenname string
  Surname string
  Birthdate time.Time
  Deceased time.Time
  Bio_mother *Person
  Bio_father *Person
  Linkedin_link string
}

func (p Person) toMap() map[string]string  {
  m := make(map[string]string)
  m["UUID"] = p.UUID
  m["Created_at"] = p.Created_at.String()
  m["Givenname"] = p.Givenname
  m["Surname"] = p.Surname
  m["Birthdate"] = p.Birthdate.String()
  m["Deceased"] = p.Deceased.String()
  m["Bio_mother_UUID"] = p.Bio_mother.UUID
  m["Bio_father_UUID"] = p.Bio_father.UUID
  m["Linkedin_link"] = p.Linkedin_link
  return m
}

//CREATE
func CreateNewPerson(fields map[string]interface{}) (*Person, error){
  person := Person{}
  if len(fields) == 0 {
    return nil, nil
  }
  UUID, ok := fields["UUID"]
  if ok {
    person.UUID = UUID.(string)
  }
  return &person, nil
}


func NewPerson(UUID string, Created_at time.Time, Givenname string, Surname string, Birthdate time.Time, Deceased time.Time, Bio_mother *Person, Bio_father *Person, Linkedin_link string)  *Person{
  person := Person{
    UUID: UUID,
    Created_at: Created_at,
    Givenname: Givenname,
    Surname: Surname,
    Birthdate: Birthdate,
    Deceased: Deceased,
    Bio_mother: Bio_mother,
    Bio_father: Bio_father,
    Linkedin_link: Linkedin_link,
  }
  return &person
}

func BuildFromTree(tree map[string]*Person, base *Person) Person {
  log.Println("BuildFromTree")
  return buildFromTree(tree, base)
}

func buildFromTree(tree map[string]*Person, base *Person)  Person{ 
  startTime := time.Now()
  log.Println("buildFromTree")
  log.Println("Building for: ", base.Givenname, base.Surname)
  var mother, father *Person
  var mok, fok bool
  if base.Bio_mother != nil {
    mother, mok = tree[base.Bio_mother.UUID]
    log.Println("Mother for ", base.Givenname, " ", base.Surname, ": ", mother.Givenname, " ", mother.Surname)
  } else {
    log.Println("No mother for: ", base.Givenname)
    mother, mok = nil, false
  }
  if base.Bio_father != nil {
    father, fok = tree[base.Bio_father.UUID]
    log.Println("father for ", base.Givenname, " ", base.Surname, ": ", father.Givenname, " ", father.Surname)
  } else {
    father, fok = nil, false
    log.Println("No father for: ", base.Givenname)
  }
  changed := *base
  if mok {
    log.Println("Finding mother for: ", changed.Givenname)
    p := buildFromTree(tree, mother)
    changed.Bio_mother = &p
  }
  log.Println(base.Givenname, " father ok: ", fok)
  if fok {
    log.Println("Finding father for: ", changed.Givenname)
    p := buildFromTree(tree, father)
    changed.Bio_father = &p
  }
  log.Println("buildFromTree,", time.Since(startTime).String())
  return changed
}

func AdvNewPerson(legend map[string]string)  *Person{
  parsedTimeValues := pgConnector.BulkTimeParser([]string{legend["Birthdate"], legend["Deceased"], legend["Created_at"], legend["mother_Birthdate"], legend["mother_Deceased"], legend["mother_Created_at"], legend["father_Birthdate"], legend["father_Deceased"], legend["father_Created_at"]})
  person := Person{
    UUID: legend["UUID"],
    Created_at: parsedTimeValues["Created_at"],
    Givenname: legend["Givenname"],
    Surname: legend["Surname"],
    Birthdate: parsedTimeValues["Birthdate"], 
    Deceased: parsedTimeValues["Deceased"],
    Bio_mother: &Person{
      UUID: legend["mother_UUID"],
      Created_at: parsedTimeValues["mother_Created_at"],
      Givenname: legend["mother_Givenname"],
      Surname: legend["mother_Surname"],
      Birthdate: parsedTimeValues["mother_Birthdate"], 
      Deceased: parsedTimeValues["mother_Deceased"],
      Bio_mother: nil,
      Bio_father: nil,
      Linkedin_link: legend["mother_Linkedin_link"],
    },
    Bio_father: &Person{
      UUID: legend["father_UUID"],
      Created_at: parsedTimeValues["father_Created_at"],
      Givenname: legend["father_Givenname"],
      Surname: legend["father_Surname"],
      Birthdate: parsedTimeValues["father_Birthdate"], 
      Deceased: parsedTimeValues["father_Deceased"],
      Bio_mother: nil,
      Bio_father: nil,
      Linkedin_link: legend["father_Linkedin_link"],
    },
    Linkedin_link: legend["Linkedin_link"],
  }
  return &person
}

func OldNewPerson(conn *pgxpool.Pool, legend map[string]string)  *Person{
  birthdate, err := time.Parse("2006-01-02", legend["Birthdate"])
  if err!=nil {
    log.Print(err)
  }
  deceased, err := time.Parse("2006-01-02", legend["Deceased"])
  if err!=nil {
    log.Print(err)
  }
  Bio_mother, err := GetPersonByUUID(conn, legend["Bio_mother_UUID"])
  if err!= nil{
    log.Print(err)
  }
  Bio_father, err := GetPersonByUUID(conn, legend["Bio_father_UUID"])
  if err!= nil{
    log.Print(err)
  }
  Created_at, err := time.Parse("2006-01-02", legend["Created_at"])
  if err!= nil{
    log.Print(err)
  }
  person := Person{
    UUID: legend["UUID"],
    Created_at: Created_at,
    Givenname: legend["Givenname"],
    Surname: legend["Surname"],
    Birthdate: birthdate, 
    Deceased: deceased,
    Bio_mother: Bio_mother,
    Bio_father: Bio_father,
    Linkedin_link: legend["Linkedin_link"],
  }
  return &person
}

func PersonFromRow(conn *pgxpool.Pool, row pgx.Row) (*Person, error) {
  startTime := time.Now()
  var UUID, Givenname, Surname, Bio_mother_UUID, Bio_father_UUID, Linkedin_link *string
  var Created_at, Birthdate, Deceased *time.Time
  err := row.Scan(&UUID, &Created_at, &Givenname, &Surname, &Birthdate, &Deceased, &Bio_mother_UUID, &Bio_father_UUID, &Linkedin_link)
  if err!= nil {
    return nil, err
  }
  m := map[string]string{
    "UUID":*UUID,
    "Created_at":*timeNilCheck(Created_at),
    "Givenname":*Givenname,
    "Surname":*Surname,
    "Birthdate":*timeNilCheck(Birthdate),
    "Deceased":*timeNilCheck(Deceased),
    "Bio_mother_UUID":*stringNilCheck(Bio_mother_UUID),
    "Bio_father_UUID":*stringNilCheck(Bio_father_UUID),
    "Linkedin_link":*stringNilCheck(Linkedin_link),
  }

  person := OldNewPerson(conn, m)
  fmt.Println("PersonScanner: ", time.Now().Sub(startTime))
  return person, nil 
}

func GetPersonByUUID(conn *pgxpool.Pool, UUIDtoSearch string) (*Person, error) {
  startTime := time.Now()
  row, err := conn.Query(context.Background(), "SELECT uuid, Created_at, givenname, surname, birthdate, deceased, bio_mother_uuid, bio_father_uuid, linkedin_link FROM person WHERE uuid=$1", UUIDtoSearch)

  person, err := PersonFromRow(conn, row)
  if err != nil{
    return nil, err
  }

  fmt.Println("getPersonByUUID: ", time.Now().Sub(startTime))
  return person, nil
}
