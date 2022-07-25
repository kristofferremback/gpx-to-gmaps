import html from '../../lib/html.js'
import { useCallback, useState, useEffect, useMemo } from '../../deps/preact/hooks.js'

import useFetch, { fetchStates } from '../../lib/hooks/use-fetch.js'
import Modal from '../Modal/Modal.component.js'

import { vehicleTypes } from './vehicle-types.js'
import GPXConverter from './GPXConverter.component.js'

const GPXConverterContainer = () => {
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
    <${GPXConverter}
      onSubmit=${onSubmit}
      requestData=${requestData}
      setRequestData=${setRequestData}
      submitAllowed=${submitAllowed}
      isLoading=${state === fetchStates.LOADING}
      googleMapsUrls=${resp?.response?.google_maps_urls ?? null}
      mapImageUrls=${resp?.response?.maps_urls ?? null}
    />
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

export default GPXConverterContainer
