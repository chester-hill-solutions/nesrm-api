package models

import (
	"context"
	"net/url"
	"time"

	//"github.com/jackc/pgx/v5/pgx"
	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)



func Connection() (*pgxpool.Pool, error){
  //Load enviroment variables
  myEnv, err := godotenv.Read(".env")
  if err != nil{
    return nil, err 
  }
  // build connection string
  dsn := url.URL{
    Scheme: myEnv["SCHEME"],
    Host: myEnv["HOST"],
    User: url.UserPassword(myEnv["USER"], myEnv["PASSWORD"]),
    Path: myEnv["DBNAME"],
  }
  q := dsn.Query()
  q.Add("sslmode", "disable")
 dsn.RawQuery = q.Encode()
  //try and connect
  conn, err := pgxpool.New(context.Background(), dsn.String())
  if err != nil{
    return nil, err
  } 
  return conn, nil
} 

func timeNilCheck(t *time.Time) *string{
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

func stringNilCheck(s *string) *string{
  if s == nil {
    n := ""
    p := &n
    return p
  } else {
    return s
  }
}
