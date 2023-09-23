import type { ServerData } from "./types.d"

export let Fixtures: ServerData = {
  groups: [
    {
      name: "Chicago",
      servers: [
        {
          name: "UFC Duel Cup CHI-1.1",
          address: "chi-1.kokuei.dev:7777",
          external_link: "https://chi-1.kokuei.dev/instance/1",
          online: true,
          player_count: {
            current: 2,
            max: 2,
          },
          players: [
            {
              name: "kokuei",
              score: 20,
              ping: 44,
            },
            {
              name: "bot",
              score: 2,
              ping: 12,
            }
          ]
        },
        {
          name: "UFC Duel Cup CHI-1.2",
          address: "chi-2.kokuei.dev:7777",
          online: true,
          player_count: {
            current: 0,
            max: 2,
          },
          players: []
        },
        {
          name: "UFC Duel Cup CHI-1.2",
          address: "chi-2.kokuei.dev:7777",
          online: false,
          player_count: {
            current: 0,
            max: 2,
          },
          players: []
        }
      ]
    },
    {
      name: "Chicago",
      servers: [
        {
          name: "UFC Duel Cup CHI-1.1",
          address: "chi-1.kokuei.dev:7777",
          online: true,
          player_count: {
            current: 2,
            max: 2,
          },
          players: [
            {
              name: "kokuei",
              score: 20,
              ping: 44,
            },
            {
              name: "bot",
              score: 2,
              ping: 12,
            }
          ]
        },
        {
          name: "UFC Duel Cup CHI-1.2",
          address: "chi-1.kokuei.dev:7797",
          online: true,
          player_count: {
            current: 2,
            max: 2,
          },
          players: [
            {
              name: "kokuei",
              score: 20,
              ping: 44,
            },
            {
              name: "bot",
              score: 2,
              ping: 12,
            }
          ]
        }
      ]
    }
  ]
}
