import { h, Component } from 'preact';
import style from './style';

import Button from 'preact-material-components/Button';
import 'preact-material-components/Button/style.css';
import Switch from 'preact-material-components/Switch';
import 'preact-material-components/Switch/style.css';
import TextField from 'preact-material-components/TextField';
import 'preact-material-components/TextField/style.css';

import { getInitialTrackStatus, getTrackState } from '../../lib/settings';
import {
  startHoming,
  startTracking,
  stopTracking,
  toggleIntervalometer,
  toggleDewController,
} from '../../lib/commands';
import { setInterval } from 'timers';

export default class Track extends Component {
  state = {
    trackingState: 'Idle',
    intervalometerEnabled: null,
    dewControllerEnabled: null,
    intervalometerState: null,
    elapsedMillis: null,
    error: null,
  };

  async componentDidMount() {
    getInitialTrackStatus()
      .then((r) => {
        this.setState({
          trackingState: r.state,
          intervalometerEnabled: r.intervalometerEnabled,
          dewControllerEnabled: r.dewControllerEnabled,
          elapsedMillis: r.elapsedMillis,
        });
        console.log('Starting Refresh Interval');
        this.timer = setInterval(this.refreshAlignmentStatus.bind(this), 1000);
      })
      .catch((e) => this.handleError(e));
  }

  handleError = (e) => {
    console.error('problem', e);
    this.setState({ error: e });
  };

  async refreshAlignmentStatus() {
    getTrackState()
      .then((r) => {
        this.setState({
          trackingState: r.state,
          intervalometerState: r.intervalometerState,
          elapsedMillis: r.elapsedMillis,
        });
      })
      .catch((e) => this.handleError(e));
  }

  componentWillUnmount() {
    console.log('Cancelling timer');
    clearInterval(this.timer._id);
  }

  msToTime(duration) {
    var seconds = Math.floor((duration / 1000) % 60),
      minutes = Math.floor((duration / (1000 * 60)) % 60),
      hours = Math.floor((duration / (1000 * 60 * 60)) % 24);

    hours = hours < 10 ? '0' + hours : hours;
    minutes = minutes < 10 ? '0' + minutes : minutes;
    seconds = seconds < 10 ? '0' + seconds : seconds;

    return hours + ':' + minutes + ':' + seconds;
  }

  errorToast() {
    if (this.state.error != null) {
      return <p>{this.state.error.toString()}</p>;
    }
  }

  onHomePressed = (e) => {
    e.preventDefault();
    startHoming()
      .then((r) => this.setState({ trackingState: r }))
      .catch((e) => this.handleError(e));
  };

  onTrackPressed = (e) => {
    e.preventDefault();
    startTracking()
      .then((r) =>
        this.setState({
          trackingState: r,
          intervalometerState: null,
          elapsedMillis: 0,
        }),
      )
      .catch((e) => this.handleError(e));
  };

  onStopPressed = (e) => {
    e.preventDefault();
    stopTracking()
      .then((r) => this.setState({ trackingState: r }))
      .catch((e) => this.handleError(e));
  };

  onIntervalometerToggled = (e) => {
    const enabled = e.target.checked;
    console.log(
      `Intervalometer toggled to: ${enabled ? 'enabled' : 'disabled'}`,
    );
    this.setState({ intervalometerEnabled: enabled });
    toggleIntervalometer(enabled)
      .then((r) => this.setState({ intervalometerEnabled: r }))
      .catch((e) => this.handleError(e));
  };

  onDewControllerEnabled = (e) => {
    const enabled = e.target.checked;
    console.log(
      `Dew controller toggled to: ${enabled ? 'enabled' : 'disabled'}`,
    );
    this.setState({ dewControllerEnabled: enabled });
    toggleDewController(enabled)
      .then((r) => this.setState({ dewControllerEnabled: r }))
      .catch((e) => this.handleError(e));
  };

  homeButton() {
    if (this.state.trackingState == 'Idle') {
      return (
        <p>
          <Button raised ripple onClick={this.onHomePressed.bind(this)}>
            Home
          </Button>
        </p>
      );
    }
  }

  trackButton() {
    if (this.state.trackingState == 'Homed') {
      return (
        <p>
          <Button raised ripple onClick={this.onTrackPressed.bind(this)}>
            Track
          </Button>
        </p>
      );
    }
  }

  stopButton() {
    if (this.state.trackingState == 'Tracking') {
      return (
        <p>
          <Button raised ripple onClick={this.onStopPressed.bind(this)}>
            Stop
          </Button>
        </p>
      );
    }
  }

  intervalometerState() {
    if (
      this.state.intervalometerEnabled === true &&
      this.state.trackingState === 'Tracking' &&
      this.state.intervalometerState != null
    ) {
      return (
        <p>
          <TextField
            label="Intervalometer"
            value={this.state.intervalometerState}
            disabled="true"
          />
        </p>
      );
    }
  }

  elapsedTime() {
    if (
      this.state.trackingState === 'Tracking' &&
      this.state.elapsedMillis != null
    ) {
      return (
        <p>
          <TextField
            label="Elapsed Time"
            value={this.msToTime(this.state.elapsedMillis)}
            disabled="true"
          />
        </p>
      );
    }
  }

  render({}, { trackingState, intervalometerEnabled }) {
    return (
      <div class={style.track}>
        <h1>Track</h1>
        <div>
          {this.errorToast()}
          <p>
            <TextField
              label="Tracking Status"
              value={trackingState}
              disabled="true"
            />
          </p>
          {this.elapsedTime()}
          {this.intervalometerState()}
          <p>
            Intervalometer:{' '}
            <Switch
              onChange={this.onIntervalometerToggled.bind(this)}
              checked={intervalometerEnabled === true}
            ></Switch>
          </p>
          {this.homeButton()}
          {this.trackButton()}
          {this.stopButton()}
        </div>
      </div>
    );
  }
}
