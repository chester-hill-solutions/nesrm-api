package models

import (
	"context"
	"fmt"
	"net/url"

	//"github.com/jackc/pgx/v5/pgx"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Connection() (*pgx.Conn, error){
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
  fmt.Println(dsn.String())
  //try and connect
  conn, err := pgx.Connect(context.Background(), dsn.String())
  if err != nil{
    return nil, err
  } 
  return conn, nil
} 
