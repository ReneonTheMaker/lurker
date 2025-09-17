package main

var xpTable = []int{
	0,    // Level 1
	100,  // Level 2
	300,  // Level 3
	600,  // Level 4
	1000, // Level 5
	1500, // Level 6
	2100, // Level 7
	2800, // Level 8
	3600, // Level 9
	4500, // Level 10
}

type Character struct {
	Level            int
	CurrentHp        int
	BaseMaxHp        int
	MaxHp            int
	CurrentSp        int
	BaseMaxSp        int
	MaxSp            int
	Strength         int
	BaseStrength     int
	Speed            int
	BaseSpeed        int
	Intelligence     int
	BaseIntelligence int
	Experience       int
	NextLevelExp     int
}

func NewCharacter() Character {
	return Character{
		Level:            1,
		CurrentHp:        20,
		BaseMaxHp:        20,
		MaxHp:            20,
		CurrentSp:        10,
		BaseMaxSp:        40,
		MaxSp:            40,
		Strength:         5,
		BaseStrength:     5,
		Speed:            5,
		BaseSpeed:        5,
		Intelligence:     5,
		BaseIntelligence: 5,
		Experience:       0,
		NextLevelExp:     100,
	}
}
