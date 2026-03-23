<script lang="ts">
  import { onMount } from 'svelte'
  import { FetchCatalog, InstallStyle, IsInstalled } from '../../bindings/github.com/alciller88/commitlore/app/marketplaceapp.js'

  type CatalogEntry = {
    name: string
    description: string
    author: string
    version: string
    tags: string[]
    preview: string
    download: string
  }

  let entries: CatalogEntry[] = []
  let loading = true
  let error = ''
  let installedMap: Record<string, boolean> = {}
  let installingMap: Record<string, boolean> = {}
  let installErrors: Record<string, string> = {}

  onMount(() => { loadCatalog() })

  async function loadCatalog() {
    loading = true
    error = ''
    try {
      entries = await FetchCatalog()
      await checkInstalled()
    } catch (e: any) {
      error = e?.message || 'Failed to load catalog. Check your internet connection.'
    } finally {
      loading = false
    }
  }

  async function checkInstalled() {
    const map: Record<string, boolean> = {}
    for (const entry of entries) {
      try {
        map[entry.name] = await IsInstalled(entry.name)
      } catch {
        map[entry.name] = false
      }
    }
    installedMap = map
  }

  async function install(entry: CatalogEntry) {
    installingMap = { ...installingMap, [entry.name]: true }
    installErrors = { ...installErrors, [entry.name]: '' }
    try {
      await InstallStyle(entry.download, entry.name)
      installedMap = { ...installedMap, [entry.name]: true }
    } catch (e: any) {
      installErrors = { ...installErrors, [entry.name]: e?.message || 'Install failed' }
    } finally {
      installingMap = { ...installingMap, [entry.name]: false }
    }
  }
</script>

<div class="marketplace">
  <div class="header">
    <h1>Marketplace</h1>
    <p class="subtitle">Browse and install community styles</p>
  </div>

  {#if loading}
    <div class="state-container">
      <div class="spinner"></div>
      <p class="state-text">Loading catalog...</p>
    </div>
  {:else if error}
    <div class="state-container">
      <div class="error-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="48" height="48">
          <circle cx="12" cy="12" r="10"/>
          <line x1="15" y1="9" x2="9" y2="15"/>
          <line x1="9" y1="9" x2="15" y2="15"/>
        </svg>
      </div>
      <p class="state-text error-text">{error}</p>
      <button class="retry-btn" on:click={loadCatalog}>Retry</button>
    </div>
  {:else}
    <div class="grid">
      {#each entries as entry}
        <div class="card">
          {#if entry.preview}
            <div class="card-preview">
              <img src={entry.preview} alt="{entry.name} preview" loading="lazy" />
            </div>
          {/if}
          <div class="card-body">
            <div class="card-header">
              <span class="card-name">{entry.name}</span>
              <span class="card-version">v{entry.version}</span>
            </div>
            <div class="card-author">by {entry.author}</div>
            <p class="card-desc">{entry.description}</p>
            {#if entry.tags && entry.tags.length > 0}
              <div class="card-tags">
                {#each entry.tags as tag}
                  <span class="tag">{tag}</span>
                {/each}
              </div>
            {/if}
            {#if installErrors[entry.name]}
              <div class="install-error">{installErrors[entry.name]}</div>
            {/if}
            <div class="card-actions">
              {#if installedMap[entry.name]}
                <button class="install-btn installed" disabled>Installed</button>
              {:else if installingMap[entry.name]}
                <button class="install-btn installing" disabled>
                  <span class="btn-spinner"></span>
                  Installing...
                </button>
              {:else}
                <button class="install-btn" on:click={() => install(entry)}>Install</button>
              {/if}
            </div>
          </div>
        </div>
      {/each}
    </div>

    {#if entries.length === 0}
      <div class="state-container">
        <p class="state-text">No styles available yet.</p>
      </div>
    {/if}
  {/if}
</div>

<style>
  .marketplace {
    display: flex;
    flex-direction: column;
    gap: 16px;
    flex: 1;
    min-height: 0;
    overflow-y: auto;
  }

  .header { display: flex; flex-direction: column; gap: 4px; }
  h1 { color: var(--cl-text, #e6edf3); font-size: 20px; margin: 0; }
  .subtitle { color: var(--cl-secondary, #8b949e); font-size: 13px; margin: 0; }

  .state-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    min-height: 200px;
  }

  .state-text {
    color: var(--cl-secondary, #8b949e);
    font-size: 14px;
    margin: 0;
    text-align: center;
    max-width: 400px;
  }
  .error-text { color: var(--cl-text, #e6edf3); }

  .error-icon { color: var(--cl-secondary, #8b949e); opacity: 0.5; }

  .spinner {
    width: 24px;
    height: 24px;
    border: 2px solid var(--cl-border, #30363d);
    border-top-color: var(--cl-accent, #58a6ff);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  .retry-btn {
    padding: 8px 20px;
    background: var(--cl-surface, #161b22);
    border: 1px solid var(--cl-border, #30363d);
    border-radius: 6px;
    color: var(--cl-text, #e6edf3);
    font-size: 13px;
    font-family: var(--cl-font-family, system-ui, sans-serif);
    cursor: pointer;
  }
  .retry-btn:hover {
    border-color: var(--cl-accent, #58a6ff);
    color: var(--cl-accent, #58a6ff);
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
  }

  .card {
    display: flex;
    flex-direction: column;
    background: var(--cl-surface, #161b22);
    border: 1px solid var(--cl-border, #30363d);
    border-radius: 8px;
    overflow: hidden;
  }

  .card-preview {
    width: 100%;
    aspect-ratio: 16 / 9;
    overflow: hidden;
    background: var(--cl-background, #0d1117);
  }
  .card-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .card-body {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 12px;
    flex: 1;
  }

  .card-header {
    display: flex;
    align-items: baseline;
    gap: 8px;
  }
  .card-name {
    font-size: 14px;
    font-weight: 600;
    color: var(--cl-text, #e6edf3);
  }
  .card-version {
    font-size: 11px;
    padding: 1px 6px;
    border-radius: 8px;
    background: var(--cl-border, #30363d);
    color: var(--cl-secondary, #8b949e);
    font-family: 'JetBrains Mono', monospace;
  }

  .card-author {
    font-size: 12px;
    color: var(--cl-secondary, #8b949e);
  }

  .card-desc {
    font-size: 13px;
    color: var(--cl-text, #e6edf3);
    line-height: 1.4;
    margin: 0;
    opacity: 0.85;
  }

  .card-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin-top: 2px;
  }
  .tag {
    font-size: 10px;
    padding: 2px 8px;
    border-radius: 10px;
    background: var(--cl-background, #0d1117);
    color: var(--cl-secondary, #8b949e);
    border: 1px solid var(--cl-border, #30363d);
  }

  .install-error {
    font-size: 12px;
    color: var(--cl-text, #e6edf3);
    background: var(--cl-background, #0d1117);
    border: 1px solid var(--cl-border, #30363d);
    border-radius: 4px;
    padding: 6px 8px;
  }

  .card-actions {
    margin-top: auto;
    padding-top: 8px;
  }

  .install-btn {
    width: 100%;
    padding: 8px;
    border-radius: 6px;
    font-size: 12px;
    font-weight: 500;
    font-family: var(--cl-font-family, system-ui, sans-serif);
    cursor: pointer;
    border: 1px solid var(--cl-accent, #58a6ff);
    background: transparent;
    color: var(--cl-accent, #58a6ff);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
  }
  .install-btn:hover {
    background: var(--cl-accent, #58a6ff);
    color: var(--cl-background, #0d1117);
  }
  .install-btn.installed {
    border-color: var(--cl-border, #30363d);
    color: var(--cl-secondary, #8b949e);
    cursor: default;
    background: transparent;
  }
  .install-btn.installing {
    border-color: var(--cl-border, #30363d);
    color: var(--cl-secondary, #8b949e);
    cursor: wait;
    background: transparent;
  }

  .btn-spinner {
    width: 12px;
    height: 12px;
    border: 2px solid var(--cl-border, #30363d);
    border-top-color: var(--cl-secondary, #8b949e);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    display: inline-block;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
