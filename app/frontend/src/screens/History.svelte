<script lang="ts">
  import { onDestroy } from 'svelte'
  import { History } from '../../bindings/github.com/alciller88/commitlore/app/gitapp.js'
  import { activeRepo } from '../lib/store'
  import type { ActiveRepo } from '../lib/store'

  let currentRepo: ActiveRepo | null = null
  const unsub = activeRepo.subscribe(v => { currentRepo = v })
  onDestroy(() => unsub())

  let author = ''
  let since = ''
  let until = ''
  let limit = 50
  let loading = false
  let error = ''
  let commits: Array<{hash: string, date: string, author: string, message: string}> = []

  async function fetchHistory() {
    if (!currentRepo) return
    loading = true
    error = ''
    commits = []
    try {
      const raw = await History(currentRepo.path, author, since, until, limit)
      commits = JSON.parse(raw)
    } catch (e: any) {
      error = e?.message || 'Failed to fetch history'
    } finally {
      loading = false
    }
  }

  function shortHash(h: string) {
    return h ? h.substring(0, 7) : ''
  }

  function copyHash(hash: string) {
    navigator.clipboard.writeText(hash)
  }

  function formatDate(d: string) {
    try {
      return new Date(d).toLocaleDateString()
    } catch {
      return d
    }
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
      <input class="author-input" type="text" bind:value={author} placeholder="Author" />
      <input class="date-input" type="text" bind:value={since} placeholder="Since" />
      <input class="date-input" type="text" bind:value={until} placeholder="Until" />
      <div class="limit-group">
        <input type="range" min="10" max="200" bind:value={limit} />
        <span class="limit-value">{limit}</span>
      </div>
      <button class="fetch-btn" on:click={fetchHistory} disabled={loading}>
        {#if loading}
          <span class="spinner"></span>
        {:else}
          Fetch
        {/if}
      </button>
    </div>

    {#if commits.length > 0}
      <div style="flex:1; min-height:0; overflow-y:auto;">
        <table style="width:100%; border-collapse:collapse;">
          <thead>
            <tr>
              <th class="col-hash">Hash</th>
              <th class="col-date">Date</th>
              <th class="col-author">Author</th>
              <th class="col-message">Message</th>
            </tr>
          </thead>
          <tbody>
            {#each commits as commit}
              <tr>
                <td class="cell-hash" on:click={() => copyHash(commit.hash)} title="Click to copy full SHA">
                  {shortHash(commit.hash)}
                </td>
                <td class="cell-date">{formatDate(commit.date)}</td>
                <td class="cell-author">{commit.author}</td>
                <td class="cell-message">{commit.message}</td>
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
    height: 36px;
  }

  .author-input {
    flex: 1;
    min-width: 100px;
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
  .author-input:focus {
    outline: none;
    border-color: var(--cl-accent);
  }

  .date-input {
    width: 100px;
    flex-shrink: 0;
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
  .date-input:focus {
    outline: none;
    border-color: var(--cl-accent);
  }

  .limit-group {
    display: flex;
    align-items: center;
    gap: var(--space-1);
    width: 120px;
    flex-shrink: 0;
  }
  .limit-group input[type="range"] {
    flex: 1;
    -webkit-appearance: none;
    appearance: none;
    height: 4px;
    background: var(--cl-border);
    border-radius: var(--radius-sm);
    outline: none;
  }
  .limit-group input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 12px;
    height: 12px;
    background: var(--cl-accent);
    border-radius: 50%;
    cursor: pointer;
  }
  .limit-value {
    font-size: var(--text-xs);
    color: var(--cl-secondary);
    min-width: 24px;
    text-align: right;
  }

  .fetch-btn {
    width: 72px;
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

  tr {
    height: 36px;
    transition: background var(--transition-fast);
  }
  tbody tr:nth-child(odd) {
    background: color-mix(in srgb, var(--cl-surface) 30%, transparent);
  }
  tbody tr:hover {
    background: var(--cl-surface);
  }

  td {
    padding: 0 var(--space-3);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .col-hash { width: 80px; }
  .col-date { width: 100px; }
  .col-author { width: 140px; }
  .col-message { }

  .cell-hash {
    font-family: 'JetBrains Mono', monospace;
    font-size: var(--text-sm);
    color: var(--cl-accent);
    cursor: pointer;
  }
  .cell-hash:hover {
    text-decoration: underline;
  }

  .cell-date {
    font-size: var(--text-xs);
    color: var(--cl-secondary);
  }

  .cell-author {
    font-size: var(--text-sm);
    color: var(--cl-text);
  }

  .cell-message {
    font-size: var(--text-base);
    color: var(--cl-text);
    max-width: 400px;
    overflow: hidden;
    text-overflow: ellipsis;
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
