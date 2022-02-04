package controllers

import (
  "os"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func getAllEntries() []map[string]string{
  ent , err := bannerdb.Query(`
    SELECT ads.fk_name, url, uri, bans.hardban, size, clicks FROM ads
    LEFT JOIN bans ON ads.fk_name = bans.fk_name
    ORDER BY id DESC
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
      SELECT * FROM ads WHERE name = ?
    ` , []interface{}{target})
  if err != nil{
    panic(err)
  }
  for _ , banner := range banners{
    err = os.Remove(controller_settings.PublicPath + banner["uri"])
    if err != nil{
      panic(err)
    }
  }
}

func removeUserFromDatabase(target string){
  _ , err := bannerdb.Query(`
      DELETE FROM ads WHERE name = ?
    ` , []interface{}{ target })
  if err != nil{
    panic(err)
  }
}

func removeIndividualBannerFromImages(uri string){
  err := os.Remove(controller_settings.PublicPath + uri)
  if err != nil{
    panic(err)
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