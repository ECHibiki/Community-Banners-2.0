<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;

use Illuminate\Http\UploadedFile;

use Illuminate\Support\Facades\DB;

class PageGenerationController extends Controller
{
	public static function CreateUserFile(string $name){
		if(preg_match('/(\\.|\\/|;)/', $name))
			return response(json_encode(["warn"=>"Name has invalid characters"]), 422)->header('Content-Type', 'text/plain');
		Storage::disk('local')->put("$name.json", '[]');
	}

	public static function StoreAdImage($img){
		$fname = Storage::putFile('public/image', $img);
		return $fname;
	}

	public static function getLimitedInfo(Request $request){
		$pool = false;
		try{
			$pool =$request->input('env');
		}
		catch(\Exception $e){}
		$data = (array)PageGenerationController::GetLimitedEntries(null, $pool);
		$data = array_reverse(array_pop($data));
		return json_encode($data);
	}

	public function GenerateAdPage(Request $request){
		$size = $request->input('size');
	  if($size == "small"){
			$rand_ad = $this->GetSmallRandomAdEntry();
		}
		else if($size == "wide"){
			$rand_ad = $this->GetWideRandomAdEntry();
		}
		else{
					$rand_ad = $this->GetRandomAdEntry();
		}
		if($rand_ad == null){
			return "asdf no ads";
		}
		else{
      $url = env('MIX_APP_URL') . '/req?s=' . urlencode($rand_ad->url) . '&f=' . substr($rand_ad->uri, strrpos($rand_ad->uri, '/') + 1);
			return view('banner', ['url'=>$url , 'uri'=>str_replace('public','storage',$rand_ad->uri), 'name'=>$rand_ad->fk_name, 'size'=>$rand_ad->size, 'clicks'=>$rand_ad->clicks]);
		}
	}

	public function GenerateAdJSON(Request $request){
		$size = $request->input('size');
		if($size == "small"){
			$rand_ad = $this->GetSmallRandomAdEntry();
		}
		else if($size == "wide"){
			$rand_ad = $this->GetWideRandomAdEntry();
		}
		else{
			$rand_ad = $this->GetRandomAdEntry();
		}
		if($rand_ad == null){
			return json_encode([['url'=>'', 'uri'=>'', 'name'=>'asdf no ads', 'size'=>'', 'clicks'=>'']]);
		}
		else{
      $url = env('MIX_APP_URL') . '/req?s=' . urlencode($rand_ad->url) . '&f=' . substr($rand_ad->uri, strrpos($rand_ad->uri, '/') + 1);
			return json_encode([['url'=>$url , 'uri'=>str_replace('public','storage',$rand_ad->uri),  'name'=>$rand_ad->fk_name, 'size'=>$rand_ad->size, 'clicks'=>$rand_ad->clicks]]);
		}
	}

	// banned users will not show up in rotation
	public static function GetRandomAdEntry(){
		try{
      $name = auth()->setToken($_COOKIE['freeadstoken'])->user()->name;
    } catch(\Exception $e){
			$name = "";
		}
    return DB::table('ads')
			->leftJoin('bans', 'ads.fk_name', '=', 'bans.fk_name')
			->whereNull('bans.hardban')
			->orWhere('ip','=', PageGenerationController::getBestIPSource())
			->orWhere('bans.fk_name','=',$name)
			->select("ads.fk_name", "uri", "url", "size", "clicks")
			->inRandomOrder()->first();

		}
		// banned users will not show up in rotation
		public static function GetSmallRandomAdEntry(){
			try{
													$name = auth()->setToken($_COOKIE['freeadstoken'])->user()->name;
									}
									catch(\Exception $e){
													$name = "";
			}
			return DB::table('ads')
				->leftJoin('bans', 'ads.fk_name', '=', 'bans.fk_name')
				->where('ads.size', '=', 'small')
				->Where(function($query) use ($name){
					$query->whereNull('bans.hardban')
					->orWhere('ip','=', PageGenerationController::getBestIPSource())
					->orWhere('bans.fk_name','=',$name);
				})
				->select("ads.fk_name", "uri", "url", "size", "clicks")
				->inRandomOrder()->first();

			}
			// banned users will not show up in rotation
			public static function GetWideRandomAdEntry(){
				try{
					$name = auth()->setToken($_COOKIE['freeadstoken'])->user()->name;
				}
				catch(\Exception $e){
					$name = "";
				}

			 return DB::table('ads')
					->leftJoin('bans', 'ads.fk_name', '=', 'bans.fk_name')
					->where('ads.size', '=', 'wide')
					->Where(function($query) use ($name){
						$query->whereNull('bans.hardban')
						->orWhere('ip','=', PageGenerationController::getBestIPSource())
						->orWhere('bans.fk_name','=',$name);
					})
					->select("ads.fk_name", "uri", "url", "size", "clicks")
					->inRandomOrder()->first();

				}

	public static function GetLimitedEntries($name = null, $env_skip = false){
		try{
			if(!$name)
				$name = auth()->setToken($_COOKIE['freeadstoken'])->user()->name;
			}
		catch(\Exception $e){
			$name = "";
		}
		return DB::table('ads')
			->leftJoin('bans', 'ads.fk_name', '=', 'bans.fk_name')
			->when(!(env('ALLOW_BANNED_USER_POOL') || $env_skip)  || ($name == "" || !PageGenerationController::checkBanned($name)),function($q) use ($name){
				return $q->whereNull('bans.hardban')
					->orWhere('ip','=', PageGenerationController::getBestIPSource())
					->orWhere('bans.fk_name','=',$name);
			})
			->select("ads.fk_name", "uri", "url", "size", "clicks")
			->orderBy('ads.created_at', 'ASC')->get();
                }

	public static function checkBanned($name){
		return  DB::table('bans')->where('fk_name', '=', $name)->count() > 0;
	}

	public static function getBestIPSource(){
		return isset($_SERVER['HTTP_X_REAL_IP']) ? $_SERVER['HTTP_X_REAL_IP'] : \Request::ip();
	}

}
