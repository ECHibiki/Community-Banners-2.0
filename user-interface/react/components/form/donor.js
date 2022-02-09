import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';


export class DonorExtendButton extends Component{
	render(){
		return (
      <div id="create-ad-start">
        <button onClick={this.props.onClickCallBack} type="button" className="btn btn-primary" >Add more features</button>
      </div>);
	}
}

export class DonorPanel extends Component{
	constructor(props){
		super(props)
    this.updateToken = this.updateToken.bind(this);
    this.checkToken = this.checkToken.bind(this);
		this.state = {info_text:"", info_class:"", token_value:""}
	}
	async checkToken(){
		var response = await APICalls.testDonorToken(this.state.token_value);
		if("error" in response){
			var key_ind = 0;
			this.setState({info_text:response['error'], info_class:"text-danger"});
		}
		else if("warn" in response){
			this.setState({info_text:response['warn'] , info_class:"text-warning bg-dark"});
		}
		else{
			this.setState({info_text:response['log'] , info_class:"text-success"});
			this.props.UpdateDonator(true)
		}
	}
  updateToken(e){
    this.setState({token_value:e.target.value});
  }
	render(){
		return <div style={{visibility: this.props.visibility, maxHeight: this.props.height, opacity: this.props.opacity}} className="donor-panel basic-form">
			<label for="donor-auth">Authorize Donor Status</label>
			<input type="url" class="form-control" id="donor-auth" value={this.state.token_value}
      onChange={this.updateToken}
      placeholder="Insert Donator Token for more options" />
			<button type="button" className="btn btn-secondary" onClick={this.checkToken}>Authorize</button>
			<p className={"err-field " + this.state.info_class}  id="cad-info-field" >{this.state.info_text}</p>
		</div>
	}
}