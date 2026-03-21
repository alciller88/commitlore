<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { OpenFolderPicker, History } from '../../bindings/github.com/alciller88/commitlore/app/gitapp.js'

  export let activeRepo = ''

  const dispatch = createEventDispatcher()

  let repo = ''
  let author = ''
  let since = ''
  let until = ''
  let limit = 50
  let loading = false
  let error = ''
  let commits: Array<{hash: string, date: string, author: string, message: string}> = []

  $: repo = activeRepo || repo

  async function pickRepo() {
    const path = await OpenFolderPicker()
    if (path) {
      repo = path
      dispatch('repoSelected', path)
    }
  }

  async function fetchHistory() {
    if (!repo) return
    loading = true
    error = ''
    commits = []
    try {
      const raw = await History(repo, author, since, until, limit)
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

<div class="screen">
  <h1>Commit History</h1>

  {#if error}
    <div class="banner error">{error}</div>
  {/if}

  <div class="form">
    <div class="field">
      <label>Repository</label>
      <div class="repo-input">
        <input type="text" bind:value={repo} placeholder="Local path or owner/repo" />
        <button class="pick-btn" on:click={pickRepo}>Browse</button>
      </div>
    </div>

    <div class="row">
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
    </div>

    <div class="row">
      <div class="field limit-field">
        <label>Limit: {limit}</label>
        <input type="range" min="10" max="200" bind:value={limit} />
      </div>
      <button class="action-btn" on:click={fetchHistory} disabled={loading || !repo}>
        {#if loading}
          <span class="spinner"></span> Loading...
        {:else}
          Fetch
        {/if}
      </button>
    </div>
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

<style>
  .screen { display: flex; flex-direction: column; gap: 20px; height: 100%; }
  h1 { color: #e6edf3; font-size: 22px; margin: 0; }
  .form { display: flex; flex-direction: column; gap: 14px; }
  .field { display: flex; flex-direction: column; gap: 4px; flex: 1; }
  .field label { color: #8b949e; font-size: 12px; text-transform: uppercase; }
  .row { display: flex; gap: 12px; align-items: flex-end; }

  .limit-field { max-width: 200px; }
  input[type="range"] {
    -webkit-appearance: none; appearance: none; height: 6px;
    background: #30363d; border-radius: 3px; outline: none;
  }
  input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none; appearance: none; width: 16px; height: 16px;
    background: #58a6ff; border-radius: 50%; cursor: pointer;
  }

  input[type="text"], select {
    padding: 8px 12px; background: #0d1117; border: 1px solid #30363d;
    border-radius: 6px; color: #e6edf3; font-size: 14px;
    font-family: 'JetBrains Mono', monospace;
  }
  input[type="text"]:focus { outline: none; border-color: #58a6ff; }

  .repo-input { display: flex; gap: 8px; }
  .repo-input input { flex: 1; }
  .pick-btn {
    padding: 8px 14px; background: #21262d; border: 1px solid #30363d;
    border-radius: 6px; color: #e6edf3; cursor: pointer; font-family: inherit; font-size: 13px;
  }
  .pick-btn:hover { border-color: #58a6ff; }

  .action-btn {
    padding: 10px 20px; background: #238636; border: none; border-radius: 6px;
    color: #fff; font-size: 14px; cursor: pointer; display: flex;
    align-items: center; gap: 8px; font-family: inherit; height: fit-content;
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
