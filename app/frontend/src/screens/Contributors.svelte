<script lang="ts">
  import { onDestroy } from 'svelte'
  import { Contributors } from '../../bindings/github.com/alciller88/commitlore/app/gitapp.js'
  import { activeRepo } from '../lib/store'
  import type { ActiveRepo } from '../lib/store'

  let currentRepo: ActiveRepo | null = null
  const unsub = activeRepo.subscribe(v => { currentRepo = v })
  onDestroy(() => unsub())

  let since = ''
  let top = 10
  let loading = false
  let error = ''
  let contribs: Array<{name: string, email: string, commits: number}> = []

  $: maxCommits = contribs.length > 0 ? contribs[0].commits : 1

  async function fetchContribs() {
    if (!currentRepo) return
    loading = true
    error = ''
    contribs = []
    try {
      const raw = await Contributors(currentRepo.path, since, top)
      contribs = JSON.parse(raw)
    } catch (e: any) {
      error = e?.message || 'Failed to fetch contributors'
    } finally {
      loading = false
    }
  }

  function barWidth(commits: number) {
    return Math.max(4, (commits / maxCommits) * 100)
  }
</script>

{#if !currentRepo}
  <div class="no-repo">
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="32" height="32"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
    <p>Select a repository in Dashboard to get started.</p>
  </div>
{:else}
  <div style="display:flex; flex-direction:column; height:100%; overflow:hidden;">
    {#if error}
      <div class="banner error" style="flex-shrink:0;">{error}</div>
    {/if}

    <div class="filters-row" style="flex-shrink:0;">
      <input class="since-input" type="text" bind:value={since} placeholder="Since (YYYY-MM-DD)" />
      <div class="top-group">
        <span class="top-label">Top</span>
        <input type="range" min="5" max="50" bind:value={top} />
        <span class="top-value">{top}</span>
      </div>
      <button class="fetch-btn" on:click={fetchContribs} disabled={loading}>
        {#if loading}
          <span class="spinner"></span>
        {:else}
          Fetch
        {/if}
      </button>
    </div>

    {#if currentRepo.type === 'github' && contribs.length > 0}
      <p class="remote-note" style="flex-shrink:0;">Note: Top files not available for remote repositories.</p>
    {/if}

    {#if contribs.length > 0}
      <div style="flex:1; min-height:0; overflow-y:auto;">
        <table style="width:100%; border-collapse:collapse;">
          <thead>
            <tr>
              <th class="col-contributor">Contributor</th>
              <th class="col-commits">Commits</th>
              <th class="col-activity">Activity</th>
            </tr>
          </thead>
          <tbody>
            {#each contribs as c}
              <tr>
                <td class="cell-contributor">
                  <span class="avatar">{c.name.charAt(0).toUpperCase()}</span>
                  <div class="contributor-info">
                    <span class="contributor-name">{c.name}</span>
                    <span class="contributor-email">{c.email}</span>
                  </div>
                </td>
                <td class="cell-commits">{c.commits}</td>
                <td class="cell-activity">
                  <div class="bar-track">
                    <div class="bar-fill" style="width: {barWidth(c.commits)}%"></div>
                  </div>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {/if}
  </div>
{/if}

<style>
  .no-repo {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    flex: 1;
    gap: var(--space-3);
    color: var(--cl-secondary);
  }
  .no-repo p {
    margin: 0;
    font-size: var(--text-md);
  }

  .filters-row {
    display: flex;
    gap: var(--space-2);
    align-items: center;
    height: 40px;
    flex-shrink: 0;
  }

  .since-input {
    flex: 1;
    max-width: 200px;
    height: 36px;
    box-sizing: border-box;
    padding: 0 var(--space-2);
    background: var(--cl-background);
    border: 1px solid var(--cl-border);
    border-radius: var(--radius-md);
    color: var(--cl-text);
    font-size: var(--text-base);
    font-family: 'JetBrains Mono', monospace;
    transition: border-color var(--transition-fast);
  }
  .since-input:focus {
    outline: none;
    border-color: var(--cl-accent);
  }

  .top-group {
    display: flex;
    align-items: center;
    gap: var(--space-1);
  }
  .top-label {
    font-size: var(--text-xs);
    color: var(--cl-secondary);
    text-transform: uppercase;
  }
  .top-group input[type="range"] {
    width: 80px;
    -webkit-appearance: none;
    appearance: none;
    height: 4px;
    background: var(--cl-border);
    border-radius: var(--radius-sm);
    outline: none;
  }
  .top-group input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 12px;
    height: 12px;
    background: var(--cl-accent);
    border-radius: 50%;
    cursor: pointer;
  }
  .top-value {
    font-size: var(--text-xs);
    color: var(--cl-secondary);
    min-width: 20px;
    text-align: right;
  }

  .fetch-btn {
    width: 80px;
    height: 36px;
    background: var(--cl-accent);
    border: none;
    border-radius: var(--radius-md);
    color: #fff;
    font-size: var(--text-base);
    font-family: var(--cl-font-family);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: var(--space-1);
    flex-shrink: 0;
    transition: opacity var(--transition-fast);
  }
  .fetch-btn:hover {
    opacity: 0.9;
  }
  .fetch-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .spinner {
    display: inline-block;
    width: 14px;
    height: 14px;
    border: 2px solid #ffffff44;
    border-top-color: #fff;
    border-radius: 50%;
    animation: spin 0.6s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .remote-note {
    margin: 0;
    font-size: var(--text-xs);
    color: var(--cl-secondary);
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  thead {
    position: sticky;
    top: 0;
  }

  th {
    text-align: left;
    padding: 0 var(--space-3);
    height: 28px;
    font-size: var(--text-xs);
    text-transform: uppercase;
    color: var(--cl-secondary);
    opacity: 0.6;
    font-weight: 500;
    border-bottom: 1px solid var(--cl-border);
    background: var(--cl-background);
  }

  .col-commits {
    text-align: right;
    width: 80px;
  }
  .col-activity {
    width: 160px;
  }

  tr {
    height: 44px;
    transition: background var(--transition-fast);
  }
  tbody tr:hover {
    background: var(--cl-surface);
  }

  td {
    padding: 0 var(--space-3);
    border-bottom: 1px solid var(--cl-border);
  }

  .cell-contributor {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    height: 44px;
  }

  .avatar {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    background: var(--cl-accent);
    color: #fff;
    font-size: var(--text-xs);
    font-weight: 600;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .contributor-info {
    display: flex;
    flex-direction: column;
    gap: 1px;
    min-width: 0;
  }

  .contributor-name {
    font-size: var(--text-base);
    font-weight: 500;
    color: var(--cl-text);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .contributor-email {
    font-size: var(--text-xs);
    color: var(--cl-secondary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .cell-commits {
    text-align: right;
    font-size: var(--text-base);
    color: var(--cl-text);
    font-family: 'JetBrains Mono', monospace;
  }

  .cell-activity {
    padding-right: var(--space-4);
  }

  .bar-track {
    width: 100%;
    height: 4px;
    background: var(--cl-surface);
    border-radius: var(--radius-sm);
    overflow: hidden;
  }

  .bar-fill {
    height: 100%;
    background: var(--cl-accent);
    border-radius: var(--radius-sm);
    transition: width var(--transition-base);
  }

  .banner.error {
    background: #da363433;
    border: 1px solid #da3634;
    color: #f85149;
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-md);
    font-size: var(--text-base);
  }
</style>
