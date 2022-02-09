import React, { Component } from 'react';
import {TopHeader} from "../information/header";
import {DonatorBox} from "../information/donator";
import {HelperText} from "../information/helper";
import {SampleBanner} from "../image/sample-banner";
import {PatreonBanner} from "../image/patreon-banner";
import {AllContainer} from "../container/all";
export class AllPage extends Component{
	render(){
			return(
				<div id="master-all" className="main-container">
					<div id="upper-master-all" className="upper-container">
					  <TopHeader />
					  <SampleBanner />
					</div>
					<hr/>
				  <div id="mid-master-all">
				    <AllContainer />
			    </div>
				  <hr/>
			    <div id="lower-master-all">
						<DonatorBox />
						<PatreonBanner />
						<HelperText />
			   </div>
			 </div>);
	}
}
