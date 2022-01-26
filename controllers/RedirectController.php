<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use DB;
use Storage;
class RedirectController extends Controller
{
    public function RedirectSiteRequest(Request $request){
        $site_query = "s";
        $site = $request->input($site_query);
        $filename_query = "f";
        $pkey = 'public/image/' . $request->input($filename_query);
        if($site == NULL){
          return response("Non existing URL", 404);
        }
        $existing = DB::table('ads')->where('uri','=',$pkey)->select('clicks', 'size', 'fk_name')->first();
        if($existing == NULL){
          return response("Non existing URI", 404);
        }
        if($existing->size == 'wide'){
          RedirectController::incrementClicksSQL($pkey, $existing->clicks);
          RedirectController::incrementClicksJSON($pkey, $existing->fk_name, $site, $existing->clicks);
        }
        return redirect($site);
    }

    public static function incrementClicksSQL(string $uri, $clicks){
        DB::table('ads')->where('uri','=',$uri)->update(['clicks' =>  $clicks + 1]);
  	}

    public static function incrementClicksJSON(string $uri, string $name, string $url, $clicks){
      $combined = json_decode(Storage::disk('local')->get("$name.json"), true);
      $new = [];
      foreach($combined as $entry){
        if($entry['uri'] == $uri && $entry['url'] == $url){
          $entry['clicks'] = $clicks + 1;
        }
          $new[] = $entry;
      }
      Storage::disk('local')->put("$name.json", json_encode($new));
    }
}
