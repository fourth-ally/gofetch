import { useState } from 'react'
import useGoFetch from './hooks/useGoFetch'

function App() {
  const { isLoading, error, gofetch, logs, clearLogs } = useGoFetch()
  const [stats, setStats] = useState({ requests: 0, successes: 0, errors: 0 })

  const updateStats = (success) => {
    setStats(prev => ({
      requests: prev.requests + 1,
      successes: success ? prev.successes + 1 : prev.successes,
      errors: success ? prev.errors : prev.errors + 1
    }))
  }

  const handleGetUsers = async () => {
    try {
      const response = await gofetch.get('/users')
      console.log('Users:', response.data)
      updateStats(true)
    } catch (err) {
      console.error('Error fetching users:', err)
      updateStats(false)
    }
  }

  const handleGetUser = async () => {
    try {
      const response = await gofetch.get('/users/1')
      console.log('User:', response.data)
      updateStats(true)
    } catch (err) {
      console.error('Error fetching user:', err)
      updateStats(false)
    }
  }

  const handleCreatePost = async () => {
    try {
      const newPost = {
        userId: 1,
        title: 'GoFetch React Demo',
        body: 'This post was created using GoFetch in a React application!'
      }
      
      const response = await gofetch.post('/posts', null, newPost)
      console.log('Created post:', response.data)
      updateStats(true)
    } catch (err) {
      console.error('Error creating post:', err)
      updateStats(false)
    }
  }

  const handleGetPosts = async () => {
    try {
      const response = await gofetch.get('/posts', { userId: 1 })
      console.log('Posts:', response.data)
      updateStats(true)
    } catch (err) {
      console.error('Error fetching posts:', err)
      updateStats(false)
    }
  }

  const handleUpdatePost = async () => {
    try {
      const updatedPost = {
        userId: 1,
        id: 1,
        title: 'Updated Title',
        body: 'Updated content from GoFetch!'
      }
      
      const response = await gofetch.put('/posts/1', null, updatedPost)
      console.log('Updated post:', response.data)
      updateStats(true)
    } catch (err) {
      console.error('Error updating post:', err)
      updateStats(false)
    }
  }

  const handleDeletePost = async () => {
    try {
      const response = await gofetch.delete('/posts/1')
      console.log('Deleted post:', response)
      updateStats(true)
    } catch (err) {
      console.error('Error deleting post:', err)
      updateStats(false)
    }
  }

  const handleWithCustomClient = async () => {
    try {
      // Create a new client instance with custom configuration
      const customClient = gofetch.newClient()
      customClient.setBaseURL('https://jsonplaceholder.typicode.com')
      customClient.setTimeout(5000)
      customClient.setHeader('X-Custom-Header', 'React-Demo')
      
      const response = await customClient.get('/todos/1')
      console.log('Custom client response:', response.data)
      updateStats(true)
    } catch (err) {
      console.error('Error with custom client:', err)
      updateStats(false)
    }
  }

  if (isLoading) {
    return (
      <div className="container">
        <div className="loading">
          <h2>🔄 Loading GoFetch WASM Module...</h2>
          <p>Please wait while we initialize the WebAssembly module.</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="container">
        <div className="error">
          <h2>❌ Failed to Load GoFetch</h2>
          <p>{error}</p>
          <p>Make sure you've built the WASM module by running: <code>npm run build:wasm</code></p>
        </div>
      </div>
    )
  }

  return (
    <div className="container">
      <h1>🚀 GoFetch React Demo</h1>
      <p className="subtitle">
        A demonstration of GoFetch HTTP client running in React via WebAssembly
      </p>

      <div className="stats">
        <div className="stat-card">
          <div className="stat-value">{stats.requests}</div>
          <div className="stat-label">Total Requests</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.successes}</div>
          <div className="stat-label">Successful</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.errors}</div>
          <div className="stat-label">Errors</div>
        </div>
      </div>

      <div className="controls">
        <button onClick={handleGetUsers}>Get All Users</button>
        <button onClick={handleGetUser}>Get User by ID</button>
        <button onClick={handleGetPosts}>Get User Posts</button>
        <button onClick={handleCreatePost}>Create Post</button>
        <button onClick={handleUpdatePost}>Update Post</button>
        <button onClick={handleDeletePost}>Delete Post</button>
        <button onClick={handleWithCustomClient}>Custom Client</button>
        <button onClick={clearLogs} style={{ marginLeft: 'auto' }}>Clear Logs</button>
      </div>

      <div className="output">
        <h3>Request Logs:</h3>
        {logs.length === 0 ? (
          <p style={{ color: '#666', fontStyle: 'italic' }}>
            No requests yet. Click a button above to make a request.
          </p>
        ) : (
          logs.map((log, index) => (
            <div key={index} className="log-entry">
              <span className="log-timestamp">{log.timestamp}</span>
              <span className={`log-${log.type}`}>{log.message}</span>
              {log.data && (
                <pre>{JSON.stringify(log.data, null, 2)}</pre>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  )
}

export default App
