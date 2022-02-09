import React, { Component } from 'react';
import {TopHeader} from "../information/header";
import {DonatorBox} from "../information/donator";
import {HelperText} from "../information/helper";
import {SampleBanner} from "../image/sample-banner";
import {PatreonBanner} from "../image/patreon-banner";
import {ModContainer} from "../container/mod";

export class ModPage extends Component{
	render(){
			return(<div id="master-mod">
				<div id="upper-master-mod">
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
						<PatreonBanner />
						<HelperText />
				   </div>
				</div>);
	}
}
