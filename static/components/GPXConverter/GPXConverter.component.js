import getUseFetch, { fetchStates } from '../../hooks/use-fetch.js'
import getModal from '../Modal/Modal.component.js'

const vehicleTypes = [
  { value: 'bike', name: 'Bike' },
  { value: 'car', name: 'Car' },
  { value: 'walking', name: 'Walking' },
]

/**
 * @param {any} html
 * @param {import('preact/hooks')} hooks
 */
const getGPXConverter = (html, hooks) => {
  const { useCallback, useState, useEffect, useMemo } = hooks
  const useFetch = getUseFetch(hooks)
  const Modal = getModal(html)

  const GPXConverter = () => {
    const [dispatchFetch, state, resp] = useFetch()
    const [requestData, setRequestData] = useState({
      vehicle_type: vehicleTypes[0].value,
      max_precision: '25',
      gpx: null,
    })

    const [errIsOpen, setErrIsOpen] = useState(false)
    const closeModal = useCallback(() => setErrIsOpen(false), [setErrIsOpen])
    useEffect(() => {
      if (state === fetchStates.ERROR) {
        setErrIsOpen(true)
      }
      return () => {
        closeModal
      }
    }, [state, closeModal])

    const submitAllowed = useMemo(() => {
      return [
        requestData.gpx != null,
        !isNaN(parseInt(requestData.max_precision)),
        requestData.vehicle_type != '',
        state != fetchStates.LOADING,
      ].every(Boolean)
    }, [requestData, state])

    const onChange = useCallback(
      (e) => {
        switch (e.target.name) {
          case 'vehicle_type':
            return setRequestData({ ...requestData, vehicle_type: e.target.value })
          case 'max_precision':
            let value = parseInt(e.target.value)
            const [min, max] = [e.target.min, e.target.max].map((v) => parseInt(v))
            if (value > max) {
              value = max
            } else if (value < min) {
              value = min
            }

            return setRequestData({ ...requestData, max_precision: value.toString() })
          case 'gpx_file':
            return setRequestData({ ...requestData, gpx: e.target.files[0] })
          default:
            break
        }
      },
      [requestData, setRequestData]
    )

    const onSubmit = useCallback(
      async (e) => {
        e.preventDefault()
        e.stopPropagation()

        const data = new FormData()
        for (const [key, value] of Object.entries(requestData)) {
          data.append(key, value)
        }

        await dispatchFetch(
          '/api/convert-gpx',
          { method: 'POST', body: data },
          { validateStatus: (status) => status === 200 }
        )
      },
      [dispatchFetch, requestData]
    )

    return html`
    <div class="gpx-converter">

      <link rel="stylesheet" href="/components/GPXConverter/GPXConverter.styles.css" />
      <div class="grid">
        <aside>
          <form onsubmit=${onSubmit}>
            <label for="vehicle_type">Vehicle type</label>
            <select
              name="vehicle_type"
              value=${requestData.vehicle_type}
              id="vehicle_type"
              onchange=${onChange}
            >
              ${vehicleTypes.map(({ name, value }) => html`<option value=${value}>${name}</option>`)}
            </select>
            <label for="max_precision">Max precision</label>
            <input
              type="number"
              min="3"
              max="30"
              name="max_precision"
              value=${requestData.max_precision}
              onchange=${onChange}
            />

            <label for="gpx_file">Select .gpx file</label>
            <input type="file" name="gpx_file" accept=".gpx" onchange=${onChange} />
            <button type="submit" disabled=${!submitAllowed} aria-busy=${state === fetchStates.LOADING}>
              Upload
            </button>
          </form>
        </aside>
        ${
          fetchStates.IDLE && resp.response != null
            ? html`
                ${resp.response.google_maps_urls.map(
                  (url, i) => html`
                    <article class="preview-map">
                      <img src=${resp.response.maps_urls[i]} />
                      <h3><a href=${url} target="_blank">Google Maps directions here</a></h3>
                    </article>
                  `
                )}
              `
            : null
        }
      </div>
    </div>
    <${Modal} isOpen=${errIsOpen} title="Something went wrong" close=${closeModal}>
    ${
      errIsOpen
        ? html`
            <details>
              <summary>An error occurred when converting the .gpx file</summary>
              <pre>${resp.error != null ? resp.error.message : 'unknown error'}</pre>
            </details>
          `
        : null
    }
    </${Modal}>
    `
  }

  return GPXConverter
}

export default getGPXConverter
