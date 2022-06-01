import React, { Component , createRef } from 'react';
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
		this.submit_ref = createRef();
	}
	render(){
		return(<div style={{visibility: this.props.visibility, opacity: this.props.opacity, maxHeight: this.props.height}} className="sign-form basic-form">
				<div className="form-group">
					<label htmlFor="name-si">UserName</label>
					<input className="form-control" id="name-si" placeholder="insert username"
						onKeyDown={
							(e) => {
								if(e.key.toLowerCase() == "enter"){
									this.submit_ref.current.click()
								}
							}
						}
					required/>
				</div>
				<div className="form-group">
					<label htmlFor="pass-si">Password</label>
					<input type="password" className="form-control" id="pass-si"
						onKeyDown={
							(e) => {
								if(e.key.toLowerCase() == "enter"){
									this.submit_ref.current.click()
								}
							}
						}
					placeholder="" required/>
				</div>
				<div className="form-group">
					<label htmlFor="funder-si">Funder Token</label>
					<input className="form-control" id="funder-si"
						onKeyDown={
							(e) => {
								if(e.key.toLowerCase() == "enter"){
									this.submit_ref.current.click()
								}
							}
						}
					placeholder="(optional)Access additional features"/>
				</div>
				<SignInAPIButton ButtonRef={this.submit_ref} swapPage={this.props.swapPage}  />
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
		var donor = document.getElementById("funder-si").value;
		this.setState({cursor:"progress"});
		var si_response = await APICalls.callSignIn(name, pass , donor);
		this.setState({cursor:"pointer"});
		if("error" in si_response){
			this.setState({
				info_text:si_response['error'],
				info_class:"text-danger"
			});
		}
		else if("warn" in si_response){
			this.setState({info_text:si_response['warn'], info_class:"text-warning bg-dark"});
		}
		else{
			this.setState({info_text:si_response['log'], info_class:"text-success"});
			// token gets stored by server response
			// DataStore.storeAuthToken(si_response['access_token']['code']);
			this.props.swapPage(si_response['donor'] == true);
		}
	}

	render(){
		return (
			<div id="sign-in-finish">
				<button type="button" className="btn btn-secondary" style={{cursor:this.state.cursor}}
				ref={this.props.ButtonRef}
				onClick={this.SendUserSignIn}>Submit</button>
				<p className={"err-field " + this.state.info_class}  id="si-info-field" >{this.state.info_text}</p>
			</div>);
	}

}
