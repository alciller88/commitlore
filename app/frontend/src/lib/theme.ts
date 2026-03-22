import { GetStyleTheme } from '../../bindings/github.com/alciller88/commitlore/app/styleapp.js'

export interface ThemeVars {
  primary: string
  secondary: string
  background: string
  surface: string
  text: string
  accent: string
  border: string
  fontFamily: string
  fontSize: string
  mode: string
  logo: string
}

const DEFAULTS: ThemeVars = {
  primary: '#58a6ff',
  secondary: '#8b949e',
  background: '#0d1117',
  surface: '#161b22',
  text: '#e6edf3',
  accent: '#58a6ff',
  border: '#30363d',
  fontFamily: 'system-ui, sans-serif',
  fontSize: '14px',
  mode: 'dark',
  logo: '',
}

let currentTheme: ThemeVars = { ...DEFAULTS }

export function getTheme(): ThemeVars {
  return currentTheme
}

export async function applyTheme(styleName: string): Promise<ThemeVars> {
  try {
    const theme = await GetStyleTheme(styleName)
    currentTheme = {
      primary: theme.primary || DEFAULTS.primary,
      secondary: theme.secondary || DEFAULTS.secondary,
      background: theme.background || DEFAULTS.background,
      surface: theme.surface || DEFAULTS.surface,
      text: theme.text || DEFAULTS.text,
      accent: theme.accent || DEFAULTS.accent,
      border: theme.border || DEFAULTS.border,
      fontFamily: theme.fontFamily || DEFAULTS.fontFamily,
      fontSize: theme.fontSize || DEFAULTS.fontSize,
      mode: theme.mode || DEFAULTS.mode,
      logo: theme.logo || '',
    }
  } catch {
    currentTheme = { ...DEFAULTS }
  }
  injectCSSVariables(currentTheme)
  return currentTheme
}

function injectCSSVariables(t: ThemeVars) {
  const root = document.documentElement
  root.style.setProperty('--cl-primary', t.primary)
  root.style.setProperty('--cl-secondary', t.secondary)
  root.style.setProperty('--cl-background', t.background)
  root.style.setProperty('--cl-surface', t.surface)
  root.style.setProperty('--cl-text', t.text)
  root.style.setProperty('--cl-accent', t.accent)
  root.style.setProperty('--cl-border', t.border)
  root.style.setProperty('--cl-font-family', t.fontFamily)
  root.style.setProperty('--cl-font-size', t.fontSize)
}
