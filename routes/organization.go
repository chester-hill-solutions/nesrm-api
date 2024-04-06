package routes

import (
	"context"
	"errors"
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

func HandleGetOrganizationByUUID(c *gin.Context){
  startTime := time.Now()
  log.Println("HandleGetOrganizationByUUID")
  //Connection
  connPool, err := pgConnector.ConnectionPool()
  if err != nil {
    log.Println(err)
    c.IndentedJSON(http.StatusInternalServerError, err.Error())
    return
  }
  defer connPool.Close()
  //Validate Request
  var requestBodyMap map[string]interface{}
  err = c.ShouldBind(&requestBodyMap)
  //if err != nil {
  //  log.Println(err)
  //  c.IndentedJSON(http.StatusBadRequest, err.Error())
  //  return
  //}
  requestIsValid, err := ValidateRequestContent(requestBodyMap)
  if requestIsValid == false{
    c.IndentedJSON(http.StatusUnprocessableEntity, err)
  }

  //Logic
  organization, err := GetOrganizationByUUID(connPool, c.Param("UUID"))
  if err!=nil {
    if err.Error() == "No resources found" {
      log.Println(err)
      c.IndentedJSON(http.StatusNotFound, organization)
      return
    }
    c.IndentedJSON(http.StatusBadRequest, err.Error())
    log.Println(err)
    return
  } else{
    log.Println("org: ", organization)
    c.IndentedJSON(http.StatusOK, organization)
  }
  simplygolog.SaveTime("HandleGetOrganizationByUUIDGetOrganizationByUUID", time.Since(startTime))
}

func GetOrganizationByUUID(connPool *pgxpool.Pool, UUID string)  (*models.Organization, error){
  startTime := time.Now()
  log.Println("enter", UUID)
  args := pgx.NamedArgs{
    "UUID":UUID,
  }
  rows, err := connPool.Query(context.Background(), "SELECT * FROM organization WHERE uuid = @UUID", args)
  if err != nil {
    log.Println("Query err:", err)
    return nil, err
  }
  organizations, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Organization])
  if err != nil {
    return nil, err
  }
  if len(organizations) == 0{
    return nil, errors.New("No resources found")
  }
  org := organizations[0]
  simplygolog.PrintSaveTime("GetOrganizationByUUID", time.Since(startTime))
  return &org, nil
}

func HandlePostOrganization(c *gin.Context)  {
  startTime := time.Now()
  log.Println("Enter HandlePostOrganization")
  //Connection
  connPool, err := pgConnector.ConnectionPool()
  if err != nil{
    log.Println(err)
    c.IndentedJSON(http.StatusInternalServerError, err.Error())
    return
  }
  defer connPool.Close()
  //validate request content
  var requestBodyMap map[string]interface{}
  err = c.ShouldBind(&requestBodyMap)
  if err != nil {
    log.Println(err)
    c.IndentedJSON(http.StatusBadRequest, err.Error())
    return
  }
  requestIsValid, err := ValidateRequestContent(requestBodyMap)
  if requestIsValid == false{
    log.Println(err)
    c.IndentedJSON(http.StatusUnprocessableEntity, err.Error())
    return
  }

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
  simplygolog.PrintSaveTime("HandlePostOrganization", time.Since(startTime))
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
