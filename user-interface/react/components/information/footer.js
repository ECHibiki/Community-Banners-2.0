import React, { Component } from 'react';
import {version_no} from '../../settings'
export class FooterInfo extends Component{
	render(){
		return(
			<div id='footer'>
				<a href="https://github.com/ECHibiki/Community-Banners-2.0">Community Banners - {version_no}</a><br/>
				Verniy - MPL-2.0, {1900 + (new Date()).getYear()}<br/>
				Concerns should be sent to Verniy @ <a href="https://kissu.moe/b/">kissu.moe</a>
			</div>);
	}
}
