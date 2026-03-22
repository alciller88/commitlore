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
