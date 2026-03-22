<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { GenerateStory } from '../../bindings/github.com/alciller88/commitlore/app/storyapp.js'
  import { GetLLMConfig } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { activeRepo, activeStyle } from '../lib/store'
  import type { ActiveRepo } from '../lib/store'

  let currentRepo: ActiveRepo | null = null
  let currentStyle = 'formal'
  const unsubRepo = activeRepo.subscribe(v => { currentRepo = v })
  const unsubStyle = activeStyle.subscribe(v => { currentStyle = v })
  onDestroy(() => { unsubRepo(); unsubStyle() })

  let from = ''
  let loading = false
  let error = ''
  let htmlResult = ''
  let llmStatus = ''

  onMount(() => {
    loadLLMStatus()
  })

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

  async function tellStory() {
    if (!currentRepo) return
    loading = true
    error = ''
    htmlResult = ''
    try {
      htmlResult = await GenerateStory(currentRepo.path, from, currentStyle)
    } catch (e: any) {
      error = e?.message || 'Story generation failed'
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
    a.download = 'story.html'
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
      <h2>Tell the Story</h2>

      {#if error}
        <div class="banner error">{error}</div>
      {/if}

      <div class="field">
        <label>From</label>
        <input type="text" bind:value={from} placeholder="Starting commit or tag" />
      </div>

      <div class="style-pill">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="12" height="12"><circle cx="12" cy="12" r="3"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
        <span>{currentStyle}</span>
      </div>

      <div class="llm-indicator">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="12" height="12"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/></svg>
        <span>LLM: {llmStatus}</span>
      </div>

      <button class="action-btn" on:click={tellStory} disabled={loading}>
        {#if loading}
          <span class="spinner"></span> Generating...
        {:else}
          Tell the story
        {/if}
      </button>
    </div>

    <div class="main-panel">
      {#if htmlResult}
        <div class="result-toolbar">
          <button class="tool-btn" on:click={copyText}>Copy text</button>
          <button class="tool-btn" on:click={saveAsFile}>Save as .html</button>
        </div>
        <iframe class="preview" srcdoc={htmlResult} sandbox="" title="Story preview"></iframe>
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

  input {
    padding: 7px 10px; background: var(--cl-background, #0d1117); border: 1px solid var(--cl-border, #30363d);
    border-radius: 6px; color: var(--cl-text, #e6edf3); font-size: 13px;
    font-family: 'JetBrains Mono', monospace;
  }
  input:focus { outline: none; border-color: var(--cl-accent, #58a6ff); }

  .style-pill {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 10px;
    background: color-mix(in srgb, var(--cl-accent, #58a6ff) 13%, transparent);
    border: 1px solid var(--cl-accent, #58a6ff);
    border-radius: 12px;
    color: var(--cl-accent, #58a6ff);
    font-size: 11px;
    font-weight: 600;
    align-self: flex-start;
  }

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
