package controllers

import (
  "time"
  "golang.org/x/crypto/bcrypt"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func validateNameBruteForce(ip string) bool{
  // lock out for N + oldest_attempt minutes if entered to table 5 times
  as , err := bannerdb.Query(`
    SELECT * FROM antispam
    WHERE name = ? AND type = "login" AND unix >= ?
  ` , []interface{}{
    ip , time.Now().Unix() - controller_settings.AttemptInterval * 60,
  })
  if err != nil{
    panic (err)
  }
  if len( as ) > controller_settings.MaxAttempts {
    return false
  }

  _ , err = bannerdb.Query(`
    INSERT INTO antispam VALUES ( ? , ? , ?)
  ` , []interface{}{
    ip , time.Now().Unix() , "login",
  })
  if err != nil{
    panic (err)
  }
  return true
}

func updateNameBruteForce(ip string){
  // clear entries on table
  _ , err := bannerdb.Query(`
    DELETE FROM antispam WHERE ip = ? AND type = "login"
  ` , []interface{}{ ip })
  if err != nil{
    panic (err)
  }
}

func checkAuthentication(name string , pass string) bool{
  acc, err := bannerdb.Query(`
    SELECT * FROM users WHERE name = ?
  ` , []interface{}{name})
  if err != nil{
    panic (err)
  }
  if len (acc) == 0{
    return false
  }
  user := acc[0]
  hashed_password := user["pass"]
  pass_err := bcrypt.CompareHashAndPassword([]byte(hashed_password) , []byte(pass))
  return pass_err == nil
}

func checkHardBanned(name string) bool{
  ban_get , err := bannerdb.Query(`
    SELECT fk_name FROM bans
    WHERE fk_name = ? AND hardban=1
  `, []interface{}{name})
  if err != nil{
    panic(err)
  }
  return len(ban_get) > 0
}