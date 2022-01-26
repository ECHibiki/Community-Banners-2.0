import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';
import {AdCreationForm, AdCreateButton} from '../form/ad-create';
import {AdDetailsTable} from '../table/user-details-table';

import {Link} from "react-router-dom";

export class UserContainer extends Component{
	constructor(props){
		super(props);
		this.AdCreateOnClick = this.AdCreateOnClick.bind(this);
		this.state = {AdCVisibility:"unset", AdCHeight:"0em", AdCOpacity:"0", AdArray:[], mod:false};
		this.UpdateDetails = this.UpdateDetails.bind(this);
	}

	componentDidMount(){
		this.UpdateDetails();
	}

	AdCreateOnClick(){
		if(this.state.AdCVisibility == "unset")
			this.setState({AdCVisibility:"initial", AdCHeight:"19.4em", AdCOpacity:"1"});
		else
			this.setState({AdCVisibility:"unset", AdCHeight:"0em", AdCOpacity:"0"});
	}

	async UpdateDetails(){
		var d_response = await APICalls.callRetrieveUserAds();
		if("message" in d_response){
			if("errors" in d_response){
				var reasons_arr = []
				for(var reason in d_response['errors']){
					reasons_arr.push(d_response['errors'][reason]);
				}
				var key_ind = 0;
				this.setState({err_text:reasons_arr.map((r) => <span key={key_ind++}>{r}<br/></span> ), war_text: "", suc_text:""});
			}
			else{
				this.setState({err_text:<span>Authorization Failed, Please Refresh<br/></span>});
			}
		}
		else if("warn" in d_response){
			this.setState({war_text:d_response['warn']});
		}
		else{
			this.setState({AdArray:d_response['ads'], mod:d_response['mod']});
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
				  <AdCreationForm visibility={this.state.AdCVisibility} opacity={this.state.AdCOpacity} height={this.state.AdCHeight} UpdateDetails={this.UpdateDetails}/>
				</div>
				<AdDetailsTable adData={this.state.AdArray} updateDetailsCallback={this.UpdateDetails}/>
			</div>)
	}

}
