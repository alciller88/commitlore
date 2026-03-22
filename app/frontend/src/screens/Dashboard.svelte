<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { OpenFolderPicker, History } from '../../bindings/github.com/alciller88/commitlore/app/gitapp.js'
  import { GetRecentRepos, AddRecentRepo } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { Events } from '@wailsio/runtime'
  import { activeRepo, repoSummary } from '../lib/store'
  import type { ActiveRepo, RepoSummary } from '../lib/store'

  let currentRepo: ActiveRepo | null = null
  let currentSummary: RepoSummary | null = null
  let recentRepos: Array<{path: string, type: string, lastUsed: string}> = []
  let githubInput = ''
  let error = ''
  let loading = false
  let unsubFileDrop: (() => void) | null = null

  const unsubRepo = activeRepo.subscribe(v => { currentRepo = v })
  const unsubSummary = repoSummary.subscribe(v => { currentSummary = v })

  onMount(() => {
    loadRecentRepos()
    unsubFileDrop = Events.On('file-dropped', (event: any) => {
      const path = event.data
      if (path && typeof path === 'string') selectRepo(path, 'local')
    })
  })
  onDestroy(() => { if (unsubFileDrop) unsubFileDrop(); unsubRepo(); unsubSummary() })

  async function loadRecentRepos() { try { recentRepos = await GetRecentRepos() } catch (e: any) { error = e?.message || 'Failed to load recent repos' } }
  async function openFolder() {
    loading = true; error = ''
    try { const p = await OpenFolderPicker(); if (p) await selectRepo(p, 'local') }
    catch (e: any) { error = e?.message || 'Failed to open folder' } finally { loading = false }
  }
  async function connectGitHub() {
    if (!githubInput.trim()) return; loading = true; error = ''
    try { await selectRepo(githubInput.trim(), 'github') }
    catch (e: any) { error = e?.message || 'Failed to connect' } finally { loading = false }
  }
  function extractName(path: string, type: string): string {
    if (type === 'github') return path
    const parts = path.replace(/\\/g, '/').split('/'); return parts[parts.length - 1] || path
  }
  async function selectRepo(path: string, type: string) {
    await AddRecentRepo(path, type)
    activeRepo.set({ path, type: type as 'local' | 'github', name: extractName(path, type) })
    await loadRecentRepos(); await loadRepoSummary(path)
  }
  async function loadRepoSummary(repo: string) {
    try {
      const raw = await History(repo, '', '', '', 0); const commits = JSON.parse(raw)
      const authors = new Set(commits.map((c: any) => c.email))
      repoSummary.set({ name: extractName(repo, repo.includes('/') && !repo.includes('\\') ? 'github' : 'local') || repo, lastCommit: commits.length > 0 ? commits[0].message : 'No commits', totalCommits: commits.length, contributors: authors.size })
    } catch { repoSummary.set(null) }
  }
  function changeRepo() { activeRepo.set(null); repoSummary.set(null) }
  function shortenPath(p: string) {
    if (p.length <= 40) return p; const parts = p.replace(/\\/g, '/').split('/')
    if (parts.length <= 3) return p; return parts[0] + '/.../' + parts.slice(-2).join('/')
  }
  function formatDate(iso: string) { try { return new Date(iso).toLocaleDateString() } catch { return iso } }
</script>

{#if error}<div class="banner-err">{error}</div>{/if}

{#if !currentRepo}
  <div class="empty">
    <h1>Welcome to CommitLore</h1>
    <p class="sub">Your repo has a story. Select one to begin.</p>

    <div class="entry-cards">
      <button class="ecard" on:click={openFolder} disabled={loading}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="24" height="24"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
        <span class="ecard-title">Open folder</span>
        <span class="ecard-sub">Local git repository</span>
      </button>

      <div class="ecard drop" data-file-drop-target="true" role="button" tabindex="0">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="24" height="24"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
        <span class="ecard-title">Drop a folder</span>
        <span class="ecard-sub">Drag from Explorer</span>
      </div>

      <div class="ecard gh-card">
        <svg viewBox="0 0 16 16" fill="currentColor" width="24" height="24"><path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/></svg>
        <div class="gh-input-row">
          <input type="text" bind:value={githubInput} placeholder="owner/repo" on:keydown={(e) => e.key === 'Enter' && connectGitHub()} />
          <button class="gh-btn" on:click={connectGitHub} disabled={loading || !githubInput.trim()}>Connect</button>
        </div>
      </div>
    </div>

    {#if recentRepos.length > 0}
      <div class="recent">
        <span class="section-label">Recent</span>
        {#each recentRepos.slice(0, 5) as repo}
          <button class="recent-row" on:click={() => selectRepo(repo.path, repo.type)}>
            <span class="rr-icon">
              {#if repo.type === 'github'}<svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14"><path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/></svg>
              {:else}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
              {/if}
            </span>
            <span class="rr-name">{repo.path.replace(/\\/g, '/').split('/').pop()}</span>
            <span class="rr-path">{shortenPath(repo.path)}</span>
            <span class="rr-date">{formatDate(repo.lastUsed)}</span>
          </button>
        {/each}
      </div>
    {/if}
  </div>
{:else}
  <div class="active-view">
    <div class="active-header">
      <div>
        <h1 class="repo-name">{currentSummary?.name || currentRepo.name}</h1>
        <span class="repo-path">{currentRepo.path}</span>
      </div>
      <button class="change-link" on:click={changeRepo}>Change</button>
    </div>

    {#if currentSummary}
      <div class="stats-row">
        <div class="stat-card"><span class="stat-val">{currentSummary.totalCommits}</span><span class="stat-lbl">Commits</span></div>
        <div class="stat-card"><span class="stat-val">{currentSummary.contributors}</span><span class="stat-lbl">Contributors</span></div>
        <div class="stat-card last-commit-card"><span class="stat-lbl">Last commit</span><span class="stat-msg">{currentSummary.lastCommit}</span></div>
      </div>
    {/if}

    {#if recentRepos.length > 1}
      <div class="recent">
        <span class="section-label">Recent</span>
        {#each recentRepos.filter(r => r.path !== currentRepo?.path).slice(0, 5) as repo}
          <button class="recent-row" on:click={() => selectRepo(repo.path, repo.type)}>
            <span class="rr-icon">
              {#if repo.type === 'github'}<svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14"><path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/></svg>
              {:else}<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
              {/if}
            </span>
            <span class="rr-name">{repo.path.replace(/\\/g, '/').split('/').pop()}</span>
            <span class="rr-path">{shortenPath(repo.path)}</span>
            <span class="rr-date">{formatDate(repo.lastUsed)}</span>
          </button>
        {/each}
      </div>
    {/if}
  </div>
{/if}

<style>
  .banner-err { background: #da363422; border: 1px solid #da3634; color: #f85149; padding: var(--space-2) var(--space-3); border-radius: var(--radius-md); font-size: var(--text-base); margin-bottom: var(--space-3); }
  .empty { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; gap: var(--space-5); }
  h1 { color: var(--cl-text); font-size: var(--text-xl); font-weight: 600; }
  .sub { color: var(--cl-secondary); font-size: var(--text-base); }

  .entry-cards { display: flex; gap: var(--space-4); flex-wrap: wrap; justify-content: center; }
  .ecard {
    width: 200px; display: flex; flex-direction: column; align-items: center; gap: var(--space-2);
    padding: var(--space-5) var(--space-4);
    background: var(--cl-surface); border: 1px solid var(--cl-border); border-radius: var(--radius-lg);
    color: var(--cl-text); font-family: inherit; font-size: var(--text-base);
    transition: border-color var(--transition-fast);
  }
  .ecard:hover { border-color: var(--cl-accent); }
  .ecard:disabled { opacity: 0.5; cursor: not-allowed; }
  .ecard-title { font-weight: 500; }
  .ecard-sub { font-size: var(--text-xs); color: var(--cl-secondary); }
  .ecard.drop { border-style: dashed; cursor: default; }
  .ecard svg { color: var(--cl-accent); }

  .gh-card { gap: var(--space-3); }
  .gh-input-row { display: flex; gap: var(--space-2); width: 100%; }
  .gh-input-row input {
    flex: 1; padding: var(--space-2); background: var(--cl-background); border: 1px solid var(--cl-border);
    border-radius: var(--radius-sm); color: var(--cl-text); font-size: var(--text-sm); font-family: 'JetBrains Mono', monospace;
  }
  .gh-input-row input:focus { outline: none; border-color: var(--cl-accent); }
  .gh-btn { padding: var(--space-2) var(--space-3); background: #238636; border: none; border-radius: var(--radius-sm); color: #fff; font-size: var(--text-sm); font-family: inherit; }
  .gh-btn:hover { background: #2ea043; }
  .gh-btn:disabled { opacity: 0.5; cursor: not-allowed; }

  .section-label { display: block; font-size: var(--text-xs); text-transform: uppercase; letter-spacing: 0.05em; color: var(--cl-secondary); opacity: 0.6; margin-bottom: var(--space-2); }
  .recent { width: 100%; max-width: 560px; }
  .recent-row {
    display: flex; align-items: center; gap: var(--space-2); width: 100%; padding: var(--space-2) var(--space-3);
    background: transparent; border: none; border-radius: var(--radius-sm); color: var(--cl-text);
    font-family: inherit; font-size: var(--text-base); text-align: left; transition: background var(--transition-fast);
  }
  .recent-row:hover { background: var(--cl-surface); }
  .rr-icon { display: flex; color: var(--cl-secondary); flex-shrink: 0; }
  .rr-name { font-weight: 500; font-size: var(--text-base); min-width: 80px; }
  .rr-path { flex: 1; font-size: var(--text-xs); color: var(--cl-secondary); font-family: 'JetBrains Mono', monospace; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .rr-date { font-size: var(--text-xs); color: var(--cl-secondary); flex-shrink: 0; }

  .active-view { display: flex; flex-direction: column; gap: var(--space-5); }
  .active-header { display: flex; justify-content: space-between; align-items: flex-start; }
  .repo-name { font-size: var(--text-xl); font-weight: 600; }
  .repo-path { font-size: var(--text-sm); color: var(--cl-secondary); font-family: 'JetBrains Mono', monospace; }
  .change-link { background: none; border: none; color: var(--cl-accent); font-size: var(--text-sm); font-family: inherit; padding: 0; }
  .change-link:hover { text-decoration: underline; }

  .stats-row { display: flex; gap: var(--space-3); flex-wrap: wrap; }
  .stat-card { display: flex; flex-direction: column; gap: var(--space-1); background: var(--cl-surface); border-radius: var(--radius-md); padding: var(--space-3) var(--space-4); min-width: 110px; }
  .stat-val { font-size: 24px; font-weight: 700; color: var(--cl-accent); }
  .stat-lbl { font-size: var(--text-xs); color: var(--cl-secondary); text-transform: uppercase; }
  .last-commit-card { flex: 1; min-width: 200px; }
  .stat-msg { font-size: var(--text-base); color: var(--cl-text); font-family: 'JetBrains Mono', monospace; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
