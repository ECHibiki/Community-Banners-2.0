import React, { Component } from 'react';
export class FooterInfo extends Component{
	render(){
		return(
			<div id='footer'>
				<a href="https://github.com/ECHibiki/Community-Banners">Community Banners - {process.env.MIX_VERSION_NO}</a><br/>
				Verniy - MPL-2.0, {1900 + (new Date()).getYear()}<br/>
				Concerns should be sent to Verniy @ <a href="https://kissu.moe/b/res/2275">kissu.moe</a>
			</div>);
	}
}
