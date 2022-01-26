import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';
import {Sort} from '../sorting/sort';
import {AllDetailsTable} from '../table/all-details-table';
import {Link} from "react-router-dom";
export class AllContainer extends Component{
	constructor(props){
		super(props);
		var preset_data = this.getLocalStoragePresets();
		this.state = {AdArray:[], filter:preset_data['filter'], sorting:preset_data['sort']};

		this.setAllDetails = this.setAllDetails.bind(this);
		this.filterOnChange = this.filterOnChange.bind(this);
		this.sortByOnChange = this.sortByOnChange.bind(this);
	}

	componentDidMount(){
		APICalls.callRetrieveAllAds(this.setAllDetails, 'AdArray');
	}

	setAllDetails(state_obj){
		this.setState(state_obj);
	}

  filterOnChange(response){
    this.setState({filter:response});
		this.setFilterPresets(response);
  }
  sortByOnChange(response){
    this.setState({sorting:response});
		this.setSortPresets(response);
  }

	setFilterPresets(filter){
			localStorage.setItem('filter', filter);
	}
	setSortPresets(sort){
			localStorage.setItem('sort', sort);
	}

	getLocalStoragePresets(){
		return {
			filter: localStorage.getItem('filter') == undefined ? 'none' : localStorage.getItem('filter'),
			sort: localStorage.getItem('sort')  == undefined ? 'none' : localStorage.getItem('sort')
		}
	}

	// possible refactor
	render(){
			return (
				<div id="all-container">
					<h2>All Banners</h2>
					<span className="all-link"><Link to="/">Back</Link></span><br/>

					<Sort sortChange={this.sortByOnChange} filterChange={this.filterOnChange} filter={this.state.filter}
						sort={this.state.sorting}/>

				  <AllDetailsTable adData={this.state.AdArray} updateDetailsCallback={this.setAllDetails}
            filterDetails={this.state.filter}
            sortingDetails={this.state.sorting}/>
			  </div>);
	}
}
