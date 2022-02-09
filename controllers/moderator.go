package controllers

import (
  "os"
  "fmt"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func init(){
  fmt.Println("")
}

func getAllEntries() []map[string]string{
  ent , err := bannerdb.Query(`
    SELECT ads.fk_name as name, url, uri, bans.hardban, size, clicks, board FROM ads
    LEFT JOIN bans ON ads.fk_name = bans.fk_name
    ORDER BY ads.id DESC
  ` , []interface{}{})
  if err != nil {
    panic(err)
  }
  return ent
}

func createNewBan(target string, hard_ban int , ip string){
  _ , err := bannerdb.Query(`
      INSERT INTO bans VALUES (NULL , ? , ? , ?)
    ` , []interface{}{target , hard_ban , ip})
  if err != nil{
    panic(err)
  }
}
func removeAllUserImages(target string){
  banners , err := bannerdb.Query(`
      SELECT * FROM ads WHERE fk_name = ?
    ` , []interface{}{target})
  if err != nil{
    panic(err)
  }
  for _ , banner := range banners{
    err = os.Remove(controller_settings.PublicPath + banner["uri"])
    if err != nil{
      fmt.Println(err)
    }
  }
}

func removeUserFromDatabase(target string){
  _ , err := bannerdb.Query(`
      DELETE FROM ads WHERE fk_name = ?
    ` , []interface{}{ target })
  if err != nil{
    panic(err)
  }
}

func removeIndividualBannerFromImages(uri string){
  err := os.Remove(controller_settings.PublicPath + uri)
  if err != nil{
    fmt.Println(err)
  }
}

func removeIndividualBannerFromDB(uri string){
  _ , err := bannerdb.Query(`
      DELETE FROM ads WHERE uri = ?
    ` , []interface{}{ uri })
  if err != nil{
    panic(err)
  }
}

func checkModLevel(name string) string{
  mod_get , err := bannerdb.Query(`
    SELECT fk_name , "admin" as level FROM mods
    WHERE fk_name = ?
  `, []interface{}{name})
  if err != nil{
    panic(err)
  }
  if len(mod_get) == 0{
    return "none"
  }
  return mod_get[0]["level"]
}