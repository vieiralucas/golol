package golol

import (
	"errors"
	"fmt"
)

type Stats struct {
	HitPoints              float64 `json:"hp"`
	HitPointsPerLevel      float64 `json:"hpperlevel"`
	HitPointsRegen         float64 `json:"hpregen"`
	HitPointsRegenPerLevel float64 `json:"hpregenperlevel"`
	Mana                   float64 `json:"mp"`
	ManaPerLevel           float64 `json:"mpperlevel"`
	ManaRegen              float64 `json:"mpregen"`
	ManaRegenPerLevel      float64 `json:"mpregenperlevel"`
	Armor                  float64 `json:"armor"`
	ArmorPerLevel          float64 `json:"armorperlevel"`
	MagicResist            float64 `json:"spellblock"`
	MagicResistPerLevel    float64 `json:"spellblockperlevel"`
	AttackDamage           float64 `json:"attackdamage"`
	AttackDamagePerLevel   float64 `json:"attackdamageperlevel"`
	CriticalChance         float64 `json:"crit"`
	CriticalChancePerLevel float64 `json:"critperlevel"`
	AttackSpeedDelay       float64 `json:"attackspeedoffset"`
	AttackSpeedPerLevel    float64 `json:"attackspeedperlevel"`
	MovementSpeed          float64 `json:"movespeed"`
	AttackRange            float64 `json:"attackrange"`
}

func (s *Stats) AttackSpeed() float64 {
	return 0.625 / (1 + s.AttackSpeedDelay)
}

type Info struct {
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
	Attack     int `json:"attack"`
}

type Passive struct {
	Name                 string `json:"name"`
	SanitizedDescription string `json:"sanitizedDescription"`
	Description          string `json:"description"`
	Image                Image  `json:"image"`
}

type Champion struct {
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Key       string   `json:"key"`
	Blurb     string   `json:"blurb"`
	Lore      string   `json:"lore"`
	PartyType string   `json:"partype"`
	Tags      []string `json:"tags"`
	EnemyTips []string `json:"enemytips"`
	AllyTips  []string `json:"allyTips"`
	Stats     Stats    `json:"stats"`
	Info      Info     `json:"info"`
	Image     Image    `json:"image"`
}

type Champions struct {
	Data map[string]Champion `json:"data"`
}

var champions = make(map[string]map[string]Champion, 11)

func GetChampions(region string) (map[string]Champion, error) {
	if CACHE && champions[region] != nil {
		return champions[region], nil
	}

	url := fmt.Sprintf(
		"%v%v/v1.2/champion?champData=allytips,altimages,blurb,enemytips,image,info,lore,partype,passive,stats,tags&api_key=%v",
		STATIC_DATA_URL,
		region,
		API_KEY,
	)

	cs := Champions{}
	err := request(url, &cs)

	if CACHE {
		champions[region] = cs.Data
	}

	return cs.Data, err
}

func GetChampionById(id int, region string) (Champion, error) {
	if _, found := champions[region]; CACHE && found {
		for _, c := range champions[region] {
			if c.Id == id {
				return c, nil
			}
		}
	}

	url := fmt.Sprintf(
		"%v%v/v1.2/champion/%v?champData=allytips,altimages,blurb,enemytips,image,info,lore,partype,passive,stats,tags&api_key=%v",
		STATIC_DATA_URL,
		region,
		id,
		API_KEY,
	)

	c := Champion{}
	err := request(url, &c)

	if err != nil && CACHE {
		if _, found := champions[region]; found {
			champions[region][c.Name] = c
		} else {
			champions[region] = make(map[string]Champion)
			champions[region][c.Name] = c
		}
	}

	return c, err
}

func GetChampionByName(n string, region string) (Champion, error) {
	if _, found := champions[region][n]; CACHE && found {
		return champions[region][n], nil
	}

	var cs map[string]Champion

	cs, err := GetChampions(region)

	if err != nil {
		return cs[n], err
	}

	if _, found := cs[n]; found {
		return cs[n], nil
	}

	return cs[n], errors.New("Couldn't find a champion named: " + n)
}
