import axios from 'axios';

const hexletAxios = axios.create({
  withCredentials: true,
});

export default hexletAxios;
