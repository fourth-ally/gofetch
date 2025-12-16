import { useState, useEffect } from 'react'
import gofetch from 'gofetch-wasm'

export default function useGoFetch() {
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState(null)
  const [client, setClient] = useState(null)
  const [logs, setLogs] = useState([])

  const addLog = (type, message, data = null) => {
    const timestamp = new Date().toLocaleTimeString()
    setLogs(prev => [...prev, { timestamp, type, message, data }])
  }

  const createWrappedClient = (gofetchClient) => {
    return {
      get: async (path, params = null) => {
        addLog('info', `GET ${path}`, params)
        try {
          const response = await gofetchClient.get(path, params)
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
          const response = await gofetchClient.post(path, params, body)
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
          const response = await gofetchClient.put(path, params, body)
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
          const response = await gofetchClient.patch(path, params, body)
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
          const response = await gofetchClient.delete(path, params)
          addLog('success', `✓ DELETE ${path} - Status: ${response.statusCode}`)
          return response
        } catch (err) {
          addLog('error', `✗ DELETE ${path} failed: ${err}`)
          throw err
        }
      },
      
      setBaseURL: (url) => gofetchClient.setBaseURL(url),
      setTimeout: (timeout) => gofetchClient.setTimeout(timeout),
      setHeader: (key, value) => gofetchClient.setHeader(key, value),
      newClient: async () => {
        const newClient = await gofetch.newClient()
        return createWrappedClient(newClient)
      }
    }
  }

  useEffect(() => {
    const initGoFetch = async () => {
      try {
        console.log('useGoFetch: Starting initialization...')
        
        // Create a new client from the npm package
        const gofetchClient = await gofetch.newClient()
        
        console.log('useGoFetch: Client created successfully')
        
        // Configure the default client
        gofetchClient.setBaseURL('https://jsonplaceholder.typicode.com')
        gofetchClient.setTimeout(10000) // 10 seconds

        // Wrap gofetch methods to add logging
        const wrappedClient = createWrappedClient(gofetchClient)

        setClient(wrappedClient)
        addLog('success', 'GoFetch loaded successfully from NPM!')
        setIsLoading(false)
      } catch (err) {
        console.error('Failed to load GoFetch:', err)
        setError(err.message)
        setIsLoading(false)
      }
    }

    initGoFetch()
  }, [])

  const clearLogs = () => {
    setLogs([])
  }

  return {
    isLoading,
    error,
    gofetch: client,
    logs,
    clearLogs,
    addLog
  }
}
