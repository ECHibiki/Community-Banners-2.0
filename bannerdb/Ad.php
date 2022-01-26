<?php

namespace App;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Foundation\Auth\User as Authenticatable;
use Illuminate\Notifications\Notifiable;
class Ad extends Model
{
    protected $fillable=['fk_name', 'uri', 'url', 'ip', 'size', 'clicks', 'hash'];
    protected $dateFormat = 'Y-m-d H:i:s';
    protected $table = 'ads';
}
