<?php

namespace App\Http\Controllers;
use Carbon\Carbon;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\DB;
use Illuminate\Http\UploadedFile;

use App\Ad;
use App\AntiSpam;
use JWTAuth;
use App\Http\Controllers\PageGenerationController;
use App\Http\Controllers\MailSendController;

class ConfidentialInfoController extends Controller
{

	public function __construct(){
		$this->middleware(['auth:api']);
		$this->middleware(['ban:api']);
	}



// can this be tested?
	public function checkDuplicateBanner($tmp_fname){
		$hash = shell_exec("blockhash " . escapeshellarg($tmp_fname));
		if($hash){
			$hash = explode(" ", $hash)[0];
		} else{
			return ["duplicate" => true, "hash" => ""];
		}
		return ["duplicate" => DB::table("ads")->where("hash", "=", $hash)->count() > 0, "hash" => $hash];
	}
	// can this be tested?
	public function doAntiSpam($name, $tmp_fname){
		// expand into cooldown and optionally set phashing algorithm.
		// return false or true
		// for phash, new column will store hash data and evaluate for simularities
		$antispam_response = [];
		if(env('USE_PERCEPTUAL_HASHING') == "1"){
			$check_arr = $this->checkDuplicateBanner($tmp_fname);
			$antispam_response['duplicate'] = $check_arr["duplicate"];
			$antispam_response['hash'] = $check_arr["hash"];
		} else{
			$antispam_response['duplicate'] = false;
			$antispam_response['hash'] = "";
		}
		$antispam_response['cooldown'] = $this->checkSubmitCooldown($name);
		return $antispam_response;
	}

// test case
	public function checkSubmitCooldown($name){
		return DB::table('antispam')
			->where('name','=',$name)
			->where('type','=','ad')
			->where('unix', '>=',
				Carbon::now()->subSeconds(intval(env('AD_CREATE_COOLDOWN',60)))->timestamp);
	}

	public function updateAntiSpam($name){
		DB::table('antispam')
			->where('unix', '<',
				Carbon::now()->subSeconds(intval(env('AD_CREATE_COOLDOWN',60)))->timestamp)
			->where('type', '=', 'ad')
			->delete();
		AntiSpam::create(['name'=>$name, 'unix' => 	Carbon::now()->timestamp, 'type'=>'ad']);
	}

	public function CreateBanner(Request $request){
		$response ="";
		$name = auth()->user()->name;
		$antispam_response = $this->doAntiSpam($name, $request->file('image')->getPathName());
	  if ($antispam_response['cooldown']->count() > 0){
			return ['warn'=>'posting too fast('.
				($antispam_response['cooldown']->first()->unix - Carbon::now()->subSeconds(intval(env('AD_CREATE_COOLDOWN',60)))->timestamp) . ' seconds)'];
		} else if ($antispam_response['duplicate']) {
			return ['warn'=> 'Duplicate detected'];
		} else{
			if($request->input('size') == "small"){
				$response = $this->createSmallInfo($request, $antispam_response['hash']);
			}
			else{
				$response = $this->createWideInfo($request, $antispam_response['hash']);
			}
		}
		$this->updateAntiSpam($name);
		return $response;
	}

	public function createSmallInfo(Request $request, $hash){
		$request->validate([
			'image'=>'required|image|dimensions:width='. env('MIX_IMAGE_DIMENSIONS_SMALL_W', '300') .',height=' . env('MIX_IMAGE_DIMENSIONS_SMALL_H', '140') .
			 '|between:0, '.env('MIX_MAX_FILE_SIZE', '2000'),
		]);
		$fname = PageGenerationController::StoreAdImage($request->file('image'));
		$this->addUserJSON($fname, env('MIX_APP_URL', 'https://kissu.moe'), 'small');
		$this->addAdSQL($fname, $hash, env('MIX_APP_URL', 'https://kissu.moe'), 'small');

		$t = MailSendController::getCooldown();
		if($t < time()){
			$err = MailSendController::sendMail(["name"=>auth()->user()->name, "time"=>date('yMd-h:i:s',time()), "url"=> $request->input('url'), 'fname'=>$fname],
				['primary_email'=>env('PRIMARY_MOD_EMAIL'), 'secondary_emails'=>env('SECONDARY_MOD_EMAIL_LIST')]);
			MailSendController::updateCooldown();
			if(!$err){
					return ['log'=>'Ad Created', 'fname'=>$fname, 'errors'=>'no email'];
			}
			if (gettype($err) != 'boolean')
				return ['log'=>'Ad Created', 'fname'=>$fname, 'errors'=>$err];

		}
		return ['log'=>'Ad Created', 'fname'=>$fname];
	}

	public function createWideInfo(Request $request, $hash){
		$request->validate([
			'image'=>'required|image|dimensions:width='. env('MIX_IMAGE_DIMENSIONS_W', '500') .',height=' . env('MIX_IMAGE_DIMENSIONS_H', '90')
				. '|between:0, ' . env('MIX_MAX_FILE_SIZE', '2000'),
			'url'=>['required','url','regex:/^http(|s):\/\/[A-Z0-9+&@#\/%?=~\-_|!:,.;]+\.[A-Z0-9:+&@#\/%=~_|?\.\-"\']+$/i']
		]);
		$fname = PageGenerationController::StoreAdImage($request->file('image'));
		$this->addUserJSON($fname, $request->input('url'), 'wide');
		$this->addAdSQL($fname, $hash, $request->input('url'), 'wide');
		$t = MailSendController::getCooldown();

		if($t < time()){
			$err = MailSendController::sendMail(["name"=>auth()->user()->name, "time"=>date('yMd-h:i:s',time()), "url"=> $request->input('url'), 'fname'=>$fname],
				['primary_email'=>env('PRIMARY_MOD_EMAIL'), 'secondary_emails'=>env('SECONDARY_MOD_EMAIL_LIST')]);
			MailSendController::updateCooldown();
			if(!$err){
			    return ['log'=>'Ad Created', 'fname'=>$fname, 'errors'=>'no email'];
			}
			if (gettype($err) != 'boolean')
				return ['log'=>'Ad Created', 'fname'=>$fname, 'errors'=>$err];

		}
		return ['log'=>'Ad Created', 'fname'=>$fname];
	}

	public static function addUserJSON(string $uri, string $url, string $size){
		$name = auth()->user()->name;
		$combined = json_decode(Storage::disk('local')->get("$name.json"), true);
		$combined[] = ['uri'=>$uri, 'url'=>$url, 'size'=> $size, 'clicks'=>'0'];
		Storage::disk('local')->put("$name.json", json_encode($combined));
	}

	public static function removeUserJSON(string $uri, string $url){
		$name = auth()->user()->name;
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

	public static function addAdSQL(string $uri, string $hash, string $url, string $size='wide'){
		$name = auth()->user()->name;
		$ad = new Ad(['fk_name'=>$name, 'hash'=>$hash, 'uri'=>$uri, 'url'=>$url, 'ip'=>ConfidentialInfoController::getBestIPSource(), 'size'=>$size]);
		$ad->save();
	}

}
