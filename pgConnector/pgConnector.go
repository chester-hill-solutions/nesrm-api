package pgConnector

import (
	"context"
	"log"
	"net/url"
	"time"

	//"github.com/jackc/pgx/v5/pgx"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)



func ConnectionPool() (*pgxpool.Pool, error){
  startTime := time.Now()
  log.Println("Establishing Connection Pool")
  //Load enviroment variables
  myEnv, err := godotenv.Read(".env")
  if err != nil{
    return nil, err 
  }
  // build connPoolection string
  dsn := url.URL{
    Scheme: myEnv["SCHEME"],
    Host: myEnv["HOST"],
    User: url.UserPassword(myEnv["USER"], myEnv["PASSWORD"]),
    Path: myEnv["DBNAME"],
  }
  q := dsn.Query()
  q.Add("sslmode", "disable")
 dsn.RawQuery = q.Encode()
  //try and connPoolect
  connPool, err := pgxpool.New(context.Background(), dsn.String())
  if err != nil{
    return nil, err
  } else {
    log.Println("Connection established")
  } 
  log.Println("Connection time: ", time.Now().Sub(startTime))
  return connPool, nil
} 

func TimeNilCheck(t *time.Time) *string{
  if t == nil {
    s := ""
    p := &s
    return p 
  } else{
    s := t.String()
    p := &s
    return p
  }
}

func StringNilCheck(s *string) *string{
  if s == nil {
    n := ""
    p := &n
    return p
  } else {
    return s
  }
}

func TimeParser(s *string) *time.Time{
  r, err := time.Parse("2006-01-02", *s)
  if err != nil{
    log.Print(err)
  }
  return &r
}

func BulkTimeParser(s []string) map[string]time.Time{
  m := make(map[string]time.Time)
  for _, v := range s{
    m[v] = *TimeParser(&v)
  }
  return m
}
