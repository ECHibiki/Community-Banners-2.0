package controllers

import(
  "os"
  "time"
  "image"
  "os/exec"
  "strings"
  "strconv"
  "mime/multipart"

  "github.com/alessio/shellescape"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)


func getUserData(name string) ([]map[string]string ){
  banners ,err := bannerdb.Query(`
      SELECT fk_name , uri , url , clicks, size FROM ads
      WHERE fk_name = ? ORDER BY id DESC
    ` , []interface{}{name})
  if err != nil{
    panic(err)
  }
  return banners;
}

func affirmImageIsOwned(name string, uri string) bool{
  banner , err := bannerdb.Query(`
    SELECT * FROM ads WHERE fk_name = ? AND uri = ?
  ` , []interface{}{name , uri})
  if err != nil{
    panic(err)
  }
  return len(banner) == 1
}

func removeAdImage(uri string){
  err := os.Remove(controller_settings.PublicPath + uri)
  if err != nil{
    panic(err)
  }
}
func removeAdSQL(uri string , name string){
  _ , err := bannerdb.Query(`
    DELETE FROM ads WHERE uri = ? AND fk_name = ?
  ` , []interface{}{uri , name})
  if err != nil{
    panic(err)
  }
}

func doCreationCooldown(ip string) int{
  cooldown := time.Now().Unix() - controller_settings.BannerInterval * 60
  rows , err := bannerdb.Query(`
    SELECT * FROM antispam WHERE ip = ? AND type="banner" AND unix >= ?
    ORDER BY unix ASC
  ` , []interface{}{ ip , cooldown} )
  if err != nil{
    panic (err)
  }
  timer := 0
  if len(rows) != 0{
    u , _ := strconv.Atoi(rows[0]["unix"])
    timer = int(cooldown) - u
  }
  return timer
}

func updateCreationCooldown(ip string){
  time_now := time.Now().Unix()
  cooldown := time_now - controller_settings.BannerInterval * 60
  _ , err := bannerdb.Query(`
    DELETE FROM antispam WHERE ip = ? AND unix < ? AND type="banner"
  ` , []interface{}{ ip , cooldown} )
  if err != nil{
    panic (err)
  }
  _ , err = bannerdb.Query(`
    INSERT INTO antispam VALUES (? , ? , ?)
  ` , []interface{}{ip , time_now , "banner"} )
  if err != nil{
    panic (err)
  }
}

func getImageDimensions(mem_image multipart.File) (int , int){
  decode_image, _, err := image.Decode(mem_image)
	if err != nil {
		panic(err)
	}
	bounds := decode_image.Bounds()
	return bounds.Max.X , bounds.Max.Y
}

func checkUniqueBanner(tmp_path string) bool{
  safe_name := shellescape.Quote(tmp_path)
  blockhash_cmd := exec.Command("blockhash", safe_name)
  hash_result, err := blockhash_cmd.Output()
  if err != nil{
    panic(err)
  }
  match , err := bannerdb.Query(`
    SELECT * FROM ads WHERE hash = ?
  ` , []interface{}{hash_result})
  if err != nil{
    panic(err)
  }
  return len(match) == 0
}

func generateHashedFileName(tmp_path string) string{
  safe_name := shellescape.Quote(tmp_path)
  sha_cmd := exec.Command("sha1sum", safe_name)
  sha1_and_file_result, err := sha_cmd.Output()
  if err != nil{
    panic(err)
  }
  sha1 := strings.Split(string(sha1_and_file_result) , " ")
  if len(sha1) != 2 {
    panic("sha1 hash issue")
  }
  return sha1[0]
}

func getBase64(path string) string{
  safe_name := shellescape.Quote(path)
  b64_cmd := exec.Command("base64", safe_name)
  b64_or_err, err := b64_cmd.Output()
  if err != nil{
    panic(err)
  }
  b64 := strings.Split(string(b64_or_err) , " ")
  if len(b64) != 1 {
    panic("Base64 encode issue")
  }
  return b64[0]
}

func addBanner(name string, file_path string, url string, ip string,  size string){
  _ , err := bannerdb.Query(`
    INSET INTO ads VALUES (? , ? , ? , ? , ?)
  ` , []interface{}{name , file_path , url , ip , size })
  if err != nil{
    panic(err)
  }
}