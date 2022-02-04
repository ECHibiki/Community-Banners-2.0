import React, { Component } from 'react';
export class AllDetailsTable extends Component{
	constructor(props){
		super(props);

		this.state = {row_data:[]}
	}

	JSXRowData(adData_const){
		var adData = [ ...adData_const];

		var JSX_var = [];
		//sort it
		if(this.props.sortingDetails == "none"){}
		else{
			for(var index_i = 0; index_i < adData.length; index_i++){
				for(var index_j = 0; index_j < adData.length; index_j++){
					if(adData[index_i]['clicks'] > adData[index_j]['clicks']){
						let ad_tmp = adData[index_i];
						adData[index_i] = adData[index_j];
					  adData[index_j] = ad_tmp;
					}
				}
			}
		}
		if(this.props.sortingDetails == "asc"){
			adData = adData.reverse();
		}
		else{}

		//filter it
		for(var index in adData){
			var entry = adData[index];
			if(this.props.filterDetails == "none" || entry['size'] == this.props.filterDetails){

				JSX_var.push(<AllDetailsEntry updateDetailsCallback={this.props.updateDetailsCallback}
					id={"banner-" + index} key={"banner-"+index} name={entry['name']}
					ad_src={entry['uri']} url={entry['size'] == "wide" ? entry['url'] : "[-]"}
					size={entry['size']}
					click_count={(entry['size'] == "wide" ? entry['clicks'] : "-")}/>);
			}
		}
		return JSX_var;
	}

	render(){
		const row_data = this.props.adData;

		return (<div id="ad-details-table" className="table table-striped table-responsive">
			<table>
				<caption>ありがとうございます!</caption>
				<thead className="thead-dark">
					<tr>
						<th className="ad-th-name">Name</th>
						<th className="ad-th-img">Image</th>
						<th className="ad-th-url">URL</th>
						<th className="ad-th-clicks">Clicks</th>
					</tr>
				</thead>
				<tbody className="" data-sorting={this.props.sortingDetails}>
					{this.JSXRowData(row_data)}
				</tbody>
			</table>
			</div>);
	}
}


export class AllDetailsEntry extends Component{
	constructor(props){
		super(props);
	}

	render(){
		return(<tr id={this.props.id} className="">
				<td className="ad-td-name"><span className="ad-td-name-text">{this.props.name}</span></td>
				<td className="ad-td-img"><a href={this.props.ad_src} ><img src={this.props.ad_src}/></a></td>
				<td className={"ad-td-url " + (this.props.size == "small" ? "url-absent" : "")}>{this.props.size == "wide" ? <a href={this.props.url}>{this.props.url}</a> : this.props.url}</td>
				<td className="ad-td-clicks"><span className="ad-td-clicks-text">{this.props.click_count}</span></td>
			</tr>);
	}
}
