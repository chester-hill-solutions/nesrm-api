package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/chester-hill-solutions/nesrm_api/models"
	"github.com/chester-hill-solutions/nesrm_api/pgConnector"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RespondGetPersonAll(c *gin.Context)  {
  log.Println("Responding GetPersonAll")
  startTime := time.Now()
  //establish Connection
  connPool, err := pgConnector.ConnectionPool()
  if err != nil{
    log.Fatal(err) 
  }
  defer connPool.Close()

  //QUERY ROWS
  rows, err := connPool.Query(context.Background(),`SELECT base.uuid base_uuid, base.created_at base_created_at, base.givenname base_givenname, base.surname base_surname, base.birthdate base_birthdate, base.deceased base_deceased, base.bio_mother_uuid base_bio_mother_uuid, base.bio_father_uuid base_bio_father_uuid, base.linkedin_link base_linkedin_link,
mother.uuid mother_uuid, mother.created_at mother_created_at, mother.givenname mother_givenname, mother.surname mother_surname, mother.birthdate mother_birthdate, mother.deceased mother_deceased, mother.bio_mother_uuid mother_bio_mother_uuid, mother.bio_father_uuid mother_bio_father_uuid, mother.linkedin_link mother_linkedin_link,
father.uuid father_uuid, father.created_at father_created_at, father.givenname father_givenname, father.surname father_surname, father.birthdate father_birthdate, father.deceased father_deceased, father.bio_mother_uuid father_bio_mother_uuid, father.bio_father_uuid father_bio_father_uuid,father.linkedin_link father_linkedin_link
FROM person AS base
LEFT JOIN person AS mother ON base.bio_mother_uuid = mother.uuid
LEFT JOIN person as father ON base.bio_father_uuid = father.uuid;`)
  if err != nil {
    log.Fatal(err)
  } 
  defer rows.Close()
  //r, err := pgx.CollectRows(rows, pgx.RowToStructByName[person])
  log.Println(len(rows.FieldDescriptions()))

  //UNMARSHALL INTO STRUCTS
  persons := []models.Person{}
  for rows.Next() {
    person, err := PersonFromRow(connPool, rows) 
    if err != nil {
      //fmt.Printf("%+v\n", person)
      log.Println(err) 
    }
    persons = append(persons, *person)
  }
  executionTime := time.Now().Sub(startTime)
  log.Println("RespondGetPersonAll Execution Time: ", executionTime)
  c.IndentedJSON(http.StatusOK, persons) 
}

func PersonFromRow(connPool *pgxpool.Pool, row pgx.Row) (*models.Person, error) {
  startTime := time.Now()
  log.Println("enter PersonFromRow")
  var base_UUID, base_givenname, base_surname, base_bio_mother_UUID, base_bio_father_UUID, base_linkedin_link *string
  var mother_UUID, mother_givenname, mother_surname, mother_bio_mother_UUID, mother_bio_father_UUID, mother_linkedin_link *string
  var father_UUID, father_givenname, father_surname, father_bio_mother_UUID, father_bio_father_UUID, father_linkedin_link *string
  var base_created_at, base_birthdate, base_deceased *time.Time
  var mother_created_at, mother_birthdate, mother_deceased *time.Time
  var father_created_at, father_birthdate, father_deceased *time.Time
  err := row.Scan(
    &base_UUID, &base_created_at, &base_givenname, &base_surname, &base_birthdate, &base_deceased, &base_bio_mother_UUID, &base_bio_father_UUID, &base_linkedin_link,
    &mother_UUID, &mother_created_at, &mother_givenname, &mother_surname, &mother_birthdate, &mother_deceased, &mother_bio_mother_UUID, &mother_bio_father_UUID, &mother_linkedin_link,
    &father_UUID, &father_created_at, &father_givenname, &father_surname, &father_birthdate, &father_deceased, &father_bio_mother_UUID, &father_bio_father_UUID, &father_linkedin_link,
  )
  if err!= nil {
    return nil, err
  }
  m := map[string]string{
    "UUID":*base_UUID,
    "Created_at":*pgConnector.TimeNilCheck(base_created_at),
    "Givenname":*base_givenname,
    "Surname":*base_surname,
    "Birthdate":*pgConnector.TimeNilCheck(base_birthdate),
    "Deceased":*pgConnector.TimeNilCheck(base_deceased),
    "Bio_mother_UUID":*pgConnector.StringNilCheck(base_bio_mother_UUID),
    "Bio_father_UUID":*pgConnector.StringNilCheck(base_bio_father_UUID),
    "Linkedin_link":*pgConnector.StringNilCheck(base_linkedin_link),
    "mother_UUID":*pgConnector.StringNilCheck(mother_UUID),
    "mother_Created_at":*pgConnector.TimeNilCheck(mother_created_at),
    "mother_Givenname":*pgConnector.StringNilCheck(mother_givenname),
    "mother_Surname":*pgConnector.StringNilCheck(mother_surname),
    "mother_Birthdate":*pgConnector.TimeNilCheck(mother_birthdate),
    "mother_Deceased":*pgConnector.TimeNilCheck(mother_deceased),
    "mother_Bio_mother_UUID":*pgConnector.StringNilCheck(mother_bio_mother_UUID),
    "mother_Bio_father_UUID":*pgConnector.StringNilCheck(mother_bio_father_UUID),
    "mother_Linkedin_link":*pgConnector.StringNilCheck(mother_linkedin_link),
    "father_UUID":*pgConnector.StringNilCheck(father_UUID),
    "father_Created_at":*pgConnector.TimeNilCheck(father_created_at),
    "father_Givenname":*pgConnector.StringNilCheck(father_givenname),
    "father_Surname":*pgConnector.StringNilCheck(father_surname),
    "father_Birthdate":*pgConnector.TimeNilCheck(father_birthdate),
    "father_Deceased":*pgConnector.TimeNilCheck(father_deceased),
    "father_Bio_mother_UUID":*pgConnector.StringNilCheck(father_bio_father_UUID),
    "father_Bio_father_UUID":*pgConnector.StringNilCheck(father_bio_father_UUID),
    "father_Linkedin_link":*pgConnector.StringNilCheck(father_linkedin_link),
  }

  person := models.NewPerson(m)
  log.Println("PersonScanner: ", time.Now().Sub(startTime))
  return person, nil 
}
