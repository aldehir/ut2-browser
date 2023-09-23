<script lang="ts">
  import type { Server, Player } from "../types.d"

  import copy from "../assets/copy.svg"
  import external_link from "../assets/external-link.svg"
  import eye from "../assets/eye.svg"

  export let server: Server

  function truncateName(s: string) {
    if (s.length > 30) {
      return s.slice(0, 30) + "..."
    }
    return s
  }

  function calcGridRows(players?: Player[]) {
    if (players == null) {
      return 1
    }

    if (players.length < 3) {
      return 2
    }

    return 3
  }

  let serverNameShort: string
  $: serverNameShort = truncateName(server.name)

  let isEmpty: boolean
  $: isEmpty = server.players == null || server.players.length == 0

  let gridRows: number
  $: gridRows = calcGridRows(server.players)
</script>

<div
  class="panel"
  style="grid-row: auto / span {gridRows};"
  class:offline={!server.online}
  class:empty={isEmpty}
  class:hasplayers={!isEmpty}>

  <div class="header">
    <div class="title">
      <h2>{serverNameShort}</h2>
    </div>
    <div class="status">
      {#if server.online}
      {server.player_count.current} / {server.player_count.max}
      {:else}
      Offline
      {/if}
    </div>
  </div>
  {#if server.players != null && server.players.length > 0}
  <div class="players">
    <div class="row">
      <div class="col-header col-name">Name</div>
      <div class="col-header col-score">Score</div>
      <div class="col-header col-ping">Ping</div>
    </div>
    {#each server.players as player}
    <div class="row">
      <div class="col col-name">{player.name}</div>
      <div class="col col-score">{player.score}</div>
      <div class="col col-ping">{player.ping}ms</div>
    </div>
    {/each}
  </div>
  {:else}
  <div class="fill"></div>
  {/if}
  <div class="footer">
    <div class="address">
      {server.address}
    </div>
    <div class="links">
      {#if server.external_link}
        <a target="_blank" href="{server.external_link}" title="External Link"><img src="{external_link}" alt="Admin"></a>
      {/if}
    </div>
  </div>
</div>

<style lang="scss">
.panel {
  box-sizing: border-box;
  border-radius: 4px;

  background: var(--color-primary-700);
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.15), 0 5px 15px rgba(0, 0, 0, 0.2);

  display: flex;
  flex-flow: column nowrap;

  h2 {
    font: 400 0.875em var(--font-primary);
    color: var(--color-primary-200);
    padding: 0;
    margin: 0;
  }
}

.offline {
  opacity: 40%;
}

.header {
  box-sizing: border-box;
  display: flex;
  flex-flow: row nowrap;
  align-items: baseline;
  padding: 1rem;

  .title {
    flex-grow: 1.0;
  }

  .status {
    font: 400 0.95em var(--font-secondary);
    color: var(--color-primary-300);
    text-transform: uppercase;
  }
}

.footer {
  display: flex;
  flex-flow: row nowrap;
  padding: 0.5rem 1rem;

  align-items: center;
  justify-content: space-between;

  font: 400 0.8rem var(--font-primary);
  color: var(--color-primary-400);

  .links {
    display: flex;
    flex-flow: row nowrap;
    gap: 1em;

    img {
      width: 16px;
      filter: invert(100%) sepia(33%) saturate(730%) hue-rotate(174deg) brightness(86%) contrast(87%);
    }
  }
}

.empty {
  .footer {
    padding: 0rem 1rem 0.5rem 1rem;
  }
}

.players {
  padding: 0.5rem 1rem;
  background: var(--color-primary-750);
  flex-grow: 1;

  .row {
    display: flex;
    flex-flow: row nowrap;
    margin-bottom: 0.2em;

    .col-header {
      font: 400 0.8em var(--font-primary);
      color: var(--color-primary-400);
      text-transform: uppercase;
    }

    .col {
      font: 400 0.8em var(--font-primary);
      color: var(--color-primary-200);
    }

    .col-name {
      width: 70%;
    }

    .col-score {
      width: 15%;
      text-align: right;
    }

    .col-ping {
      width: 15%;
      text-align: right;
    }
  }
}

.fill {
  flex-grow: 1;
}
</style>
