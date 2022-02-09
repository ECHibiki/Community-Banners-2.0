import React, { Component } from 'react';
export class Sort extends Component{
	constructor(props){
		super(props);

		this.filterOnChange = this.filterOnChange.bind(this);
		this.sortByOnChange = this.sortByOnChange.bind(this);
		
	}

	filterOnChange(response){
		this.props.filterChange(response.target.value);
	}
	sortByOnChange(response){
		this.props.sortChange(response.target.value);
	}

	render(){
		return (
			<div id="select-fields">
	      <span id="filter-selections">Banner Size:&nbsp;
	        <select onChange={this.filterOnChange} defaultValue={this.props.filter}>
	          <option value='none'>All</option>
	          <option value='small'>Small</option>
	          <option value='wide'>Wide</option>
	        </select>
	      </span>
	      <span id="sortby-selections"> Sort By:&nbsp;
	        <select onChange={this.sortByOnChange} defaultValue={this.props.sort}>
	          <option value='none'>Chronological</option>
	          <option value='asc'>Ascending</option>
	          <option value='desc'>Descending</option>
	        </select>
	      </span>
	    </div>
		);
	}
}
