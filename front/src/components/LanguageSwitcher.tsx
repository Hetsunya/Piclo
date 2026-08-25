import React from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate, useLocation } from 'react-router-dom';

interface LanguageSwitcherProps {
  className?: string;
}

const SUPPORTED_LANGUAGES = ['en', 'ru'];

const LanguageSwitcher: React.FC<LanguageSwitcherProps> = ({ className }) => {
  const { i18n } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();

  const toggleLang = (newLang: string) => {
    // Сохраняем в localStorage для детектора
    localStorage.setItem('piclo_lang', newLang);
    
    // Меняем язык в i18next
    i18n.changeLanguage(newLang);
    
    // Обновляем URL, заменяя префикс языка
    const pathParts = location.pathname.split('/').filter(Boolean);
    const restOfPath = pathParts.length > 1 
      ? '/' + pathParts.slice(1).join('/') 
      : '/';
    
    // Если мы на странице с языковым префиксом, заменяем его
    if (pathParts.length > 0 && SUPPORTED_LANGUAGES.includes(pathParts[0])) {
      navigate(`/${newLang}${restOfPath}${location.search}`, { replace: true });
    } else {
      // Если нет префикса, просто добавляем новый язык
      navigate(`/${newLang}${location.pathname}${location.search}`, { replace: true });
    }
  };

  return (
    <div className={`language-switcher ${className || ''}`}>
      <button 
        className={`lang-btn ${i18n.language === 'en' ? 'active' : ''}`}
        onClick={() => toggleLang(i18n.language === 'en' ? 'ru' : 'en')}
        aria-label="Switch language"
      >
        <span className="lang-option">EN</span>
        <span className="lang-divider">/</span>
        <span className="lang-option">RU</span>
      </button>
    </div>
  );
};

export default LanguageSwitcher;
