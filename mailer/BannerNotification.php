<?php

namespace App\Mail;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Mail\Mailable;
use Illuminate\Queue\SerializesModels;

class BannerNotification extends Mailable // implements ShouldQueue
{
    use Queueable, SerializesModels;

	protected $name, $time, $url, $fname, $error;

    /**
     * Create a new message instance.
     *
     * @return void
     */
    public function __construct($data_arr)
    {
	    $this->name = $data_arr['name'];
	    $this->time = $data_arr['time'];
	    $this->url = $data_arr['url'];
	    $this->fname = $data_arr['fname'];
	    $this->error = $data_arr['err'];
    }

    /**
     * Build the message.
     *
     * @return $this
     */
    public function build()
    {
      $this->fname =  public_path() . "/" .  str_replace("public", "storage", $this->fname);
	    return $this->from('banner@mail.notice')
		    	->view('banner-mail-notice')->with(["name"=>$this->name, "time"=>$this->time, "url"=>$this->url, "err"=>$this->error, "fname"=>$this->fname]);
    }
}
