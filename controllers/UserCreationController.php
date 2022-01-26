<?php

namespace App\Http\Controllers;
use Carbon\Carbon;
use Illuminate\Http\Request;
use App\User;
use App\AntiSpam;
use Symfony\Component\HttpFoundation\Response;
use Illuminate\Support\Facades\Validator;
use Illuminate\Support\Facades\Hash;
use DB;
class UserCreationController extends Controller
{
	public function rejectUserCreation(){
			return response(json_encode(["warn"=>"Pool Closed - Come back later"]), 200)->header('Content-Type', 'text/plain');
	}

	public function createNewUser(Request $request){
		$request->validate([
			'name' => 'required|max:30',
			'pass' => 'required|confirmed|min:5'
		]);

		$cr_resp = $this->validateIPCreation();
		if($cr_resp->count() > 0){
			return response()->json(['warn'=>'Too many accounts'], 401);
		}
		$this->updateIPCreation();

		$response_msg = UserCreationController::addNewUserToDB($request->input('name'), $request->input('pass'));
		if($response_msg === 1){
			if($err = PageGenerationController::CreateUserFile($request->input('name')))
				return $err;
			return response(json_encode(["log"=>"Successfully Created"]), 200)->header('Content-Type', 'text/plain');
		}
		else{
			return $response_msg;
		}

	}

	public static function addNewUserToDB(string $name, string $hashedpass){
		if(preg_match('/(\\.|\\/|;)/', $name)){
			return response(json_encode(["warn"=>"Name has invalid characters"]), 422)->header('Content-Type', 'text/plain');
		}
		if(!User::where('name', $name)->exists()){
			$user = new User (['name' => $name, 'pass' => bcrypt($hashedpass)]);
			$user->save();
			return 1;
		}
		else
			return response(json_encode(["warn"=>"Username Already Exists"]), 401)->header('Content-Type', 'text/plain');

	}
	// needs test case
	public function validateIPCreation(){
		return DB::table('antispam')
			->where('name','=',ConfidentialInfoController::getBestIPSource())
			->where('type','=','create')
			->where('unix', '>=',
				Carbon::now()->subSeconds(intval(env('IP_CREATE_COOLDOWN',60)))->timestamp);
	}
	// needs test case
	public function updateIPCreation(){
		DB::table('antispam')
			->where('unix', '<',
				Carbon::now()->subSeconds(intval(env('IP_CREATE_COOLDOWN',60)))->timestamp)
			->where('type', '=', 'create')
			->delete();
		AntiSpam::create(['name'=>ConfidentialInfoController::getBestIPSource(), 'unix' => 	Carbon::now()->timestamp, 'type'=>'create']);
	}
}
