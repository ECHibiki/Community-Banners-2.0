<?php

namespace App;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Foundation\Auth\User as Authenticatable;
use Illuminate\Notifications\Notifiable;

use Tymon\JWTAuth\Contracts\JWTSubject;

class User extends Authenticatable implements JWTSubject
{
	protected $fillable=['name','pass'];
	protected $table = 'users';

	use Notifiable;

	public function getAuthPassword() {
       		return $this->pass;
    	}

	public function username() {
       		return "name";
    	}

     /**
     * Get the identifier that will be stored in the subject claim of the JWT.
     *
     * @return mixed
     */
    public function getJWTIdentifier()
    {
        return $this->getKey();
    }

    /**
     * Return a key value array, containing any custom claims to be added to the JWT.
     *
     * @return array
     */
    public function getJWTCustomClaims()
    {
	    return
	    [
		    'is_mod'=> $this->isMod()
	    ];
    }

    public function isMod(){
    	return Mod::where('fk_name','=', $this->name)->count() > 0;
    }

    public function isBanned(){
    	return Ban::where('fk_name','=', $this->name)->where('hardban','=', '1')->count() > 0;

    }

}
