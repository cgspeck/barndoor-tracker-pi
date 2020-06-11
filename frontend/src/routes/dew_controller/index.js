import { h, Component } from "preact";
import linkState from "linkstate";

import Button from "preact-material-components/Button";
import "preact-material-components/Button/style.css";
import Switch from "preact-material-components/Switch";
import "preact-material-components/Switch/style.css";
import TextField from 'preact-material-components/TextField';
import "preact-material-components/TextField/style.css";

import Radio from 'preact-material-components/Radio';
import "preact-material-components/Radio/style.css";

import FormField from 'preact-material-components/FormField';
import 'preact-material-components/FormField/style.css';

import { getDewControllerStatus, setTargetTemperature } from '../../lib/settings';
import { toggleDewController } from '../../lib/commands';
import { setInterval } from "timers";

import style from "./style";

export default class DewController extends Component {
  state = {
    currentTemperature: null,
    currentlyHeating: null,
    dewControllerEnabled: null,
    targetTemperature: null,
    info: null,
    error: null
  };

  async componentDidMount() {
    getDewControllerStatus()
      .then(r => {
        this.setState({
          currentTemperature: r.currentTemperature,
          currentlyHeating: r.currentlyHeating,
          targetTemperature: r.targetTemperature,
          dewControllerEnabled: r.dewControllerEnabled
        })
        console.log("Starting Refresh Interval");
        this.timer = setInterval(this.refreshStatus.bind(this), 1000);
      })
      .catch(e => this.handleError(e));
  }

  handleError = e => {
    console.error('problem', e)
    this.setState({error: e});
  }

  async refreshStatus() {
    getDewControllerStatus()
      .then(r => {
        this.setState({
          currentTemperature: r.currentTemperature,
          currentlyHeating: r.currentlyHeating
        })
      })
      .catch(e => this.handleError(e));
  }

  componentWillUnmount() {
    console.log("Cancelling timer");
    clearInterval(this.timer._id);
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

  onSubmit = e => {
    e.preventDefault();
    this.setState({error: null, info: null});
    const { targetTemperature } = this.state;
    setTargetTemperature(targetTemperature)
      .then(r => this.setState(
        {
          info: "Intervalometer Settings Updated",
          targetTemperature: r.targetTemperature
        }))
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

  render({}, { currentTemperature, currentlyHeating, dewControllerEnabled, targetTemperature }) {
    return(
      <div class={style.main}>
        <h1>Dew Controller</h1>
        <div>
          Enabled: <Switch id="dewControllerEnable" onChange={this.onDewControllerEnabled.bind(this)} checked={dewControllerEnabled === true}></Switch>
          <form onSubmit={this.onSubmit.bind(this)}>
            {this.infoToast()}
            {this.errorToast()}
            <FormField>
              <Radio id="heating" name="Basic Options" disabled="true" checked={currentlyHeating === true}/>
              <label for="heating">Heating</label>
            </FormField>
            <p>
              <TextField
                label="Current Temperature °C"
                value={currentTemperature}
                disabled="true"
              ></TextField>
            </p>
            <p>
              <TextField
                label="Target Temperature °C"
                value={targetTemperature}
                onInput={linkState(this, 'targetTemperature')}
              ></TextField>
            </p>
            <Button raised ripple onClick={e => e.prevent_default}>
              Update
            </Button>
          </form>
        </div>
      </div>
    )
  }
}
