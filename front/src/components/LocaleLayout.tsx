import { useEffect } from 'react';
import { Outlet, useNavigate, useParams, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Header } from '../App';

const SUPPORTED_LANGUAGES = ['en', 'ru'];
const DEFAULT_LANGUAGE = 'ru';

export default function LocaleLayout() {
  const { lng } = useParams<{ lng: string }>();
  const { i18n } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    // Если язык не указан или не поддерживается, редирект на дефолтный
    if (!lng || !SUPPORTED_LANGUAGES.includes(lng)) {
      // Сохраняем текущий путь без языкового префикса
      const pathWithoutLng = location.pathname.replace(/^\/[^\/]+/, '') || '/';
      navigate(`/${DEFAULT_LANGUAGE}${pathWithoutLng}${location.search}`, { replace: true });
      return;
    }

    // Синхронизируем язык из URL с i18next и localStorage
    if (i18n.language !== lng) {
      i18n.changeLanguage(lng);
      localStorage.setItem('piclo_lang', lng);
    }
  }, [lng, i18n, navigate, location.pathname, location.search]);

  return (
    <>
      <Header />
      <Outlet />
    </>
  );
}
