import axios from 'axios';
import config from '../../config';

async function getLogList() {
  return axios
    .get(`${config.endpoint}/backend/log_list`)
    .then((r) => r.data.files);
}

export { getLogList };
