<script lang="ts">
  import { createEventDispatcher, onMount, onDestroy } from 'svelte'
  import { OpenFolderPicker, History } from '../../bindings/github.com/alciller88/commitlore/app/gitapp.js'
  import { GetRecentRepos, AddRecentRepo } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { Events } from '@wailsio/runtime'

  export let activeRepo = ''

  const dispatch = createEventDispatcher()

  let recentRepos: Array<{path: string, type: string, lastUsed: string}> = []
  let githubInput = ''
  let error = ''
  let loading = false
  let repoSummary: {name: string, lastCommit: string, totalCommits: number, contributors: number} | null = null
  let unsubFileDrop: (() => void) | null = null

  onMount(() => {
    loadRecentRepos()
    unsubFileDrop = Events.On('file-dropped', (event: any) => {
      const path = event.data
      if (path && typeof path === 'string') {
        selectRepo(path, 'local')
      }
    })
  })

  onDestroy(() => {
    if (unsubFileDrop) unsubFileDrop()
  })

  async function loadRecentRepos() {
    try {
      recentRepos = await GetRecentRepos()
    } catch (e: any) {
      error = e?.message || 'Failed to load recent repos'
    }
  }

  async function openFolder() {
    loading = true
    error = ''
    try {
      const path = await OpenFolderPicker()
      if (path) {
        await selectRepo(path, 'local')
      }
    } catch (e: any) {
      error = e?.message || 'Failed to open folder'
    } finally {
      loading = false
    }
  }

  async function connectGitHub() {
    if (!githubInput.trim()) return
    loading = true
    error = ''
    try {
      await selectRepo(githubInput.trim(), 'github')
    } catch (e: any) {
      error = e?.message || 'Failed to connect to GitHub repo'
    } finally {
      loading = false
    }
  }

  async function selectRepo(path: string, type: string) {
    await AddRecentRepo(path, type)
    activeRepo = path
    dispatch('repoSelected', path)
    await loadRecentRepos()
    await loadRepoSummary(path)
  }

  async function loadRepoSummary(repo: string) {
    try {
      const raw = await History(repo, '', '', '', 0)
      const commits = JSON.parse(raw)
      const authors = new Set(commits.map((c: any) => c.email))
      const name = repo.includes('/') ? repo.split('/').pop() : repo.split('\\').pop()
      repoSummary = {
        name: name || repo,
        lastCommit: commits.length > 0 ? commits[0].message : 'No commits',
        totalCommits: commits.length,
        contributors: authors.size,
      }
    } catch {
      repoSummary = null
    }
  }

  function shortenPath(p: string) {
    if (p.length <= 40) return p
    const parts = p.replace(/\\/g, '/').split('/')
    if (parts.length <= 3) return p
    return parts[0] + '/.../' + parts.slice(-2).join('/')
  }

  function formatDate(iso: string) {
    try {
      return new Date(iso).toLocaleDateString()
    } catch {
      return iso
    }
  }
</script>

{#if error}
  <div class="banner error">{error}</div>
{/if}

{#if !activeRepo}
  <div class="empty-state">
    <h1>Welcome to CommitLore</h1>
    <p class="subtitle">Your repo has a story. Select one to begin.</p>

    <div class="entry-points">
      <button class="entry-btn" on:click={openFolder} disabled={loading}>
        <svg class="entry-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="28" height="28">
          <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>
        <span class="entry-label">Open folder</span>
      </button>

      <div
        class="drop-zone"
        data-file-drop-target="true"
        role="button"
        tabindex="0"
      >
        <span class="drop-text">Drag & drop a folder here</span>
      </div>

      <div class="github-input">
        <input
          type="text"
          bind:value={githubInput}
          placeholder="owner/repo"
          on:keydown={(e) => e.key === 'Enter' && connectGitHub()}
        />
        <button class="connect-btn" on:click={connectGitHub} disabled={loading || !githubInput.trim()}>
          Connect
        </button>
      </div>
    </div>

    {#if recentRepos.length > 0}
      <div class="recent-section">
        <h3>Recent repositories</h3>
        <div class="recent-list">
          {#each recentRepos.slice(0, 5) as repo}
            <button class="recent-item" on:click={() => selectRepo(repo.path, repo.type)}>
              <span class="recent-icon">
                {#if repo.type === 'github'}
                  <svg viewBox="0 0 16 16" fill="currentColor" width="14" height="14"><path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/></svg>
                {:else}
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" /></svg>
                {/if}
              </span>
              <span class="recent-path">{shortenPath(repo.path)}</span>
              <span class="recent-date">{formatDate(repo.lastUsed)}</span>
            </button>
          {/each}
        </div>
      </div>
    {/if}
  </div>
{:else}
  <div class="repo-active">
    {#if repoSummary}
      <h1>{repoSummary.name}</h1>
      <div class="stats">
        <div class="stat">
          <span class="stat-value">{repoSummary.totalCommits}</span>
          <span class="stat-label">Commits</span>
        </div>
        <div class="stat">
          <span class="stat-value">{repoSummary.contributors}</span>
          <span class="stat-label">Contributors</span>
        </div>
      </div>
      <div class="last-commit">
        <span class="label">Last commit:</span>
        <span class="value">{repoSummary.lastCommit}</span>
      </div>
    {:else}
      <p class="loading-text">Loading repo info...</p>
    {/if}
    <button class="change-repo" on:click={() => { activeRepo = ''; repoSummary = null; dispatch('repoSelected', '') }}>
      Change repository
    </button>
  </div>
{/if}

<style>
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    gap: 24px;
  }

  h1 { color: #e6edf3; font-size: 28px; margin: 0; }
  .subtitle { color: #8b949e; margin: 0; }

  .entry-points {
    display: flex;
    gap: 16px;
    align-items: stretch;
    flex-wrap: wrap;
    justify-content: center;
    margin-top: 16px;
  }

  .entry-btn {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 24px 32px;
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    color: #e6edf3;
    cursor: pointer;
    transition: border-color 0.2s, background 0.2s;
    font-size: 14px;
    font-family: inherit;
  }
  .entry-btn:hover { border-color: #58a6ff; background: #1c2333; }
  .entry-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .entry-icon { display: flex; align-items: center; color: #58a6ff; }

  .drop-zone {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 24px 32px;
    border: 2px dashed #30363d;
    border-radius: 8px;
    color: #8b949e;
    transition: border-color 0.2s, background 0.2s;
    min-width: 160px;
    cursor: default;
  }
  .drop-zone:hover { border-color: #58a6ff; background: #1c233322; }

  .github-input {
    display: flex;
    gap: 8px;
    align-items: center;
  }
  .github-input input {
    padding: 10px 14px;
    background: #0d1117;
    border: 1px solid #30363d;
    border-radius: 6px;
    color: #e6edf3;
    font-size: 14px;
    width: 180px;
    font-family: 'JetBrains Mono', monospace;
  }
  .github-input input:focus { outline: none; border-color: #58a6ff; }
  .connect-btn {
    padding: 10px 16px;
    background: #238636;
    border: none;
    border-radius: 6px;
    color: #fff;
    cursor: pointer;
    font-size: 14px;
    font-family: inherit;
  }
  .connect-btn:hover { background: #2ea043; }
  .connect-btn:disabled { opacity: 0.5; cursor: not-allowed; }

  .recent-section { margin-top: 24px; width: 100%; max-width: 500px; }
  .recent-section h3 { color: #8b949e; font-size: 13px; text-transform: uppercase; margin-bottom: 8px; }

  .recent-list { display: flex; flex-direction: column; gap: 4px; }
  .recent-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 12px;
    background: transparent;
    border: 1px solid transparent;
    border-radius: 6px;
    color: #e6edf3;
    cursor: pointer;
    text-align: left;
    width: 100%;
    font-family: inherit;
    font-size: 13px;
    transition: background 0.15s;
  }
  .recent-item:hover { background: #161b22; border-color: #30363d; }
  .recent-path { flex: 1; font-family: 'JetBrains Mono', monospace; font-size: 12px; }
  .recent-date { color: #8b949e; font-size: 11px; }

  .banner.error {
    background: #da363433;
    border: 1px solid #da3634;
    color: #f85149;
    padding: 8px 12px;
    border-radius: 6px;
    margin-bottom: 16px;
    font-size: 13px;
  }

  .repo-active {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }
  .stats {
    display: flex;
    gap: 24px;
  }
  .stat {
    display: flex;
    flex-direction: column;
    background: #161b22;
    border: 1px solid #30363d;
    border-radius: 8px;
    padding: 16px 24px;
    min-width: 120px;
  }
  .stat-value { font-size: 28px; font-weight: 700; color: #58a6ff; }
  .stat-label { font-size: 12px; color: #8b949e; text-transform: uppercase; }
  .last-commit { color: #8b949e; font-size: 13px; }
  .last-commit .label { margin-right: 6px; }
  .last-commit .value { color: #e6edf3; font-family: 'JetBrains Mono', monospace; }
  .change-repo {
    align-self: flex-start;
    padding: 6px 12px;
    background: transparent;
    border: 1px solid #30363d;
    border-radius: 6px;
    color: #8b949e;
    cursor: pointer;
    font-size: 13px;
    font-family: inherit;
  }
  .change-repo:hover { border-color: #58a6ff; color: #58a6ff; }
  .loading-text { color: #8b949e; }
</style>
