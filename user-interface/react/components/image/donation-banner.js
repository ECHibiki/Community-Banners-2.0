import React, { Component } from 'react';
import {PatreonBanner} from "./patreon-banner";
import {KofiBanner} from "./ko-fi-banner";
export class DonationBanner extends Component{
	render(){
		return (<KofiBanner />);
	}
}
