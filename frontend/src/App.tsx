import { useState } from 'react'
import './App.css'

function App() {
  const [status, setStatus] = useState<string>('')
  const [loading, setLoading] = useState(false)

  const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080').replace(/\/+$/, '')

  const handlePing = async () => {
    setLoading(true)
    setStatus('Pinging...')
    try {
      const res = await fetch(`${apiBaseUrl}/api/ping`)
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      const data = await res.json()
      setStatus(`OK: ${JSON.stringify(data)}`)
    } catch (err: unknown) {
      setStatus(`Error: ${err instanceof Error ? err.message : 'failed to connect'}`)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="form-card">
      <h1 className="form-title">JobShop Ping</h1>
      <p className="form-desc">{apiBaseUrl}/api/ping</p>
      <button className="ping-btn" onClick={handlePing} disabled={loading}>
        {loading ? 'Sending...' : 'Ping'}
      </button>
      {status && <div className="form-result">{status}</div>}
    </div>
  )
}

export default App
