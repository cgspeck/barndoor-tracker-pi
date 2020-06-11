import axios from "axios";
import config from "../../config";

async function getAllSettings() {
  return axios.get(`${config.endpoint}/status/debug`).then(r => r.data);
}

async function getAPSettings() {
  return axios.get(`${config.endpoint}/settings/network/ap`).then(r => r.data);
}

async function setAPSettings(ssid, key) {
  // remove when frontend actually has wireless management features
  await axios.post(`${config.endpoint}/settings/network`, {
    accessPointMode: true
  });
  return axios
    .post(`${config.endpoint}/settings/network/ap`, {
      ssid: ssid,
      key: key
    })
    .then(r => r.data);
}

async function getIntervalometerSettings() {
  return axios.get(`${config.endpoint}/settings/intervalometer`).then(r => r.data);
}

async function setIntervalometerSettings(bulbInterval, restInterval) {
  return axios
    .post(`${config.endpoint}/settings/intervalometer`, {
      bulbInterval,
      restInterval
    })
    .then(r => r.data);
}

async function getLocationSettings() {
  return axios.get(`${config.endpoint}/settings/location`).then(r => r.data);
}

async function setLocationSettings(
  latitude,
  magDeclination,
  azError,
  altError,
  xOffset,
  yOffset,
  zOffset
) {
  return axios
    .post(
      `${config.endpoint}/settings/location`,
      {
        latitude,
        magDeclination,
        azError,
        altError,
        xOffset,
        yOffset,
        zOffset
      },
      {
        transformRequest: [
          function(data, _) {
            const res = {
              latitude: parseFloat(data.latitude),
              magDeclination: parseFloat(data.magDeclination),
              azError: parseFloat(data.azError),
              altError: parseFloat(data.altError),
              xOffset: parseInt(data.xOffset),
              yOffset: parseInt(data.yOffset),
              zOffset: parseInt(data.zOffset)
            };
            return JSON.stringify(res);
          }
        ]
      }
    )
    .then(r => r.data);
}

async function getFlags() {
  return axios.get(`${config.endpoint}/status/flags`).then(r => r.data);
}

async function getAlignStatus() {
  return axios.get(`${config.endpoint}/status/align`).then(r => r.data);
}

async function getInitialTrackStatus() {
  return axios.get(`${config.endpoint}/track/all`).then(r => r.data);
}

async function getTrackState() {
  return axios.get(`${config.endpoint}/track`).then(r => r.data);
}

async function getDewControllerStatus() {
  return axios.get(`${config.endpoint}/status/dew_controller`).then(r => r.data);
}

async function setTargetTemperature(targetTemp) {
  return axios.post(
    `${config.endpoint}/settings/dew_controller`,
    {
      targetTemp
    }
    ).then(r => r.data);
}

export {
  getAllSettings,
  getAlignStatus,
  getAPSettings,
  getFlags,
  getLocationSettings,
  getInitialTrackStatus,
  getTrackState,
  setAPSettings,
  setLocationSettings,
  getIntervalometerSettings,
  setIntervalometerSettings,
  getDewControllerStatus,
  setTargetTemperature
};
