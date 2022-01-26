<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

use Storage;
use App\Mail\BannerNotification;
use Illuminate\Support\Facades\Mail;
use GuzzleHttp\Client;

class MailSendController extends Controller
{

	public static function updateCooldown(){
		$json = [];
		if(Storage::disk('local')->exists('mail.json')){
			$json = json_decode(Storage::disk('local')->get('mail/mail.json'), true);
		}
		$json['cooldown'] = time() + env('MAIL_COOLDOWN', 5);
		Storage::disk('local')->put('mail/mail.json', json_encode($json));
	}

	public static function getCooldown(){
		if(Storage::disk('local')->exists('mail/mail.json')){
			$json = json_decode(Storage::disk('local')->get('mail/mail.json'), true);
			return $json['cooldown'];
		}
		else return 0;
	}


	public static function sendMail($data, $emails){
	$bcc_errors = '';
	 try{
	    if(isset($emails['primary_email'])){
		    $msg = Mail::to(str_replace("\n", "", $emails['primary_email']));
		   if(isset($emails['secondary_emails'])){
			$bcc = explode(',', $emails['secondary_emails']);
			foreach($bcc as $mail){
				$mail = str_replace("\n", "", $mail);

				if(preg_match('/[^@]+@[^@\.]+\.[^@]+/', $mail)){
					$msg->bcc($mail);
				}
				else if($mail == ""){}
				else{
					$errors .= '|E: ' . $mail;
				}
			}
		}

		$msg->send(new BannerNotification(['name'=>$data['name'], 'time'=>$data['time'], 'url'=>$data['url'], 'fname'=>$data['fname'], 'err'=>$bcc_errors]));
		return true;
	    }
	}catch(\Exception $e){
		return $e->getMessage();
	}
	return false;
   }
}
