import { MasterPage } from './components/page/master';
import { AllPage } from './components/page/all';
import { ModPage } from './components/page/mod';
import { FooterInfo } from './components/information/footer';
import ReactDOM from 'react-dom';
import React from 'react';

import {
  BrowserRouter as Router,
  Switch,
  Route
} from "react-router-dom";

if(document.getElementById("index")){
	ReactDOM.render(
		<Router>
		 <Switch>
  		  <Route path="/all" component={AllPage} />
  		  <Route path="/mod" component={ModPage} />
  		  <Route path="/" component={MasterPage} />
		 </Switch>
		</Router>,
	document.getElementById("index"));
}
if(document.getElementById("footer")){
	ReactDOM.render(
		<FooterInfo />,
	document.getElementById("footer"));
}
