<script lang="ts">
  import type { Server, Player } from "../types.d"

  import play from "../assets/play-circle.svg"
  import external_link from "../assets/external-link.svg"
  import eye from "../assets/eye.svg"

  export let server: Server

  function truncateName(s: string, len: number) {
    if (s.length > len) {
      return s.slice(0, len) + "..."
    }
    return s
  }

  function pickAddress(addr: string, resolved?: string) {
    if (resolved != null && resolved.length > 0) {
      return resolved
    }
    return addr
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
  $: serverNameShort = truncateName(server.name, 30)

  let serverAddress: string
  $: serverAddress = pickAddress(server.address, server.resolved_address)

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
      <span>{server.map}</span>
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
      <div class="col col-name">{truncateName(player.name,30)}</div>
      <div class="col col-score">{player.score}</div>
      <div class="col col-ping">{player.ping}</div>
    </div>
    {/each}
  </div>
  {:else}
  <div class="fill"></div>
  {/if}
  <div class="footer">
    <div class="address">
      {serverAddress}
    </div>
    <div class="links">
      <a href="ut2004://{serverAddress}" title="Join"><img src="{play}" alt="Join"></a>
      <a href="ut2004://{serverAddress}?SpectatorOnly=1" title="Spectate"><img src="{eye}" alt="Spectate"></a>
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
    display: flex;
    flex-flow: column nowrap;
    flex-grow: 1.0;

    h2 {
      font: 400 0.875em var(--font-primary);
      color: var(--color-primary-200);
      padding: 0;
      margin: 0;
    }

    span {
      font: 400 0.75em var(--font-secondary);
      color: var(--color-primary-400);
    }
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
  color: var(--color-primary-300);

  .links {
    display: flex;
    flex-flow: row nowrap;
    gap: 1em;

    img {
      width: 16px;
      filter: invert(100%) sepia(33%) saturate(730%) hue-rotate(174deg) brightness(65%) contrast(87%);
    }

    a:hover, a:active {
      img {
        filter: invert(100%) sepia(33%) saturate(730%) hue-rotate(174deg) brightness(86%) contrast(87%);
      }
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
      width: 60%;
    }

    .col-score {
      width: 20%;
      text-align: right;
    }

    .col-ping {
      width: 20%;
      text-align: right;
    }
  }
}

.fill {
  flex-grow: 1;
}
</style>
