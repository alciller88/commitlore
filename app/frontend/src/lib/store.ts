import { writable } from 'svelte/store'

export interface ActiveRepo {
  path: string
  type: 'local' | 'github'
  name: string
}

export interface RepoSummary {
  name: string
  lastCommit: string
  totalCommits: number
  contributors: number
}

export const activeRepo = writable<ActiveRepo | null>(null)
export const repoSummary = writable<RepoSummary | null>(null)
export const activeStyle = writable<string>('formal')

export interface UILabelsType {
  dashboard: string
  generate: string
  story: string
  history: string
  contributors: string
  styles: string
  settings: string
  generateButton: string
  storyButton: string
}

export const uiLabels = writable<UILabelsType>({
  dashboard: 'Dashboard',
  generate: 'Generate',
  story: 'Story',
  history: 'History',
  contributors: 'Contributors',
  styles: 'Styles',
  settings: 'Settings',
  generateButton: 'Generate',
  storyButton: 'Tell the story',
})
