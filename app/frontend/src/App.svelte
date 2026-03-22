<script lang="ts">
  import { onMount } from 'svelte'
  import Dashboard from './screens/Dashboard.svelte'
  import Generate from './screens/Generate.svelte'
  import Story from './screens/Story.svelte'
  import History from './screens/History.svelte'
  import Contributors from './screens/Contributors.svelte'
  import Styles from './screens/Styles.svelte'
  import Settings from './screens/Settings.svelte'
  import { activeRepo, repoSummary, activeStyle } from './lib/store'
  import type { ActiveRepo } from './lib/store'
  import { GetRecentRepos, GetActiveStyle } from '../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { History as HistoryBinding } from '../bindings/github.com/alciller88/commitlore/app/gitapp.js'
  import { applyTheme, getTheme } from './lib/theme'
  import type { ThemeVars } from './lib/theme'
  import { Window } from '@wailsio/runtime'

  const viewScreens = [
    { name: 'Dashboard', icon: 'dashboard' },
    { name: 'Generate', icon: 'generate' },
    { name: 'Story', icon: 'story' },
    { name: 'History', icon: 'history' },
    { name: 'Contributors', icon: 'contributors' },
  ]
  const systemScreens = [
    { name: 'Styles', icon: 'styles' },
    { name: 'Settings', icon: 'settings' },
  ]

  let activeScreen = 'Dashboard'
  let currentRepo: ActiveRepo | null = null
  let theme: ThemeVars = getTheme()

  activeRepo.subscribe(value => { currentRepo = value })
  activeStyle.subscribe(async (styleName) => {
    theme = await applyTheme(styleName)
  })

  onMount(async () => {
    try {
      const styleName = await GetActiveStyle()
      activeStyle.set(styleName)
    } catch {}
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
      repoSummary.set({ name, lastCommit: commits.length > 0 ? commits[0].message : 'No commits', totalCommits: commits.length, contributors: authors.size })
    } catch { repoSummary.set(null) }
  }

  function truncatePath(p: string, maxLen: number): string {
    if (p.length <= maxLen) return p
    const parts = p.replace(/\\/g, '/').split('/')
    if (parts.length <= 3) return p.substring(0, maxLen - 1) + '\u2026'
    return parts[0] + '/\u2026/' + parts.slice(-2).join('/')
  }

  function winMinimize() { Window.Minimise() }
  function winMaximize() { Window.ToggleMaximise() }
  function winClose() { Window.Close() }

  const components: Record<string, any> = { Dashboard, Generate, Story, History, Contributors, Styles, Settings }
</script>

<div class="layout">
  <nav class="sidebar" style="--wails-draggable: drag">
    <div class="sidebar-header" style="--wails-draggable: drag">
      <div class="brand" style="--wails-draggable: no-drag">
        {#if theme.logo}
          <span class="brand-logo">
            <img src={theme.logo} alt="" width="24" height="24" />
          </span>
        {/if}
        <span class="brand-text"><span class="wm-c">Commit</span><span class="wm-l">Lore</span></span>
      </div>
      <div class="win-controls" style="--wails-draggable: no-drag">
        <button class="wc wc-close" on:click={winClose} title="Close"><svg viewBox="0 0 12 12" width="10" height="10"><path d="M1 1l10 10M11 1L1 11" stroke="currentColor" stroke-width="1.5" fill="none"/></svg></button>
        <button class="wc wc-min" on:click={winMinimize} title="Minimize"><svg viewBox="0 0 12 12" width="10" height="10"><path d="M2 6h8" stroke="currentColor" stroke-width="1.5" fill="none"/></svg></button>
        <button class="wc wc-max" on:click={winMaximize} title="Maximize"><svg viewBox="0 0 12 12" width="10" height="10"><rect x="2" y="2" width="8" height="8" rx="1" stroke="currentColor" stroke-width="1.2" fill="none"/></svg></button>
      </div>
    </div>

    <div class="nav-section" style="--wails-draggable: no-drag">
      <span class="nav-label-section">Views</span>
      {#each viewScreens as screen}
        <button class="nav-item" class:active={activeScreen === screen.name} on:click={() => activeScreen = screen.name}>
          <span class="nav-icon">
            {#if screen.icon === 'dashboard'}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/></svg>
            {:else if screen.icon === 'generate'}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><path d="M12 5v14M5 12h14"/></svg>
            {:else if screen.icon === 'story'}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><path d="M4 19.5A2.5 2.5 0 016.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 014 19.5v-15A2.5 2.5 0 016.5 2z"/></svg>
            {:else if screen.icon === 'history'}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
            {:else if screen.icon === 'contributors'}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 00-3-3.87"/><path d="M16 3.13a4 4 0 010 7.75"/></svg>
            {/if}
          </span>
          <span class="nav-text">{screen.name}</span>
        </button>
      {/each}

      <span class="nav-label-section">System</span>
      {#each systemScreens as screen}
        <button class="nav-item" class:active={activeScreen === screen.name} on:click={() => activeScreen = screen.name}>
          <span class="nav-icon">
            {#if screen.icon === 'styles'}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><circle cx="12" cy="12" r="3"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
            {:else if screen.icon === 'settings'}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
            {/if}
          </span>
          <span class="nav-text">{screen.name}</span>
        </button>
      {/each}
    </div>

  </nav>

  <div class="main-area">
    <div class="topbar">
      {#if currentRepo}
        <span class="tb-icon">
          {#if currentRepo.type === 'github'}<svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14"><path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/></svg>
          {:else}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
          {/if}
        </span>
        <span class="tb-name">{currentRepo.name}</span>
        <span class="tb-sep">&middot;</span>
        <span class="tb-path">{truncatePath(currentRepo.path, 60)}</span>
      {:else}
        <span class="tb-empty">No repository selected</span>
      {/if}
    </div>
    <main class="content">
      <svelte:component this={components[activeScreen]} />
    </main>
  </div>
</div>

<style>
  .layout {
    display: flex; height: 100vh;
    background: var(--cl-background, #0d1117);
    color: var(--cl-text, #e6edf3);
    font-family: var(--cl-font-family, system-ui, sans-serif);
    font-size: var(--text-base, 13px);
  }

  .sidebar {
    width: 220px; flex-shrink: 0;
    background: var(--cl-background, #0d1117);
    border-right: 1px solid var(--cl-surface, #161b22);
    display: flex; flex-direction: column;
  }

  .sidebar-header {
    height: 52px; display: flex; align-items: center;
    justify-content: space-between;
    padding: 0 var(--space-3, 12px);
    border-bottom: 1px solid var(--cl-surface, #161b22);
    flex-shrink: 0;
  }

  .brand { display: flex; align-items: center; gap: var(--space-2, 8px); }
  .brand-logo { display: flex; color: var(--cl-accent, #58a6ff); flex-shrink: 0; }
  .brand-logo img { border-radius: var(--radius-sm, 4px); }
  .brand-text { font-size: 15px; font-weight: 600; line-height: 1; font-family: var(--cl-font-family, system-ui, sans-serif); }
  .wm-c { color: var(--cl-primary, #58a6ff); }
  .wm-l { color: var(--cl-accent, #58a6ff); }

  .win-controls { display: flex; gap: 6px; align-items: center; }
  .wc {
    width: 12px; height: 12px; border-radius: 50%; border: none;
    background: var(--cl-win-default, #666666); padding: 0;
    display: flex; align-items: center; justify-content: center;
    color: transparent; transition: var(--transition-fast, 120ms ease);
  }
  .wc svg { opacity: 0; transition: var(--transition-fast, 120ms ease); }
  .wc-close:hover { background: var(--cl-win-close, #FF5F57); color: #fff; }
  .wc-close:hover svg { opacity: 1; }
  .wc-min:hover { background: var(--cl-win-minimize, #FEBC2E); color: #1a1a1a; }
  .wc-min:hover svg { opacity: 1; }
  .wc-max:hover { background: var(--cl-win-maximize, #28C840); color: #1a1a1a; }
  .wc-max:hover svg { opacity: 1; }

  .nav-section { flex: 1; overflow-y: auto; padding: var(--space-2, 8px) 0; }
  .nav-label-section {
    display: block; padding: var(--space-4, 16px) var(--space-3, 12px) var(--space-1, 4px);
    font-size: var(--text-xs, 11px); text-transform: uppercase; letter-spacing: 0.05em;
    color: var(--cl-secondary, #8b949e); opacity: 0.6;
  }

  .nav-item {
    display: flex; align-items: center; gap: var(--space-2, 8px);
    height: 32px; padding: 0 var(--space-3, 12px);
    border: none; border-left: 2px solid transparent;
    background: transparent;
    color: var(--cl-secondary, #8b949e);
    font-size: var(--text-base, 13px);
    width: 100%; text-align: left;
    font-family: var(--cl-font-family, system-ui, sans-serif);
    transition: background var(--transition-fast, 120ms ease), color var(--transition-fast, 120ms ease);
  }
  .nav-item:hover {
    background: color-mix(in srgb, var(--cl-surface, #161b22) 60%, transparent);
    color: var(--cl-text, #e6edf3);
  }
  .nav-item.active {
    background: var(--cl-surface, #161b22);
    color: var(--cl-accent, #58a6ff);
    border-left-color: var(--cl-accent, #58a6ff);
  }
  .nav-icon { display: flex; align-items: center; width: 16px; height: 16px; flex-shrink: 0; }
  .nav-text { white-space: nowrap; }

  .main-area { flex: 1; display: flex; flex-direction: column; min-width: 0; }

  .topbar {
    height: 32px; flex-shrink: 0; display: flex; align-items: center; gap: var(--space-2, 8px);
    padding: 0 var(--space-4, 16px);
    background: var(--cl-surface, #161b22);
    border-bottom: 1px solid var(--cl-background, #0d1117);
    font-size: var(--text-base, 13px);
  }
  .tb-icon { display: flex; color: var(--cl-secondary, #8b949e); flex-shrink: 0; }
  .tb-name { color: var(--cl-text, #e6edf3); font-weight: 500; white-space: nowrap; }
  .tb-sep { color: var(--cl-secondary, #8b949e); opacity: 0.5; }
  .tb-path {
    color: var(--cl-secondary, #8b949e); font-size: var(--text-sm, 12px);
    font-family: 'JetBrains Mono', monospace;
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 400px;
  }
  .tb-empty { color: var(--cl-secondary, #8b949e); font-style: italic; opacity: 0.5; }

  .content { flex: 1; overflow-y: auto; padding: var(--space-5, 24px); }
</style>
