<?php

namespace App\Http\Middleware;

use Closure;
use App\User;
use JWTAuth;
class CheckIsMod
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
	    if(!User::query()->where("name", "=", auth()->user()->name)->first()->isMod()){
		    return response(json_encode(['warn'=>'You are not a moderator']), 401);
	    }
	    return $next($request);
    }
}
