// @ts-check

// TODO: move to global lib
// eslint-disable-next-line import/no-unresolved
import { locale } from 'gon'
import * as routes from '../../routes.js'

// NOTE: for en using default path
if (locale !== 'en') {
  routes.configure({ default_url_options: { locale } })
}

export default routes
