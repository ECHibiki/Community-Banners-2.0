package bannerdb

import (
    "fmt"
    "sync"
    "database/sql"
    "encoding/json"
    "io/ioutil"

    _ "github.com/go-sql-driver/mysql"
)

// written for documentation purposes or future uses.
type DBSettings struct{
  User string
  Pass string
  Address string
  Database string
  CrossDatabase string
}

type AntiSpam struct {
  IP string
  Unix int32
  Type string
}
type User struct {
  Id int32
  Name string
  Pass string
}
type Ban struct {
  Id int32
  Name string
  HardBan bool
}
type Mod struct {
  Id int32
  Name string
}
type Ad struct {
  Id int32
  Name string
  URI string
  URL string
  IP string
  Size string
  Clicks int32
}

var lock = &sync.Mutex{}
var db_connection *sql.DB
var db_settings DBSettings

func init(){
    fmt.Println("\nDatabase connection initialization...")

    setting_json_bytes, err := ioutil.ReadFile("./settings/db-settings.json")
    if err != nil{
      panic(err)
    }
    err = json.Unmarshal(setting_json_bytes, &db_settings)
    if err != nil{
      panic(err)
    }
    initMysqlDBConnection()

    fmt.Println("...Database connection initialized")
}

// golang auto sanitizes
func Query(request string , parameters []interface{}) ([]map[string]string, error){
  // get similar entries of base thread id
   items,err := db_connection.Query(request , parameters...)
   if err != nil{
     return []map[string]string{}, err
   }
   columns, _ := items.Columns()
   var return_map []map[string]string
   for items.Next(){
     given_columns := make([]sql.NullString, len(columns))
     given_column_pointers := make([]interface{}, len(columns))
     for column_index, _ := range columns{
         given_column_pointers[column_index] = &given_columns[column_index]
     }
     err := items.Scan(given_column_pointers...)
     if err != nil{
       return []map[string]string{}, err
     }
     current_map_data := make(map[string]string)
     for index, column_name := range columns{
       if given_columns[index].Valid{
         current_map_data[column_name] = given_columns[index].String
       }
     }
     return_map = append(return_map, current_map_data)
   }
   return return_map, nil
}
func CrossDBQuery(request string , parameters []interface{}) ([]map[string]string, error){
  // get similar entries of base thread id
   items,err := db_connection.Query(fmt.Sprintf(request , db_settings.CrossDatabase) , parameters...)
   if err != nil{
     return []map[string]string{}, err
   }
   columns, _ := items.Columns()
   var return_map []map[string]string
   for items.Next(){
     given_columns := make([]sql.NullString, len(columns))
     given_column_pointers := make([]interface{}, len(columns))
     for column_index, _ := range columns{
         given_column_pointers[column_index] = &given_columns[column_index]
     }
     err := items.Scan(given_column_pointers...)
     if err != nil{
       return []map[string]string{}, err
     }
     current_map_data := make(map[string]string)
     for index, column_name := range columns{
       if given_columns[index].Valid{
         current_map_data[column_name] = given_columns[index].String
       }
     }
     return_map = append(return_map, current_map_data)
   }
   return return_map, nil
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
