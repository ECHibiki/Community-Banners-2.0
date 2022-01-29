import React, { Component } from 'react';
import {extra_info} from '../../settings'

export class DonatorBox extends Component{
	render(){
		//safe because build variable
		var html = {__html: extra_info}
		return(<div id="donation" dangerouslySetInnerHTML={html}></div>);
	}
}
