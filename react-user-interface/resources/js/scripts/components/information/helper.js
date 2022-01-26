import React, { Component } from 'react';
import {dimensions_w, dimensions_h, dimensions_small_w,dimensions_small_h} from '../../settings'

export class HelperText extends Component{
	render(){
		var code_string_iframe = '<iframe id="kissu-banner" src="/banner" scrolling="no" width="' + dimensions_w + '" height="' + dimensions_h + '" style="margin:auto;display:block;max-width:100%;border:none"></iframe>';
		var code_string_ajax = 	'<script>var fetch = new XMLHttpRequest();fetch.open("GET", "{{ config.banner_src }}api/banner");fetch.addEventListener("load", function(){var info = JSON.parse(this.responseText)[0];window.ban_url = document.createElement("A");window.ban_url.setAttribute("href", info["url"]);ban_url.setAttribute("style", "display:contents");var ban_img = document.createElement("IMG");ban_img.setAttribute("src", "{{ config.banner_src }}" + info["uri"]);ban_img.setAttribute("style", "margin:auto;display:block;max-width:100%;border:none;");	window.ban_url.appendChild(ban_img);if(document.getElementById("banner-container") != undefined);document.getElementById("banner-container").appendChild(window.ban_url);});fetch.send();</script>;'
		+ '\n\n<body><div id="banner-container"></div></body>';
		return (<div id="helper"><h2>How To Use</h2><p>To easily embed on a website use:<br/>
							<textarea className="code" value={code_string_iframe} readOnly/>
							A more sophisticated method is to preload the banner and then place into a container:<br/>
							<textarea className="code" value={code_string_ajax} readOnly/>
							This is a slight bit faster than iframes and removes difficulties with resizing.
							<br/>Uploaded images must be {dimensions_w}x{dimensions_h} or {dimensions_small_w}x{dimensions_small_h} and safe for work
						</p>
					</div>);
	};
}
