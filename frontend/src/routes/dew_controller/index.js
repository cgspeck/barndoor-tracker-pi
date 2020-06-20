import { h, Component } from 'preact';
import linkState from 'linkstate';

import Button from 'preact-material-components/Button';
import 'preact-material-components/Button/style.css';
import Switch from 'preact-material-components/Switch';
import 'preact-material-components/Switch/style.css';
import TextField from 'preact-material-components/TextField';
import 'preact-material-components/TextField/style.css';

import Radio from 'preact-material-components/Radio';
import 'preact-material-components/Radio/style.css';

import FormField from 'preact-material-components/FormField';
import 'preact-material-components/FormField/style.css';

import Dialog from 'preact-material-components/Dialog';
import 'preact-material-components/Dialog/style.css';

import List from 'preact-material-components/List';
import 'preact-material-components/List/style.css';

import Slider from 'preact-material-components/Slider';
import 'preact-material-components/Slider/style.css';

import {
  getDewControllerStatus,
  setTargetTemperature,
  setPID,
  setDutyCycle,
} from '../../lib/settings';
import {
  toggleDewController,
  toggleDewControllerLogging,
} from '../../lib/commands';
import { getLogList } from '../../lib/information';
import { setInterval } from 'timers';

import style from './style';

export default class DewController extends Component {
  constructor() {
    super();
    this.state = {
      currentTemperature: null,
      currentlyHeating: null,
      dewControllerEnabled: null,
      targetTemperature: null,
      P: null,
      I: null,
      D: null,
      vPFormValue: null,
      desiredP: null,
      desiredI: null,
      desiredD: null,
      PIDFormDirty: false,
      loggingEnabled: null,
      info: null,
      error: null,
      logList: null,
      sensorOk: null,
      lastSliderChanged: Date.now(),
      dutyCycle: null,
    };
  }

  async componentDidMount() {
    getLogList().then((r) => {
      this.setState({ logList: r });
    });
    getDewControllerStatus()
      .then((r) => {
        const { dutyCycle } = r;

        this.setState({
          currentTemperature: r.currentTemperature,
          currentlyHeating: r.currentlyHeating,
          targetTemperature: r.targetTemperature,
          dewControllerEnabled: r.dewControllerEnabled,
          P: r.p,
          I: r.i,
          D: r.d,
          loggingEnabled: r.loggingEnabled,
          dutyCycle: dutyCycle,
          sensorOk: r.sensorOk,
        });

        const slider = this.fallbackSlider;
        if (slider) {
          slider.setValue(dutyCycle);
        }
        this.timer = setInterval(this.refreshStatus.bind(this), 2000);
      })
      .catch((e) => this.handleError(e));
  }

  handleError = (e) => {
    console.error('problem', e);
    this.setState({ error: e });
  };

  async refreshStatus() {
    getDewControllerStatus()
      .then((r) => {
        const { dutyCycle } = r;
        this.setState({
          currentTemperature: r.currentTemperature,
          currentlyHeating: r.currentlyHeating,
          P: r.p,
          I: r.i,
          D: r.d,
          loggingEnabled: r.loggingEnabled,
          dutyCycle: dutyCycle,
          sensorOk: r.sensorOk,
        });

        const slider = this.fallbackSlider;
        if (slider) {
          slider.setValue(dutyCycle);
        }
      })
      .catch((e) => this.handleError(e));
  }

  componentWillUnmount() {
    console.log('Cancelling timer');
    clearInterval(this.timer._id);
  }

  errorToast() {
    if (this.state.error != null) {
      return <p>{this.state.error.toString()}</p>;
    }
  }

  onSubmit = (e) => {
    e.preventDefault();
    this.setState({ error: null, info: null });
    const { targetTemperature } = this.state;
    setTargetTemperature(targetTemperature)
      .then((r) =>
        this.setState({
          info: 'Intervalometer Settings Updated',
          targetTemperature: r.targetTemperature,
        }),
      )
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

  errorToast() {
    if (this.state.error != null) {
      return <p>{this.state.error.toString()}</p>;
    }
  }

  infoToast() {
    if (this.state.info != null) {
      return <p>{this.state.info.toString()}</p>;
    }
  }

  onPIDSubmit() {
    const { desiredP, desiredI, desiredD } = this.state;
    console.log('Updating PID to ', desiredP, desiredI, desiredD);
    setPID(desiredP, desiredI, desiredD)
      .then(() => this.setState({ info: 'PID settings updated' }))
      .catch((e) => this.handleError(e));
  }

  PIDDialogTags() {
    return (
      <Dialog
        ref={(pidDialog) => {
          this.pidDialog = pidDialog;
        }}
      >
        <Dialog.Header>PID Values</Dialog.Header>
        <Dialog.Body>
          <form>
            <p>
              <TextField
                label="P"
                value={this.state.desiredP}
                onInput={linkState(this, 'desiredP')}
              />
            </p>
            <p>
              <TextField
                label="I"
                value={this.state.desiredI}
                onInput={linkState(this, 'desiredI')}
              />
            </p>
            <p>
              <TextField
                label="D"
                value={this.state.desiredD}
                onInput={linkState(this, 'desiredD')}
              />
            </p>
          </form>
        </Dialog.Body>
        <Dialog.Footer>
          <Dialog.FooterButton cancel={true}>Cancel</Dialog.FooterButton>
          <Dialog.FooterButton accept={true} onClick={() => this.onPIDSubmit()}>
            Update
          </Dialog.FooterButton>
        </Dialog.Footer>
      </Dialog>
    );
  }

  onLoggingEnabled = (e) => {
    const enabled = e.target.checked;
    console.log(`Logging toggled to: ${enabled ? 'enabled' : 'disabled'}`);
    this.setState({ dewControllerEnabled: enabled });
    toggleDewControllerLogging(enabled)
      .then((r) => this.setState({ dewControllerEnabled: r }))
      .catch((e) => this.handleError(e));
  };

  logListTags() {
    const { logList } = this.state;

    if (logList === null) {
      return;
    }

    if (logList === undefined) {
      return;
    }

    return (
      <List>
        {logList.map((logFile) => (
          <List.Item>
            <a href={`pid-log-viewer/?file=${logFile.escapedFilename}`}>
              {logFile.filename}
            </a>
          </List.Item>
        ))}
      </List>
    );
  }

  CSVDialogTags() {
    return (
      <Dialog
        ref={(csvDialog) => {
          this.csvDialog = csvDialog;
        }}
      >
        <Dialog.Header>Download Logs</Dialog.Header>
        <Dialog.Body>
          <p>
            Logging Enabled:{' '}
            <Switch
              id="loggingEnabled"
              onChange={this.onLoggingEnabled.bind(this)}
              checked={this.state.loggingEnabled === true}
            />
          </p>
          <p>{this.logListTags()}</p>
        </Dialog.Body>
        <Dialog.Footer>
          <Dialog.FooterButton accept={true}>Close</Dialog.FooterButton>
        </Dialog.Footer>
      </Dialog>
    );
  }

  showPIDDialog() {
    const { P, I, D } = this.state;
    this.setState({
      PIDFormDirty: false,
      desiredP: P,
      desiredI: I,
      desiredD: D,
    });
    this.pidDialog.MDComponent.show();
  }

  onFallbackSliderChange = (e) => {
    // filter out spurious events: https://github.com/material-components/material-components-web/issues/5789
    const now = Date.now();
    const { lastSliderChanged } = this.state;

    if (now - lastSliderChanged > 100 || lastSliderChanged === undefined) {
      this.setState({ lastSliderChanged: now });
      const dutyCycle = e.detail.value;
      console.log(`Setting duty cycle to ${dutyCycle}`);
      setDutyCycle(dutyCycle);
    }
  };

  FallbackControlTags() {
    const { sensorOk } = this.state;
    if (sensorOk !== true) {
      return (
        <div>
          <Slider
            ref={(fallbackSlider) => {
              this.fallbackSlider = fallbackSlider;
            }}
            step={1}
            min={0}
            max={10}
            discrete={true}
            onChange={(e) => this.onFallbackSliderChange(e)}
          />
        </div>
      );
    }
  }

  PIDControlsTags() {
    const { currentTemperature, sensorOk, targetTemperature } = this.state;

    if (sensorOk === true) {
      return (
        <div>
          <form onSubmit={this.onSubmit.bind(this)}>
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
            <p>
              <Button raised ripple onClick={(e) => e.prevent_default}>
                Update
              </Button>
            </p>
          </form>
          <p>
            <Button raised ripple onClick={() => this.showPIDDialog()}>
              PID Values
            </Button>
          </p>
        </div>
      );
    }
  }

  render({}, { currentlyHeating, dewControllerEnabled }) {
    return (
      <div class={style.main}>
        <h1>Dew Controller</h1>
        <div>
          Enabled:{' '}
          <Switch
            id="dewControllerEnable"
            onChange={this.onDewControllerEnabled.bind(this)}
            checked={dewControllerEnabled === true}
          ></Switch>
          {this.infoToast()}
          {this.errorToast()}
          <p>
            <FormField>
              <Radio
                id="heating"
                name="Basic Options"
                disabled="true"
                checked={currentlyHeating === true}
              />
              <label for="heating">Heating</label>
            </FormField>
          </p>
          {this.FallbackControlTags()}
          {this.PIDControlsTags()}
          <p>
            <Button
              raised
              ripple
              onClick={() => this.csvDialog.MDComponent.show()}
            >
              Logs
            </Button>
          </p>
          {this.PIDDialogTags()}
          {this.CSVDialogTags()}
        </div>
      </div>
    );
  }
}
