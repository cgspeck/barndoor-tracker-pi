import axios from "axios";
import config from "../../config";

async function getAllSettings() {
  return axios.get(`${config.endpoint}/settings/debug`).then(r => r.data);
}

async function getAPSettings() {
  return axios.get(`${config.endpoint}/settings/ap`).then(r => r.data);
}

async function setAPSettings(ssid, key) {
  return axios
    .post(`${config.endpoint}/settings/ap`, {
      ssid: ssid,
      key: key
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
    .post(`${config.endpoint}/settings/location`, {
      latitude,
      magDeclination,
      azError,
      altError,
      xOffset,
      yOffset,
      zOffset
    })
    .then(r => r.data);
}

async function getFlags() {
  return axios.get(`${config.endpoint}/flags`).then(r => r.data);
}

async function getAlignStatus() {
  return axios.get(`${config.endpoint}/status/align`).then(r => r.data);
}

export {
  getAllSettings,
  getAlignStatus,
  getAPSettings,
  getFlags,
  getLocationSettings,
  setAPSettings,
  setLocationSettings
};
