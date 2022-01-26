import React, { Component } from 'react';
import {SignInForm, SignInButton} from '../form/sign-in';
import {CreationForm, CreateButton} from '../form/create';
import {Link} from "react-router-dom";
export class LoginContainer extends Component{
	constructor(props){
		super(props);
		this.state = {si_visibility:"unset", si_height:"0em", si_opacity:"0", c_visibility:"unset", c_height:"0em", c_opacity:"0"};

		this.SignInOnClick = this.SignInOnClick.bind(this);
		this.CreateOnClick = this.CreateOnClick.bind(this);

	}

	SignInOnClick(){
		if(this.state.si_visibility=="unset"){
			this.setState({si_visibility:"initial", si_height:"17em", si_opacity:"1"});
		}
		else{
			this.setState({si_visibility:"unset", si_height:"0em", si_opacity:"0"});
		}
	}

	CreateOnClick(){
		if(this.state.c_visibility=="unset"){
			this.setState({c_visibility:"initial", c_height:"22.5em", c_opacity:"1"});
		}
		else{
			this.setState({c_visibility:"unset", c_height:"0em", c_opacity:"0"});
		}
	}

	render(){
		    return (<div id="login-container">
			    <div className="mid-header-container">
			    <h2>Authentication</h2>
			     <span className="all-link"><Link to="/all">View All</Link></span>
			    </div>
			    <div id="si-button-container">
  			    <SignInButton onClickCallBack={this.SignInOnClick}/>
    				<SignInForm  swapPage={this.props.swapPage} opacity={this.state.si_opacity} visibility={this.state.si_visibility} height={this.state.si_height} />
    				<CreateButton onClickCallBack={this.CreateOnClick}/>
    				<CreationForm swapPage={this.props.swapPage} opacity={this.state.c_opacity} visibility={this.state.c_visibility} height={this.state.c_height} />
			    </div>
			</div>)
	}
}
