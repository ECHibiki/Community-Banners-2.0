package controllers

import (
  "time"
  "golang.org/x/crypto/bcrypt"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func validateNameBruteForce(ip string) bool{
  // lock out for N + oldest_attempt minutes if entered to table 5 times
  timer := time.Now().Unix()
  cooldown := time.Now().Unix() + controller_settings.AttemptInterval * 60
  as , err := bannerdb.Query(`
    SELECT * FROM antispam
    WHERE ip = ? AND type = "login" AND unix >= ?
  ` , []interface{}{
    ip , timer ,
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
    ip , cooldown , "login",
  })
  if err != nil{
    panic (err)
  }
  return true
}

func updateNameBruteForce(){
  // clear entries on table
  timer := time.Now().Unix()
  _ , err := bannerdb.Query(`
    DELETE FROM antispam WHERE unix < ?
  ` , []interface{}{ timer })
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

func checkIsMod(name string) bool{
  mod_get , err := bannerdb.Query(`
    SELECT fk_name FROM mods
    WHERE fk_name = ?
  `, []interface{}{name})
  if err != nil{
    panic(err)
  }
  return len(mod_get) > 0
}
func checkIsDonor(token string) bool{
  token_get , err := bannerdb.CrossDBQuery(`
    SELECT token FROM %[1]s.whitelist_tokens
    WHERE token = ?
  `, []interface{}{token})
  if err != nil{
    panic(err)
  }
  return len(token_get) > 0
}

func checkHardBanned(name string , ip string) bool{
  ban_get , err := bannerdb.Query(`
    SELECT fk_name FROM bans
    WHERE (fk_name = ? OR ip = ?) AND hardban=1
  `, []interface{}{name, ip})
  if err != nil{
    panic(err)
  }
  return len(ban_get) > 0
}