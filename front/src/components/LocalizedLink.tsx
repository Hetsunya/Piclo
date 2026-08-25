import { Link, useLocation } from 'react-router-dom';
import type { To } from 'react-router-dom';

interface LocalizedLinkProps {
  to: To;
  children?: React.ReactNode;
  className?: string;
  [key: string]: any;
}

export default function LocalizedLink({ to, children, className, ...rest }: LocalizedLinkProps) {
  const location = useLocation();
  
  // Извлекаем текущий язык из pathname
  const pathParts = location.pathname.split('/').filter(Boolean);
  const currentLng = pathParts.length > 0 && ['en', 'ru'].includes(pathParts[0]) 
    ? pathParts[0] 
    : 'ru';

  // Обрабатываем to как строку или объект
  let localizedTo: To;
  
  if (typeof to === 'string') {
    // Если путь абсолютный (начинается с /), добавляем префикс языка
    if (to.startsWith('/')) {
      localizedTo = `/${currentLng}${to}`;
    } else {
      localizedTo = to;
    }
  } else {
    // Если to - объект (например, { pathname: '/about', search: '?q=1' })
    localizedTo = {
      ...to,
      pathname: to.pathname?.startsWith('/') 
        ? `/${currentLng}${to.pathname}` 
        : to.pathname,
    };
  }

  return (
    <Link to={localizedTo} className={className} {...rest}>
      {children}
    </Link>
  );
}
