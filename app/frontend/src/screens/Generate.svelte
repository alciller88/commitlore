<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { Generate } from '../../bindings/github.com/alciller88/commitlore/app/changelogapp.js'
  import { ListStyles } from '../../bindings/github.com/alciller88/commitlore/app/styleapp.js'
  import { GetLLMConfig } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { activeRepo } from '../lib/store'
  import type { ActiveRepo } from '../lib/store'

  let currentRepo: ActiveRepo | null = null
  const unsub = activeRepo.subscribe(v => { currentRepo = v })
  onDestroy(() => unsub())

  let since = ''
  let until = ''
  let styleName = 'formal'
  let styles: Array<{name: string}> = []
  let loading = false
  let error = ''
  let htmlResult = ''
  let llmStatus = ''

  onMount(() => {
    loadStyles()
    loadLLMStatus()
  })

  async function loadStyles() {
    try {
      const raw = await ListStyles()
      styles = JSON.parse(raw)
    } catch {}
  }

  async function loadLLMStatus() {
    try {
      const cfg = await GetLLMConfig()
      if (cfg.provider && cfg.keyConfigured) {
        llmStatus = `${cfg.provider} (${cfg.model || 'default'})`
      } else if (cfg.provider) {
        llmStatus = `${cfg.provider} — key not configured`
      } else {
        llmStatus = 'not configured'
      }
    } catch {
      llmStatus = 'not configured'
    }
  }

  async function generate() {
    if (!currentRepo) return
    loading = true
    error = ''
    htmlResult = ''
    try {
      htmlResult = await Generate(currentRepo.path, since, until, styleName)
    } catch (e: any) {
      error = e?.message || 'Generation failed'
    } finally {
      loading = false
    }
  }

  function copyText() {
    const tmp = document.createElement('div')
    tmp.innerHTML = htmlResult
    navigator.clipboard.writeText(tmp.innerText)
  }

  function saveAsFile() {
    const blob = new Blob([htmlResult], { type: 'text/html' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'changelog.html'
    a.click()
    URL.revokeObjectURL(url)
  }
</script>

{#if !currentRepo}
  <div class="no-repo">
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="32" height="32"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
    <p>Select a repository in Dashboard to get started.</p>
  </div>
{:else}
  <div class="two-col">
    <div class="sidebar-form">
      <h2>Generate Changelog</h2>

      {#if error}
        <div class="banner error">{error}</div>
      {/if}

      <div class="field">
        <label>Since</label>
        <input type="text" bind:value={since} placeholder="tag, SHA, or YYYY-MM-DD" />
      </div>
      <div class="field">
        <label>Until</label>
        <input type="text" bind:value={until} placeholder="default: HEAD" />
      </div>
      <div class="field">
        <label>Style</label>
        <select bind:value={styleName}>
          {#each styles as s}
            <option value={s.name}>{s.name}</option>
          {/each}
        </select>
      </div>

      <div class="llm-indicator">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="12" height="12"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/></svg>
        <span>LLM: {llmStatus}</span>
      </div>

      <button class="action-btn" on:click={generate} disabled={loading}>
        {#if loading}
          <span class="spinner"></span> Generating...
        {:else}
          Generate
        {/if}
      </button>
    </div>

    <div class="main-panel">
      {#if htmlResult}
        <div class="result-toolbar">
          <button class="tool-btn" on:click={copyText}>Copy text</button>
          <button class="tool-btn" on:click={saveAsFile}>Save as .html</button>
        </div>
        <iframe class="preview" srcdoc={htmlResult} sandbox="" title="Changelog preview"></iframe>
      {:else}
        <div class="empty-preview">
          <span class="empty-text">Preview will appear here</span>
        </div>
      {/if}
    </div>
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
    color: var(--cl-secondary, #8b949e);
  }
  .no-repo p { margin: 0; font-size: 14px; }

  .two-col {
    display: flex;
    gap: 16px;
    height: 100%;
  }

  .sidebar-form {
    width: 280px;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  h2 { color: var(--cl-text, #e6edf3); font-size: 18px; margin: 0 0 4px 0; }

  .field { display: flex; flex-direction: column; gap: 4px; }
  .field label { color: var(--cl-secondary, #8b949e); font-size: 11px; text-transform: uppercase; }

  input, select {
    padding: 7px 10px; background: var(--cl-background, #0d1117); border: 1px solid var(--cl-border, #30363d);
    border-radius: 6px; color: var(--cl-text, #e6edf3); font-size: 13px;
    font-family: 'JetBrains Mono', monospace;
  }
  input:focus, select:focus { outline: none; border-color: var(--cl-accent, #58a6ff); }
  select { cursor: pointer; }

  .llm-indicator {
    display: flex; align-items: center; gap: 6px;
    color: var(--cl-secondary, #8b949e); font-size: 11px;
  }

  .action-btn {
    padding: 10px 20px; background: #238636; border: none; border-radius: 6px;
    color: #fff; font-size: 14px; cursor: pointer;
    display: flex; align-items: center; justify-content: center; gap: 8px; font-family: inherit;
    margin-top: auto;
  }
  .action-btn:hover { background: #2ea043; }
  .action-btn:disabled { opacity: 0.5; cursor: not-allowed; }

  .spinner {
    display: inline-block; width: 14px; height: 14px;
    border: 2px solid #ffffff44; border-top-color: #fff;
    border-radius: 50%; animation: spin 0.6s linear infinite;
  }
  @keyframes spin { to { transform: rotate(360deg); } }

  .main-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 0;
  }

  .result-toolbar { display: flex; gap: 8px; }
  .tool-btn {
    padding: 5px 10px; background: var(--cl-surface, #161b22); border: 1px solid var(--cl-border, #30363d);
    border-radius: 6px; color: var(--cl-text, #e6edf3); cursor: pointer; font-size: 12px; font-family: inherit;
  }
  .tool-btn:hover { border-color: var(--cl-accent, #58a6ff); }

  .preview {
    flex: 1; border: 1px solid var(--cl-border, #30363d); border-radius: 6px;
    background: #fff; min-height: 200px;
  }

  .empty-preview {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px dashed var(--cl-border, #30363d);
    border-radius: 6px;
  }
  .empty-text { color: var(--cl-secondary, #8b949e); font-size: 13px; }

  .banner.error {
    background: #da363433; border: 1px solid #da3634; color: #f85149;
    padding: 8px 12px; border-radius: 6px; font-size: 13px;
  }
</style>
