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
  <div class="screen">
    <h1>Commit History</h1>

    {#if error}
      <div class="banner error">{error}</div>
    {/if}

    <div class="filters-row">
      <div class="field">
        <label>Author</label>
        <input type="text" bind:value={author} placeholder="name or email" />
      </div>
      <div class="field">
        <label>Since</label>
        <input type="text" bind:value={since} placeholder="YYYY-MM-DD" />
      </div>
      <div class="field">
        <label>Until</label>
        <input type="text" bind:value={until} placeholder="YYYY-MM-DD" />
      </div>
      <div class="field limit-field">
        <label>Limit: {limit}</label>
        <input type="range" min="10" max="200" bind:value={limit} />
      </div>
      <button class="action-btn" on:click={fetchHistory} disabled={loading}>
        {#if loading}
          <span class="spinner"></span>
        {:else}
          Fetch
        {/if}
      </button>
    </div>

    {#if commits.length > 0}
      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>Hash</th>
              <th>Date</th>
              <th>Author</th>
              <th>Message</th>
            </tr>
          </thead>
          <tbody>
            {#each commits as commit}
              <tr>
                <td>
                  <button class="hash-btn" on:click={() => copyHash(commit.hash)} title="Click to copy full SHA">
                    {shortHash(commit.hash)}
                  </button>
                </td>
                <td class="date">{formatDate(commit.date)}</td>
                <td class="author">{commit.author}</td>
                <td class="message">{commit.message}</td>
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
    height: 100%;
    gap: 12px;
    color: #484f58;
  }
  .no-repo p { margin: 0; font-size: 14px; }

  .screen { display: flex; flex-direction: column; gap: 16px; height: 100%; }
  h1 { color: #e6edf3; font-size: 22px; margin: 0; }

  .filters-row {
    display: flex;
    gap: 10px;
    align-items: flex-end;
    flex-wrap: wrap;
  }

  .field { display: flex; flex-direction: column; gap: 4px; flex: 1; min-width: 100px; }
  .field label { color: #8b949e; font-size: 11px; text-transform: uppercase; }
  .limit-field { max-width: 160px; }

  input[type="text"] {
    padding: 7px 10px; background: #0d1117; border: 1px solid #30363d;
    border-radius: 6px; color: #e6edf3; font-size: 13px;
    font-family: 'JetBrains Mono', monospace;
  }
  input[type="text"]:focus { outline: none; border-color: #58a6ff; }

  input[type="range"] {
    -webkit-appearance: none; appearance: none; height: 6px;
    background: #30363d; border-radius: 3px; outline: none;
  }
  input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none; appearance: none; width: 14px; height: 14px;
    background: #58a6ff; border-radius: 50%; cursor: pointer;
  }

  .action-btn {
    padding: 8px 16px; background: #238636; border: none; border-radius: 6px;
    color: #fff; font-size: 13px; cursor: pointer; display: flex;
    align-items: center; gap: 6px; font-family: inherit; height: fit-content;
    flex-shrink: 0;
  }
  .action-btn:hover { background: #2ea043; }
  .action-btn:disabled { opacity: 0.5; cursor: not-allowed; }

  .spinner {
    display: inline-block; width: 14px; height: 14px;
    border: 2px solid #ffffff44; border-top-color: #fff;
    border-radius: 50%; animation: spin 0.6s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .table-container { overflow: auto; flex: 1; }
  table { width: 100%; border-collapse: collapse; font-size: 13px; }
  thead { position: sticky; top: 0; }
  th {
    text-align: left; padding: 8px 12px; color: #8b949e;
    border-bottom: 1px solid #30363d; background: #0d1117;
    font-size: 11px; text-transform: uppercase;
  }
  td {
    padding: 6px 12px; border-bottom: 1px solid #21262d;
    white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 400px;
  }
  .hash-btn {
    background: none; border: none; color: #58a6ff; cursor: pointer;
    font-family: 'JetBrains Mono', monospace; font-size: 13px; padding: 0;
  }
  .hash-btn:hover { text-decoration: underline; }
  .date { color: #8b949e; }
  .author { color: #e6edf3; }
  .message { color: #e6edf3; font-family: 'JetBrains Mono', monospace; }

  .banner.error {
    background: #da363433; border: 1px solid #da3634; color: #f85149;
    padding: 8px 12px; border-radius: 6px; font-size: 13px;
  }
</style>
