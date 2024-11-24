import axios from 'axios';
import { sha256 } from 'js-sha256'
import Cookies from 'js-cookie';

import React from 'react';

import {host_addr, host_name} from '../settings';

var error_404 = {"error":"Server error with request"}

export class APICalls{
	static hashPass(name, pass){
		var fake_salt_pass = name + "V" + pass;
		if(name == undefined || pass == undefined || fake_salt_pass == undefined){
			console.log("HASHPASS UNDEFINED VARIABLE");
			return false;
		}
		var hash = sha256(fake_salt_pass);
		return hash;
	}

	static callCreate(name, pass, pass_confirmation){
		var post_data = {"name":name, "pass":pass, "pass_confirmation":pass_confirmation};
		return axios.post(host_addr + '/api/create', post_data, {headers:
			{
				"accept":"application/json", "content-type":"application/json"
			}
			})
			.then(function(res){
				return res.data ;
			})
			.catch(function(err){
				if(!err.response){
					console.log(err);
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});
	}
	static callSignIn(name, pass, donor){
		var post_data = {"name":name, "pass":pass, "token":donor};
		return axios.post(host_addr + '/api/login', post_data, {headers:
			{
				"accept":"application/json", "content-type":"application/json"
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					console.log(err);
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});
	}
	static testDonorToken(token){
		var post_data = {"token":token};
		return axios.post(host_addr + '/api/user/token', post_data, {headers:
			{
				"accept":"application/json", "content-type":"application/json"
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					console.log(err);
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});
	}
	static callCreateNewAd(imagefile, url, hidden_url , board){
		var post_data = new FormData();
		post_data.append("image", imagefile);
		post_data.append("url", url);
		post_data.append("size", hidden_url);
		post_data.append("board", board);
		return axios.post(host_addr + '/api/user/details', post_data, {headers:
			{
				"accept":"application/json", "content-type":"multipart/form-data",
				"authorization": "bearer " + DataStore.getAuthToken()
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					console.log(err);
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});

	}
	static callRetrieveUserAds(){
		return axios.get(host_addr + '/api/user/details', {headers:
			{
				"accept":"application/json",
				"authorization": "bearer " + DataStore.getAuthToken()
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;

			});

	}
	static callRetrieveAllAds(setterCallBack, key){
		axios.get(host_addr + '/api/all', {headers:
			{
				"accept":"application/json"
			}
			})
			.then(function(res){
				setterCallBack({[key]: res.data});
			})
			.catch(function(err){
				console.log(err);
				if(!err.response){
					setterCallBack({[key] : [{fk_name:'404', uri:'' , url:'server out of order'}]});
					return;
				}
				else{
					setterCallBack({[key] : [{fk_name:err.response.status, uri:'', url:JSON.stringify(err.response.data ? err.response.data : error_404)}]});
				}
			});

	}
	static callRetrieveModAds(setterCallBack, key){
		axios.get(host_addr + '/api/user/mod/all', {headers:
			{
				"accept":"application/json",
				"authorization": "bearer " + DataStore.getAuthToken()
			}
			})
			.then(function(res){
				setterCallBack({[key]: res.data});
			})
			.catch(function(err){
				console.log(err);
				if(!err.response){
					setterCallBack({[key] : [{fk_name:'404', uri:'' , url:'server out of order'}]});
					return;
				}
				else{
					setterCallBack({[key] : [{fk_name:err.response.status, uri:'', url:JSON.stringify(err.response.data ? err.response.data : error_404)}]});
				}
			});

	}

	static callRemoveUserAds(uri, url){
		var post_data = {"uri":uri, "url":url};
		return axios.post(host_addr + '/api/user/removal', post_data, {headers:
			{
				"accept":"application/json",
				"authorization": "bearer " + DataStore.getAuthToken()
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					console.log(err);
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});
	}
	static callModRemoveIndividualAds(name, uri, url){
		var post_data = {"target": name, "uri":uri, "url":url};
		return axios.post(host_addr + '/api/user/mod/individual', post_data, {headers:
			{
				"accept":"application/json",
				"authorization": "bearer " + DataStore.getAuthToken()
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					console.log(err);
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});
	}
	static callModRemoveAllUserAds(name){
		var post_data = {"target": name};
		return axios.post(host_addr + '/api/user/mod/purge', post_data, {headers:
			{
				"accept":"application/json",
				"authorization": "bearer " + DataStore.getAuthToken()
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});
	}

	static callEditUserAds(uri, url){
		var post_data = {"uri":uri, "url":url};
		return axios.post(host_addr + '/api/user/edit', post_data, {headers:
			{
				"accept":"application/json",
				"authorization": "bearer " + DataStore.getAuthToken()
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					console.log(err);
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});
	}
	static callModEditIndividualAds(name, uri, url){
		var post_data = {"target": name, "uri":uri, "url":url};
		return axios.post(host_addr + '/api/user/mod/edit', post_data, {headers:
			{
				"accept":"application/json",
				"authorization": "bearer " + DataStore.getAuthToken()
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					console.log(err);
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});
	}

	static callModBanUser(name,hard){
		var post_data = {"target": name,"hard": hard};
		return axios.post(host_addr + '/api/user/mod/ban', post_data, {headers:
			{
				"accept":"application/json",
				"authorization": "bearer " + DataStore.getAuthToken()
			}
			})
			.then(function(res){
				return res.data;
			})
			.catch(function(err){
				if(!err.response){
					console.log(err);
					return error_404;
				}
				return err.response.data ? err.response.data : error_404;
			});
	}

}

export class DataStore{
	static getAuthToken(){
		if(Cookies.get("freeadstoken") != undefined && this.token == undefined){
			this.token = Cookies.get("freeadstoken");
		}
		return this.token;
	}
}
