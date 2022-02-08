import React, { Component } from 'react';
import {AdRemovalForm, AdRemovalButton} from '../form/ad-remove';

export class AdDetailsTable extends Component{
	constructor(props){
		super(props);
		this.state = {row_data:[]}
	}

	JSXRowData(adData){
		var JSX_var = [];
		for(var index in adData){
			var entry = adData[index];
			entry['uri'] = entry['uri'].replace('public/image/', 'storage/image/');
			entry['clicks'] = entry['size'] == 'small' ? "-" :  entry['clicks'];
			entry['clicks'] = entry['clicks'] == undefined ? "0" : entry['clicks'];
			JSX_var.push(<AdDetailsEntry updateDetailsCallback={this.props.updateDetailsCallback}
				id={"banner-" + index} key={"banner-"+index} ad_src={entry['uri']} url={entry['url']}
				clicks={entry['clicks']} board={entry['board'] ? entry['board'] : "*"}
				/>);
		}
		return JSX_var;
	}

	render(){
		return (<div id="details-table" className="table table-striped table-responsive">
			<table>
				<caption>ありがとうございます!</caption>
				<thead className="thead-dark">
					<tr>
						<th className="ad-th-del">Remove</th>
						<th className="ad-th-img">Image</th>
						<th className="ad-th-url">URL</th>
						<th className="ad-th-clicks">Clicks</th>
						<th className="ad-th-board">Board</th>
					</tr>
				</thead>
				<tbody className="">
				{this.JSXRowData(this.props.adData)}
				</tbody>
			</table>
			</div>);
	}
}



export class AdDetailsEntry extends Component{
	constructor(props){
		super(props);
	}

	render(){
		return(
			<tr id={this.props.id} className="">
				<td className="ad-td-del"><AdRemovalButton updateDetailsCallback={this.props.updateDetailsCallback}  ad_src={this.props.ad_src} url={this.props.url} /></td>
				<td className="ad-td-img"><a href={this.props.ad_src} ><img src={this.props.ad_src}/></a></td>
				<td className={"ad-td-url " + (this.props.url ? "" : "url-absent")}>{this.props.url ? <a href={this.props.url}>{this.props.url}</a> : "[-]"}</td>
				<td className="ad-td-clicks">{this.props.clicks}</td>
				<td className="ad-td-board">{this.props.board}</td>
			</tr>);
	}
}
