package main

import (
	"net/http"

	r "github.com/gin-gonic/gin"

	"github.com/chester-hill-solutions/nesrm_api/routes"
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
func getTestPersons(c *r.Context)  {
  c.IndentedJSON(http.StatusOK, testPersons)
}

func main()  {
  router := r.Default()
  router.GET("/ping", func(ctx *r.Context) {
    ctx.IndentedJSON(http.StatusOK, map[string]string{"Hello":"From Sai",})
  })
  // /person
  router.GET("/person/all", routes.RespondGetPersonAll)
  router.GET("/person/:UUID", routes.RespondGetPersonByUUID)

  // /organization
  router.GET("/organization/:UUID", routes.RespondGetOrganizationByUUID)
  router.POST("/organization", routes.RespondPostOrganization)

  router.Run("localhost:8000")
}
