import React, { Component , createRef } from 'react';
import {DataStore, APICalls} from '../../network/api';
import {free_mode} from '../../settings';
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
		this.submit_ref = createRef();
	}
	render(){
		return(<div style={{visibility: this.props.visibility, opacity:this.props.opacity, maxHeight: this.props.height}} className="create-form basic-form">
				<div className="form-group">
					<label htmlFor="name-c">UserName</label>
					<input className="form-control" id="name-c" placeholder="insert username"
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
					<label htmlFor="pass-c">Password</label>
					<input type="password" className="form-control" id="pass-c"
						onKeyDown={
							(e) => {
								if(e.key.toLowerCase() == "enter"){
									this.submit_ref.current.click()
								}
							}
						}
					placeholder="5 character min" required/>
				</div>
				<div className="form-group">
					<label htmlFor="pass-c-conf">Confirm Password</label>
					<input type="password" className="form-control" id="pass-c-conf"
					onKeyDown={
						(e) => {
							if(e.key.toLowerCase() == "enter"){
								this.submit_ref.current.click()
							}
						}
					}
					placeholder="confirmation" required/>
				</div>
				<div style={{visibility: free_mode ? "hidden" : "unset"}} className="form-group">
						<label htmlFor="funder-c">Funder Token</label>
						<input className="form-control" id="funder-c"
							onKeyDown={
								(e) => {
									if(e.key.toLowerCase() == "enter"){
										this.submit_ref.current.click()
									}
								}
							}
						placeholder="(optional)Access additional features"/>
					</div>

				<CreateAPIButton ButtonRef={this.submit_ref} swapPage={this.props.swapPage}  />
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
		var donor = document.getElementById("funder-c").value;
		var pass_confirmation = document.getElementById("pass-c-conf").value;
		this.setState({cursor:"progress"});
		var response = await APICalls.callCreate(name, pass, pass_confirmation);
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
			var response = await APICalls.callSignIn(name, pass, donor);
			if("error" in response){
				var reasons_arr = []
				for(var reason in response['error']){
					reasons_arr.push(response['error'][reason]);
				}
				var key_ind = 0;
				this.setState({
					info_text:reasons_arr.map((r) => <span key={key_ind++}>{r}<br/></span> ),
					info_class:"text-danger"
				});
			}
			else if("warn" in response){
				this.setState({info_text:response['warn'], info_class:"text-warning bg-dark"});
			}
			else{
				this.setState({info_text:response['log'], info_class:"text-success"});
				// Token gets stored by server response
				// DataStore.storeAuthToken(response['access_token']['code']);
				this.props.swapPage();
			}
		}
	}

	render(){
		return (
			<div id="create-finish">
				<button type="button" className="btn btn-secondary"
					style={{cursor:this.state.cursor}}
					onClick={this.SendUserCreate}
					ref={this.props.ButtonRef}
				>Create</button>
				<p className={"err-field " + this.state.info_class}  id="c-info-field" >{this.state.info_text}</p>
			</div>);
	}

}
