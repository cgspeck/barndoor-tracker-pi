import { h, Component } from "preact";
import linkState from 'linkstate';

import style from "./style";

import Button from "preact-material-components/Button";
import "preact-material-components/Button/style.css";
import TextField from 'preact-material-components/TextField';
import "preact-material-components/TextField/style.css";

import { getIntervalometerSettings, setIntervalometerSettings } from '../../lib/settings';

export default class IntervalometerSettings extends Component {
  state = {
    intervalometerSettings: {
      bulbInterval: null,
      restInterval: null
    },
    error: null
  };

  async componentDidMount() {
		getIntervalometerSettings()
			.then(r => {
        this.setState({ intervalometerSettings: {...r} })
      })
			.catch(e => this.handleError(e));
	}

  handleError = e => {
		console.error('problem', e)
		this.setState({error: e});
  }

  onSubmit = e => {
		e.preventDefault();
		this.setState({error: null, info: null});
		const { bulbInterval, restInterval } = this.state.intervalometerSettings;
		setIntervalometerSettings(bulbInterval, restInterval)
			.then(r => this.setState(
        {
          info: "Intervalometer Settings Updated",
          intervalometerSettings: {...r}
        }))
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

  render({}, { intervalometerSettings }) {
    return(
      <div class={style.main}>
        <h1>Intervalometer Settings</h1>
				{this.infoToast()}
				{this.errorToast()}
        <form onSubmit={this.onSubmit.bind(this)}>
          <p>
            <TextField label="Bulb Interval (seconds)" value={intervalometerSettings.bulbInterval} onInput={linkState(this, 'intervalometerSettings.bulbInterval')}></TextField>
          </p>
          <p>
            <TextField label="Rest Interval (seconds)" value={intervalometerSettings.restInterval} onInput={linkState(this, 'intervalometerSettings.restInterval')}></TextField>
          </p>
          <Button raised ripple onClick={e => e.prevent_default}>
            Update
          </Button>
        </form>
      </div>
    )
  }
}
