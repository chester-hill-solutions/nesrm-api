package models

import (
	"time"

	"github.com/chester-hill-solutions/nesrm_api/pgConnector"
)

type Organization struct{
  UUID string
  Created_at time.Time
  Name string
  Linkedin_link string
}

func NewOrganization(m map[string]string)  *Organization{
  parsedTimeValues := pgConnector.BulkTimeParser([]string{m["Created_at"]}) 
  organization := Organization{
    UUID: m["UUID"],
    Created_at: parsedTimeValues["Created_at"],
    Name: m["Name"],
    Linkedin_link: m["Linkedin_link"],
  }
  return &organization
}


