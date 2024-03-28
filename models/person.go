package models

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type person struct{
  UUID string
  Created_at time.Time
  Givenname string
  Surname string
  Birthdate time.Time
  Deceased time.Time
  Bio_mother_UUID string
  Bio_father_UUID string
  Linkedin_link string
}
func (p *person) setProperty(propName string, propValue string) *person {
	reflect.ValueOf(p).Elem().FieldByName(propName).Set(reflect.ValueOf(propValue))
	return p
}
func RepondGetPersons(c *gin.Context){
  persons, err := getPersons()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("persons:", persons)
  c.IndentedJSON(http.StatusOK, persons) 
}

func getPersons() (*[]person, error) {
  fmt.Println("enter getPersons")
  conn, err := Connection() 
  if err != nil{
    return nil, err
  }
  defer conn.Close(context.Background())
  rows, err := conn.Query(context.Background(),"select * from person")
  if err != nil {
    return nil, err
  } 
  defer rows.Close()
  r, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[person])
  return &r, err
}
