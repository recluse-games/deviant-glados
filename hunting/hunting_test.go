package hunting

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"

	deviant "github.com/recluse-games/deviant-protobuf/genproto/go"
)

func generateHandLiterals(size int32, class deviant.Classes) *deviant.Hand {
	deckLiteral := &deviant.Hand{
		Id:    "grass_0000",
		Cards: generateCardLiterals(size, class),
	}

	return deckLiteral
}

func generateDiscardLiteral(size int32, class deviant.Classes) *deviant.Discard {
	deckLiteral := &deviant.Discard{
		Id:    "grass_0000",
		Cards: generateCardLiterals(size, class),
	}

	return deckLiteral
}

func generateDeckLiterals(size int32, class deviant.Classes) *deviant.Deck {
	deckLiteral := &deviant.Deck{
		Id:    "grass_0000",
		Cards: generateCardLiterals(size, class),
	}

	return deckLiteral
}

func generateCardLiterals(size int32, class deviant.Classes) []*deviant.Card {
	var cardLiterals = []*deviant.Card{}

	switch class {
	case deviant.Classes_WARRIOR:
		for i := int32(0); i < size; i++ {
			card := &deviant.Card{
				Id:          "attack_bash_0000",
				BackId:      "back_0000",
				InstanceId:  uuid.New().String(),
				Cost:        3,
				Damage:      2,
				Title:       "Bash",
				Flavor:      "OP Area Move",
				Description: "Something Too Broken to Be Real",
				Type:        deviant.CardType_ATTACK,
				Action: &deviant.CardAction{
					Id: uuid.New().String(),
					Pattern: []*deviant.Pattern{
						{
							Direction: deviant.Direction_DOWN,
							Distance:  3,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_LEFT,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_DOWN,
									Distance:  3,
								},
							},
						},
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_RIGHT,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_DOWN,
									Distance:  3,
								},
							},
						},
					},
				},
			}

			cardLiterals = append(cardLiterals, card)

			card = &deviant.Card{
				Id:          "attack_slash_0000",
				BackId:      "back_0000",
				InstanceId:  uuid.New().String(),
				Cost:        2,
				Damage:      2,
				Title:       "Slash",
				Flavor:      "Downward Dog",
				Description: "A Simple Slash",
				Type:        deviant.CardType_ATTACK,
				Action: &deviant.CardAction{
					Id: uuid.New().String(),
					Pattern: []*deviant.Pattern{
						{
							Direction: deviant.Direction_DOWN,
							Distance:  3,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
							},
						},
					},
				},
			}

			cardLiterals = append(cardLiterals, card)

			card = &deviant.Card{
				Id:          "block_wall_0000",
				BackId:      "back_0000",
				InstanceId:  uuid.New().String(),
				Cost:        1,
				Damage:      2,
				Title:       "Block",
				Flavor:      "A Simple Block",
				Description: "The most beautiful block",
				Type:        deviant.CardType_BLOCK,
				Action: &deviant.CardAction{
					Id: uuid.New().String(),
					Pattern: []*deviant.Pattern{
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
							},
						},
					},
				},
			}

			cardLiterals = append(cardLiterals, card)
		}
	case deviant.Classes_PRIEST:
		for i := int32(0); i < size; i++ {
			card := &deviant.Card{
				Id:          "cast_radience_0000",
				BackId:      "back_0000",
				InstanceId:  uuid.New().String(),
				Cost:        1,
				Damage:      1,
				Title:       "Radience",
				Flavor:      "A Basic Clearing Spell",
				Description: "A Basic Clearing Spell",
				Type:        deviant.CardType_ATTACK,
				Action: &deviant.CardAction{
					Id: uuid.New().String(),
					Pattern: []*deviant.Pattern{
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_RIGHT,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_LEFT,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_UP,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_UP,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_LEFT,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_UP,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_RIGHT,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_RIGHT,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_LEFT,
									Distance:  1,
								},
							},
						},
					},
				},
			}

			cardLiterals = append(cardLiterals, card)

			card = &deviant.Card{
				Id:          "cast_heal_0000",
				BackId:      "back_0000",
				InstanceId:  uuid.New().String(),
				Cost:        1,
				Damage:      2,
				Title:       "Heal",
				Flavor:      "A Basic Heal",
				Description: "A Simple Heal",
				Type:        deviant.CardType_HEAL,
				Action: &deviant.CardAction{
					Id: uuid.New().String(),
					Pattern: []*deviant.Pattern{
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
							},
						},
					},
				},
			}

			cardLiterals = append(cardLiterals, card)

			card = &deviant.Card{
				Id:          "cast_healing_ray_0000",
				BackId:      "back_0000",
				InstanceId:  uuid.New().String(),
				Cost:        3,
				Damage:      2,
				Title:       "Healing Ray",
				Flavor:      "A Basic Heal",
				Description: "A Ranged Heal",
				Type:        deviant.CardType_HEAL,
				Action: &deviant.CardAction{
					Id: uuid.New().String(),
					Pattern: []*deviant.Pattern{
						{
							Direction: deviant.Direction_DOWN,
							Distance:  3,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
							},
						},
					},
				},
			}

			cardLiterals = append(cardLiterals, card)
		}
	case deviant.Classes_MAGE:
		for i := int32(0); i < size; i++ {
			card := &deviant.Card{
				Id:          "attack_fireball_0000",
				BackId:      "back_0000",
				InstanceId:  uuid.New().String(),
				Cost:        2,
				Damage:      2,
				Title:       "Fireball",
				Flavor:      "Dunking",
				Description: "A Simple Fireball",
				Type:        deviant.CardType_ATTACK,
				Action: &deviant.CardAction{
					Id: uuid.New().String(),
					Pattern: []*deviant.Pattern{
						{
							Direction: deviant.Direction_DOWN,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  3,
								},
							},
						},
					},
				},
			}

			cardLiterals = append(cardLiterals, card)

			card = &deviant.Card{
				Id:          "attack_searing_touch_0000",
				BackId:      "back_0000",
				InstanceId:  uuid.New().String(),
				Cost:        2,
				Damage:      3,
				Title:       "Searing Touch",
				Flavor:      "Burn Baby",
				Description: "A Cross Attack",
				Type:        deviant.CardType_ATTACK,
				Action: &deviant.CardAction{
					Id: uuid.New().String(),
					Pattern: []*deviant.Pattern{
						{
							Direction: deviant.Direction_RIGHT,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_UP,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_RIGHT,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_RIGHT,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_RIGHT,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_LEFT,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_UP,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_LEFT,
									Distance:  1,
								},
							},
						},
						{
							Direction: deviant.Direction_LEFT,
							Distance:  1,
							Offset: []*deviant.Offset{
								{
									Direction: deviant.Direction_DOWN,
									Distance:  1,
								},
								{
									Direction: deviant.Direction_LEFT,
									Distance:  1,
								},
							},
						},
					},
				},
			}

			cardLiterals = append(cardLiterals, card)
		}
	}

	// Suffle the Deck
	rand.Seed(time.Now().UnixNano())
	for i := len(cardLiterals) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		cardLiterals[i], cardLiterals[j] = cardLiterals[j], cardLiterals[i]
	}

	return cardLiterals
}

// GenerateMatch Generates a new match
func generateMatch() *deviant.Encounter {
	test := &deviant.Encounter{
		Id:        "encounter_0000",
		Completed: false,
		Board: &deviant.Board{
			Tiles: &deviant.Tiles{
				Tiles: []*deviant.TilesRow{
					{
						Tiles: []*deviant.Tile{
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
						},
					},
					{
						Tiles: []*deviant.Tile{
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
						},
					},
					{
						Tiles: []*deviant.Tile{
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
						},
					},
					{
						Tiles: []*deviant.Tile{
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
						},
					},
					{
						Tiles: []*deviant.Tile{
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
						},
					},
					{
						Tiles: []*deviant.Tile{
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
						},
					},
					{
						Tiles: []*deviant.Tile{
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
						},
					},
					{
						Tiles: []*deviant.Tile{
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
						},
					},
					{
						Tiles: []*deviant.Tile{
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
							{Id: "grass_0000"},
						},
					},
				},
			},
			OverlayTiles: []*deviant.Tile{},
			Entities: &deviant.Entities{
				Entities: []*deviant.EntitiesRow{
					{
						Entities: []*deviant.Entity{
							{
								Id:         "0001",
								Name:       "Ian",
								Hp:         10,
								MaxHp:      10,
								Ap:         5,
								MaxAp:      5,
								Alignment:  deviant.Alignment_FRIENDLY,
								Class:      deviant.Classes_WARRIOR,
								Hand:       generateHandLiterals(0, deviant.Classes_WARRIOR),
								Deck:       generateDeckLiterals(10, deviant.Classes_WARRIOR),
								Discard:    generateDiscardLiteral(0, deviant.Classes_WARRIOR),
								Initiative: 5,
								OwnerId:    "0001",
								Rotation:   deviant.EntityRotationNames_WEST,
							},
							{},
							{},
							{
								Id:         "0003",
								Name:       "Zach",
								Hp:         10,
								MaxHp:      10,
								Ap:         5,
								MaxAp:      5,
								Alignment:  deviant.Alignment_FRIENDLY,
								Class:      deviant.Classes_MAGE,
								Hand:       generateHandLiterals(0, deviant.Classes_MAGE),
								Deck:       generateDeckLiterals(10, deviant.Classes_MAGE),
								Discard:    generateDiscardLiteral(0, deviant.Classes_MAGE),
								Initiative: 5,
								OwnerId:    "0001",
								Rotation:   deviant.EntityRotationNames_NORTH,
							},
							{
								Id:         "0006",
								Name:       "Chris",
								Hp:         10,
								MaxHp:      10,
								Ap:         5,
								MaxAp:      5,
								Alignment:  deviant.Alignment_FRIENDLY,
								Class:      deviant.Classes_PRIEST,
								Hand:       generateHandLiterals(0, deviant.Classes_PRIEST),
								Deck:       generateDeckLiterals(10, deviant.Classes_PRIEST),
								Discard:    generateDiscardLiteral(0, deviant.Classes_PRIEST),
								Initiative: 5,
								OwnerId:    "0001",
								Rotation:   deviant.EntityRotationNames_NORTH,
							},
							{},
							{},
							{},
						},
					},
					{
						Entities: []*deviant.Entity{
							{},
							{},
							{},
							{},
							{},
							{},
							{},
							{},
						},
					},
					{
						Entities: []*deviant.Entity{
							{},
							{},
							{},
							{},
							{},
							{},
							{},
							{},
						},
					},
					{
						Entities: []*deviant.Entity{
							{},
							{},
							{},
							{},
							{},
							{},
							{},
							{},
						},
					},
					{
						Entities: []*deviant.Entity{
							{
								Id:        uuid.New().String(),
								Name:      "Wall",
								Hp:        2,
								MaxHp:     2,
								Class:     deviant.Classes_WALL,
								State:     deviant.EntityStateNames_IDLE,
								Alignment: deviant.Alignment_NEUTRAL,
							},
							{
								Id:        uuid.New().String(),
								Name:      "Wall",
								Hp:        2,
								MaxHp:     2,
								Class:     deviant.Classes_WALL,
								State:     deviant.EntityStateNames_IDLE,
								Alignment: deviant.Alignment_NEUTRAL,
							},
							{
								Id:        uuid.New().String(),
								Name:      "Wall",
								Hp:        2,
								MaxHp:     2,
								Class:     deviant.Classes_WALL,
								State:     deviant.EntityStateNames_IDLE,
								Alignment: deviant.Alignment_NEUTRAL,
							},
							{
								Id:        uuid.New().String(),
								Name:      "Wall",
								Hp:        2,
								MaxHp:     2,
								State:     deviant.EntityStateNames_IDLE,
								Class:     deviant.Classes_WALL,
								Alignment: deviant.Alignment_NEUTRAL,
							},
							{
								Id:        uuid.New().String(),
								Name:      "Wall",
								Hp:        2,
								MaxHp:     2,
								Class:     deviant.Classes_WALL,
								State:     deviant.EntityStateNames_IDLE,
								Alignment: deviant.Alignment_NEUTRAL,
							},
							{
								Id:        uuid.New().String(),
								Name:      "Wall",
								Hp:        2,
								MaxHp:     2,
								Class:     deviant.Classes_WALL,
								State:     deviant.EntityStateNames_IDLE,
								Alignment: deviant.Alignment_NEUTRAL,
							},
							{
								Id:        uuid.New().String(),
								Name:      "Wall",
								Hp:        2,
								MaxHp:     2,
								Class:     deviant.Classes_WALL,
								State:     deviant.EntityStateNames_IDLE,
								Alignment: deviant.Alignment_NEUTRAL,
							},
							{
								Id:        uuid.New().String(),
								Name:      "Wall",
								Hp:        2,
								MaxHp:     2,
								Class:     deviant.Classes_WALL,
								State:     deviant.EntityStateNames_IDLE,
								Alignment: deviant.Alignment_NEUTRAL,
							},
						},
					},
					{
						Entities: []*deviant.Entity{
							{},
							{},
							{},
							{},
							{},
							{},
							{},
							{},
						},
					},
					{
						Entities: []*deviant.Entity{
							{},
							{},
							{},
							{},
							{},
							{},
							{},
							{},
						},
					},
					{
						Entities: []*deviant.Entity{
							{},
							{},
							{},
							{},
							{},
							{},
							{},
							{},
						},
					},
					{
						Entities: []*deviant.Entity{
							{
								Id:         "0002",
								Name:       "Cameron",
								Hp:         10,
								MaxHp:      10,
								Ap:         5,
								MaxAp:      5,
								Alignment:  deviant.Alignment_UNFRIENDLY,
								Class:      deviant.Classes_WARRIOR,
								Hand:       generateHandLiterals(0, deviant.Classes_WARRIOR),
								Deck:       generateDeckLiterals(10, deviant.Classes_WARRIOR),
								Discard:    generateDiscardLiteral(0, deviant.Classes_WARRIOR),
								Initiative: 5,
								OwnerId:    "0001",
								Rotation:   deviant.EntityRotationNames_SOUTH,
							},
							{},
							{
								Id:         "0004",
								Name:       "Matt",
								Hp:         10,
								MaxHp:      10,
								Ap:         5,
								MaxAp:      5,
								Alignment:  deviant.Alignment_UNFRIENDLY,
								Class:      deviant.Classes_MAGE,
								Hand:       generateHandLiterals(0, deviant.Classes_MAGE),
								Deck:       generateDeckLiterals(10, deviant.Classes_MAGE),
								Discard:    generateDiscardLiteral(0, deviant.Classes_MAGE),
								Initiative: 5,
								OwnerId:    "0001",
								Rotation:   deviant.EntityRotationNames_SOUTH,
							},
							{
								Id:         "0005",
								Name:       "Ben",
								Hp:         10,
								MaxHp:      10,
								Ap:         5,
								MaxAp:      5,
								Alignment:  deviant.Alignment_UNFRIENDLY,
								Class:      deviant.Classes_PRIEST,
								Hand:       generateHandLiterals(0, deviant.Classes_PRIEST),
								Deck:       generateDeckLiterals(10, deviant.Classes_PRIEST),
								Discard:    generateDiscardLiteral(0, deviant.Classes_PRIEST),
								Initiative: 5,
								OwnerId:    "0001",
								Rotation:   deviant.EntityRotationNames_SOUTH,
							},
							{},
							{},
							{},
							{},
						},
					},
				},
			},
		},
		ActiveEntity: &deviant.Entity{
			Id:         "0001",
			Name:       "Ian",
			Hp:         10,
			MaxHp:      10,
			Ap:         5,
			MaxAp:      5,
			Alignment:  deviant.Alignment_FRIENDLY,
			Class:      deviant.Classes_WARRIOR,
			Hand:       generateHandLiterals(1, deviant.Classes_WARRIOR),
			Deck:       generateDeckLiterals(10, deviant.Classes_WARRIOR),
			Discard:    generateDiscardLiteral(0, deviant.Classes_WARRIOR),
			Initiative: 5,
			OwnerId:    "0001",
			Rotation:   deviant.EntityRotationNames_WEST,
		},
		ActiveEntityOrder: []string{"0001", "0002", "0003", "0004", "0006", "0005"},
		Turn: &deviant.Turn{
			Id:    "turn_0000",
			Phase: deviant.TurnPhaseNames_PHASE_POINT,
		},
	}

	return test
}

func TestGenerateCardVertexPairs(t *testing.T) {
	match := generateMatch()
	startingVertex := &gridNode{
		X: 0,
		Y: 0,
	}

	cardVertexPairsNorth := GenerateCardVertexPairs(startingVertex, match.ActiveEntity, match.Board.Entities.Entities, deviant.EntityRotationNames_NORTH)
	cardVertexPairsSouth := GenerateCardVertexPairs(startingVertex, match.ActiveEntity, match.Board.Entities.Entities, deviant.EntityRotationNames_SOUTH)
	cardVertexPairsEast := GenerateCardVertexPairs(startingVertex, match.ActiveEntity, match.Board.Entities.Entities, deviant.EntityRotationNames_EAST)
	cardVertexPairsWest := GenerateCardVertexPairs(startingVertex, match.ActiveEntity, match.Board.Entities.Entities, deviant.EntityRotationNames_WEST)

	// This is a super low value test this should be replaced with one verifying the contents of each rotation based on a static match.
	if len(cardVertexPairsNorth) <= 0 {
		t.Fail()
	}
	if len(cardVertexPairsSouth) <= 0 {
		t.Fail()
	}
	if len(cardVertexPairsEast) <= 0 {
		t.Fail()
	}
	if len(cardVertexPairsWest) <= 0 {
		t.Fail()
	}
}

func TestGenerateAllLocationMoveCombinations(t *testing.T) {
	match := generateMatch()

	GenerateAllLocationMoveCombinations(match.ActiveEntity, match.Board.Entities.Entities, match)
}

func TestFilterCardPlaysToHits(t *testing.T) {
	match := generateMatch()

	GenerateAllLocationMoveCombinations(match.ActiveEntity, match.Board.Entities.Entities, match)
	FilterCardPlaysToHits(match.Board.Entities.Entities, match, deviant.Alignment_NEUTRAL)

}

func TestSortCardPlaysByDamageInflicted(t *testing.T) {
	match := generateMatch()

	allHittingMoveCombinations := FilterCardPlaysToHits(match.Board.Entities.Entities, match, deviant.Alignment_NEUTRAL)
	bestMovesInDamageOrder := SortCardPlaysByDamageInflicted(allHittingMoveCombinations, match.Board.Entities)

	for _, vertex := range bestMovesInDamageOrder {
		t.Log(vertex.cardVertexPair.card.Id)
		t.Log(vertex.damage)
		t.Log(vertex.origin)
		t.Log(vertex.rotation)
	}
}

func TestGetPlayThatDealsTheMostDamageToTheLowestHealthTargets(t *testing.T) {
	match := generateMatch()

	allHittingMoveCombinations := FilterCardPlaysToHits(match.Board.Entities.Entities, match, deviant.Alignment_NEUTRAL)
	bestMovesInDamageOrder := SortCardPlaysByDamageInflicted(allHittingMoveCombinations, match.Board.Entities)
	entityLocationVertexPairs := GenerateEntityLocationPairs(deviant.Alignment_NEUTRAL, match.Board.Entities.Entities)
	theBestPlay := GetPlayThatDealsTheMostDamageToTheLowestHealthTargets(bestMovesInDamageOrder, entityLocationVertexPairs)

	t.Log(theBestPlay.cardVertexPair.card.Id)
	t.Log(theBestPlay.damage)
	t.Log(theBestPlay.origin)
	t.Log(theBestPlay.rotation)
}
