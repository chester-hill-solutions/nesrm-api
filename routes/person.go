package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	dbConnector "github.com/chester-hill-solutions/nesrm_api/conn"
	"github.com/chester-hill-solutions/nesrm_api/models"
	"github.com/gin-gonic/gin"
)

func RespondGetPersonAll(c *gin.Context)  {
  startTime := time.Now()
  //establish Connection
  conn, err := dbConnector.Connection()
  if err != nil{
    log.Fatal(err) 
  }
  defer conn.Close()

  //QUERY ROWS
  rows, err := conn.Query(context.Background(),"SELECT * FROM person")
  if err != nil {
    log.Fatal(err)
  } 
  defer rows.Close()
  //r, err := pgx.CollectRows(rows, pgx.RowToStructByName[person])
  
  //UNMARSHALL INTO STRUCTS
  persons := []models.Person{}
  for rows.Next() {
    person, err := models.PersonFromRow(conn, rows) 
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
