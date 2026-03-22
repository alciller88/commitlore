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
  <div class="screen">
    <h1>Contributors</h1>

    {#if error}
      <div class="banner error">{error}</div>
    {/if}

    <div class="filters-row">
      <div class="field">
        <label>Since</label>
        <input type="text" bind:value={since} placeholder="YYYY-MM-DD" />
      </div>
      <div class="field limit-field">
        <label>Top: {top}</label>
        <input type="range" min="5" max="50" bind:value={top} />
      </div>
      <button class="action-btn" on:click={fetchContribs} disabled={loading}>
        {#if loading}
          <span class="spinner"></span>
        {:else}
          Fetch
        {/if}
      </button>
    </div>

    {#if currentRepo.type === 'github' && contribs.length > 0}
      <div class="note">
        <span class="info">Note: Top files not available for remote repositories.</span>
      </div>
    {/if}

    {#if contribs.length > 0}
      <div class="table-container">
        <table>
          <thead>
            <tr>
              <th>Contributor</th>
              <th>Commits</th>
              <th>Activity</th>
            </tr>
          </thead>
          <tbody>
            {#each contribs as c}
              <tr>
                <td>
                  <div class="contributor-name">{c.name}</div>
                  <div class="contributor-email">{c.email}</div>
                </td>
                <td class="commits-cell">{c.commits}</td>
                <td>
                  <div class="bar-container">
                    <div class="bar" style="width: {barWidth(c.commits)}%"></div>
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

  .note { font-size: 12px; }
  .info { color: #8b949e; }

  .table-container { overflow: auto; flex: 1; }
  table { width: 100%; border-collapse: collapse; font-size: 13px; }
  thead { position: sticky; top: 0; }
  th {
    text-align: left; padding: 8px 12px; color: #8b949e;
    border-bottom: 1px solid #30363d; background: #0d1117;
    font-size: 11px; text-transform: uppercase;
  }
  td { padding: 8px 12px; border-bottom: 1px solid #21262d; }
  .contributor-name { color: #e6edf3; }
  .contributor-email { color: #8b949e; font-size: 11px; }
  .commits-cell { color: #e6edf3; font-family: 'JetBrains Mono', monospace; text-align: center; }

  .bar-container {
    background: #21262d; border-radius: 3px; height: 8px; width: 100%; min-width: 80px;
  }
  .bar {
    height: 100%; background: #58a6ff; border-radius: 3px;
    transition: width 0.3s ease;
  }

  .banner.error {
    background: #da363433; border: 1px solid #da3634; color: #f85149;
    padding: 8px 12px; border-radius: 6px; font-size: 13px;
  }
</style>
