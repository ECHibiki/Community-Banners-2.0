package controllers

import(
  "os"

  "github.com/ECHibiki/Community-Banners-2.0/bannerdb"
)


func getUserData(name string) ([]map[string]string , bool){
  banners ,err := bannerdb.Query(`
      SELECT fk_name , uri , url , clicks, size FROM ads
      WHERE fk_name = ? ORDER BY id DESC
    ` , []interface{}{name})
  if err != nil{
    panic(err)
  }
  return banners;
}

// func doAntiSpam(){
//   // can this be tested?
//   public function doAntiSpam($name, $tmp_fname){
//     // expand into cooldown and optionally set phashing algorithm.
//     // return false or true
//     // for phash, new column will store hash data and evaluate for simularities
//     $antispam_response = [];
//     if(env('USE_PERCEPTUAL_HASHING') == "1"){
//       $check_arr = $this->checkDuplicateBanner($tmp_fname);
//       $antispam_response['duplicate'] = $check_arr["duplicate"];
//       $antispam_response['hash'] = $check_arr["hash"];
//     } else{
//       $antispam_response['duplicate'] = false;
//       $antispam_response['hash'] = "";
//     }
//     $antispam_response['cooldown'] = $this->checkSubmitCooldown($name);
//     return $antispam_response;
//   }
// }

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
func removeAdSQL(uri string){
  _ , err := bannerdb.Query(`
    DELETE FROM ads WHERE uri = ?
  ` , []interface{}{uri})
  if err != nil{
    panic(err)
  }
}