package controllers

import (
  "fmt"
  "time"
  "regexp"
  "golang.org/x/crypto/bcrypt"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func init(){
  fmt.Print("")
}

func validateIPCreation(ip string) bool{
  cooldown := time.Now().Unix() - controller_settings.AccountInterval * 60
  rows , err := bannerdb.Query(`
    SELECT * FROM antispam WHERE ip = ? AND type="create" AND unix < ?
  ` , []interface{}{ ip , cooldown} )
  if err != nil{
    panic (err)
  }
  fmt.Println(rows, cooldown)
  return len(rows) == 0
}

func updateIPCreation(ip string){
  time_now := time.Now().Unix()
  cooldown := time_now - controller_settings.AccountInterval * 60
  _ , del_err := bannerdb.Query(`
    DELETE FROM antispam WHERE unix < ? AND type="create"
  ` , []interface{}{ cooldown } )
  if del_err != nil{
    panic (del_err)
  }
  _ , ins_err := bannerdb.Query(`
    INSERT INTO antispam VALUES ( ? , ? , ? )
  ` , []interface{}{ ip , time_now , "create"} )
  if ins_err != nil{
    panic (ins_err)
  }
}

func addNewUserToDB(name string, pass string) string{
  invalid_reg := regexp.MustCompile("/[^a-zA-Z0-9_\\-]/")
  if invalid_reg.MatchString(name) {
    return "Name has invalid characters(alphanumeric names only)"
  }
  hashed_bytes , berr := bcrypt.GenerateFromPassword([]byte(pass), 10)
  if berr != nil{
    panic (berr)
  }
  rows, err := bannerdb.Query(`
    SELECT name FROM users WHERE name = ?
  `,  []interface{}{ name })
  if err != nil{
    panic (err)
  }
  if len(rows) == 0{
    _, ins_err := bannerdb.Query(`
      INSERT INTO users VALUES (NULL , ? , ? )
    `,  []interface{}{ name , string(hashed_bytes) })
    if ins_err != nil{
      panic (ins_err)
    }
    return ""
  } else{
    return "Username Already Exists"
  }
}