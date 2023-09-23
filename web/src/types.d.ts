export interface ServerData {
  groups: Group[]
}

export interface Group {
  name: string
  servers: Server[]
}

export interface Server {
  name: string
  address: string
  resolved_address?: string
  external_link?: string
  online: boolean
  map: string
  player_count: PlayerCount
  players: Player[]
}

export interface PlayerCount {
  current: number
  max: number
}

export interface Player {
  name: string
  score: number
  ping: number
}
