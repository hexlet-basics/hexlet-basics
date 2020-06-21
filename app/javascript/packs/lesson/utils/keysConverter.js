import humps from 'humps';

export const camelize = (obj) => humps.camelizeKeys(obj);
export const decamelize = (obj) => humps.decamelizeKeys(obj);
