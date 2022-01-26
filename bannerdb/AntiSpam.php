<?php

namespace App;

use Illuminate\Database\Eloquent\Model;

class AntiSpam extends Model
{
    protected $fillable=['name', 'unix', 'type'];
    protected $table = 'antispam';
}
