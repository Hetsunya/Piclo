import { useState, useRef, useEffect, type DragEvent } from 'react'
import { BrowserRouter, Routes, Route, useParams, Link } from 'react-router-dom'
import './index.css'


function UploadPage() {
  const [file, setFile] = useState<File | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string>('')
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFile = (selected: File | undefined) => {
    if (!selected) return
    setFile(selected)
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

      window.location.href = `/image/${data.id}`
      
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="container">
      <h1>📸 Piclo Upload</h1>
      <p style={{textAlign: 'center', color: '#666', marginBottom: '24px'}}>
        Быстрый хостинг изображений. Максимум 10MB.
      </p>

      <div
        className="drop-zone"
        onClick={() => fileInputRef.current?.click()}
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
      >
        <p>Нажми или перетащи картинку</p>
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
    </div>
  )
}

function ImageViewer() {
  const { id } = useParams<{ id: string }>()
  const [imageUrl, setImageUrl] = useState<string>('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string>('')

  // ИСПРАВЛЕНО: useState заменен на useEffect
  useEffect(() => {
    const img = new Image()
    img.onload = () => {
      setImageUrl(`/api/v1/image/${id}/raw`)
      setLoading(false)
    }
    img.onerror = () => {
      setError('not_found')
      setLoading(false)
    }
    img.src = `/api/v1/image/${id}/raw`
  }, [id])

  if (loading) {
    return (
      <div className="container" style={{textAlign: 'center', padding: '60px 20px'}}>
        <div className="spinner">Загрузка...</div>
      </div>
    )
  }

  if (error === 'not_found') {
    return (
      <div className="container" style={{textAlign: 'center', padding: '60px 20px'}}>
        <h1>🖼️ 404</h1>
        <p style={{color: '#666', marginBottom: '24px'}}>
          Содержимое не найдено или срок его жизни истек.
        </p>
        <Link to="/" className="btn" style={{display: 'inline-block', width: 'auto', padding: '12px 32px', textDecoration: 'none'}}>
          Загрузить новое изображение
        </Link>
      </div>
    )
  }

  return (
    <div className="image-viewer">
      <div className="image-header">
        <Link to="/" className="logo">📸 Piclo</Link>
        <div className="image-actions">
          <button 
            className="btn-copy" 
            onClick={() => {
              navigator.clipboard.writeText(window.location.href)
              // Можно заменить alert на красивый тост, но для MVP alert ок
            }}
          >
            Копировать ссылку
          </button>
        </div>
      </div>
      <div className="image-container">
        <img src={imageUrl} alt="Uploaded" />
      </div>
      <div className="image-footer">
        <p style={{color: '#666', fontSize: '14px'}}>
          Изображение доступно по ссылке: <br/>
          <code>{window.location.href}</code>
        </p>
      </div>
    </div>
  )
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<UploadPage />} />
        <Route path="/image/:id" element={<ImageViewer />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App