package mocks

import (
	"context"
	"time"

	"goal-tracker/api/pkg/steam"
)

type MockSteamClient struct {
}

func NewMockSteamClient() steam.Client {
	return MockSteamClient{}
}

func (client MockSteamClient) GetOwnedGames(
	_ context.Context,
	_ string,
) (*steam.OwnedGamesResponse, error) {
	response := steam.OwnedGamesResponse{
		Response: steam.OwnedGamesResponseData{
			GameCount: 1,
			Games: []steam.Game{
				{
					AppID:                    1,
					Name:                     "test",
					ImgIconURL:               "",
					ImgLogoURL:               "",
					HasCommunityVisibleStats: true,
				},
			},
		},
	}
	return &response, nil
}

func (client MockSteamClient) GetPlayerAchievements(
	_ context.Context,
	steamID string,
	_ int,
) (*steam.AchievementsResponse, error) {
	response := steam.AchievementsResponse{
		PlayerStats: steam.PlayerStats{
			SteamID:  steamID,
			GameName: "test",
			Achievements: []steam.Achievement{
				{
					APIName:     "TEST",
					Achieved:    1,
					UnlockTime:  int64(time.Now().UTC().Second()),
					Name:        "test",
					Description: "Hello, World!",
				},
			},
		},
	}
	return &response, nil
}
