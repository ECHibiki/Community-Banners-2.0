import React, { Component } from 'react';
import {TopHeader} from "../information/header";
import {DonatorBox} from "../information/donator";
import {HelperText} from "../information/helper";
import {SampleBanner} from "../image/sample-banner";
import {PatreonBanner} from "../image/patreon-banner";
import {ModContainer} from "../container/mod";

export class ThemeToggle extends Component{
	constructor(props){
		super(props);
	}

	componentDidMount(){
		var ds = document.getElementById("dark-sheet");
		if(window.localStorage &&
			(window.localStorage.getItem("theme") && window.localStorage.getItem("theme") == "dark" )){
				ds.href="/css/dark.css";
			} else{
				ds.href="";
			}
	}

	render(){
		return(<div className="theme-toggle">
					<i className="fa-solid fa-cloud-moon"
						onClick={function(e){
							var ds = document.getElementById("dark-sheet");
							console.log(ds.href)
							if (ds.href.indexOf("/css/dark.css") != -1){
								window.localStorage ? window.localStorage.setItem("theme" , "light"): "";
								ds.href="";
							} else{
								ds.href="/css/dark.css";
								window.localStorage ? window.localStorage.setItem("theme" , "dark") : "";
							}
						}}
					></i>
				</div>);
	}
}
