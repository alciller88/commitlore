import { GetStyleTheme } from '../../bindings/github.com/alciller88/commitlore/app/styleapp.js'
import { uiLabels } from './store'

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
  winDefault: string
  winClose: string
  winMinimize: string
  winMaximize: string
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
  winDefault: '#666666',
  winClose: '#FF5F57',
  winMinimize: '#FEBC2E',
  winMaximize: '#28C840',
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
      winDefault: theme.winDefault || DEFAULTS.winDefault,
      winClose: theme.winClose || DEFAULTS.winClose,
      winMinimize: theme.winMinimize || DEFAULTS.winMinimize,
      winMaximize: theme.winMaximize || DEFAULTS.winMaximize,
    }
    const labels = (theme as any).uiLabels
    if (labels) {
      uiLabels.set({
        dashboard: labels.dashboard || 'Dashboard',
        generate: labels.generate || 'Generate',
        story: labels.story || 'Story',
        history: labels.history || 'History',
        contributors: labels.contributors || 'Contributors',
        styles: labels.styles || 'Styles',
        settings: labels.settings || 'Settings',
        generateButton: labels.generateButton || 'Generate',
        storyButton: labels.storyButton || 'Tell the story',
      })
    } else {
      uiLabels.set({
        dashboard: 'Dashboard', generate: 'Generate', story: 'Story',
        history: 'History', contributors: 'Contributors', styles: 'Styles', settings: 'Settings',
        generateButton: 'Generate', storyButton: 'Tell the story',
      })
    }
  } catch {
    currentTheme = { ...DEFAULTS }
    uiLabels.set({
      dashboard: 'Dashboard', generate: 'Generate', story: 'Story',
      history: 'History', contributors: 'Contributors', styles: 'Styles', settings: 'Settings',
      generateButton: 'Generate', storyButton: 'Tell the story',
    })
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
  root.style.setProperty('--cl-win-default', t.winDefault)
  root.style.setProperty('--cl-win-close', t.winClose)
  root.style.setProperty('--cl-win-minimize', t.winMinimize)
  root.style.setProperty('--cl-win-maximize', t.winMaximize)
}
