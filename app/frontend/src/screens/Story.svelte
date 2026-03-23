<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { GenerateStory } from '../../bindings/github.com/alciller88/commitlore/app/storyapp.js'
  import { GetLLMConfig } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { activeRepo, activeStyle, uiLabels } from '../lib/store'
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
      <span class="section-label">Filters</span>

      {#if error}
        <div class="banner error">{error}</div>
      {/if}

      <div class="fields">
        <div class="field">
          <label>From</label>
          <input type="text" bind:value={from} placeholder="Starting commit or tag" />
        </div>
      </div>

      <div class="style-pill">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="12" height="12"><circle cx="12" cy="12" r="3"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>
        <span>{currentStyle}</span>
      </div>

      <div class="llm-indicator">
        <span class="llm-dot" class:configured={llmStatus !== 'not configured'}></span>
        <span>{llmStatus}</span>
      </div>

      <button class="action-btn" on:click={tellStory} disabled={loading}>
        {#if loading}
          <span class="spinner"></span> Generating...
        {:else}
          {$uiLabels.storyButton}
        {/if}
      </button>
    </div>

    <div class="main-panel">
      {#if htmlResult}
        <div class="result-toolbar">
          <button class="tool-btn" on:click={copyText}>Copy text</button>
          <button class="tool-btn" on:click={saveAsFile}>Save as .html</button>
        </div>
        <iframe class="preview" srcdoc={htmlResult} sandbox="allow-scripts" title="Story preview"></iframe>
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
    flex: 1;
    gap: var(--space-3);
    color: var(--cl-secondary);
  }
  .no-repo p {
    margin: 0;
    font-size: var(--text-md);
  }

  .two-col {
    display: flex;
    gap: var(--space-4);
    flex: 1;
    min-height: 0;
  }

  .sidebar-form {
    width: 240px;
    flex-shrink: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }

  .section-label {
    font-size: var(--text-xs);
    text-transform: uppercase;
    color: var(--cl-secondary);
    opacity: 0.6;
    letter-spacing: 0.05em;
  }

  .fields {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
  }
  .field label {
    font-size: var(--text-xs);
    text-transform: uppercase;
    color: var(--cl-secondary);
  }

  input {
    width: 100%;
    box-sizing: border-box;
    padding: var(--space-2) var(--space-2);
    background: var(--cl-background);
    border: 1px solid var(--cl-border);
    border-radius: var(--radius-md);
    color: var(--cl-text);
    font-size: var(--text-base);
    font-family: 'JetBrains Mono', monospace;
    transition: border-color var(--transition-fast);
  }
  input:focus {
    outline: none;
    border-color: var(--cl-accent);
  }

  .style-pill {
    display: inline-flex;
    align-items: center;
    gap: var(--space-1);
    padding: var(--space-1) var(--space-2);
    background: color-mix(in srgb, var(--cl-accent) 13%, transparent);
    border: 1px solid var(--cl-accent);
    border-radius: 12px;
    color: var(--cl-accent);
    font-size: var(--text-xs);
    font-weight: 600;
    align-self: flex-start;
  }

  .llm-indicator {
    display: flex;
    align-items: center;
    gap: var(--space-1);
    font-size: var(--text-xs);
    color: var(--cl-secondary);
  }
  .llm-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--cl-secondary);
    flex-shrink: 0;
  }
  .llm-dot.configured {
    background: #3fb950;
  }

  .action-btn {
    width: 100%;
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
    gap: var(--space-2);
    margin-top: auto;
    transition: opacity var(--transition-fast);
  }
  .action-btn:hover {
    opacity: 0.9;
  }
  .action-btn:disabled {
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

  .main-panel {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .result-toolbar {
    display: flex;
    gap: var(--space-2);
    justify-content: flex-end;
    height: 36px;
    align-items: center;
    flex-shrink: 0;
  }
  .tool-btn {
    padding: var(--space-1) var(--space-3);
    background: transparent;
    border: 1px solid var(--cl-border);
    border-radius: var(--radius-md);
    color: var(--cl-secondary);
    cursor: pointer;
    font-size: var(--text-sm);
    font-family: var(--cl-font-family);
    transition: border-color var(--transition-fast), color var(--transition-fast);
  }
  .tool-btn:hover {
    border-color: var(--cl-accent);
    color: var(--cl-text);
  }

  .preview {
    flex: 1;
    border: 1px solid var(--cl-border);
    border-radius: var(--radius-lg);
    background: #fff;
    min-height: 200px;
  }

  .empty-preview {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px dashed var(--cl-border);
    border-radius: var(--radius-lg);
  }
  .empty-text {
    color: var(--cl-secondary);
    font-size: var(--text-base);
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
