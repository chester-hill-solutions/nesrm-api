package models

import (
	"errors"
	"time"
)

type Organization struct{
  UUID string
  Created_at time.Time
  Name string
  Linkedin_link string
}
func NewOrganization(m map[string]interface{})  (*Organization, error){
  UUID, _ := m["UUID"].(string)
  Created_at_string, _ := m["UUID"].(string)
  Created_at, err := time.Parse("2006-01-02", Created_at_string) 
  if err != nil {
    Created_at = time.Time{}
  }
  Name, ok := m["name"].(string)
  if !ok {
    return nil, errors.New("Missing a name to add") 
  }
  Linkedin_link, ok := m["linkedin_link"].(string)
  if !ok {
    return nil, errors.New("Missing a linkedin_link to add. Orgainzation and experience tracking is relies heavily on linkedin, so please provide a linkedin list")
  }
  organization := Organization{
    UUID: UUID,
    Created_at: Created_at,
    Name: Name,
    Linkedin_link: Linkedin_link,
  }
  return &organization, nil
}


