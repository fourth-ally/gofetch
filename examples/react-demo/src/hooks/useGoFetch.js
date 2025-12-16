import { useState, useEffect } from 'react'

export default function useGoFetch() {
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState(null)
  const [gofetch, setGoFetch] = useState(null)
  const [logs, setLogs] = useState([])

  useEffect(() => {
    const loadWasm = async () => {
      try {
        // Load wasm_exec.js
        const wasmExecScript = document.createElement('script')
        wasmExecScript.src = '/wasm_exec.js'
        
        await new Promise((resolve, reject) => {
          wasmExecScript.onload = resolve
          wasmExecScript.onerror = reject
          document.head.appendChild(wasmExecScript)
        })

        // Initialize Go WASM
        const go = new window.Go()
        const result = await WebAssembly.instantiateStreaming(
          fetch('/gofetch.wasm'),
          go.importObject
        )
        
        // Run the WASM module
        go.run(result.instance)

        // Check if gofetch is available
        if (!window.gofetch) {
          throw new Error('gofetch not found in window object')
        }

        // Configure the default client
        window.gofetch.setBaseURL('https://jsonplaceholder.typicode.com')
        window.gofetch.setTimeout(10000) // 10 seconds

        // Wrap gofetch methods to add logging
        const wrappedGoFetch = createWrappedClient(window.gofetch)

        setGoFetch(wrappedGoFetch)
        addLog('success', 'GoFetch WASM module loaded successfully!')
        setIsLoading(false)
      } catch (err) {
        console.error('Failed to load GoFetch WASM:', err)
        setError(err.message)
        setIsLoading(false)
      }
    }

    loadWasm()
  }, [])

  const addLog = (type, message, data = null) => {
    const timestamp = new Date().toLocaleTimeString()
    setLogs(prev => [...prev, { timestamp, type, message, data }])
  }

  const createWrappedClient = (client) => {
    return {
      get: async (path, params = null) => {
        addLog('info', `GET ${path}`, params)
        try {
          const response = await client.get(path, params)
          addLog('success', `✓ GET ${path} - Status: ${response.statusCode}`, response.data)
          return response
        } catch (err) {
          addLog('error', `✗ GET ${path} failed: ${err}`)
          throw err
        }
      },
      
      post: async (path, params = null, body = null) => {
        addLog('info', `POST ${path}`, body)
        try {
          const response = await client.post(path, params, body)
          addLog('success', `✓ POST ${path} - Status: ${response.statusCode}`, response.data)
          return response
        } catch (err) {
          addLog('error', `✗ POST ${path} failed: ${err}`)
          throw err
        }
      },
      
      put: async (path, params = null, body = null) => {
        addLog('info', `PUT ${path}`, body)
        try {
          const response = await client.put(path, params, body)
          addLog('success', `✓ PUT ${path} - Status: ${response.statusCode}`, response.data)
          return response
        } catch (err) {
          addLog('error', `✗ PUT ${path} failed: ${err}`)
          throw err
        }
      },
      
      patch: async (path, params = null, body = null) => {
        addLog('info', `PATCH ${path}`, body)
        try {
          const response = await client.patch(path, params, body)
          addLog('success', `✓ PATCH ${path} - Status: ${response.statusCode}`, response.data)
          return response
        } catch (err) {
          addLog('error', `✗ PATCH ${path} failed: ${err}`)
          throw err
        }
      },
      
      delete: async (path, params = null) => {
        addLog('info', `DELETE ${path}`, params)
        try {
          const response = await client.delete(path, params)
          addLog('success', `✓ DELETE ${path} - Status: ${response.statusCode}`)
          return response
        } catch (err) {
          addLog('error', `✗ DELETE ${path} failed: ${err}`)
          throw err
        }
      },
      
      setBaseURL: (url) => client.setBaseURL(url),
      setTimeout: (timeout) => client.setTimeout(timeout),
      setHeader: (key, value) => client.setHeader(key, value),
      newClient: () => createWrappedClient(client.newClient())
    }
  }

  const clearLogs = () => {
    setLogs([])
  }

  return {
    isLoading,
    error,
    gofetch,
    logs,
    clearLogs,
    addLog
  }
}
