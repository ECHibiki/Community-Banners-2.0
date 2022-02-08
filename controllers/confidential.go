package controllers

import(
  "os"
  "io"
  "fmt"
  "time"
  "image"
  "os/exec"
  "strings"
  "strconv"
  "mime/multipart"
  "net/http"
  _ "image/jpeg"
  _ "image/png"
  _ "image/gif"

  "github.com/alessio/shellescape"
  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)

func init(){
  fmt.Print("")
}

func getUserData(name string) ([]map[string]string ){
  banners ,err := bannerdb.Query(`
      SELECT fk_name , uri , url , clicks, size , board FROM ads
      WHERE fk_name = ? ORDER BY id DESC
    ` , []interface{}{name})
  if err != nil{
    panic(err)
  }
  return banners;
}


func getMimeType(image multipart.File) string{
  file_header := make([]byte, 512)
  // Copy the headers into the FileHeader buffer
  _, err := image.Read(file_header)
  if err != nil {
    panic(err)
  }
  image.Seek(0, io.SeekStart)
  return http.DetectContentType(file_header)
}
func writeTmpFile( name string , image multipart.File) string{
  fname := strings.ReplaceAll("tmp-" + strconv.Itoa(int(time.Now().Unix()))  + "-" + name , "." , "-", )
  fname = strings.ReplaceAll(fname , "/" , "-")
  tmp_location := controller_settings.PublicPath + "storage/tmp/" + fname
  handler, err := os.OpenFile(tmp_location,os.O_WRONLY|os.O_CREATE, 0775)
	if err != nil {
		return ""
	}
  defer handler.Close()
  _ , err = io.Copy(handler , image )
  if err != nil {
    panic(err)
  }
  return tmp_location
}

func getImageDimensions(mem_image multipart.File) (int , int){
  decode_image, _, err := image.Decode(mem_image)
	if err != nil {
		panic(err)
	}
	bounds := decode_image.Bounds()
  mem_image.Seek(0, io.SeekStart)
	return bounds.Max.X , bounds.Max.Y
}

func checkUniqueBanner(tmp_path string) (string , string){
  safe_name := shellescape.Quote(tmp_path)
  blockhash_cmd := exec.Command("blockhash", safe_name)
  hash_result, err := blockhash_cmd.Output()

  // var err_buffer , output_buffer bytes.Buffer
  // blockhash_cmd.Stdout = &output_buffer
  // blockhash_cmd.Stderr = &err_buffer
  // hash_result := output_buffer.String()
  // err_result := err_buffer.String()

  if string(hash_result) == "" || err != nil{
    return "Issue reading banner - 1" , ""
  }
  hash_result_arr := strings.Split(string(hash_result) , " ")
  if len(hash_result_arr) != 3{
    return "Issue reading banner - 2", ""
  }
  hash_result_str := hash_result_arr[0]
  fmt.Println("hash_result_str" , hash_result_str)

  match , err := bannerdb.Query(`
    SELECT * FROM ads WHERE hash = ?
  ` , []interface{}{hash_result_str})
  if err != nil{
    panic(err)
  }
  if len(match) != 0{
    return "This banner is not unique" , ""
  } else{
    return "", hash_result_str
  }
}

func generateHashedFileName(tmp_path string) string{
  safe_name := shellescape.Quote(tmp_path)
  sha_cmd := exec.Command("sha1sum", safe_name)
  sha1_and_file_result, err := sha_cmd.Output()
  if err != nil{
    panic(err)
  }
  sha1 := strings.Split(string(sha1_and_file_result) , " ")
  if len(sha1) != 3 {
    panic("sha1 hash issue" )
  }
  if sha1[0] == ""{
    panic("sha1 hash issue" )
  }
  return sha1[0]
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
    // panic(err)
  }
}
func removeAdSQL(name string, uri string){
  _ , err := bannerdb.Query(`
    DELETE FROM ads WHERE uri = ? AND fk_name = ?
  ` , []interface{}{uri , name})
fmt.Println(`
  DELETE FROM ads WHERE uri = ? AND fk_name = ?
` , []interface{}{uri , name})
  if err != nil{
    panic(err)
  }
}

func doCreationCooldown(ip string) int{
  time_now := time.Now().Unix()
  rows , err := bannerdb.Query(`
    SELECT * FROM antispam WHERE ip = ? AND type="banner" AND unix >= ?
    ORDER BY unix ASC
  ` , []interface{}{ ip , time_now} )
  if err != nil{
    panic (err)
  }
  timer := 0
  if len(rows) != 0{
    u , _ := strconv.Atoi(rows[0]["unix"])
    timer = u - int(time_now)
  }
  return timer / 60
}

func insertHaltCooldown(ip string){
  // freeze uploads for 5 sec to prevent duplicates
  cooldown := time.Now().Unix() + 5
  _ , err := bannerdb.Query(`
    INSERT INTO antispam VALUES (? , ? , ?)
  ` , []interface{}{ip , cooldown , "banner"} )
  if err != nil{
    panic (err)
  }
}

func updateCreationCooldown(ip string){
  time_now := time.Now().Unix()
  cooldown := time_now + controller_settings.BannerInterval * 60
  _ , err := bannerdb.Query(`
    DELETE FROM antispam WHERE unix < ? AND type="banner"
  ` , []interface{}{ time_now} )
  if err != nil{
    panic (err)
  }
  _ , err = bannerdb.Query(`
    INSERT INTO antispam VALUES (? , ? , ?)
  ` , []interface{}{ip , cooldown , "banner"} )
  if err != nil{
    panic (err)
  }
}


func getBase64(path string ) string{
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
  return  b64[0]
}

func addBanner(name string, file_path string, url string, ip string,  size string,
  hash string , board string){
  _ , err := bannerdb.Query(`
    INSERT INTO ads VALUES (? , ? , ? , ? , ? , 0 , NULL , ? , ? )
  ` , []interface{}{name , file_path , url , ip , size , hash , board })
  if err != nil{
    panic(err)
  }
}