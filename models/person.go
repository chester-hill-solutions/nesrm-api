package models

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
func newPerson(conn *pgxpool.Pool, legend map[string]string)  *Person{
  birthdate, err := time.Parse("2006-01-02", legend["Birthdate"])
  if err!=nil {
    log.Print(err)
  }
  deceased, err := time.Parse("2006-01-02", legend["Deceased"])
  if err!=nil {
    log.Print(err)
  }
  Bio_mother, err := getPersonByUUID(conn, legend["Bio_mother_UUID"])
  if err!= nil{
    log.Print(err)
  }
  Bio_father, err := getPersonByUUID(conn, legend["Bio_father_UUID"])
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

func personScanner(conn *pgxpool.Pool, row pgx.Row) (*Person, error) {
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

  person := newPerson(conn, m)
  fmt.Println("personScanner: ", time.Now().Sub(startTime))
  return person, nil 
}

func getPersonByUUID(conn *pgxpool.Pool, UUIDtoSearch string) (*Person, error) {
  startTime := time.Now()
  row := conn.QueryRow(context.Background(), "SELECT uuid, Created_at, givenname, surname, birthdate, deceased, bio_mother_uuid, bio_father_uuid, linkedin_link FROM person WHERE uuid=$1", UUIDtoSearch)

  person, err := personScanner(conn, row)
  if err != nil{
    return nil, err
  }

  fmt.Println("getPersonByUUID: ", time.Now().Sub(startTime))
  return person, nil
}

func RepondGetPersonAll(c *gin.Context){
  conn, err := Connection() 
  if err != nil{
    log.Fatal(err) 
  }
  defer conn.Close()
  persons, err := getPersonAll(conn)
  if err != nil {
    log.Fatal(err)
  }
  c.IndentedJSON(http.StatusOK, persons) 
}

func getPersonAll(conn *pgxpool.Pool) (*[]Person, error) {
  startTime := time.Now()
  rows, err := conn.Query(context.Background(),"SELECT * FROM person")
  if err != nil {
    return nil, err
  } 
  defer rows.Close()
  //r, err := pgx.CollectRows(rows, pgx.RowToStructByName[person])
  persons := []Person{}
  for rows.Next() {
    person, err := personScanner(conn, rows) 
    if err != nil {
      //fmt.Printf("%+v\n", person)
      log.Println(err) 
    }
    persons = append(persons, *person)
  }
  executionTime := time.Now().Sub(startTime)
  fmt.Println("enter persons ended: ", executionTime)
  return &persons, err
}
