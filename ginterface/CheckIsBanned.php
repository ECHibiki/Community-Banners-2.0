<?php

namespace App\Http\Middleware;

use Closure;
use App\User;
use JWTAuth;
class CheckIsBanned
{
    /**
     * Handle an incoming request.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  \Closure  $next
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
	    if(User::query()->where("name", "=", auth()->user()->name)->first()->isBanned()){
		    auth()->invalidate();
		    return response(json_encode(['warn'=>'You\'ve been banned...']),401);
	    }
	    return $next($request);
    }
}
