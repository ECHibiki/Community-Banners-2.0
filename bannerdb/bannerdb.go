package bannerdb

import (
    "fmt"
    "sync"
    "database/sql"
    "encoding/json"
    "io/ioutil"

    _ "github.com/go-sql-driver/mysql"
)

type DBSettings struct{
  User string
  Pass string
  Address string
  Database string
}

var lock = &sync.Mutex{}
var db_connection *sql.DB
var db_settings DBSettings

func Init(){
    fmt.Println("\nDatabase connection initialization...")

    setting_json_bytes, err := ioutil.ReadFile("./settings/.db-settings.json")
    if err != nil{
      panic(err)
    }
    json.Unmarshal(setting_json_bytes, &db_settings)

    initMysqlDBConnection()

    fmt.Println("...Database connection initialized")
}

func initMysqlDBConnection(){
  if db_connection == nil {
  	lock.Lock()
    if db_connection == nil {
      fmt.Println("\tBegining DB connection...")
      db , err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", db_settings.User, db_settings.Pass, db_settings.Address, db_settings.Database))
      if err != nil {
        panic(err)
      }
      db_connection = db
      db_connection.SetConnMaxLifetime(0)
      db_connection.SetConnMaxIdleTime(0)
      db_connection.SetMaxOpenConns(0)
      db_connection.SetMaxIdleConns(5)
      fmt.Println("\t" , db_connection.Stats())
    }
    lock.Unlock()
  }
}
