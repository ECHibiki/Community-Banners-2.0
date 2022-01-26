import React, { Component } from 'react';
export class DonatorBox extends Component{
	render(){
		//safe because build variable
		var html = {__html: process.env.MIX_EXTRA_INFO}
		return(<div id="donation" dangerouslySetInnerHTML={html}></div>);
	}
}
