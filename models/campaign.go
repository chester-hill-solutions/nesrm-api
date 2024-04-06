package models

import (
  "time"
  "github.com/shopspring/decimal"
)

type Campaign struct{
  UUID string 
  Created_at time.Time 
  Start_date time.Time 
  Election_date time.Time 
  Governance_level string 
  Name string
  Campaign_type string 
  Candidate_uuid string 
  Candidate *Person
  Points_accrued decimal.Decimal 
  Total_points_cast decimal.Decimal 
}

func NewCampaign(UUID string, Created_at time.Time, Start_date time.Time, Election_date time.Time, Governance_level string, Name string, Campaign_type string, Candidate_uuid string, Points_accrued decimal.Decimal, Total_points_cast decimal.Decimal) *Campaign {
 return &Campaign{
    UUID: UUID,
    Created_at: Created_at,
    Start_date: Start_date,
    Election_date: Election_date,
    Governance_level: Governance_level,
    Name: Name,
    Campaign_type: Campaign_type,
    Candidate_uuid: Candidate_uuid,
    Points_accrued: Points_accrued,
    Total_points_cast: Total_points_cast,
  }
}
