import htm from 'https://cdn.skypack.dev/htm'
import * as preact from 'https://cdn.skypack.dev/preact@10.5.15'
import * as hooks from 'https://cdn.skypack.dev/preact@10.5.15/hooks'
import getGPXConverter from './components/GPXConverter/index.js'

window.addEventListener('load', () => {
  const html = htm.bind(preact.h)
  const GPXConverter = getGPXConverter(html, hooks)

  preact.render(html`<${GPXConverter} />`, document.querySelector('#gpx-converter'))
})
