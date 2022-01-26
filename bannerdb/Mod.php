<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class Mod extends Model {
	protected $fillable=['fk_name'];
		protected $table = 'mods';
}
