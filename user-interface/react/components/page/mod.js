import React, { Component } from 'react';
import {TopHeader} from "../information/header";
import {DonatorBox} from "../information/donator";
import {HelperText} from "../information/helper";
import {SampleBanner} from "../image/sample-banner";
import {DonationBanner} from "../image/donation-banner";
import {ModContainer} from "../container/mod";
import {ThemeToggle } from "../page/themes";

export class ModPage extends Component{
	render(){
			return(<div id="master-mod">
				<div id="upper-master-mod">
					<ThemeToggle />
				  <TopHeader />
				  <SampleBanner />
				</div>
				<hr/>
				  <div id="mid-master-mod">
				    <ModContainer />
				   </div>
				  <hr/>
				   <div id="lower-master-mod">
					 	<DonatorBox />
						<DonationBanner />
						<HelperText />
				   </div>
				</div>);
	}
}
