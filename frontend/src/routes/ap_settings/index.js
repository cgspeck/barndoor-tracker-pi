import { h, Component } from 'preact';
import linkState from 'linkstate';

import Button from 'preact-material-components/Button';
import TextField from 'preact-material-components/TextField';
import 'preact-material-components/Button/style.css';
import 'preact-material-components/TextField/style.css';
import style from './style';

import { getAPSettings, setAPSettings } from '../../lib/settings';

export default class APSettings extends Component {
	state = {
		apSettings: {
			ssid: null,
			key: null
		},
		error: null,
		info: null,
	};

	async componentDidMount() {
		getAPSettings()
			.then(r => this.setState({ apSettings: {...r} }))
			.catch(e => this.handleError(e));
	}

	handleError = e => {
		console.error('problem', e)
		this.setState({error: e});
	}

	onSubmit = e => {
		e.preventDefault();
		this.setState({error: null, info: null});
		const { ssid, key } = this.state.apSettings;
		setAPSettings(ssid, key)
			.then(r => this.setState({ info: "AP settings updated, device will restart soon" }))
			.catch(e => this.handleError(e));
  }

	errorToast() {
		if (this.state.error != null) {
			return(
				<p>
					{ this.state.error.toString() }
				</p>
			)
		}
	}

	infoToast() {
		if (this.state.info != null) {
			return(
				<p>
					{ this.state.info.toString() }
				</p>
			)
		}
	}

	render({}, { apSettings }) {
		return(
			<div class={style.ap}>
				<h1>Access Point Settings</h1>
				{this.infoToast()}
				{this.errorToast()}
				<form onSubmit={this.onSubmit.bind(this)}>
					<p>
						<TextField label="SSID" value={apSettings.ssid} onInput={linkState(this, 'apSettings.ssid')}></TextField>
					</p>
					<p>
						<TextField label="Key" value={apSettings.key} onInput={linkState(this, 'apSettings.key')}></TextField>
					</p>
					<Button raised ripple onClick={this.onSubmit.bind(this)}>Update</Button>
				</form>
			</div>
		)
	}
};
