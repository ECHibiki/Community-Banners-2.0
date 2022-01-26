package main

import (
  "fmt"

  "github.com/ECHibiki/Community-Banners-2.0/gin_interface"
)

func main(){
  fmt.Println("Community-Banners-2.0 - ECVerniy 2022")
  port := ":4200"
  gin_interface.Init(port)

}