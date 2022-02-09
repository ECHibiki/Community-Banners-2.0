import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';
import {TopHeader} from "../information/header";
import {DonatorBox} from "../information/donator";
import {HelperText} from "../information/helper";
import {SampleBanner} from "../image/sample-banner";
import {PatreonBanner} from "../image/patreon-banner";
import {LoadingSpinner} from "../image/loading-spinner";
import {LoginContainer} from "../container/login";
import {UserContainer} from "../container/user";
export class MasterPage extends Component{
	constructor(props){
		super(props);
		this.state = {auth: undefined , is_donor:false};
		this.swapPage = this.swapPage.bind(this);
	}

	componentDidMount(){
		this.checkLoggedIn();
	}

	async checkLoggedIn(){
		var instant_login = await APICalls.callRetrieveUserAds();
		if(!("name" in instant_login)){
			this.setState({auth: false});
		}
		else{
			let is_donor = false;
			if ("donor" in instant_login){
				is_donor = instant_login["donor"]
			}
			this.setState({auth: true , is_donor: is_donor});
		}
	}
	swapPage(donor){
		this.setState({auth: !this.state.auth , is_donor:donor});
	}

	render(){
		if(this.state.auth == undefined){
			return(<div id="master-waiting"  className="main-container">
							<div id="upper-master-login" className="upper-container">
							 <TopHeader />
							 <SampleBanner />
							</div>
							<div id="lower-master-waiting">
							 <hr/>
							 <LoadingSpinner />
							</div>
			       </div>);
		}
		else if(!this.state.auth){
			// possible context API on swapPage
			return(<div id="master-login" className="main-container">
							<div id="upper-master-login" className="upper-container">
							  <TopHeader />
							  <SampleBanner />
							</div>
							<hr/>
							  <div id="mid-master-login">
							    <LoginContainer swapPage={this.swapPage}/>
							   </div>
							  <hr/>
							   <div id="lower-master-login">
								 	<DonatorBox />
 							  	<PatreonBanner />
									<HelperText />
							   </div>
							</div>
					);
		}
		else{
			return(<div id="master-user" className="main-container">
							<div id="upper-master-user" className="upper-container">
							<TopHeader />
							<SampleBanner />
							</div>
							<hr/>
							<div id="mid-master-user">
							  <UserContainer isDonor={this.state.is_donor}/>
							</div>
							<div id="lower-master-user">
								{!this.state.is_donor && <><DonatorBox />
							  <PatreonBanner /></> }
							  <HelperText />
							</div>
						</div>
					);
		}
	}
}
