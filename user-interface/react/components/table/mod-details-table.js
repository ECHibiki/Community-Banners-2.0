import React, { Component } from 'react';
import Popup from "reactjs-popup";
import {DataStore, APICalls} from '../../network/api';
export class ModDetailsTable extends Component{
	constructor(props){
		super(props);
		this.state = {row_data:[]}
	}

	JSXRowData(adData){
		var JSX_var = [];
		for(var index in adData){
			var entry = adData[index];
			entry['uri'] = entry['uri'].replace('public/image/', 'storage/image/');
			JSX_var.push(<ModDetailsEntry updateDetailsCallback={this.props.updateDetailsCallback} id={"banner-" + index}
			key={"banner-"+index} ad_src={entry['uri']} url={entry['url']} name={entry['name']} ban={entry['hardban']}/>);
		}
		return JSX_var;
	}

	render(){
		return (<div id="mod-details-table" className="table table-striped table-responsive">
			<table>
				<caption>ありがとうございます!</caption>
				<thead className="thead-dark">
					<tr>
						<th className="ad-th-del">Delete</th>
						<th className="ad-th-ban">Ban</th>
						<th className="ad-th-name">Name</th>
						<th className="ad-th-img">Image</th>
						<th className="ad-th-url">URL</th>
						<th className="ad-th-ban">Ban State</th>
					</tr>
				</thead>
				<tbody className="">
				{this.JSXRowData(this.props.adData)}
				</tbody>
			</table>
			</div>);
	}
}

export class ModDetailsEntry extends Component{
	constructor(props){
		super(props);
	}

	render(){
		var ban_str = "-";;
		if(this.props.ban == 1){
			var ban_str =  "Hardban";
		}
		else if(this.props.ban == 0){
			var ban_str = "Softban";
		}
		return(<tr id={this.props.id} className="">
				<td className="ad-td-ban"><ModRemovalButton updateDetailsCallback={this.props.updateDetailsCallback} ad_src={this.props.ad_src} url={this.props.url} name={this.props.name}/></td>
				<td className="ad-td-del"><ModBanButton updateDetailsCallback={this.props.updateDetailsCallback} ad_src={this.props.ad_src} url={this.props.url} name={this.props.name}/></td>

				<td className="ad-td-name">{this.props.name}</td>
				<td className="ad-td-img"><a href={this.props.ad_src} ><img src={this.props.ad_src}/></a></td>
				<td className="ad-td-url"><a href={this.props.url}>{this.props.url}</a></td>
				<td className="ad-td-ban">{ban_str}</td>
			</tr>);
	}
}


export class ModRemovalButton extends Component{
	constructor(props){
		super(props);
	}

	render(){
		return (<div className="ad-remove">
			<Popup trigger={<button type="button" className="btn btn-danger btn-sm">Remove</button>}>
			{close => (
				<div>
				<p className="text-danger">Delete all by name or delete individual?</p>
				<ModIndividualRemovalAPIButton updateDetailsCallback={this.props.updateDetailsCallback}  onClickCallBack={close} ad_src={this.props.ad_src} url={this.props.url} name={this.props.name}/>&nbsp;
				<ModCompleteRemovalAPIButton updateDetailsCallback={this.props.updateDetailsCallback}  onClickCallBack={close} ad_src={this.props.ad_src} url={this.props.url} name={this.props.name}/>

				</div>
			)}
			</Popup></div>);
	}
}
export class ModBanButton extends Component{
	constructor(props){
		super(props);
	}

	render(){
		return (<div className="ad-remove">
			<Popup trigger={<button type="button" className="btn btn-info btn-sm">Ban</button>}>
			{close => (
				<div>
				<p className="text-info">Shadow realm this user or hard ban user?</p>
				<ModSoftBanAPIButton updateDetailsCallback={this.props.updateDetailsCallback}  onClickCallBack={close} ad_src={this.props.ad_src} url={this.props.url} name={this.props.name}/>
				<ModHardBanAPIButton updateDetailsCallback={this.props.updateDetailsCallback}  onClickCallBack={close} ad_src={this.props.ad_src} url={this.props.url} name={this.props.name}/>
				</div>
			)}
			</Popup></div>);
	}
}


export class ModIndividualRemovalAPIButton extends Component{
	constructor(props){
		super(props);
		this.state = {info_text:"", info_class:""};

		this.RemoveAd = this.RemoveAd.bind(this);
	}

	async RemoveAd(){
		const uri=this.props.ad_src;
		const url= this.props.url;
		const name = this.props.name;
		var response = await APICalls.callModRemoveIndividualAds(name,uri,url);
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
		return (<div className="ad-remove-soft"><button type="button" className="btn btn-danger btn-sm" onClick={this.RemoveAd}>Remove Selected</button>
			<p className={this.state.info_class}  id="mi-info-field" >{this.state.info_text}</p>
			</div>);
	}
}
export class ModCompleteRemovalAPIButton extends Component{
	constructor(props){
		super(props);
		this.state = {info_text:"", info_class:""};

		this.RemoveAd = this.RemoveAd.bind(this);
	}

	async RemoveAd(){
		const name = this.props.name;

		var response = await APICalls.callModRemoveAllUserAds(name);
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
			this.setState({info_text:response['log'], info_class:"text-success"});			this.props.updateDetailsCallback();
			this.props.onClickCallBack();

		}
	}

	render(){
		return (<div id="ad-remove-hard"><button type="button" className="btn btn-danger btn-sm" onClick={this.RemoveAd}>Remove All</button>
			<p className={this.state.info_class}  id="mc-info-field" >{this.state.info_text}</p>
			</div>);
	}
}
export class ModSoftBanAPIButton extends Component{
	constructor(props){
		super(props);
		this.state = {info_text:"", info_class:""};

		this.SoftBan = this.SoftBan.bind(this);
	}

	async SoftBan(){
		const name = this.props.name;
		var response = await APICalls.callModBanUser(name, 0);
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
			this.setState({info_text:response['log'], info_class:"text-success"});			this.props.updateDetailsCallback();
			this.props.onClickCallBack();

		}
	}

	render(){
		return (<div id="ban-soft"><button type="button" className="btn btn-info btn-sm" onClick={this.SoftBan}>Soft Ban</button>
			<p className={this.state.info_class}  id="sb-info-field" >{this.state.info_text}</p>
			</div>);
	}
}
export class ModHardBanAPIButton extends Component{
	constructor(props){
		super(props);
		this.state = {info_text:"", info_class:""};

		this.HardBan = this.HardBan.bind(this);
	}

	async HardBan(){
		const name = this.props.name;

		var response = await APICalls.callModBanUser(name, 1);
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
			this.props.updateDetailsCallback();
			this.props.onClickCallBack();

		}
	}

	render(){
		return (<div id="ban-hard"><button type="button" className="btn btn-info btn-sm" onClick={this.HardBan}>Hard Ban</button>
			<p className={this.state.info_class}  id="hb-info-field" >{this.state.info_text}</p>
			</div>);
	}
}
