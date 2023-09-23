<script lang="ts">
  import type { ServerData } from "./types.d"

  import { onMount } from "svelte"

  import { load } from "./data"
  import Server from "./lib/Server.svelte"

  let data: ServerData

  onMount(async () => {
    data = await performLoad()
  })

  async function performLoad() {
    data = await load()
    setTimeout(performLoad, 15*1000)
    return data
  }
</script>

<main>
{#if data != null && data.groups != null}
{#each data.groups as group (group.name) }
  <div class="group">
    <h1>{group.name}</h1>
    <div class="panels">
    {#if group.servers != null}
    {#each group.servers as server (server.address) }
      <Server server={server}></Server>
    {/each}
    {/if}
    </div>
  </div>
{/each}
{/if}
</main>

<style lang="scss">
.panels {
  padding: 1rem 2rem 2rem 2rem;

  display: grid;
  grid-gap: 1rem;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
}

h1 {
  font: 400 1.2rem var(--font-secondary);
  margin: 1rem 2rem 0 2rem;
  color: var(--color-primary-300);
  text-transform: uppercase;
}
</style>
