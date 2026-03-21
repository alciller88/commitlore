<script lang="ts">
  import Dashboard from './screens/Dashboard.svelte'
  import Generate from './screens/Generate.svelte'
  import Story from './screens/Story.svelte'
  import History from './screens/History.svelte'
  import Contributors from './screens/Contributors.svelte'
  import Styles from './screens/Styles.svelte'
  import Settings from './screens/Settings.svelte'

  const screens = [
    { name: 'Dashboard', icon: '~' },
    { name: 'Generate', icon: '+' },
    { name: 'Story', icon: '#' },
    { name: 'History', icon: '>' },
    { name: 'Contributors', icon: '@' },
    { name: 'Styles', icon: '*' },
    { name: 'Settings', icon: '=' },
  ] as const

  let activeScreen = 'Dashboard'
  let activeRepo = ''

  function handleRepoSelected(event: CustomEvent<string>) {
    activeRepo = event.detail
  }

  const components: Record<string, any> = {
    Dashboard, Generate, Story, History, Contributors, Styles, Settings,
  }
</script>

<div class="layout">
  <nav class="sidebar">
    <div class="sidebar-header">CommitLore</div>
    {#each screens as screen}
      <button
        class="nav-item"
        class:active={activeScreen === screen.name}
        on:click={() => activeScreen = screen.name}
      >
        <span class="nav-icon">{screen.icon}</span>
        <span class="nav-label">{screen.name}</span>
      </button>
    {/each}
  </nav>
  <main class="content">
    <svelte:component
      this={components[activeScreen]}
      {activeRepo}
      on:repoSelected={handleRepoSelected}
    />
  </main>
</div>

<style>
  .layout {
    display: flex;
    height: 100vh;
    background: #0d1117;
    color: #e6edf3;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  }

  .sidebar {
    width: 220px;
    background: #161b22;
    border-right: 1px solid #30363d;
    display: flex;
    flex-direction: column;
    padding: 0;
    flex-shrink: 0;
  }

  .sidebar-header {
    padding: 20px 16px;
    font-size: 18px;
    font-weight: 700;
    color: #f0f6fc;
    border-bottom: 1px solid #30363d;
    font-family: 'JetBrains Mono', 'Courier New', monospace;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 16px;
    border: none;
    background: transparent;
    color: #8b949e;
    font-size: 14px;
    cursor: pointer;
    text-align: left;
    width: 100%;
    transition: background 0.15s, color 0.15s;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  }

  .nav-item:hover {
    background: #1f2937;
    color: #e6edf3;
  }

  .nav-item.active {
    background: #1f6feb22;
    color: #58a6ff;
    border-left: 3px solid #58a6ff;
  }

  .nav-icon {
    font-family: 'JetBrains Mono', monospace;
    font-size: 16px;
    width: 20px;
    text-align: center;
  }

  .content {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
  }
</style>
