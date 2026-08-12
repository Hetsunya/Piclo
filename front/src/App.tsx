import { useState, useRef, type DragEvent } from 'react'
import './index.css'

interface UploadResponse {
  id: string
  url: string
}

function App() {
  const [ttl, setTtl] = useState<number>(1) // по умолчанию 1 час
  const [file, setFile] = useState<File | null>(null)
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<UploadResponse | null>(null)
  const [error, setError] = useState<string>('')
  const [copied, setCopied] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFile = (selected: File | undefined) => {
    if (!selected) return
    setFile(selected)
    setResult(null)
    setError('')
  }

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.currentTarget.classList.remove('active')
    handleFile(e.dataTransfer.files[0])
  }

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.currentTarget.classList.add('active')
  }

  const handleDragLeave = (e: DragEvent<HTMLDivElement>) => {
    e.currentTarget.classList.remove('active')
  }

  const upload = async () => {
    if (!file) return
    setLoading(true)
    setError('')
    try {
      const formData = new FormData()
      formData.append('file', file)

      const res = await fetch('/api/v1/upload', {
        method: 'POST',
        body: formData,
      })

      const data = await res.json()

      if (!res.ok) {
        throw new Error(data.error || 'Upload failed')
      }

      setResult(data)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  const copyUrl = async () => {
    if (!result) return
    await navigator.clipboard.writeText(result.url)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="container">
      <h1>📸 Piclo Upload</h1>

      <div
        className="drop-zone"
        onClick={() => fileInputRef.current?.click()}
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
      >
        <p>Нажми или перетащи картинку (max 10MB)</p>
        <input
          ref={fileInputRef}
          type="file"
          accept="image/png, image/jpeg, image/gif, image/webp"
          onChange={(e) => handleFile(e.target.files?.[0])}
        />
      </div>

      {file && (
        <p className="file-info">
          {file.name} — {(file.size / 1024 / 1024).toFixed(2)} MB
        </p>
      )}
      <button className="btn" onClick={upload} disabled={!file || loading}>
        {loading ? 'Загрузка...' : 'Загрузить'}
      </button>

      {error && <div className="error">❌ {error}</div>}

      {result && (
        <div className="result">
          <p>✅ Успешно!</p>
          <img src={result.url} alt="Preview" />
          <div className="url-box">
            <input type="text" readOnly value={result.url} />
            <button onClick={copyUrl}>{copied ? 'Скопировано!' : 'Копировать'}</button>
          </div>
        </div>
      )}
    </div>
  )
}

export default App