<script lang="ts">
  import { onMount, onDestroy } from 'svelte'
  import { GetLLMConfig, SetLLMConfig, ClearLLMKey, GetActiveStyle, SetActiveStyle, GetLanguage, SetLanguage } from '../../bindings/github.com/alciller88/commitlore/app/configapp.js'
  import { ListStyles, GetStyleTheme, GetAvailableLanguagesForStyle } from '../../bindings/github.com/alciller88/commitlore/app/styleapp.js'
  import { activeStyle, activeLanguage } from '../lib/store'
  import { applyTheme } from '../lib/theme'

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

  let language = 'en'
  let availableLanguages: string[] = ['en']
  let languageError = ''

  onMount(() => {
    loadConfig()
    loadStyles()
    loadThemePreview(styleName)
    loadLanguage()
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
    await loadAvailableLanguages(name)
  }

  async function loadLanguage() {
    try {
      language = await GetLanguage()
      activeLanguage.set(language)
      await loadAvailableLanguages(styleName)
    } catch {}
  }

  async function loadAvailableLanguages(name: string) {
    try {
      availableLanguages = await GetAvailableLanguagesForStyle(name)
    } catch {
      availableLanguages = ['en']
    }
  }

  async function changeLanguage(lang: string) {
    languageError = ''
    try {
      await SetLanguage(lang)
      language = lang
      activeLanguage.set(lang)
      await applyTheme(styleName)
    } catch (e: any) {
      languageError = e?.message || 'Failed to change language'
    }
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
    <h2>Language</h2>
    <div class="form">
      <div class="field">
        <label>App language</label>
        <div class="lang-selector">
          <label class="lang-option">
            <input type="radio" name="lang" value="en"
              checked={language === 'en'}
              on:change={() => changeLanguage('en')}
              disabled={!availableLanguages.includes('en')} />
            <span>English</span>
          </label>
          <label class="lang-option">
            <input type="radio" name="lang" value="es"
              checked={language === 'es'}
              on:change={() => changeLanguage('es')}
              disabled={!availableLanguages.includes('es')} />
            <span>Español</span>
          </label>
        </div>
        {#if languageError}
          <div class="lang-error">{languageError}</div>
        {/if}
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
          <button class="btn-outline" on:click={() => { showKeyModal = true; keyInput = '' }} disabled={!provider}>
            Set key
          </button>
          {#if keyConfigured}
            <button class="btn-outline danger" on:click={clearKey}>Clear key</button>
          {/if}
        </div>
      </div>

      <button class="btn-primary" on:click={saveConfig} disabled={loading}>
        {loading ? 'Saving...' : 'Save configuration'}
      </button>
    </div>
  </section>

  <section>
    <h2>About</h2>
    <div class="about">
      <div class="about-row"><span class="about-label">Version:</span> 0.0.0 (Phase 11)</div>
      <div class="about-row"><span class="about-label">Tagline:</span> Your repo has a story. CommitLore tells it.</div>
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
        <button class="btn-outline" on:click={() => showKeyModal = false}>Cancel</button>
        <button class="btn-primary" on:click={saveKey} disabled={savingKey || !keyInput}>
          {savingKey ? 'Saving...' : 'Save key'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  .screen {
    display: flex;
    flex-direction: column;
    gap: var(--space-5);
    flex: 1;
    overflow-y: auto;
    max-width: 600px;
  }

  section {
    padding: var(--space-4);
    background: var(--cl-surface);
    border: 1px solid var(--cl-border);
    border-radius: var(--radius-lg);
  }

  h2 {
    font-size: var(--text-xs);
    text-transform: uppercase;
    color: var(--cl-secondary);
    font-weight: 500;
    letter-spacing: 0.05em;
    margin: 0 0 var(--space-3) 0;
    padding-bottom: var(--space-2);
    border-bottom: 1px solid var(--cl-surface);
  }

  h3 {
    color: var(--cl-text);
    font-size: var(--text-lg);
    margin: 0 0 var(--space-2) 0;
  }

  .form {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .field label {
    font-size: var(--text-xs);
    text-transform: uppercase;
    color: var(--cl-secondary);
  }

  input, select {
    height: 36px;
    box-sizing: border-box;
    padding: 0 var(--space-3);
    background: var(--cl-background);
    border: 1px solid var(--cl-border);
    border-radius: var(--radius-md);
    color: var(--cl-text);
    font-size: var(--text-md);
    font-family: 'JetBrains Mono', monospace;
    transition: border-color var(--transition-fast);
  }
  input:focus, select:focus {
    outline: none;
    border-color: var(--cl-accent);
  }
  select {
    cursor: pointer;
  }

  .key-status {
    display: flex;
    align-items: center;
    gap: var(--space-3);
  }
  .status {
    font-size: var(--text-md);
  }
  .configured {
    color: #3fb950;
  }
  .not-configured {
    color: #f85149;
  }

  .btn-outline {
    padding: var(--space-1) var(--space-3);
    background: transparent;
    border: 1px solid var(--cl-accent);
    border-radius: var(--radius-md);
    color: var(--cl-accent);
    cursor: pointer;
    font-size: var(--text-base);
    font-family: var(--cl-font-family);
    transition: background var(--transition-fast), color var(--transition-fast);
  }
  .btn-outline:hover {
    background: color-mix(in srgb, var(--cl-accent) 10%, transparent);
  }
  .btn-outline:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .btn-outline.danger {
    border-color: #da3634;
    color: #f85149;
  }
  .btn-outline.danger:hover {
    background: color-mix(in srgb, #da3634 10%, transparent);
  }

  .btn-primary {
    height: 36px;
    padding: 0 var(--space-4);
    background: var(--cl-accent);
    border: none;
    border-radius: var(--radius-md);
    color: #fff;
    font-size: var(--text-base);
    font-family: var(--cl-font-family);
    cursor: pointer;
    align-self: flex-start;
    transition: opacity var(--transition-fast);
  }
  .btn-primary:hover {
    opacity: 0.9;
  }
  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .about {
    font-size: var(--text-base);
    color: var(--cl-text);
  }
  .about-row {
    margin-bottom: var(--space-1);
  }
  .about-label {
    color: var(--cl-secondary);
  }

  .banner {
    padding: var(--space-2) var(--space-3);
    border-radius: var(--radius-md);
    font-size: var(--text-base);
  }
  .banner.error {
    background: #da363433;
    border: 1px solid #da3634;
    color: #f85149;
  }
  .banner.success {
    background: #23863633;
    border: 1px solid #238636;
    color: #3fb950;
  }

  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
  }
  .modal {
    background: var(--cl-surface);
    border: 1px solid var(--cl-border);
    border-radius: var(--radius-lg);
    padding: var(--space-5);
    max-width: 400px;
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
  }
  .modal-hint {
    color: var(--cl-secondary);
    font-size: var(--text-base);
    margin: 0;
  }
  .modal-actions {
    display: flex;
    gap: var(--space-2);
    justify-content: flex-end;
    margin-top: var(--space-2);
  }

  .lang-selector {
    display: flex;
    gap: var(--space-4);
  }
  .lang-option {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    cursor: pointer;
    font-size: var(--text-md);
    color: var(--cl-text);
  }
  .lang-option input[type="radio"] {
    accent-color: var(--cl-accent);
    cursor: pointer;
    width: 16px;
    height: 16px;
  }
  .lang-option input[type="radio"]:disabled + span {
    opacity: 0.4;
    cursor: not-allowed;
  }
  .lang-error {
    margin-top: var(--space-1);
    padding: var(--space-2) var(--space-3);
    background: #da363433;
    border: 1px solid #da3634;
    border-radius: var(--radius-md);
    color: #f85149;
    font-size: var(--text-base);
  }

  .style-selector {
    display: flex;
    align-items: center;
    gap: var(--space-2);
  }
  .style-selector select {
    flex: 1;
  }
  .swatch-preview {
    display: flex;
    gap: var(--space-1);
  }
  .swatch {
    width: 20px;
    height: 20px;
    border-radius: 50%;
    border: 1px solid var(--cl-surface);
  }
</style>
