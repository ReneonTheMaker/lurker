package main

const (
	// Item types
	ITEM_DAMAGE = iota
	ITEM_HEAL
	ITEM_SP_RESTORE
	ITEM_BASE_SP_INCREASE
	ITEM_MAX_SP_INCREASE
	ITEM_SANITY
	ITEM_BUFF_DAMAGE
	ITEM_DAMAGE_FIRE
	ITEM_DAMAGE_HOLY
	ITEM_DAMAGE_DARK
)

type Item struct {
	Name        string
	Description string
	Effects     map[int]int // e.g., {ITEM_HEAL: 10, ITEM_MANA: 5, ITEM_SANITY: 1}
	Consumable  bool
	CantDrop    bool
}
