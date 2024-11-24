import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';
import Popup from "reactjs-popup";

// NOTE: This thing needs to get isEditing to work
// What we want to happen is that when isEditing changes to true, the URL changes from a standard hyperlink to an input form where we want to grab the contents for editing the URL
export class AdEditForm extends Component{
	constructor(props) {
		super(props);
		this.state = { inputValue: this.props.url };
		this.handleInputChange = this.handleInputChange.bind(this);
	}

	async handleInputChange(event) {
		this.setState({ inputValue: event.target.value });
   		this.props.onInputChange(event.target.value);
	}

	async componentDidUpdate(prevProps) {
		if (prevProps.isEditing && !this.props.isEditing) {
			// Reset inputValue when editing is canceled
			this.setState({ inputValue: this.props.url });
		}
	}

	render(){
		const { isEditing } = this.props;
		const { inputValue } = this.state;
		return(<div>
			{!isEditing ? (<a href={this.props.url}>{this.props.url}</a>) : (
				<input type="url" pattern="/^http(|s):\/\/[-A-Z0-9+&amp;@#\/%?=~_|!:,.;]+\.[A-Z0-9+&amp;@#\/%=~_|]+$/i" class="form-control" placeholder="http/https urls" value={inputValue} onChange={this.handleInputChange} required></input>
			)} 
  			</div>);
	}
}

export class AdEditButton extends Component{
	constructor(props) {
		super(props);

		this.ToggleEditAd = this.props.ToggleEditAd.bind(this);
	}

	render(){
		const { isEditing, ToggleEditAd, inputValue, url } = this.props;
		return (<div>
			{!isEditing ? (<div className="ad-edit"><button type="button" className="btn btn-secondary btn-sm" onClick={() => this.ToggleEditAd(url)}>Edit</button></div>) : (
				<div className="ad-edit-pair">
					<AdEditAPIButton updateDetailsCallback={this.props.updateDetailsCallback} ad_src={this.props.ad_src} inputValue={inputValue} ToggleEditAd={ToggleEditAd} url={this.props.url} />
					<div className="ad-cancel"><button type="button" class="btn btn-danger btn-sm" onClick={() => this.ToggleEditAd(url)}>Cancel</button></div>
				</div>
			)} 
  			</div>);
	}
}

export class AdEditAPIButton extends Component{
	constructor(props){
		super(props);
		this.state = {inputValue: this.props.inputValue };
		this.EditAd = this.EditAd.bind(this);
	}

	async EditAd(){
		const uri = this.props.ad_src;
		const url = this.props.inputValue; // Use the current input value
		this.setState({cursor:"progress"});
		var response = await APICalls.callEditUserAds(uri, url);
		//console.log("URI: " + uri + ", URL: " + url);
		this.setState({cursor:"pointer"});
		if("error" in response){
			this.setState({
				info_text:response['error'],
				info_class:"text-danger"
			});
		}
		else if("warn" in response){
			this.setState({info_text:response['warn'], info_class:"text-warning bg-dark"});
		}
		else{
			this.setState({info_text:response['log'], info_class:"text-success"});
			this.props.updateDetailsCallback();
			this.props.ToggleEditAd(url);
		}
		
	}

	render(){
		return (
		<div className="ad-save"><button type="button" class="btn btn-success btn-sm" onClick={this.EditAd}>Save</button></div>
		);
	}
}
