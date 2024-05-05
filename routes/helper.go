package routes

import (
	"errors"
	"log"
	"net/http"
)


func ValidateRequestContent(headers http.Header ) (bool, error) {
  if (len(headers["User"])>0)&&(len(headers["Service"])>0) {
    log.Println("Request Validated")
    return true, nil
  } else {
    return false, errors.New("Missing required headers: please include the user and the service this request is for")
  }
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
