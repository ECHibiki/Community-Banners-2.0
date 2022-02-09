package controllers

import (
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func checkURIExists(file string) bool{
  check , err := bannerdb.Query("SELECT id FROM ads WHERE uri = ?" , []interface{}{file})
  if err != nil{
    panic(err)
  }
  return len(check) != 0
}

func incrementClicksSQL(file string){
  _, err := bannerdb.Query("UPDATE ads SET clicks = ( clicks + 1 ) WHERE uri = ?", []interface{}{file})
  if err != nil{
    panic(err)
  }
}

