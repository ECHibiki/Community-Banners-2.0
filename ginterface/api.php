<?php

use Illuminate\Http\Request;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/


	// standard api routes
	Route::post("create", "UserCreationController@createNewUser");
	// Route::post("create", "UserCreationController@rejectUserCreation");
	Route::get("all", "PageGenerationController@getLimitedInfo");
	Route::get("banner", "PageGenerationController@GenerateAdJSON");

	// credential level
	Route::post("login", "UserSignInController@loginUser");
	//Route::post("login", "UserCreationController@rejectUserCreation");
	Route::get("details", "ConfidentialInfoController@accessInfo");
	Route::post("details", "ConfidentialInfoController@createInfo");
	Route::post("removal", "ConfidentialInfoController@removeInfo");

	// moderator api routes
	Route::post("mod/login", "UserSignInController@loginMod");
	Route::get("mod/all", "ModeratorActivityController@getAllInfo");
	Route::post("mod/ban", "ModeratorActivityController@banUser");
	Route::post("mod/purge", "ModeratorActivityController@deleteAll");
	Route::post("mod/individual", "ModeratorActivityController@deleteIndividual");
