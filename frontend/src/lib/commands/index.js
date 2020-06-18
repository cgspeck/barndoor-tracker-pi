import axios from 'axios';
import config from '../../config';

async function toggleIgnoreAz(enabled) {
  return axios
    .post(`${config.endpoint}/backend/toggle/ignoreAz`, {
      enabled,
    })
    .then((r) => {
      return r.data.enabled;
    });
}

async function toggleIgnoreAlt(enabled) {
  return axios
    .post(`${config.endpoint}/backend/toggle/ignoreAlt`, {
      enabled,
    })
    .then((r) => {
      return r.data.enabled;
    });
}

async function toggleIntervalometer(enabled) {
  return axios
    .post(`${config.endpoint}/backend/toggle/intervalometer`, {
      enabled,
    })
    .then((r) => {
      return r.data.enabled;
    });
}

async function toggleDewController(enabled) {
  return axios
    .post(`${config.endpoint}/backend/toggle/dewcontroller`, {
      enabled,
    })
    .then((r) => {
      return r.data.enabled;
    });
}

async function toggleDewControllerLogging(enabled) {
  return axios
    .post(`${config.endpoint}/backend/toggle/dewcontroller/logging`, {
      enabled,
    })
    .then((r) => {
      return r.data.enabled;
    });
}

async function startHoming() {
  return axios
    .post(`${config.endpoint}/backend/track`, {
      command: 'home',
    })
    .then((r) => r.data.state);
}

async function startTracking() {
  return axios
    .post(`${config.endpoint}/backend/track`, {
      command: 'track',
    })
    .then((r) => r.data.state);
}

async function stopTracking() {
  return axios
    .post(`${config.endpoint}/backend/track`, {
      command: 'stop',
    })
    .then((r) => r.data.state);
}

export {
  toggleIgnoreAz,
  toggleIgnoreAlt,
  toggleIntervalometer,
  toggleDewController,
  toggleDewControllerLogging,
  startHoming,
  startTracking,
  stopTracking,
};
