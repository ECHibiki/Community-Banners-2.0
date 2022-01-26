<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use DB;
use App\Ban;
use App\Ad;
use App\Mod;
use App\Http\Controllers\ConfidentialInfoController;
use App\Http\Controllers\PageGenerationController;
use Storage;
class ModeratorActivityController extends Controller
{
	public function __construct(){
		$this->middleware(['auth:api']);
		$this->middleware(['mod:api']);
	}

	public static function getAllInfo(){
		$data = (array)ModeratorActivityController::GetAllEntries();
		$data = array_reverse(array_pop($data));
		return json_encode($data);
	}

	public function banUser(Request $request){
		$request->validate([
			"target"=>"required|string",
			"hard"=>"required|boolean"
		]);
		$hard = ($request->input("hard"));
		$target = ($request->input("target"));

		$this->createNewBan($request->input("target"), $request->input("hard"));
		if($hard == 1)
			$hard = "hard";
		else
			$hard = "soft";
		return response(json_encode(["log"=>"user $target was $hard banned"]), 200);
	}

	public function deleteAll(Request $request){
		$request->validate([
			"name"=>"required|string"
		]);
		$target = $request->input("name");

		//images must be first
		ModeratorActivityController::removeAllUserImages($target);
		ModeratorActivityController::removeUserFromDatabase($target);
		ModeratorActivityController::truncateUserJSON($target);

		return response(json_encode(["log"=>"user $target was purged"]), 200);

	}

	public function deleteIndividual(Request $request){
		$request->validate([
			"name"=>"required|string",
			"uri"=>"required|string",
			"url"=>"required|string"
		]);
		$name = $request->input("name");
		$uri = str_replace("storage/image/", "public/image/", $request->input("uri"));
		ModeratorActivityController::removeIndividualBannerFromImages($uri);
		ModeratorActivityController::removeIndividualBannerFromJSON($name, $uri, $request->input("url"));
		ModeratorActivityController::removeIndividualBannerFromDB($uri);
		return response(json_encode(["log"=>"$name's image was pruned"]), 200);
	}

	public static function GetAllEntries(){
		return DB::table('ads')->leftJoin('bans', 'ads.fk_name', '=', 'bans.fk_name')->select('ads.fk_name', 'url', 'uri', 'bans.hardban', 'size', 'clicks')->orderBy('ads.created_at', 'ASC')->get();
	}

	public function createNewBan($target, $hard){
		if($ban = ModeratorActivityController::GetBanInfo($target)){}
		else{
			$ban = new Ban();
		}
		$ban->fk_name = $target;
		$ban->hardban = $hard;
		$ban->save();
	}
	public static function removeIndividualBannerFromJSON($name, $uri, $url){
		ModeratorActivityController::removeUserJSON($name,$uri,$url);
	}

	public static function removeIndividualBannerFromDB($uri){
		Ad::where('uri', '=', $uri)->delete();
	}

	public static function removeIndividualBannerFromImages($uri){
		Storage::delete($uri);
	}

	public static function GetBanInfo($name){
		return Ban::query()->where('fk_name', '=', $name)->first();
	}

	public static function createMod($name){
		$mod = new Mod(['fk_name'=>$name]);
		$mod->save();
	}

	public static function removeUserFromDatabase($name){
		Ad::where('fk_name','=', $name)->delete();
	}

	public static function truncateUserJSON($target){
		Storage::disk('local')->put("$target.json","[]");
	}

	public static function removeAllUserImages($target){
		$ads = json_decode(ModeratorActivityController::getSelectJSON($target), true);
		foreach($ads as $ad){
			Storage::delete($ad['uri']);
		}
	}

	public static function getSelectJSON($name){
		return Storage::disk('local')->get("$name.json");
	}

	public static function removeUserJSON(string $name, string $uri, string $url){
		$combined = json_decode(Storage::disk('local')->get("$name.json"), true);
		$reduced = [];
		foreach($combined as $entry){
			if($entry['uri'] == $uri && $entry['url'] == $url){
				continue;
			}
			else{
				$reduced[] = $entry;
			}
		}
		Storage::disk('local')->put("$name.json", json_encode($reduced));
	}
}
