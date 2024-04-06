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
	"github.com/shopspring/decimal"
)

func HandleGetCampaignByUUID(c *gin.Context)  {
  startTime := time.Now()
  log.Println("HandleGetCampaignByUUID")
  //Connection
  connPool, err := pgConnector.ConnectionPool()
  if err != nil{
    log.Println(err)
    c.IndentedJSON(http.StatusInternalServerError, err.Error())
    return
  }
  defer connPool.Close()
  //Valuidte Request
  var requestBodyMap map[string]interface{}
  err = c.ShouldBind(&requestBodyMap)
  requestIsValid, err := ValidateRequestContent(requestBodyMap)
  if !requestIsValid {
    log.Println(err)
    c.IndentedJSON(http.StatusUnprocessableEntity, err.Error())
  }

  //Logic
  campaign, err := GetCampaignByUUID(connPool, c.Param("UUID"))
  if err != nil {
    if err.Error() == "No resources found" {
      log.Println(err)
      c.IndentedJSON(http.StatusNotFound, campaign)
      return
    }
    log.Println(err)
    c.IndentedJSON(http.StatusBadRequest, err.Error())
    return
  }
  c.IndentedJSON(http.StatusOK, *campaign)
  simplygolog.PrintSaveTime("HandleGetCampaignByUUID",time.Since(startTime))
}

func GetCampaignByUUID(connPool *pgxpool.Pool, UUID string) (*models.Campaign, error)  {
  startTime := time.Now()
  log.Println("GetCampaignByUUID")
  query := "SELECT uuid, created_at, start_date, election_date, name, candidate_uuid, governance_level, campaign_type, points_accrued, total_points_cast FROM campaign WHERE uuid = @UUID"
  args := pgx.NamedArgs{
    "UUID":UUID,
  }
  rows, err := connPool.Query(context.Background(), query, args)
  if err != nil {
    log.Println("Query err:", err)
    return nil, err
  }
  defer rows.Close()
  campaigns := []models.Campaign{}
  for rows.Next(){
    campaign, err := RowToCampaign(connPool, rows)
    if err != nil {
      log.Println(err)
      return nil, err
    }
    campaigns = append(campaigns, *campaign)
  }
  if len(campaigns) == 0 {
    return nil, errors.New("No resources found")
  }
  simplygolog.PrintSaveTime("GetCampaignByUUID", time.Since(startTime))
  campaign := campaigns[0]
  return &campaign, nil
}

func RowToCampaign(connPool *pgxpool.Pool, rows pgx.Row) (*models.Campaign, error) {
  var UUID, Name, Candidate_uuid, Governance_level, Campaign_type *string
  var Created_at, Start_date, Election_date *time.Time
  var Points_accrued, Total_points_cast *decimal.Decimal
  err := rows.Scan(&UUID, &Created_at, &Start_date, &Election_date, &Name, &Candidate_uuid, &Governance_level, &Campaign_type, &Points_accrued, &Total_points_cast)
  if err != nil {
    log.Println(err)
    return nil, err
  }
  campaign := models.Campaign{
    UUID: *UUID,
    Created_at: *pgConnector.TimeNilCheck(Created_at),
    Start_date: *pgConnector.TimeNilCheck(Start_date),
    Election_date: *pgConnector.TimeNilCheck(Election_date),
    Name: *pgConnector.StringNilCheck(Name),
    Candidate_uuid: *pgConnector.StringNilCheck(Candidate_uuid),
    Governance_level: *pgConnector.StringNilCheck(Governance_level),
    Campaign_type: *pgConnector.StringNilCheck(Campaign_type),
    Points_accrued: *pgConnector.DecimalNilCheck(Points_accrued),
    Total_points_cast: *pgConnector.DecimalNilCheck(Total_points_cast),
  }
  candidate, err := GetPersonByUUID(connPool, campaign.Candidate_uuid)
  if err != nil {
    log.Println(err)
  } else {
    campaign.Candidate = candidate
  }
  return &campaign, nil
}
