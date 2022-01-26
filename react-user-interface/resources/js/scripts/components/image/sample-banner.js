import React, { Component } from 'react';
import {dimensions_w, dimensions_h, dimensions_small_w,dimensions_small_h} from '../../settings'

export class SampleBanner extends Component{
	render(){
		return (
			<div id="sample-banner">
				<iframe src="/banner?size=wide" scrolling="no" width="500" height="90" className="iframe-banner"></iframe>
			</div>);
	}
}
