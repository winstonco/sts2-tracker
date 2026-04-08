package main

import (
	"strings"
)

type CharacterWinLossObject struct {
	Character string `json:"character"`
	Losses    int    `json:"losses"`
	Wins      int    `json:"wins"`
}

func (c CharacterWinLossObject) Name() string {
	s, _ := strings.CutPrefix(c.Character, "CHARACTER.")
	return s[:1] + strings.ToLower(s[1:])
}

type AncientStatsObject struct {
	AncientId      string                   `json:"ancient_id"`
	CharacterStats []CharacterWinLossObject `json:"character_stats"`
}

type CardStatsObject struct {
	Id           string `json:"id"`
	TimesLost    int    `json:"times_lost"`
	TiemsPicked  int    `json:"times_picked"`
	TimesSkipped int    `json:"times_skipped"`
	TimesWon     int    `json:"times_won"`
}

type Badge struct {
	Count  int    `json:"count"`
	Id     string `json:"id"`
	Rarity string `json:"rarity"`
}

type CharacterStatsObject struct {
	Badges             []Badge `json:"badges"`
	BestWinStreak      int     `json:"best_win_streak"`
	CurrentStreak      int     `json:"current_streak"`
	FastestWinTime     int     `json:"fastest_win_time"`
	Id                 string  `json:"id"`
	MaxAscension       int     `json:"max_ascension"`
	Playtime           int     `json:"playtime"`
	PreferredAscension int     `json:"preferred_ascension"`
	TotalLosses        int     `json:"total_losses"`
	TotalWins          int     `json:"total_wins"`
}

func (c CharacterStatsObject) Name() string {
	s, _ := strings.CutPrefix(c.Id, "CHARACTER.")
	return s[:1] + strings.ToLower(s[1:])
}

type EncounterStatsObject struct {
	EncounterId string                   `json:"encounter_id"`
	FightStats  []CharacterWinLossObject `json:"fight_stats"`
}

type EnemyStatsObject struct {
	EnemyId    string                   `json:"enemy_id"`
	FightStats []CharacterWinLossObject `json:"fight_stats"`
}

type EpochsObject struct {
	Id         string `json:"id"`
	ObtainDate int    `json:"obtain_date"`
	State      string `json:"state"`
}

type ProgressSave struct {
	AncientStats      []AncientStatsObject   `json:"ancient_stats"`
	ArchitectDamage   int                    `json:"architect_damage"`
	CardStats         []CardStatsObject      `json:"card_stats"`
	CharacterStats    []CharacterStatsObject `json:"character_stats"`
	CurrentScore      int                    `json:"current_score"`
	DiscoveredActs    []string               `json:"discovered_acts"`
	DiscoveredCards   []string               `json:"discovered_cards"`
	DiscoveredEvents  []string               `json:"discovered_events"`
	DiscoveredPotions []string               `json:"discovered_potions"`
	DiscoveredRelics  []string               `json:"discovered_relics"`
	EnableFtues       bool                   `json:"enable_ftues"`
	EncounterStats    []EncounterStatsObject `json:"encounter_stats"`
	EnemyStats        []EnemyStatsObject     `json:"enemy_stats"`
	Epochs            []EpochsObject         `json:"epochs"`
	FloorsClimbed     int                    `json:"floors_climbed"`
	// FtueCompleted []string `json:"ftue_completed"`
	// MaxMultiplayerAscension int `json:"max_multiplayer_ascension"`
	// PendingCharacterUnlock string `json:"pending_character_unlock"`
	// PreferredMultiplayerAscension int                    `json:"preferred_multiplayer_ascension"`
	SchemaVersion    int    `json:"schema_version"`
	TestSubjectKills int    `json:"test_subject_kills"`
	TotalPlaytime    int    `json:"total_playtime"`
	TotalUnlocks     int    `json:"total_unlocks"`
	UniqueId         string `json:"unique_id"`
	// UnlockedAchievements
	WongoPoints int `json:"wongo_points"`
}
