package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"sort"
)

func NewBrowserServer(root fs.FS) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/servers", HandleServerQuery)
	mux.Handle("/", http.FileServer(http.FS(root)))

	server := &http.Server{
		Handler: mux,
	}

	return server
}

type groupJSON struct {
	Name    string       `json:"name"`
	Servers []serverJSON `json:"servers"`
}

type serverJSON struct {
	Name            string          `json:"name"`
	Address         string          `json:"address"`
	ResolvedAddress string          `json:"resolved_address"`
	ExternalLink    string          `json:"external_link"`
	Online          bool            `json:"online"`
	Map             string          `json:"map"`
	PlayerCount     playerCountJSON `json:"player_count"`
	Players         []playerJSON    `json:"players"`
}

type playerCountJSON struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}

type playerJSON struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
	Ping  int    `json:"ping"`
}

func HandleServerQuery(w http.ResponseWriter, r *http.Request) {
	servers := state.Servers()

	// Organize servers by group
	byGroup := make(map[string][]ServerState)
	for _, svr := range servers {
		// Ignore offline servers for now
		if !svr.Online {
			continue
		}

		list, exists := byGroup[svr.Registration.Group]
		if !exists {
			list = make([]ServerState, 0, 10)
		}

		list = append(list, svr)
		byGroup[svr.Registration.Group] = list
	}

	groups := make([]groupJSON, 0, len(byGroup))
	for name, groupServers := range byGroup {
		group := groupJSON{
			Name:    name,
			Servers: make([]serverJSON, 0, len(groupServers)),
		}

		for _, svr := range groupServers {
			var jsonObj serverJSON

			// Fallback values
			jsonObj.Name = svr.Registration.Address
			jsonObj.Address = svr.Registration.Address
			jsonObj.ResolvedAddress = svr.ResolvedAddress
			jsonObj.ExternalLink = svr.Registration.ExternalLink
			jsonObj.Online = svr.Online

			if svr.Online && svr.Details != nil {
				jsonObj.Name = svr.Details.Info.ServerName.String()
				jsonObj.Map = svr.Details.Info.MapName.String()
				jsonObj.PlayerCount = playerCountJSON{
					Current: int(svr.Details.Info.CurrentPlayers),
					Max:     int(svr.Details.Info.MaxPlayers),
				}

				for _, player := range svr.Details.Players {
					jsonObj.Players = append(jsonObj.Players, playerJSON{
						Name:  player.Name.String(),
						Score: int(player.Score),
						Ping:  int(player.Ping),
					})
				}

				// Sort players by score, then by name
				sort.SliceStable(jsonObj.Players, func(i, j int) bool {
					left := jsonObj.Players[i].Score
					right := jsonObj.Players[j].Score

					if left > right {
						return true
					} else if left == right {
						return jsonObj.Players[i].Name < jsonObj.Players[j].Name
					}

					return false
				})
			}

			group.Servers = append(group.Servers, jsonObj)
		}

		sort.SliceStable(group.Servers, func(i, j int) bool {
			return group.Servers[i].Name < group.Servers[j].Name
		})

		sort.SliceStable(group.Servers, func(i, j int) bool {
			return len(group.Servers[i].Players) > len(group.Servers[j].Players)
		})

		groups = append(groups, group)
	}

	// Sort groups by name (Dallas, Chicago, New Jersey, UK)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	var doc struct {
		Groups []groupJSON `json:"groups"`
	}

	doc.Groups = groups

	w.Header().Add("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(doc)

	if err != nil {
		logger.Error("failed to encode document", "err", err)
	}
}
