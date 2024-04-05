package routes

import (
	"github.com/gin-gonic/gin"
)

func ValidateRequestContent(c *gin.Context) (bool, error) {
 return true, nil
}

func WHERE(params map[string]interface{}) string{
  var output string
  conditionCount := 0
  for key, value := range params {
    conditionCount++
    if conditionCount <1{
      output = "WHERE"
    } else {
      output += " AND"
    }
    output = output + " " + key + " = " + value.(string)
  }
  return output
}
