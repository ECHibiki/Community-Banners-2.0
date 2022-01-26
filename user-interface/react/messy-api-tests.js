import { APICalls } from '../../../resources/js/components/api';
import { DataStore } from '../../../resources/js/components/api';
import axios from 'axios';
import polyfill from "babel-polyfill";

// These tests exist because axios on jest seems to have cors issues on localhost and mocking doesn't seem worth

async function case1(){

	var create_re = await APICalls.callCreate("test" + Math.ceil(Math.random() * 1000), "2xhashpass");
	console.log("Create User: " + (JSON.stringify(create_re) == JSON.stringify({'log':'Successfully Created'})));
	console.log("    " + JSON.stringify(create_re));
}; 

async function case1f1(){

	var create_re = await APICalls.callCreate("test" + Math.ceil(Math.random() * 1000), "");
	console.log("Create User: " + (create_re['message'] == 'The given data was invalid.'));
	console.log("    " + JSON.stringify(create_re));
}; 


async function case2(){

	var create_re = await APICalls.callSignIn("test296", "2xhashpass");
	console.log("Sign in: " + (create_re['log'] == 'Successfully Logged In'));
	console.log("    " + JSON.stringify(create_re));
}; 

async function case3(){
	DataStore.storeAuthToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC9sb2NhbGhvc3Q6ODAwMFwvYXBpXC9sb2dpbiIsImlhdCI6MTU4Mjc4ODQyNCwiZXhwIjoxNTgyNzkyMDI0LCJuYmYiOjE1ODI3ODg0MjQsImp0aSI6InhYYXh5dXd1M3lEUURWQnkiLCJzdWIiOjMsInBydiI6Ijg3ZTBhZjFlZjlmZDE1ODEyZmRlYzk3MTUzYTE0ZTBiMDQ3NTQ2YWEifQ.S8PP48Vq1MnULHYxzLJHZqMe96IlfTex92HiyxHHu0Q");
	var create_re = await APICalls.callCreateNewAd(null, "https://test.com");
	console.log("Create new ad: " + (create_re['message'] == 'Successfully Logged In'));
	console.log("    " + JSON.stringify(create_re));
}
async function case4(){
	DataStore.storeAuthToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC9sb2NhbGhvc3Q6ODAwMFwvYXBpXC9sb2dpbiIsImlhdCI6MTU4Mjc4ODQyNCwiZXhwIjoxNTgyNzkyMDI0LCJuYmYiOjE1ODI3ODg0MjQsImp0aSI6InhYYXh5dXd1M3lEUURWQnkiLCJzdWIiOjMsInBydiI6Ijg3ZTBhZjFlZjlmZDE1ODEyZmRlYzk3MTUzYTE0ZTBiMDQ3NTQ2YWEifQ.S8PP48Vq1MnULHYxzLJHZqMe96IlfTex92HiyxHHu0Q");
	var create_re = await APICalls.callRetrieveUserAds();
	console.log("get info: " + (create_re == 'Successfully Logged In'));
	console.log("    " + JSON.stringify(create_re));
}
async function case5(){
	DataStore.storeAuthToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOlwvXC9sb2NhbGhvc3Q6ODAwMFwvYXBpXC9sb2dpbiIsImlhdCI6MTU4Mjc4ODQyNCwiZXhwIjoxNTgyNzkyMDI0LCJuYmYiOjE1ODI3ODg0MjQsImp0aSI6InhYYXh5dXd1M3lEUURWQnkiLCJzdWIiOjMsInBydiI6Ijg3ZTBhZjFlZjlmZDE1ODEyZmRlYzk3MTUzYTE0ZTBiMDQ3NTQ2YWEifQ.S8PP48Vq1MnULHYxzLJHZqMe96IlfTex92HiyxHHu0Q");
	var create_re = await APICalls.callRemoveUserAds("asd", "https://ki.com");
	console.log("get info: " + (create_re == 'Successfully Logged In'));
	console.log("    " + JSON.stringify(create_re));
}

//case1();
//case1f1();
//case2();
//case3();
//case4();
//case5();
