<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class Ban extends Model
{
    protected $fillable=['fk_name', 'hardban'];
  	protected $table = 'bans';
}
