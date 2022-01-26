import React, { Component } from 'react';
import {DataStore, APICalls} from '../../network/api';

export class CreateButton extends Component{
	render(){
		return (
			<div id="create-start">
				<button onClick={this.props.onClickCallBack} type="button" className="btn btn-outline-dark" >New User</button>
			</div>);
	}
}

export class CreationForm extends Component{
	constructor(props){
		super(props);
	}
	render(){
		return(<div style={{visibility: this.props.visibility, opacity:this.props.opacity, height: this.props.height}} id="create-form">
				<div className="form-group">
					<label htmlFor="name-c">UserName</label>
					<input className="form-control" id="name-c" placeholder="insert username" required/>
				</div>
				<div className="form-group">
					<label htmlFor="pass-c">Password</label>
					<input type="password" className="form-control" id="pass-c" placeholder="5 character min" required/>
				</div>
				<div className="form-group">
					<label htmlFor="pass-c-conf">Confirm Password</label>
					<input type="password" className="form-control" id="pass-c-conf" placeholder="confirmation" required/>
				</div>

				<CreateAPIButton swapPage={this.props.swapPage}  />
			</div>);
	}
}

export class CreateAPIButton extends Component{
	constructor(props){
		super(props);
		this.SendUserCreate = this.SendUserCreate.bind(this);
		this.state = {info_text:"", info_class:"", cursor:"pointer"};
	}

	async SendUserCreate(e){
		var name = document.getElementById("name-c").value;
		var pass = document.getElementById("pass-c").value;
		var pass_confirmation = document.getElementById("pass-c-conf").value;
		this.setState({cursor:"progress"});
		var response = await APICalls.callCreate(name, pass, pass_confirmation);
		this.setState({cursor:"pointer"});
		if("message" in response){
			if("errors" in response){
				var reasons_arr = []
				for(var reason in response['errors']){
					reasons_arr.push(response['errors'][reason]);
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
		else if("warn" in response){
			this.setState({info_text:response['warn'], info_class:"text-warning bg-dark"});
		}
		else{
			this.setState({info_text:response['log'], info_class:"text-success"});
			var response = await APICalls.callSignIn(name, pass);
			if("message" in response){
				if("errors" in response){
					var reasons_arr = []
					for(var reason in response['errors']){
						reasons_arr.push(response['errors'][reason]);
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
			else if("warn" in response){
				this.setState({info_text:response['warn'], info_class:"text-warning bg-dark"});
			}
			else{
				this.setState({info_text:response['log'], info_class:"text-success"});
				DataStore.storeAuthToken(response['access_token']);
				this.props.swapPage();
			}
		}
	}

	render(){
		return (
			<div id="create-finish">
				<button type="button" className="btn btn-secondary" style={{cursor:this.state.cursor}} onClick={this.SendUserCreate}>Create</button>
				<p className={this.state.info_class}  id="c-info-field" >{this.state.info_text}</p>
			</div>);
	}

}
