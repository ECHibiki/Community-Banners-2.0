import React, { Component } from 'react';
import {AdRemovalForm, AdRemovalButton} from '../form/ad-remove';
import {AdEditForm, AdEditButton} from '../form/ad-edit';

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
						<th className="ad-th-edit">Edit</th>
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
		this.state = {isEditing: false, inputValue: this.props.url, url: this.props.url};
		this.ToggleEditAd = this.ToggleEditAd.bind(this);
		this.handleInputChange = this.handleInputChange.bind(this);
	}

	async ToggleEditAd(newURL){
		console.log("new URL: " + newURL);
		this.setState({ url: newURL });
		this.setState({ isEditing: !this.state.isEditing });
	}

	async handleInputChange(newValue) {
		this.setState({ inputValue: newValue });
	}

	render(){
		const { isEditing, inputValue, url } = this.state;
		return(
			<tr id={this.props.id} className="">
				<td className="ad-td-del"><AdRemovalButton updateDetailsCallback={this.props.updateDetailsCallback}  ad_src={this.props.ad_src} url={url} /></td>
				<td className="ad-td-img"><a href={this.props.ad_src} ><img src={this.props.ad_src}/></a></td>
				<td className={"ad-td-url" + (url ? "" : " url-absent")}>{url ? <AdEditForm updateDetailsCallback={this.props.updateDetailsCallback}  ad_src={this.props.ad_src} url={url} isEditing={isEditing} inputValue={inputValue} onInputChange={this.handleInputChange} /> : "[-]"}</td>
				<td className="ad-td-clicks">{this.props.clicks}</td>
				<td className="ad-td-board">{this.props.board}</td>
				<td className={"ad-td-edit" + (url ? "" : " url-absent")}> {url ? <AdEditButton updateDetailsCallback={this.props.updateDetailsCallback}  ad_src={this.props.ad_src} url={url} isEditing={isEditing} ToggleEditAd={this.ToggleEditAd} inputValue={inputValue} /> : "[-]"} </td>
			</tr>);
	}
}
