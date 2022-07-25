import { useState, useCallback } from '../deps/preact/hooks.js'

export const fetchStates = {
  IDLE: 'IDLE',
  LOADING: 'LOADING',
  ERROR: 'ERROR',
}

export class RequestError extends Error {
  /**
   * @param {string} message
   * @param {Response} response
   */
  constructor(message, response) {
    super(`${message}. Status ${response.status}, ${response.statusText}`)
    this.name = RequestError.name
    this.response = response
  }
}

const useFetch = () => {
  const [state, setState] = useState(fetchStates.IDLE)
  const [results, setResults] = useState({ response: null, error: null })

  const update = useCallback(
    (state, response, error) => {
      setState(state)
      setResults({ error, response })
    },
    [setState, setResults]
  )

  const dispatchFetch = useCallback(
    /**
     * @param {RequestInfo | URL} input
     * @param {RequestInit} [init]
     * @param {Object} [opts]
     * @param {(status: number) => boolean} [opts.validateStatus]
     */
    async (input, init, { validateStatus = (_) => true } = {}) => {
      setState(fetchStates.LOADING)
      try {
        const resp = await fetch(input, init)
        if (!validateStatus(resp.status)) {
          throw new RequestError(`Status error`, resp)
        }

        const respJson = await resp.json()
        update(fetchStates.IDLE, respJson, null)
      } catch (error) {
        update(fetchStates.ERROR, null, error)
        console.error('useFetch:dispatchFetch', { input, init, error })
      }
    },
    []
  )

  return [dispatchFetch, state, results]
}

export default useFetch
