import { createContext, useContext, useState, useEffect, ReactNode } from 'react'
import { 
  Language, 
  Translations, 
  getTranslations, 
  getStoredLanguage, 
  setStoredLanguage 
} from '../i18n'

interface LanguageContextType {
  language: Language
  setLanguage: (lang: Language) => void
  t: Translations
}

const LanguageContext = createContext<LanguageContextType | undefined>(undefined)

interface LanguageProviderProps {
  children: ReactNode
}

export function LanguageProvider({ children }: LanguageProviderProps) {
  const [language, setLanguageState] = useState<Language>(getStoredLanguage())
  const [translations, setTranslations] = useState<Translations>(getTranslations(language))

  useEffect(() => {
    setTranslations(getTranslations(language))
  }, [language])

  const setLanguage = (lang: Language) => {
    setLanguageState(lang)
    setStoredLanguage(lang)
  }

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t: translations }}>
      {children}
    </LanguageContext.Provider>
  )
}

export function useLanguage() {
  const context = useContext(LanguageContext)
  if (context === undefined) {
    throw new Error('useLanguage must be used within a LanguageProvider')
  }
  return context
}


