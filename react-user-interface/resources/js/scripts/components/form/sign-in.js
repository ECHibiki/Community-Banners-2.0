import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';

export class SignInButton extends Component{
	constructor(props){
		super(props);
	}

	render(){
		return (<div id="sign-in-start"><button onClick={this.props.onClickCallBack} type="button" className="btn btn-primary" >Sign In</button></div>);
	}

}


export class SignInForm extends Component{
	constructor(props){
		super(props);
	}
	render(){
		return(<div style={{visibility: this.props.visibility, opacity: this.props.opacity, height: this.props.height}} id="sign-form">
				<div className="form-group">
					<label htmlFor="name-si">UserName</label>
					<input className="form-control" id="name-si" placeholder="insert username" required/>
				</div>
				<div className="form-group">
					<label htmlFor="pass-si">Password</label>
					<input type="password" className="form-control" id="pass-si" placeholder="" required/>
				</div>
				<SignInAPIButton swapPage={this.props.swapPage}  />
			</div>);
	}
}

export class SignInAPIButton extends Component{
	constructor(props){
		super(props);
		this.SendUserSignIn = this.SendUserSignIn.bind(this);
		this.state = {info_text:"", info_class:"", cursor:"pointer"};
	}

	async SendUserSignIn(e){
		var name = document.getElementById("name-si").value;
		var pass = document.getElementById("pass-si").value;
		this.setState({cursor:"progress"});
		var si_response = await APICalls.callSignIn(name, pass);
		this.setState({cursor:"pointer"});
		if("message" in si_response){
			if("errors" in si_response){
				var reasons_arr = []
				for(var reason in si_response['errors']){
					reasons_arr.push(si_response['errors'][reason]);
				}
				var key_ind = 0;
				this.setState({
					info_text:reasons_arr.map((r) => <span key={key_ind++}>{r}<br/></span> ),
					info_class:"text-danger"
				});
			}
			else{
				this.setState({
					info_text:<span>{response['message']}<br/></span>,
					info_class:"text-danger"
				});
			}
		}
		else if("warn" in si_response){
			this.setState({info_text:si_response['warn'], info_class:"text-warning bg-dark"});
		}
		else{
			this.setState({info_text:si_response['log'], info_class:"text-success"});
			DataStore.storeAuthToken(si_response['access_token']);
			this.props.swapPage();
		}
	}

	render(){
		return (
			<div id="sign-in-finish">
				<button type="button" className="btn btn-secondary" style={{cursor:this.state.cursor}} onClick={this.SendUserSignIn}>Submit</button>
				<p className={this.state.info_class}  id="si-info-field" >{this.state.info_text}</p>
			</div>);
	}

}
