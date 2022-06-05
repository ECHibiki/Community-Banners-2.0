import React, { Component } from 'react';
import {dimensions_w, dimensions_h, dimensions_small_w,dimensions_small_h, free_mode} from '../../settings'

export class HelperText extends Component{
	render(){
		return (<div id="helper">
			<h2>How To Use</h2>
						<p>
							Uploaded images must be {dimensions_w}x{dimensions_h} or {dimensions_small_w}x{dimensions_small_h} and safe for work.
							Wide banners should not be used to promote other existing social platforms or communities. Small banners must relate to Kissu in some way.
							<em>{free_mode ? "Banners for NSFW boards may be NSFW*." : "Certain exemptions exist for donators who may submit board specific banners."}</em><br/>
							{free_mode ? "" : "If you wish to upload banners specific to a given board* you need to become a donator(see the above links for more information).\n\
							Upon donating you will be emailed a donation reward token granting you access to this feature."}
						</p>
						<small>*Note that banners to NSFW boards will not show up on /all/, generic boards, or the public banner listing.
						Banners to SFW boards are shown on /all/ and other applicable locations.
						This is done to keep the site friendly to those who do not like NSFW content.</small>
				</div>);
	};
}
