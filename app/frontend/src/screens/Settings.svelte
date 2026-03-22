<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { GetLLMConfig, SetLLMConfig, ClearLLMKey, GetActiveStyle, SetActiveStyle } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { ListStyles, GetStyleTheme } from '../../bindings/github.com/alciller88/commitlore/app/styleapp.js'
  import { activeStyle } from '../lib/store'

  let styleName = 'formal'
  let styles: Array<{name: string}> = []
  let themePreview = { primary: '', accent: '' }

  const unsubStyle = activeStyle.subscribe(v => { styleName = v })
  onDestroy(() => unsubStyle())

  let provider = ''
  let model = ''
  let keyConfigured = false
  let loading = false
  let error = ''
  let success = ''
  let showKeyModal = false
  let keyInput = ''
  let savingKey = false

  onMount(() => {
    loadConfig()
    loadStyles()
    loadThemePreview(styleName)
  })

  async function loadConfig() {
    loading = true
    error = ''
    try {
      const cfg = await GetLLMConfig()
      provider = cfg.provider
      model = cfg.model
      keyConfigured = cfg.keyConfigured
    } catch (e: any) {
      error = e?.message || 'Failed to load config'
    } finally {
      loading = false
    }
  }

  async function loadStyles() {
    try {
      const raw = await ListStyles()
      styles = JSON.parse(raw)
    } catch {}
  }

  async function loadThemePreview(name: string) {
    try {
      const t = await GetStyleTheme(name)
      themePreview = { primary: t.primary, accent: t.accent }
    } catch {
      themePreview = { primary: '#58a6ff', accent: '#58a6ff' }
    }
  }

  async function changeStyle(name: string) {
    styleName = name
    activeStyle.set(name)
    await SetActiveStyle(name)
    await loadThemePreview(name)
  }

  async function saveConfig() {
    loading = true
    error = ''
    success = ''
    try {
      await SetLLMConfig(provider, model, '')
      success = 'Configuration saved.'
      setTimeout(() => success = '', 3000)
    } catch (e: any) {
      error = e?.message || 'Failed to save config'
    } finally {
      loading = false
    }
  }

  async function saveKey() {
    if (!keyInput || !provider) return
    savingKey = true
    error = ''
    try {
      await SetLLMConfig(provider, model, keyInput)
      keyConfigured = true
      showKeyModal = false
      keyInput = ''
      success = 'API key saved to OS keychain.'
      setTimeout(() => success = '', 3000)
    } catch (e: any) {
      error = e?.message || 'Failed to save key'
    } finally {
      savingKey = false
    }
  }

  async function clearKey() {
    if (!provider) return
    error = ''
    try {
      await ClearLLMKey(provider)
      keyConfigured = false
      success = 'API key removed from keychain.'
      setTimeout(() => success = '', 3000)
    } catch (e: any) {
      error = e?.message || 'Failed to clear key'
    }
  }
</script>

<div class="screen">
  <h1>Settings</h1>

  {#if error}
    <div class="banner error">{error}</div>
  {/if}
  {#if success}
    <div class="banner success">{success}</div>
  {/if}

  <section>
    <h2>Appearance</h2>
    <div class="form">
      <div class="field">
        <label>App style</label>
        <div class="style-selector">
          <select bind:value={styleName} on:change={() => changeStyle(styleName)}>
            {#each styles as s}
              <option value={s.name}>{s.name}</option>
            {/each}
          </select>
          <span class="swatch-preview">
            <span class="swatch" style="background: {themePreview.primary}"></span>
            <span class="swatch" style="background: {themePreview.accent}"></span>
          </span>
        </div>
      </div>
    </div>
  </section>

  <section>
    <h2>LLM Provider</h2>

    <div class="form">
      <div class="field">
        <label>Provider</label>
        <select bind:value={provider}>
          <option value="">None</option>
          <option value="anthropic">Anthropic</option>
          <option value="openai">OpenAI</option>
          <option value="groq">Groq</option>
          <option value="ollama">Ollama</option>
        </select>
      </div>

      <div class="field">
        <label>Model</label>
        <input type="text" bind:value={model} placeholder="e.g., claude-haiku-4-5-20251001" />
      </div>

      <div class="field">
        <label>API Key</label>
        <div class="key-status">
          {#if keyConfigured}
            <span class="status configured">&#10003; Configured</span>
          {:else}
            <span class="status not-configured">&#10007; Not configured</span>
          {/if}
          <button class="tool-btn" on:click={() => { showKeyModal = true; keyInput = '' }} disabled={!provider}>
            Set key
          </button>
          {#if keyConfigured}
            <button class="tool-btn danger" on:click={clearKey}>Clear key</button>
          {/if}
        </div>
      </div>

      <button class="action-btn" on:click={saveConfig} disabled={loading}>
        {loading ? 'Saving...' : 'Save configuration'}
      </button>
    </div>
  </section>

  <section>
    <h2>About</h2>
    <div class="about">
      <div class="about-row"><span class="label">Version:</span> 0.0.0 (Phase 11)</div>
      <div class="about-row"><span class="label">Tagline:</span> Your repo has a story. CommitLore tells it.</div>
    </div>
  </section>
</div>

{#if showKeyModal}
  <div class="modal-overlay" on:click={() => showKeyModal = false} role="dialog">
    <div class="modal" on:click|stopPropagation role="document">
      <h3>Set API Key for {provider}</h3>
      <p class="modal-hint">The key will be stored securely in your OS keychain.</p>
      <input
        type="password"
        bind:value={keyInput}
        placeholder="Paste your API key"
        on:keydown={(e) => e.key === 'Enter' && saveKey()}
      />
      <div class="modal-actions">
        <button class="tool-btn" on:click={() => showKeyModal = false}>Cancel</button>
        <button class="action-btn" on:click={saveKey} disabled={savingKey || !keyInput}>
          {savingKey ? 'Saving...' : 'Save key'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .screen { display: flex; flex-direction: column; gap: 24px; max-width: 600px; }
  h1 { color: var(--cl-text, #e6edf3); font-size: 22px; margin: 0; }
  h2 { color: var(--cl-text, #e6edf3); font-size: 16px; margin: 0 0 12px; }
  h3 { color: var(--cl-text, #e6edf3); font-size: 16px; margin: 0 0 8px; }

  section {
    padding: 20px; background: var(--cl-surface, #161b22); border: 1px solid var(--cl-border, #30363d); border-radius: 8px;
  }

  .form { display: flex; flex-direction: column; gap: 14px; }
  .field { display: flex; flex-direction: column; gap: 4px; }
  .field label { color: var(--cl-secondary, #8b949e); font-size: 12px; text-transform: uppercase; }

  input, select {
    padding: 8px 12px; background: var(--cl-background, #0d1117); border: 1px solid var(--cl-border, #30363d);
    border-radius: 6px; color: var(--cl-text, #e6edf3); font-size: 14px;
    font-family: 'JetBrains Mono', monospace;
  }
  input:focus, select:focus { outline: none; border-color: var(--cl-accent, #58a6ff); }
  select { cursor: pointer; }

  .key-status { display: flex; align-items: center; gap: 12px; }
  .status { font-size: 14px; }
  .configured { color: #3fb950; }
  .not-configured { color: #f85149; }

  .tool-btn {
    padding: 6px 14px; background: var(--cl-surface, #161b22); border: 1px solid var(--cl-border, #30363d);
    border-radius: 6px; color: var(--cl-text, #e6edf3); cursor: pointer; font-size: 13px; font-family: inherit;
  }
  .tool-btn:hover { border-color: var(--cl-accent, #58a6ff); }
  .tool-btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .tool-btn.danger:hover { border-color: #da3634; color: #f85149; }

  .action-btn {
    padding: 10px 20px; background: #238636; border: none; border-radius: 6px;
    color: #fff; font-size: 14px; cursor: pointer; align-self: flex-start; font-family: inherit;
  }
  .action-btn:hover { background: #2ea043; }
  .action-btn:disabled { opacity: 0.5; cursor: not-allowed; }

  .about { font-size: 13px; }
  .about-row { margin-bottom: 4px; }
  .about-row .label { color: var(--cl-secondary, #8b949e); }

  .banner {
    padding: 8px 12px; border-radius: 6px; font-size: 13px;
  }
  .banner.error {
    background: #da363433; border: 1px solid #da3634; color: #f85149;
  }
  .banner.success {
    background: #23863633; border: 1px solid #238636; color: #3fb950;
  }

  .modal-overlay {
    position: fixed; top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0, 0, 0, 0.6); display: flex;
    align-items: center; justify-content: center; z-index: 100;
  }
  .modal {
    background: var(--cl-surface, #161b22); border: 1px solid var(--cl-border, #30363d); border-radius: 12px;
    padding: 24px; max-width: 400px; width: 100%;
    display: flex; flex-direction: column; gap: 12px;
  }
  .modal-hint { color: var(--cl-secondary, #8b949e); font-size: 13px; margin: 0; }
  .modal-actions { display: flex; gap: 8px; justify-content: flex-end; margin-top: 8px; }

  .style-selector { display: flex; align-items: center; gap: 10px; }
  .style-selector select { flex: 1; }
  .swatch-preview { display: flex; gap: 4px; }
  .swatch {
    width: 16px; height: 16px; border-radius: 3px;
    border: 1px solid var(--cl-border, #30363d);
  }
</style>
