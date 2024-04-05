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
	"github.com/sai-sy/simplygolog"
)

func RespondGetOrganizationByUUID(c *gin.Context){
  startTime := time.Now()
  //Connection
  connPool, err := pgConnector.ConnectionPool()
  if err != nil {
    log.Fatal(err)
  }
  defer connPool.Close()
  //Validate Request
  requestIsValid, err := ValidateRequestContent(c)
  if requestIsValid == false{
    c.IndentedJSON(http.StatusUnprocessableEntity, err)
  }

  log.Println("enter responsd get organization by uuid")
  //Logic
  organization, err := GetPersonByUUID(connPool, c.Param("UUID"))
  if err!=nil {
    c.IndentedJSON(http.StatusBadRequest, err)
    log.Println(err)
  } else{
    log.Println("org: ", organization)
    c.IndentedJSON(http.StatusOK, organization)
  }
  simplygolog.SaveTime("RespondGetPerson", time.Since(startTime))
}

func GetOrganizationByUUID(connPool *pgxpool.Pool, UUID string)  (*models.Organization, error){
  startTime := time.Now()
  log.Println("enter", UUID)
  rows, err := connPool.Query(context.Background(), "SELECT * FROM organization WHERE uuid = $1", UUID)
  if err != nil {
    return nil, err
  }
  organizations, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Organization])
  if err != nil {
    return nil, err
  }
  org := organizations[0]
  simplygolog.PrintSaveTime("GetOrganizationByUUID", time.Since(startTime))
  return &org, nil
  
}
