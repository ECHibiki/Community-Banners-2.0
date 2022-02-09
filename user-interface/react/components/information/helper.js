import React, { Component } from 'react';
import {dimensions_w, dimensions_h, dimensions_small_w,dimensions_small_h} from '../../settings'

export class HelperText extends Component{
	render(){
		return (<div id="helper">
		<h2>How To Use</h2>
						<p>
							Uploaded images must be {dimensions_w}x{dimensions_h} or {dimensions_small_w}x{dimensions_small_h} and safe for work.
							Wide banners should not be used to promote other existing social platforms or communities and small banners must relate to Kissu in some way.
							Certain exemptions exist for donators who may submit board specific banners. <br/>
							If you wish to upload banners specific to a given board* you need to become a donator(see the above links for more information).
							Using your donator token will give you a new option on banner creation.
						</p>
						<small>*Note that banners to NSFW boards will not show up on /all/, generic boards, or the public banner listing.
						Banners to SFW boards will show up on /all/ and generic pages. This is done to keep the site friendly to those who do not like NSFW content.</small>
					</div>);
	};
}
