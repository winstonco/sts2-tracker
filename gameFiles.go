package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const STEAM_USERDATA_PATH = `C:\Program Files (x86)\Steam\userdata\`
const STEAM_USER_ID = "286123118"
const GAME_ID = "2868840"
const PROFILE = "profile1"

var PROFILE_SAVES_PATH = filepath.Join(STEAM_USERDATA_PATH, STEAM_USER_ID, GAME_ID, "remote", PROFILE, "saves")
var PROGRESS_SAVE_PATH = filepath.Join(PROFILE_SAVES_PATH, "progress.save")
var CURRENT_RUN_SAVE_PATH = filepath.Join(PROFILE_SAVES_PATH, "current_run.save")

func prettifyName(name string, prefix string) string {
	caser := cases.Title(language.English)
	s, _ := strings.CutPrefix(name, prefix)
	s = strings.ReplaceAll(s, "_", " ")
	return caser.String(s)
}

type CharacterWinLossObject struct {
	Character string `json:"character"`
	Losses    int    `json:"losses"`
	Wins      int    `json:"wins"`
}

func (c CharacterWinLossObject) Name() string {
	return prettifyName(c.Character, "CHARACTER.")
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
	Count  *int   `json:"count"`
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
	return prettifyName(c.Id, "ID.")
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

type ProgressSaveFile struct {
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

func readProgressSave() (ProgressSaveFile, error) {
	file, err := os.ReadFile(PROGRESS_SAVE_PATH)
	if err != nil {
		return ProgressSaveFile{}, err
	}

	progressSave := ProgressSaveFile{}
	if err := json.Unmarshal(file, &progressSave); err != nil {
		return ProgressSaveFile{}, err
	}
	return progressSave, nil
}

type RoomsObject struct {
	AncientId               string   `json:"ancient_id"`
	BossEncountersVisited   int      `json:"boss_encounters_visited"`
	BossId                  string   `json:"boss_id"`
	EliteEncounterIds       []string `json:"elite_encounter_ids"`
	EliteEncountersVisited  int      `json:"elite_encounters_visited"`
	EventIds                []string `json:"event_ids"`
	EventsVisited           int      `json:"events_visited"`
	NormalEncounterIds      []string `json:"normal_encounter_ids"`
	NormalEncountersVisited int      `json:"normal_encounters_visited"`
	SecondBossId            *string  `json:"second_boss_id"`
}

type Coord struct {
	Col int `json:"col"`
	Row int `json:"row"`
}

type MapNodeType string

const (
	MapNodeTypeAncient  MapNodeType = "ancient"
	MapNodeTypeMonster  MapNodeType = "monster"
	MapNodeTypeElite    MapNodeType = "elite"
	MapNodeTypeRestSite MapNodeType = "rest_site"
	MapNodeTypeShop     MapNodeType = "shop"
	MapNodeTypeTreasure MapNodeType = "treasure"
	MapNodeTypeUnknown  MapNodeType = "unknown"
	MapNodeTypeBoss     MapNodeType = "boss"
)

func (t *MapNodeType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case
		"ancient",
		"monster",
		"elite",
		"rest_site",
		"shop",
		"treasure",
		"unknown",
		"boss":
		*t = MapNodeType(s)
		return nil
	}
	return fmt.Errorf("invalid MapNodeType: %s", s)
}

func (t MapNodeType) String() string {
	return string(t)
}

func (t MapNodeType) Name() string {
	return prettifyName(t.String(), "")
}

type MapNode struct {
	CanModify bool        `json:"can_modify"`
	Children  []Coord     `json:"children"`
	Coord     Coord       `json:"coord"`
	Type      MapNodeType `json:"type"`
}

type SavedMap struct {
	Boss        MapNode   `json:"boss"`
	Height      int       `json:"height"`
	Points      []MapNode `json:"points"`
	Start       MapNode   `json:"start"`
	StartCoords []Coord   `json:"start_coords"`
	Width       int       `json:"width"`
}

type Act struct {
	Id       string      `json:"id"`
	Rooms    RoomsObject `json:"rooms"`
	SavedMap *SavedMap   `json:"saved_map"`
}

type Title struct {
	Key   string `json:"key"`
	Table string `json:"table"`
}

type AncientChoice struct {
	TextKey   string `json:"TextKey"`
	Title     Title  `json:"title"`
	WasChosen bool   `json:"was_chosen"`
}

func (c AncientChoice) Name() string {
	return prettifyName(c.TextKey, "")
}

type Enchantment struct {
	Amount int    `json:"amount"`
	Id     string `json:"id"`
}

type Card struct {
	CurrentUpgradeLevel *int         `json:"current_upgrade_level"`
	Enchantment         *Enchantment `json:"enchantment"`
	FloorAddedToDeck    *int         `json:"floor_added_to_deck"`
	Id                  string       `json:"id"`
}

func (c Card) Name() string {
	return prettifyName(c.Id, "CARD.")
}

type CardChoice struct {
	Card      Card `json:"card"`
	WasPicked bool `json:"was_picked"`
}

type EventChoice struct {
	Title Title `json:"title"`
	// TODO
	Variables *any `json:"variables"`
}

type PotionChoice struct {
	Choice    string `json:"choice"`
	Waspicked bool   `json:"was_picked"`
}

type RelicChoice struct {
	Choice    string `json:"choice"`
	Waspicked bool   `json:"was_picked"`
}

type CardGained struct {
	CurrentUpgradeLevel *int   `json:"current_upgrade_level"`
	Id                  string `json:"id"`
}

type CardEnchanted struct {
	Card        Card   `json:"card"`
	Enchantment string `json:"enchantment"`
}

type PlayerStatsObject struct {
	CurrentGold int `json:"current_gold"`
	CurrentHP   int `json:"current_hp"`
	DamageTaken int `json:"damage_taken"`
	GoldGained  int `json:"gold_gained"`
	GoldLost    int `json:"gold_lost"`
	GoldSpent   int `json:"gold_spent"`
	GoldStolen  int `json:"gold_stolen"`
	HPHealed    int `json:"hp_healed"`
	MaxHP       int `json:"max_hp"`
	MaxHPGained int `json:"max_hp_gained"`
	MaxHPLost   int `json:"max_hp_lost"`
	PlayerId    int `json:"player_id"`

	BoughtPotions *[]string `json:"bought_potions"`
	BoughtRelics  *[]string `json:"bought_relics"`

	AncientChoice  *[]AncientChoice `json:"ancient_choice"`
	CardChoices    *[]CardChoice    `json:"card_choices"`
	EventChoices   *[]EventChoice   `json:"event_choice"`
	PotionChoices  *[]PotionChoice  `json:"potion_choices"`
	RelicChoices   *[]RelicChoice   `json:"relic_choices"`
	CardsGained    *[]CardGained    `json:"cards_gained"`
	CardsEnchanted *[]CardEnchanted `json:"cards_enchanted"`
}

type RoomInfo struct {
	ModelId    string `json:"model_id"`
	RoomType   string `json:"room_type"`
	TurnsTaken int    `json:"turns_taken"`

	MonsterIds *[]string `json:"monster_ids"`
}

type MapPoint struct {
	MapPointType MapNodeType         `json:"map_point_type"`
	PlayerStats  []PlayerStatsObject `json:"player_stats"`
	Rooms        []RoomInfo          `json:"rooms"`
}

type Potion struct {
	Id        string `json:"id"`
	SlotIndex int    `json:"slot_index"`
}

func (p Potion) Name() string {
	return prettifyName(p.Id, "POTION.")
}

type RelicIdLists struct {
	Common   []string `json:"common"`
	Uncommon []string `json:"uncommon"`
	Rare     []string `json:"rare"`
	Shop     []string `json:"shop"`
	Event    []string `json:"event"`
	Ancient  []string `json:"ancient"`
}

type RelicGrabBag struct {
	RelicIdLists RelicIdLists `json:"relic_id_lists"`
}

type Relic struct {
	FloorAddedToDeck int    `json:"floor_added_to_deck"`
	Id               string `json:"id"`
}

func (r Relic) Name() string {
	return prettifyName(r.Id, "RELIC.")
}

type RNG struct {
	Counters map[string]int `json:"counters"` // TODO
	Seed     any            `json:"seed"`
}

type UnlockState struct {
	EncountersSeen []string `json:"encounters_seen"`
	NumberOfRuns   int      `json:"number_of_runs"`
	UnlockedEpochs []string `json:"unlocked_epochs"`
}

type Player struct {
	BaseOrbSlotCount   int                `json:"base_orb_slot_count"`
	CharacterId        string             `json:"character_id"`
	CurrentHP          int                `json:"current_hp"`
	Deck               []Card             `json:"deck"`
	ExtraFields        map[string]any     `json:"extra_fields"`
	Gold               int                `json:"gold"`
	MaxEnergy          int                `json:"max_energy"`
	MaxHP              int                `json:"max_hp"`
	MaxPotionSlotCount int                `json:"max_potion_slot_count"`
	NetId              int                `json:"net_id"`
	Odds               map[string]float64 `json:"odds"` // TODO
	Potions            []Potion           `json:"potions"`
	RelicGrabBag       RelicGrabBag       `json:"relic_grab_bag"`
	Relics             []Relic            `json:"relics"`
	RNG                RNG                `json:"rng"`
	UnlockState        UnlockState        `json:"unlock_state"`
}

func (p Player) Name() string {
	return prettifyName(p.CharacterId, "CHARACTER.")
}

func (p Player) PotionChance() float64 {
	return p.Odds["potion_reward_odds_value"]
}

type PreFinishedRoom struct {
	EncounterID             *string `json:"encounter_id"`
	EventID                 *string `json:"event_id"`
	IsPreFinished           bool    `json:"is_pre_finished"`
	ParentEventID           *string `json:"parent_event_id"`
	RewardProportion        int     `json:"reward_proportion"`
	RoomType                string  `json:"room_type"`
	ShouldResumeParentEvent bool    `json:"should_resume_parent_event"`
}

type CurrentRunSaveFile struct {
	Acts               []Act              `json:"acts"`
	Ascension          int                `json:"ascension"`
	CurrentActIndex    int                `json:"current_act_index"`
	ExtraFields        map[string]any     `json:"extra_fields"`
	GameMode           string             `json:"game_mode"`
	MapDrawings        string             `json:"map_drawings"`
	MapPointHistory    [][]MapPoint       `json:"map_point_history"`
	Modifiers          []any              `json:"modifiers"`
	Odds               map[string]float64 `json:"odds"` // TODO
	PlatformType       string             `json:"platform_type"`
	Players            []Player           `json:"players"`
	PreFinishedRoom    PreFinishedRoom    `json:"pre_finished_room"`
	RNG                RNG                `json:"rng"`
	RunTime            int                `json:"run_time"`
	SaveTime           int64              `json:"save_time"`
	SchemaVersion      int                `json:"schema_version"`
	SharedRelicGrabBag RelicGrabBag       `json:"shared_relic_grab_bag"`
	StartTime          int64              `json:"start_time"`
	VisitedMapCoords   []Coord            `json:"visited_map_coords"`
	WinTime            int                `json:"win_time"`
}

func readCurrentSave() (CurrentRunSaveFile, error) {
	file, err := os.ReadFile(CURRENT_RUN_SAVE_PATH)
	if err != nil {
		return CurrentRunSaveFile{}, err
	}

	currentRunSave := CurrentRunSaveFile{}
	if err := json.Unmarshal(file, &currentRunSave); err != nil {
		return CurrentRunSaveFile{}, err
	}
	return currentRunSave, nil
}

func (s CurrentRunSaveFile) UnknownEliteChance() float64 {
	return s.Odds["unknown_map_point_elite_odds_value"]
}

func (s CurrentRunSaveFile) UnknownMonsterChance() float64 {
	return s.Odds["unknown_map_point_monster_odds_value"]
}

func (s CurrentRunSaveFile) UnknownShopChance() float64 {
	return s.Odds["unknown_map_point_shop_odds_value"]
}

func (s CurrentRunSaveFile) UnknownTreasureChance() float64 {
	return s.Odds["unknown_map_point_treasure_odds_value"]
}

type PastRunPlayer struct {
	Badges             []Badge  `json:"badges"`
	Character          string   `json:"character"`
	Deck               []Card   `json:"deck"`
	Id                 int      `json:"id"`
	MaxPotionSlotCount int      `json:"max_potion_slot_count"`
	Potions            []Potion `json:"potions"`
	Relics             []Relic  `json:"relics"`
}

func (p PastRunPlayer) Name() string {
	return prettifyName(p.Character, "CHARACTER.")
}

type PastRunFile struct {
	Acts              []string        `json:"acts"`
	Ascension         int             `json:"ascension"`
	BuildId           string          `json:"build_id"`
	GameMode          string          `json:"game_mode"`
	KilledByEncounter string          `json:"killed_by_encounter"`
	KilledByEvent     string          `json:"killed_by_event"`
	MapPointHistory   [][]MapPoint    `json:"map_point_history"`
	Modifiers         []any           `json:"modifiers"`
	PlatformType      string          `json:"platform_type"`
	Players           []PastRunPlayer `json:"players"`
	RunTime           int             `json:"run_time"`
	SchemaVersion     int             `json:"schema_version"`
	Seed              string          `json:"seed"`
	StartTime         int             `json:"start_time"`
	WasAbandoned      bool            `json:"was_abandoned"`
	Win               bool            `json:"win"`
}

func readPastRunFile(file []byte) (PastRunFile, error) {
	pastRun := PastRunFile{}
	if err := json.Unmarshal(file, &pastRun); err != nil {
		return PastRunFile{}, err
	}
	return pastRun, nil
}
