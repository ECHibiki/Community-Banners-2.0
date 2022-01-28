package controllers

import (
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func getAllEntries() []map[string]string{
  ent , err := bannerdb.Query(`
    SELECT ads.fk_name, url, uri, bans.hardban, size, clicks FROM ads
    LEFT JOIN bans ON ads.fk_name = bans.fk_name
    ORDER BY id DESC
  ` , []interface{})
  if err != nil {
    panic(err)
  }
  return ent
}