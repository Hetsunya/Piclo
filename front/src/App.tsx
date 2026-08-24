import { useState, useRef, useEffect, type DragEvent } from 'react'
import { BrowserRouter, Routes, Route, useParams, Link } from 'react-router-dom'
import { useTranslation, Trans } from 'react-i18next'
import './index.css'
import './i18n/i18n'
import LanguageSwitcher from './components/LanguageSwitcher'

function UploadPage() {
  const { t } = useTranslation()
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
    <div className="page-wrapper">
      <LanguageSwitcher />
      
      <div className="container">
        <h1>{t('title')}</h1>
        <p className="subtitle">{t('subtitle')}</p>

        <div
          className="drop-zone"
          onClick={() => fileInputRef.current?.click()}
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
        >
          <div className="dropzone-content">
            <div className="dropzone-icon">📁</div>
            <p>{t('dropzone')}</p>
          </div>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/png, image/jpeg, image/gif, image/webp"
            onChange={(e) => handleFile(e.target.files?.[0])}
          />
        </div>

        {file && (
          <div className="file-preview">
            <div className="file-info">
              <span className="file-name">{file.name}</span>
              <span className="file-size">{(file.size / 1024 / 1024).toFixed(2)} {t('fileSize')}</span>
            </div>
            <div className="file-thumbnail">
              <img src={URL.createObjectURL(file)} alt="Preview" />
            </div>
          </div>
        )}

        <button className="btn" onClick={upload} disabled={!file || loading}>
          {loading ? t('uploading') : t('uploadBtn')}
        </button>

        {error && <div className="error">❌ {t('error')}: {error}</div>}

        <div className="terms-notice">
          <Trans i18nKey="termsAgreement">
            <Link to="/terms" target="_blank" rel="noopener noreferrer" />
            <Link to="/privacy" target="_blank" rel="noopener noreferrer" />
          </Trans>
        </div>
      </div>
    </div>
  )
}

function ImageViewer() {
  const { id } = useParams<{ id: string }>()
  const { t } = useTranslation()
  const [imageUrl, setImageUrl] = useState<string>('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string>('')

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
        <div className="spinner">{t('uploading')}</div>
      </div>
    )
  }

  if (error === 'not_found') {
    return (
      <div className="container" style={{textAlign: 'center', padding: '60px 20px'}}>
        <h1>🖼️ 404</h1>
        <p style={{color: '#666', marginBottom: '24px'}}>
          {t('notFound')}
        </p>
        <Link to="/" className="btn" style={{display: 'inline-block', width: 'auto', padding: '12px 32px', textDecoration: 'none'}}>
          {t('uploadNew')}
        </Link>
      </div>
    )
  }

  return (
    <div className="image-viewer">
      <div className="image-header">
        <Link to="/" className="logo">{t('logo')}</Link>
        <div className="header-actions">
          <button 
            className="btn-copy" 
            onClick={() => {
              navigator.clipboard.writeText(window.location.href)
            }}
          >
            {t('copyLink')}
          </button>
          <LanguageSwitcher />
        </div>
      </div>
      <div className="image-container">
        <img src={imageUrl} alt="Uploaded" />
      </div>
      <div className="image-footer">
        <p style={{color: '#999', fontSize: '14px'}}>
          {t('imageUrl')} <br/>
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
        <Route path="/terms" element={<TermsPage />} />
        <Route path="/privacy" element={<PrivacyPage />} />
      </Routes>
    </BrowserRouter>
  )
}

function TermsPage() {
  const { t } = useTranslation()
  
  return (
    <div className="legal-page">
      <div className="container legal-container">
        <Link to="/" className="back-link">← {t('uploadNew')}</Link>
        <h1>{t('termsTitle')}</h1>
        <p className="last-updated">{t('lastUpdated')}</p>
        
        <p>{t('terms.intro')}</p>
        
        <h2>{t('terms.section1Title')}</h2>
        <p>{t('terms.section1Content')}</p>
        
        <h2>{t('terms.section2Title')}</h2>
        <p>{t('terms.section2Intro')}</p>
        <ul>
          {(t('terms.section2Points', { returnObjects: true }) as string[]).map((point: string, index: number) => (
            <li key={index}>{point}</li>
          ))}
        </ul>
        <p>{t('terms.section2Prohibited')}</p>
        <ul>
          {(t('terms.section2ProhibitedPoints', { returnObjects: true }) as string[]).map((point: string, index: number) => (
            <li key={index}>{point}</li>
          ))}
        </ul>
        
        <h2>{t('terms.section3Title')}</h2>
        <p>{t('terms.section3Content')}</p>
        
        <h2>{t('terms.section4Title')}</h2>
        <p>{t('terms.section4Content')}</p>
        
        <h2>{t('terms.section5Title')}</h2>
        <p>{t('terms.section5Content')}</p>
        
        <h2>{t('terms.section6Title')}</h2>
        <p>{t('terms.section6Content')}</p>
      </div>
    </div>
  )
}

function PrivacyPage() {
  const { t } = useTranslation()
  
  return (
    <div className="legal-page">
      <div className="container legal-container">
        <Link to="/" className="back-link">← {t('uploadNew')}</Link>
        <h1>{t('privacyTitle')}</h1>
        <p className="last-updated">{t('lastUpdated')}</p>
        
        <p>{t('privacy.intro')}</p>
        
        <h2>{t('privacy.section1Title')}</h2>
        <p>{t('privacy.section1Content')}</p>
        <ul>
          {(t('privacy.section1Points', { returnObjects: true }) as string[]).map((point: string, index: number) => (
            <li key={index}>{point}</li>
          ))}
        </ul>
        <p>{t('privacy.section1Future')}</p>
        
        <h2>{t('privacy.section2Title')}</h2>
        <p>{t('privacy.section2Content')}</p>
        <ul>
          {(t('privacy.section2Points', { returnObjects: true }) as string[]).map((point: string, index: number) => (
            <li key={index}>{point}</li>
          ))}
        </ul>
        
        <h2>{t('privacy.section3Title')}</h2>
        <p>{t('privacy.section3Content')}</p>
        
        <h2>{t('privacy.section4Title')}</h2>
        <p>{t('privacy.section4Content')}</p>
        
        <h2>{t('privacy.section5Title')}</h2>
        <p>{t('privacy.section5Content')}</p>
        
        <h2>{t('privacy.section6Title')}</h2>
        <p>{t('privacy.section6Content')}</p>
        
        <h2>{t('privacy.section7Title')}</h2>
        <p>{t('privacy.section7Content')}</p>
      </div>
    </div>
  )
}

export default App
