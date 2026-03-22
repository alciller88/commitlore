<script lang="ts">
  import { onMount } from 'svelte'
  import Dashboard from './screens/Dashboard.svelte'
  import Generate from './screens/Generate.svelte'
  import Story from './screens/Story.svelte'
  import History from './screens/History.svelte'
  import Contributors from './screens/Contributors.svelte'
  import Styles from './screens/Styles.svelte'
  import Settings from './screens/Settings.svelte'
  import { activeRepo } from './lib/store'
  import type { ActiveRepo } from './lib/store'
  import { GetRecentRepos } from '../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { History as HistoryBinding } from '../bindings/github.com/alciller88/commitlore/app/gitapp.js'
  import { repoSummary } from './lib/store'

  const screens = [
    { name: 'Dashboard', icon: 'dashboard' },
    { name: 'Generate', icon: 'generate' },
    { name: 'Story', icon: 'story' },
    { name: 'History', icon: 'history' },
    { name: 'Contributors', icon: 'contributors' },
    { name: 'Styles', icon: 'styles' },
    { name: 'Settings', icon: 'settings' },
  ] as const

  let activeScreen = 'Dashboard'
  let currentRepo: ActiveRepo | null = null

  activeRepo.subscribe(value => { currentRepo = value })

  onMount(async () => {
    try {
      const recents = await GetRecentRepos()
      if (recents && recents.length > 0) {
        const r = recents[0]
        const name = extractName(r.path, r.type)
        activeRepo.set({ path: r.path, type: r.type as 'local' | 'github', name })
        loadSummaryForRepo(r.path)
      }
    } catch {}
  })

  function extractName(path: string, type: string): string {
    if (type === 'github') return path
    const parts = path.replace(/\\/g, '/').split('/')
    return parts[parts.length - 1] || path
  }

  async function loadSummaryForRepo(repo: string) {
    try {
      const raw = await HistoryBinding(repo, '', '', '', 0)
      const commits = JSON.parse(raw)
      const authors = new Set(commits.map((c: any) => c.email))
      const name = extractName(repo, repo.includes('/') && !repo.includes('\\') ? 'github' : 'local')
      repoSummary.set({
        name,
        lastCommit: commits.length > 0 ? commits[0].message : 'No commits',
        totalCommits: commits.length,
        contributors: authors.size,
      })
    } catch {
      repoSummary.set(null)
    }
  }

  function truncatePath(p: string, maxLen: number): string {
    if (p.length <= maxLen) return p
    const parts = p.replace(/\\/g, '/').split('/')
    if (parts.length <= 3) return p.substring(0, maxLen - 1) + '\u2026'
    const first = parts[0]
    const last = parts.slice(-2).join('/')
    const middle = '\u2026'
    const result = first + '/' + middle + '/' + last
    if (result.length <= maxLen) return result
    return p.substring(0, maxLen - 1) + '\u2026'
  }

  const components: Record<string, any> = {
    Dashboard, Generate, Story, History, Contributors, Styles, Settings,
  }
</script>

<div class="layout">
  <nav class="sidebar">
    <div class="sidebar-header">CommitLore</div>
    <div class="nav-items">
      {#each screens as screen}
        <button
          class="nav-item"
          class:active={activeScreen === screen.name}
          on:click={() => activeScreen = screen.name}
        >
          <span class="nav-icon">
            {#if screen.icon === 'dashboard'}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
            {:else if screen.icon === 'generate'}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><path d="M12 5v14M5 12h14"/></svg>
            {:else if screen.icon === 'story'}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><path d="M4 19.5A2.5 2.5 0 016.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z"/></svg>
            {:else if screen.icon === 'history'}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            {:else if screen.icon === 'contributors'}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 00-3-3.87"/><path d="M16 3.13a4 4 0 010 7.75"/></svg>
            {:else if screen.icon === 'styles'}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><circle cx="12" cy="12" r="3"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
            {:else if screen.icon === 'settings'}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
            {/if}
          </span>
          <span class="nav-label">{screen.name}</span>
        </button>
      {/each}
    </div>

    <div class="sidebar-repo-indicator">
      <div class="indicator-separator"></div>
      {#if currentRepo}
        <div class="indicator-content">
          <span class="indicator-icon">
            {#if currentRepo.type === 'github'}
              <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14"><path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/></svg>
            {:else}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
            {/if}
          </span>
          <div class="indicator-text">
            <span class="indicator-name">{currentRepo.name}</span>
            <span class="indicator-path">{truncatePath(currentRepo.path, 28)}</span>
          </div>
        </div>
      {:else}
        <div class="indicator-empty">No repo selected</div>
      {/if}
    </div>
  </nav>
  <main class="content">
    <svelte:component this={components[activeScreen]} />
  </main>
</div>

<style>
  .layout {
    display: flex;
    height: 100vh;
    background: #0d1117;
    color: #e6edf3;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  }

  .sidebar {
    width: 220px;
    background: #161b22;
    border-right: 1px solid #30363d;
    display: flex;
    flex-direction: column;
    padding: 0;
    flex-shrink: 0;
  }

  .sidebar-header {
    padding: 20px 16px;
    font-size: 18px;
    font-weight: 700;
    color: #f0f6fc;
    border-bottom: 1px solid #30363d;
    font-family: 'JetBrains Mono', 'Courier New', monospace;
  }

  .nav-items {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 16px;
    border: none;
    background: transparent;
    color: #8b949e;
    font-size: 14px;
    cursor: pointer;
    text-align: left;
    width: 100%;
    transition: background 0.15s, color 0.15s;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  }

  .nav-item:hover {
    background: #1f2937;
    color: #e6edf3;
  }

  .nav-item.active {
    background: #1f6feb22;
    color: #58a6ff;
    border-left: 3px solid #58a6ff;
  }

  .nav-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 16px;
  }

  .content {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
  }

  .sidebar-repo-indicator {
    padding: 0 12px 12px;
  }

  .indicator-separator {
    border-top: 1px solid #30363d;
    margin-bottom: 10px;
  }

  .indicator-content {
    display: flex;
    gap: 8px;
    align-items: flex-start;
  }

  .indicator-icon {
    display: flex;
    align-items: center;
    color: #8b949e;
    flex-shrink: 0;
    margin-top: 1px;
  }

  .indicator-text {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }

  .indicator-name {
    color: #e6edf3;
    font-size: 12px;
    font-weight: 600;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .indicator-path {
    color: #8b949e;
    font-size: 10px;
    font-family: 'JetBrains Mono', monospace;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .indicator-empty {
    color: #484f58;
    font-size: 12px;
  }
</style>
