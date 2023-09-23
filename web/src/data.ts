import type { ServerData } from "./types.d"

export let load = async function() {
  let resp = await fetch("api/v1/servers")
  return resp.json()
}

