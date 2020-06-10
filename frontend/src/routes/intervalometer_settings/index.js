import { h, Component } from "preact";
import style from "./style";

import Button from "preact-material-components/Button";
import "preact-material-components/Button/style.css";
import Switch from "preact-material-components/Switch";
import "preact-material-components/Switch/style.css";
import TextField from 'preact-material-components/TextField';
import "preact-material-components/TextField/style.css";

import { getInitialTrackStatus, getTrackState } from '../../lib/settings';
import { startHoming, startTracking, stopTracking, toggleIntervalometer, toggleDewController } from '../../lib/commands';
import { setInterval } from "timers";

export default class IntervalometerSettings extends Component {
  state = {
    trackingState: 'Idle',
    intervalometerEnabled: null,
    dewControllerEnabled: null,
    error: null
  };

  async componentDidMount() {
		getInitialTrackStatus()
			.then(r => {
        this.setState({
          trackingState: r.trackingState,
          intervalometerEnabled: r.intervalometerEnabled,
          dewControllerEnabled: r.dewControllerEnabled
        })
      })
			.catch(e => this.handleError(e));
	}

  handleError = e => {
		console.error('problem', e)
		this.setState({error: e});
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

  onHomePressed = e => {
    e.preventDefault();
    startHoming()
      .then(r => this.setState({ trackingState: r }))
      .catch(e => this.handleError(e));
  }

  onTrackPressed = e => {
    e.preventDefault();
    startTracking()
      .then(r => this.setState({ trackingState: r }))
      .catch(e => this.handleError(e));
  }

  onStopPressed = e => {
    e.preventDefault();
    stopTracking()
      .then(r => this.setState({ trackingState: r }))
      .catch(e => this.handleError(e));
  }

  onIntervalometerToggled = e => {
    const enabled = e.target.checked;
    console.log(`Intervalometer toggled to: ${enabled ? "enabled" : "disabled"}`);
    this.setState({intervalometerEnabled: enabled})
    toggleIntervalometer(enabled)
      .then(r => this.setState({intervalometerEnabled: r}))
      .catch(e => this.handleError(e));
  }

  onDewControllerEnabled = e => {
    const enabled = e.target.checked;
    console.log(`Dew controller toggled to: ${enabled ? "enabled" : "disabled"}`);
    this.setState({dewControllerEnabled: enabled})
    toggleDewController(enabled)
      .then(r => this.setState({dewControllerEnabled: r}))
      .catch(e => this.handleError(e));
  }

  homeButton() {
    if (this.state.trackingState == 'Idle') {
      return(
        <p>
          <Button raised ripple onClick={this.onHomePressed.bind(this)}>
            Home
          </Button>
        </p>
      )
    }
  }

  trackButton() {
    if (this.state.trackingState == 'Homed') {
      return(
        <p>
          <Button raised ripple onClick={this.onTrackPressed.bind(this)}>
            Track
          </Button>
        </p>
      )
    }
  }

  stopButton() {
    if (this.state.trackingState == 'Tracking') {
      return(
        <p>
          <Button raised ripple onClick={this.onStopPressed.bind(this)}>
            Stop
          </Button>
        </p>
      )
    }
  }

  render({}, { trackingState, intervalometerEnabled, dewControllerEnabled }) {
    return(
      <div class={style.main}>
        <h1>Intervalometer Settings</h1>
        <div>
          <p>
            <TextField label="Bulb (seconds)" value="30"></TextField>
          </p>
          <p>
            <TextField label="Rest (seconds)" value="30"></TextField>
          </p>
          <Button raised ripple onClick={e => e.prevent_default}>
            Update
          </Button>
        </div>
      </div>
    )
  }
}
