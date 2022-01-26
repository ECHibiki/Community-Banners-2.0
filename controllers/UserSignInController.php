<?php

namespace App\Http\Controllers;
use Carbon\Carbon;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use App\User;
use App\AntiSpam;
use App\Mod;
use App\Ban;
use DB;

use JWTAuth;
use Tymon\JWTAuth\Exceptions\TokenExpiredException;
use Tymon\JWTAuth\Exceptions\TokenInvalidException;
use Tymon\JWTAuth\Exceptions\JWTException;

class UserSignInController extends Controller
{

	public function username()
	{
    		return 'name';
	}

	public function loginUser(Request $request){
		$request->validate([
			'name'=>'required',
			'pass'=>'required'
		]);

		$bf_resp = $this->validateNameBruteForce($request->input("name"));
		if($bf_resp->count() > intval(env('MAX_PASS_ATTEMPTS_PER_CYCLE'))){
			return response()->json(['warn'=>'Too many password attempts'], 401);
		}
		$this->updateNameBruteForce($request->input("name"));

		$token = $this->returnJWT($request->input("name"), $request->input("pass"), false);
		if(!$token){
			return response()->json(['warn'=>'Username or Password Incorrect'], 401);
		}
		else if($this->checkIfBanned($request->input("name"))){
			return response()->json(['warn'=>'You\'ve been banned...'], 401);
		}
		else{
			$token_arr = [
     		'access_token' => $token,
  			'token_type' => 'bearer',
				'expires_in' => auth()->factory()->getTTL() * 60,
				'log' => "Successfully Logged In"
			];
			return response()->json($token_arr);
		}
	}

	//probably will be unused, but worthwhile fallback
	public function loginMod(Request $request){
		$request->validate([
			'name'=>'required',
			'pass'=>'required'
		]);

		if(!User::query()->where("name", "=", $request->input("name"))->first()->isMod()){
			return response()->json(['warn'=>'You are not a moderator'], 401);
		}

		$token = $this->returnJWT($request->input("name"), $request->input("pass"), true);
		if(!$token){
			return response()->json(['warn'=>'Username or Password Incorrect'], 401);
		}
		else{
			$token_arr = [
		       		'access_token' => $token,
            			'token_type' => 'bearer',
				'expires_in' => auth()->factory()->getTTL() * 60,
				'log' => "Successfully Logged In"
			];
			return response()->json($token_arr);
		}
	}

	public static function returnJWT($name, $pass){
		$token = auth()->attempt(["name"=>$name, "password"=>$pass]);
		return $token;
	}

	public function checkIfBanned($name){
		return Ban::query()
			->where("hardban", "=", "1")
			->where("fk_name", "=", $name)->count() > 0;
	}

	// needs test case
	public function validateNameBruteForce($name){
		return DB::table('antispam')
			->where('name','=',$name)
			->where('type','=','login')
			->where('unix', '>=',
				Carbon::now()->subSeconds(intval(env('NAME_LOGIN_COOLDOWN',60)))->timestamp);
	}
	// needs test case
	public function updateNameBruteForce($name){
		DB::table('antispam')
			->where('unix', '<',
				Carbon::now()->subSeconds(intval(env('NAME_LOGIN_COOLDOWN',60)))->timestamp)
			->where('type', '=', 'login')
			->delete();
		AntiSpam::create(['name'=>$name, 'unix' => 	Carbon::now()->timestamp, 'type'=>'login']);
	}
}
