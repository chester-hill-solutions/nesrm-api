package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

  "github.com/chester-hill-solutions/nesrm_api/models"
)

type testPerson struct {
  UUID string `json:"id"`
  Givenname string `json:"givenname"`
  Surname string `json:"surname"`  
}
var testPersons = []testPerson{
  {UUID: "3fdvbh4578e", Givenname: "Saihaan", Surname: "Syed"},  {UUID: "jhb436bhfbd", Givenname: "Ish", Surname: "Dur"},
  {UUID: "sdfgy4378wf", Givenname: "David", Surname: "Attenborough"},
}
func getTestPersons(c *gin.Context)  {
  c.IndentedJSON(http.StatusOK, testPersons)
}

func main()  {
  ginRouter := gin.Default()
  ginRouter.GET("/testPersons", getTestPersons)

  ginRouter.GET("/persons/all", models.RepondGetPersonAll)

  ginRouter.Run("localhost:8000")
}
