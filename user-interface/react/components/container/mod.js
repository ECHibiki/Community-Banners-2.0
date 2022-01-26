import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';
import {ModDetailsTable} from '../table/mod-details-table';
import {Link} from "react-router-dom";
export class ModContainer extends Component{
	constructor(props){
		super(props);
		this.state = {AdArray:[]}
		this.setAllDetails = this.setAllDetails.bind(this);
		this.UpdateDetails = this.UpdateDetails.bind(this);

	}

	componentDidMount(){
		APICalls.callRetrieveModAds(this.setAllDetails, 'AdArray');
	}
	// set warnings from afar
	setAllDetails(state_obj){
		this.setState(state_obj);
	}
	async UpdateDetails(){
		APICalls.callRetrieveModAds(this.setAllDetails, 'AdArray');
	}

	render(){
		return (<div id="mod-container">
			<h2>Mod Mode</h2>
			<span className="mod-link"><Link to="/">Back</Link></span>
			   <ModDetailsTable adData={this.state.AdArray} updateDetailsCallback={this.UpdateDetails}/>
			</div>);
	}
}
