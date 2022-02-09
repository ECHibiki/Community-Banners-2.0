import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';
import Popup from "reactjs-popup";

export class AdRemovalForm extends Component{
	constructor(props){
		super(props);
		this.state = {display: "none", height: "10em"};
	}
	render(){
		return(<div style={{display: this.state.display, maxHeight: this.state.height}} className="ad-remove-form basic-form">
			<form>
				<input type="hidden" id="r-uri" required/>
				<input type="hidden" id="r-url" required/>
				<small className="form-text text-muted text-danger">Confirm Delete for this Ad</small>
				<AdRemovalAPIButton />
			</form></div>);
	}

}

export class AdRemovalButton extends Component{
	constructor(props){
		super(props);
	}

	render(){
		return (<div className="ad-remove">
			<Popup trigger={<button type="button" className="btn btn-secondary btn-sm">Remove</button>}>
  			{close => (
  				<div>
    				<p className="text-danger">Are you sure you want to delete this?</p>
    				<AdRemovalAPIButton updateDetailsCallback={this.props.updateDetailsCallback}  onClickCallBack={close} ad_src={this.props.ad_src} url={this.props.url}/>
  				</div>
  			)}
			</Popup></div>);
	}
}

export class AdRemovalAPIButton extends Component{
	constructor(props){
		super(props);
		this.state = {info_text:"", info_class:"", cursor:"pointer", click_num:0 };
		this.RemoveAd = this.RemoveAd.bind(this);
	}

	async RemoveAd(){
		if(this.state.click_num != 1){
			this.setState({click_num : 1});
			return;
		}
		const uri=this.props.ad_src;
		const url= this.props.url;
		this.setState({cursor:"progress"});
		var response = await APICalls.callRemoveUserAds(uri,url);
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
			this.props.onClickCallBack();
		}
	}

	render(){
		return (
      <div className="ad-remove">
				<button type="button" className="btn btn-danger" style={{cursor:this.state.cursor}} onClick={this.RemoveAd}>Remove</button><br/>
				<small>Click Twice ( {this.state.click_num} )</small>
  			<p className={"err-field " + this.state.info_class}  id="cr-info-field" >{this.state.info_text}</p>
			</div>);
	}
}
