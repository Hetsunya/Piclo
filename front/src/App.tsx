import { useState, useRef, useEffect, type DragEvent } from 'react'
import { BrowserRouter, Routes, Route, useParams, Link } from 'react-router-dom'
import './index.css'

// Типы для локализации
type Lang = 'en' | 'ru'

interface Translations {
  [key: string]: {
    en: string
    ru: string
  }
}

const translations: Translations = {
  title: { en: '📸 Piclo Upload', ru: '📸 Piclo Upload' },
  subtitle: { en: 'Fast image hosting. Max 10MB.', ru: 'Быстрый хостинг изображений. Максимум 10MB.' },
  dropzone: { en: 'Click or drag image here', ru: 'Нажми или перетащи картинку' },
  fileSize: { en: 'MB', ru: 'МБ' },
  uploadBtn: { en: 'Upload', ru: 'Загрузить' },
  uploading: { en: 'Uploading...', ru: 'Загрузка...' },
  error: { en: 'Error', ru: 'Ошибка' },
  copyLink: { en: 'Copy Link', ru: 'Копировать ссылку' },
  notFound: { en: 'Content not found or expired', ru: 'Содержимое не найдено или срок его жизни истек' },
  uploadNew: { en: 'Upload new image', ru: 'Загрузить новое изображение' },
  imageUrl: { en: 'Image available at:', ru: 'Изображение доступно по ссылке:' },
  logo: { en: '📸 Piclo', ru: '📸 Piclo' },
}

// Хук для управления языком
function useLanguage() {
  const [lang, setLang] = useState<Lang>(() => {
    // Проверяем localStorage
    const saved = localStorage.getItem('piclo_lang') as Lang | null
    if (saved && (saved === 'en' || saved === 'ru')) {
      return saved
    }
    // Определяем язык браузера
    const browserLang = navigator.language.toLowerCase()
    if (browserLang.startsWith('ru')) {
      return 'ru'
    }
    return 'en'
  })

  useEffect(() => {
    localStorage.setItem('piclo_lang', lang)
    document.documentElement.lang = lang
  }, [lang])

  const t = (key: string): string => {
    return translations[key]?.[lang] || key
  }

  const toggleLang = () => {
    setLang(prev => prev === 'en' ? 'ru' : 'en')
  }

  return { lang, t, toggleLang }
}

function LanguageSwitcher({ lang, toggleLang }: { lang: Lang; toggleLang: () => void }) {
  return (
    <div className="language-switcher">
      <button 
        className={`lang-btn ${lang === 'en' ? 'active' : ''}`}
        onClick={toggleLang}
        aria-label="Switch language"
      >
        <span className="lang-option">EN</span>
        <span className="lang-divider">/</span>
        <span className="lang-option">RU</span>
      </button>
    </div>
  )
}

function UploadPage() {
  const { lang, t, toggleLang } = useLanguage()
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
      <LanguageSwitcher lang={lang} toggleLang={toggleLang} />
      
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
      </div>
    </div>
  )
}

function ImageViewer() {
  const { id } = useParams<{ id: string }>()
  const { lang, t, toggleLang } = useLanguage()
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
          <LanguageSwitcher lang={lang} toggleLang={toggleLang} />
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
      </Routes>
    </BrowserRouter>
  )
}

export default App
