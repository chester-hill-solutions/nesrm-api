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
  var requestBodyMap map[string]interface{}
  err = c.BindJSON(&requestBodyMap)
  if err != nil {
    c.IndentedJSON(http.StatusBadRequest, err)
    return
  }
  requestIsValid, err := ValidateRequestContent(requestBodyMap)
  if requestIsValid == false{
    c.IndentedJSON(http.StatusUnprocessableEntity, err)
  }

  log.Println("enter responsd get organization by uuid")
  //Logic
  organization, err := GetOrganizationByUUID(connPool, c.Param("UUID"))
  if err!=nil {
    c.IndentedJSON(http.StatusBadRequest, err)
    log.Println(err)
  } else{
    log.Println("org: ", organization)
    c.IndentedJSON(http.StatusOK, organization)
  }
  simplygolog.SaveTime("RespondGetOrganizationByUUID", time.Since(startTime))
}

func GetOrganizationByUUID(connPool *pgxpool.Pool, UUID string)  (*models.Organization, error){
  startTime := time.Now()
  log.Println("enter", UUID)
  args := pgx.NamedArgs{
    "UUID":UUID,
  }
  rows, err := connPool.Query(context.Background(), "SELECT * FROM organization WHERE uuid = @UUID", args)
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

func RespondPostOrganization(c *gin.Context)  {
  startTime := time.Now()
  log.Println("Enter RespondPostOrganization")
  //Connection
  connPool, err := pgConnector.ConnectionPool()
  if err != nil{
    c.IndentedJSON(http.StatusInternalServerError, err)
    return
  }
  defer connPool.Close()
  //validate request content
  var requestBodyMap map[string]interface{}
  err = c.BindJSON(&requestBodyMap)
  if err != nil {
    c.IndentedJSON(http.StatusBadRequest, err)
    return
  }
  requestIsValid, err := ValidateRequestContent(requestBodyMap)
  if requestIsValid == false{
    c.IndentedJSON(http.StatusUnprocessableEntity, err)
    return
  }

  //unmarshall body into map
  log.Println(requestBodyMap)
  //Logic
  organization, err := PostOrganization(connPool, requestBodyMap)
  if err!=nil {
    c.IndentedJSON(http.StatusBadRequest, err)
    log.Println(err)
    return
  } else{
    log.Println("org: ", organization)
    c.IndentedJSON(http.StatusOK, organization)
  }
  simplygolog.PrintSaveTime("RespondPostOrganization", time.Since(startTime))
}
func PostOrganization(connPool *pgxpool.Pool, fields map[string]interface{}) (*models.Organization, error) {
  startTime := time.Now()
  log.Println("Enter PostOrganization")
  organization, err := models.NewOrganization(fields)
  if err != nil {
    return nil, err
  }
  query := `INSERT INTO organization(name, linkedin_link) VALUES(@organizationName, @organizationLinkedin_link)`
  args := pgx.NamedArgs{
    "organizationName": organization.Name,
    "organizationLinkedin_link":organization.Linkedin_link,
  }
  _, err = connPool.Exec(context.Background(), query, args)
  if err!=nil {
    log.Println(err)
    return nil, err
  }
  simplygolog.PrintSaveTime("PostOrganization", time.Since(startTime))
  return organization, nil
}
