import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';
import {dimensions_w, dimensions_h,dimensions_small_w, dimensions_small_h,rules,rules_small} from '../../settings'

export class AdCreateButton extends Component{
	render(){
		return (
      <div id="create-ad-start">
        <button onClick={this.props.onClickCallBack} type="button" className="btn btn-primary" >New Banner</button>
      </div>);
	}
}

export class AdCreationForm extends Component{
	constructor(props){
		super(props);
		this.state = {file_input: "", url_input: "", board_input: "", hide_url:false}

		this.unsetFormFields = this.unsetFormFields.bind(this);
		this.handleFileChange = this.handleFileChange.bind(this);
		this.handleURLChange = this.handleURLChange.bind(this);
		this.handleBoardChange = this.handleBoardChange.bind(this);

	}

	unsetFormFields(){
		this.setState({file_input:"", url_input: ""});
	}

	handleFileChange(e){
		this.setState({file_input:e.target.value});
	}
	handleURLChange(e){
		this.setState({url_input:e.target.value});
	}
	handleBoardChange(e){
		this.setState({board_input:e.target.value});
	}

	render(){
		return(<div style={{visibility: this.props.visibility, maxHeight: this.props.height, opacity: this.props.opacity}} className="ad-create-form basic-form">
		<label>Banner Type:&nbsp;</label>
			<label><input id="img-size-wide" type="radio" name="size" onClick={()=>this.setState({hide_url:false})} defaultChecked/>&nbsp;Wide</label>&nbsp;
			<label><input id="img-size-small" type="radio" name="size" onClick={()=>this.setState({hide_url:true})}/>&nbsp;Small</label><br/>
				<input type="hidden" id="img-size" value={(this.state.hide_url ? "small" : "wide")}/>
				<div className="form-group">
					<label htmlFor="image-ad-c">Image</label>
					<input onChange={this.handleFileChange} value={this.state.file_input}  type="file" className="form-control-file" id="image-ad-c" accept="image/*" />
					<small  className="form-text text-muted">Must be&nbsp;
						{ (this.state.hide_url ? dimensions_small_w : dimensions_w) }x
						{ (this.state.hide_url ? dimensions_small_h : dimensions_h) }<br/>
						Rules: { (this.state.hide_url ? rules_small : rules) }</small>
				</div>
				 <div className="form-group" style={{display:(this.state.hide_url ? "none" : "initial")}}>
					<label htmlFor="ad-url-c">URL</label>
					<input onChange={this.handleURLChange} value={this.state.url_input} type="url" pattern="/^http(|s):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]+\.[A-Z0-9+&@#\/%=~_|]+$/i" className="form-control" id="ad-url-c" placeholder="http/https urls only" />
				</div>
				 <div className="form-group">
					<label htmlFor="ad-board-c">Board {!this.props.isDonor ? "(Funders only)" : ""}</label>
					<select disabled={!this.props.isDonor} onChange={this.handleBoardChange} value={this.state.board_input} className="" id="ad-board-c" name="board">
						<option value="" default>/all/</option>
						<option value="qa" default>/qa/</option>
						<option value="jp" default>/jp/</option>
						<option value="win" default>/win/</option>
						<option value="cry" default>/cry/</option>
						<option value="poll" default>/poll/</option>
						<option value="b" default>/b/</option>
						<option value="f" default>/f/</option>
						<option value="ec" default>/ec/(NSFW)</option>
						<option value="trans" default>/trans/(NSFW)</option>
						<option value="test" default>/test/(NSFW)</option>
					</select>
					<small>NSFW banners are not viewable on the public banner listings</small>
				</div>
				<AdCreateAPIButton UnsetFormFields={this.unsetFormFields} UpdateDetails={this.props.UpdateDetails}/>
			</div>);
	}
}

export class AdCreateAPIButton extends Component{
	constructor(props){
		super(props);
		this.SendNewBanner = this.SendNewBanner.bind(this);
		this.state = {info_text:"", info_class:"", cursor:"pointer"};
	}

	async SendNewBanner(e){
		var image = document.getElementById("image-ad-c").files[0];
		var url = document.getElementById("ad-url-c").value;
		var size = document.getElementById("img-size").value;
		var board = document.getElementById("ad-board-c").value;
		this.setState({cursor:"progress"});
		var response = await APICalls.callCreateNewAd(image, url,size , board);
		this.setState({cursor:"pointer"});
		if("error" in response){
			this.setState({
				info_text:response['error'],
				info_class:"text-danger"
			});
		}
		else if("warn" in response){
			this.setState({info_text:response['warn'], info_class:"text-warning bg-dark"});
			this.props.UnsetFormFields();
			this.props.UpdateDetails();
		}
		else{
			this.setState({info_text:response['log'], info_class:"text-success"});
			this.props.UnsetFormFields();
			this.props.UpdateDetails();
		}
	}


	render(){
		return (
      <div id="ad-create-finish">
				<button type="button" className="btn btn-secondary" style={{cursor:this.state.cursor}} onClick={this.SendNewBanner}>Create</button>
			  <p className={"err-field " + this.state.info_class}  id="cad-info-field" >{this.state.info_text}</p>
			</div>);
	}
}