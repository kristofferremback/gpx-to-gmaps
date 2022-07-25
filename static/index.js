import html from './lib/html.js'
import { render } from './deps/preact.js'

import GPXConverter from './components/GPXConverter/index.js'

window.addEventListener('load', () => {
  render(html`<${GPXConverter} />`, document.querySelector('#gpx-converter'))
})
