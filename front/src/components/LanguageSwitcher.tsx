import React from 'react';
import { useTranslation } from 'react-i18next';

interface LanguageSwitcherProps {
  className?: string;
}

const LanguageSwitcher: React.FC<LanguageSwitcherProps> = ({ className }) => {
  const { i18n } = useTranslation();

  const toggleLang = () => {
    const newLang = i18n.language === 'en' ? 'ru' : 'en';
    i18n.changeLanguage(newLang);
  };

  return (
    <div className={`language-switcher ${className || ''}`}>
      <button 
        className={`lang-btn ${i18n.language === 'en' ? 'active' : ''}`}
        onClick={toggleLang}
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
