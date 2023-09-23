<script lang="ts">
  import Server from "./lib/Server.svelte";
  import { Fixtures }  from "./data";

  let data = Fixtures
</script>

<main>
{#each data.groups as group }
  <div class="group">
    <h1>{group.name}</h1>
    <div class="panels">
    {#each group.servers as server }
      <Server server={server}></Server>
    {/each}
    </div>
  </div>
{/each}
</main>

<style lang="scss">
.panels {
  padding: 1rem 2rem 2rem 2rem;
  display: flex;
  flex-flow: column wrap;
  gap: 1rem;

  max-height: 200px;
}

h1 {
  font: 400 1rem var(--font-secondary);
  margin: 1rem 2rem 0 2rem;
  color: var(--color-primary-300);
}
</style>
