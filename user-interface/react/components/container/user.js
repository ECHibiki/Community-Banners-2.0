import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';
import {AdCreationForm, AdCreateButton } from '../form/ad-create';
import {DonorPanel , DonorExtendButton} from '../form/donor';
import {AdDetailsTable} from '../table/user-details-table';

import {Link} from "react-router-dom";

export class UserContainer extends Component{
	constructor(props){
		super(props);
		this.AdCreateOnClick = this.AdCreateOnClick.bind(this);
		this.DonorAuthOnClick = this.DonorAuthOnClick.bind(this);
		this.state = {
			AdCVisibility:"unset", AdCHeight:"0em", AdCOpacity:"0",
			DonorCVisibility:"unset", DonorCHeight:"0em", DonorCOpacity:"0",
		AdArray:[], mod:false , donor: this.props.isDonor
	};
		this.UpdateDetails = this.UpdateDetails.bind(this);
		this.UpdateDonator = this.UpdateDonator.bind(this);
	}

	componentDidMount(){
		this.UpdateDetails();
	}

	AdCreateOnClick(){
		if(this.state.AdCVisibility == "unset")
			this.setState({AdCVisibility:"initial", AdCHeight:"max-content", AdCOpacity:"1"});
		else
			this.setState({AdCVisibility:"unset", AdCHeight:"0em", AdCOpacity:"0"});
	}
	DonorAuthOnClick(){
		if(this.state.DonorCVisibility == "unset")
			this.setState({DonorCVisibility:"initial", DonorCHeight:"max-content", DonorCOpacity:"1"});
		else
			this.setState({DonorCVisibility:"unset", DonorCHeight:"0", DonorCOpacity:"0"});
	}

	async UpdateDetails(){
		var d_response = await APICalls.callRetrieveUserAds();
		if("error" in d_response){
			var key_ind = 0;
			this.setState({err_text:d_response['error'], war_text: "", suc_text:""});
		}
		else if("warn" in d_response){
			this.setState({err_text: "" , war_text:d_response['warn'] , suc_text:""});
		}
		else{
			this.setState({AdArray:d_response['ads'], mod:d_response['mod']});
		}
	}

	async UpdateDonator(auth){
		if(auth){
			this.setState({donor:true})
		}
	}

	render(){
		if(this.state.mod){
			var mod_button = (<span className="mod-link"><Link to="/mod">Mod Mode</Link></span>);
		}
		return (<div id="user-container">
				<h2>Your Banners</h2>
				{mod_button}
				<span className="all-link"><Link to="/all">View All</Link></span>

				<div id="ad-button-container">
				  <AdCreateButton onClickCallBack={this.AdCreateOnClick}/>
				  <AdCreationForm visibility={this.state.AdCVisibility} isDonor={this.state.donor} opacity={this.state.AdCOpacity} height={this.state.AdCHeight} UpdateDetails={this.UpdateDetails}/>
				</div>
				{!this.state.donor && 	<div id="donor-button-container">
						<DonorExtendButton onClickCallBack={this.DonorAuthOnClick}/>
						<DonorPanel visibility={this.state.DonorCVisibility} opacity={this.state.DonorCOpacity} height={this.state.DonorCHeight} UpdateDonator={this.UpdateDonator}/>
					</div>
				}
				<AdDetailsTable adData={this.state.AdArray} updateDetailsCallback={this.UpdateDetails}/>
			</div>)
	}

}
