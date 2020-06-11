package hunting

import (
	"fmt"
	"log"
	"math"
	"sort"

	deviant "github.com/recluse-games/deviant-protobuf/genproto/go"
)

type manhattenPair struct {
	distance int
	X        int32
	Y        int32
}

type gridNode struct {
	Id     string
	X      int32
	Y      int32
	apCost int
}

type gridLocation struct {
	X float64
	Y float64
}

// Vertex a vertex represents a point in our grid.
type Vertex struct {
	X      int
	Y      int
	apCost int
}

// EntityVertexPair A mapping of entities and vertexs that they exist at.
type EntityVertexPair struct {
	entity *deviant.Entity
	vertex *Vertex
}

// CardVertexPair A mapping of cards and the vertex they would be played at.
type CardVertexPair struct {
	card   *deviant.Card
	vertex *Vertex
}

type CardVertexRotationPair struct {
	deaths         int
	damage         int
	origin         *Vertex
	cardVertexPair *CardVertexPair
	rotation       deviant.EntityRotationNames
}

// GenerateEntityLocationPairs Returns the entities of a certain alignment and their locations in points.
func GenerateEntityLocationPairs(alignment deviant.Alignment, entities []*deviant.EntitiesRow) []*EntityVertexPair {
	entityVertexPairs := []*EntityVertexPair{}

	for y, entityRow := range entities {
		for x, entity := range entityRow.Entities {
			if entity.Id != "" && entity.Alignment == alignment {
				entityVertexPair := &EntityVertexPair{
					entity: entity,
					vertex: &Vertex{
						X: y,
						Y: x,
					},
				}

				entityVertexPairs = append(entityVertexPairs, entityVertexPair)
			}
		}
	}

	return entityVertexPairs
}

// GetEntityVertex Get the location of yourself
func GetEntityVertex(desiredEntity *deviant.Entity, entities []*deviant.EntitiesRow) *Vertex {
	for y, entityRow := range entities {
		for x, entity := range entityRow.Entities {
			if desiredEntity.Id == entity.Id {
				entityVertexPair := &Vertex{
					X: y,
					Y: x,
				}

				return entityVertexPair
			}
		}
	}
	return nil
}

// Generate a list of all avaliable locations based on AP cost - FloodFill
func floodFill(startx int32, starty int32, x int32, y int32, filledID string, blockedID string, limit int32, tiles []*[]*gridNode) {
	if (*tiles[x])[y].Id != blockedID && (*tiles[x])[y].Id != filledID {
		var apCostX int32
		var apCostY int32

		if startx > x {
			apCostX = startx - x
		} else if startx < x {
			apCostX = x - startx
		} else {
			apCostX = 0
		}

		if starty > y {
			apCostY = starty - y
		} else if starty < y {
			apCostY = y - starty
		} else {
			apCostY = 0
		}

		newTile := &gridNode{}
		newTile.X = int32(x)
		newTile.Y = int32(y)
		newTile.Id = filledID
		newTile.apCost = int(apCostX + apCostY)
		(*tiles[x])[y] = newTile

		if limit-apCostX-apCostY >= 0 {
			if x+1 <= 8 {
				floodFill(startx, starty, x+1, y, filledID, blockedID, limit, tiles)
			}

			if y+1 <= 7 {
				floodFill(startx, starty, x, y+1, filledID, blockedID, limit, tiles)
			}

			if x-1 >= 0 {
				floodFill(startx, starty, x-1, y, filledID, blockedID, limit, tiles)
			}

			if y-1 >= 0 {
				floodFill(startx, starty, x, y-1, filledID, blockedID, limit, tiles)
			}
		}
	}
}

// GeneratePermissableMoves Generate a list of permissable moves.
func GeneratePermissableMoves(origin *gridNode, avaliableAp int32, entities *deviant.Entities) []*gridNode {
	finalTiles := []*gridNode{}
	moveTargetTiles := []*[]*gridNode{}

	for y := 0; y < len(entities.Entities); y++ {
		newRow := []*gridNode{}

		for x := 0; x < len(entities.Entities[y].Entities); x++ {
			newTile := gridNode{}
			newTile.X = int32(y)
			newTile.Y = int32(x)

			if entities.Entities[y].Entities[x].Id != "" {
				newTile.Id = "select_0002"
			}

			if int32(y) == origin.X && int32(x) == origin.Y {
				newTile.Id = "select_0001"
			}

			newRow = append(newRow, &newTile)
		}

		moveTargetTiles = append(moveTargetTiles, &newRow)
	}

	floodFill(origin.X, origin.Y, origin.X, origin.Y, "select_0000", "select_0002", avaliableAp, moveTargetTiles)

	for _, row := range moveTargetTiles {
		for _, tile := range *row {
			if (*tile).Id == "select_0000" {
				finalTiles = append(finalTiles, tile)
			}
		}
	}

	return finalTiles
}

// GenerateValidMoveVertexes Generate list of vertexs
func GenerateValidMoveVertexes(entity *deviant.Entity, entities []*deviant.EntitiesRow, encounter *deviant.Encounter) []*gridNode {
	entityVertex := GetEntityVertex(entity, entities)
	entityGraphNode := &gridNode{
		X: int32(entityVertex.X),
		Y: int32(entityVertex.Y),
	}

	validTiles := GeneratePermissableMoves(entityGraphNode, entity.Ap, encounter.Board.Entities)

	return validTiles
}

// Generate list of all avaliable cards in hand in each rotation at current location based on avaliable AP
func rotateTilePatterns(ocx float64, ocy float64, px float64, py float64, rotationAngle float64) *gridLocation {
	var radians = (math.Pi / 180) * rotationAngle
	var s = math.Sin(radians)
	var c = math.Cos(radians)

	// translate point back to origin:
	px -= ocx
	py -= ocy

	// rotate point
	var xnew = px*c - py*s
	var ynew = px*s + py*c

	// translate point back:
	px = xnew + ocx
	py = ynew + ocy

	return &gridLocation{px, py}
}

func convertDirectionToDegree(characterRotation deviant.EntityRotationNames) float64 {
	switch characterRotation {
	case deviant.EntityRotationNames_NORTH:
		return 180.00
	case deviant.EntityRotationNames_SOUTH:
		return 0.00
	case deviant.EntityRotationNames_EAST:
		return 270.00
	case deviant.EntityRotationNames_WEST:
		return 90.00
	}

	return 0.00
}

func generateOffSetVertex(startingLocation *Vertex, offSet *deviant.Offset) *Vertex {
	switch offSet.Direction {
	case deviant.Direction_UP:
		newVertex := &Vertex{
			X: startingLocation.X,
			Y: startingLocation.Y,
		}

		for i := 0; int32(i) < offSet.Distance; i++ {
			newVertex.X = newVertex.X + 1
		}

		return newVertex
	case deviant.Direction_DOWN:
		newVertex := &Vertex{
			X: startingLocation.X,
			Y: startingLocation.Y,
		}

		for i := 0; int32(i) < offSet.Distance; i++ {
			newVertex.X = newVertex.X - 1
		}

		return newVertex
	case deviant.Direction_LEFT:
		newVertex := &Vertex{
			X: startingLocation.X,
			Y: startingLocation.Y,
		}

		for i := 0; int32(i) < offSet.Distance; i++ {
			newVertex.Y = newVertex.Y + 1
		}

		return newVertex
	case deviant.Direction_RIGHT:
		newVertex := &Vertex{
			X: startingLocation.X,
			Y: startingLocation.Y,
		}

		for i := 0; int32(i) < offSet.Distance; i++ {
			newVertex.Y = newVertex.Y - 1
		}

		return newVertex
	}

	return nil
}

func generatePatternVertex(startingLocation *gridNode, offSet *Vertex, direction deviant.Direction, distance int32) []*Vertex {
	patternVertexes := []*Vertex{}

	switch direction {
	case deviant.Direction_UP:
		newVertex := &Vertex{}
		for i := 0; int32(i) < distance; i++ {
			newVertex = &Vertex{
				X: int(startingLocation.X) + offSet.X,
				Y: int(startingLocation.Y) + offSet.Y,
			}

			patternVertexes = append(patternVertexes, newVertex)
			offSet.X = offSet.X + 1
		}

		return patternVertexes
	case deviant.Direction_DOWN:
		newVertex := &Vertex{}
		for i := 0; int32(i) < distance; i++ {
			newVertex = &Vertex{
				X: int(startingLocation.X) + offSet.X,
				Y: int(startingLocation.Y) + offSet.Y,
			}

			patternVertexes = append(patternVertexes, newVertex)
			offSet.X = offSet.X + -1
		}

		return patternVertexes
	case deviant.Direction_LEFT:
		newVertex := &Vertex{}
		for i := 0; int32(i) < distance; i++ {
			newVertex = &Vertex{
				X: int(startingLocation.X) + offSet.X,
				Y: int(startingLocation.Y) + offSet.Y,
			}

			patternVertexes = append(patternVertexes, newVertex)
			offSet.Y = offSet.Y + 1
		}

		return patternVertexes
	case deviant.Direction_RIGHT:
		newVertex := &Vertex{}
		for i := 0; int32(i) < distance; i++ {
			newVertex = &Vertex{
				X: int(startingLocation.X) + offSet.X,
				Y: int(startingLocation.Y) + offSet.Y,
			}

			patternVertexes = append(patternVertexes, newVertex)
			offSet.Y = offSet.Y + -1
		}

		return patternVertexes
	}

	return nil
}

func GenerateCardVertexPair(card *deviant.Card, location *gridNode, entity *deviant.Entity, entities []*deviant.EntitiesRow, rotation deviant.EntityRotationNames) []*CardVertexRotationPair {
	nonRotatedCardVertexPairs := []*CardVertexPair{}
	rotatedCardVertexPairs := []*CardVertexRotationPair{}

	for _, play := range card.Action.Pattern {
		offsetVertex := &Vertex{
			X: 0,
			Y: 0,
		}

		for _, offSet := range play.Offset {
			offsetVertex = generateOffSetVertex(offsetVertex, offSet)
		}

		playPatternVertexes := generatePatternVertex(location, offsetVertex, play.Direction, play.Distance)

		for _, vertex := range playPatternVertexes {
			cardVertexPair := &CardVertexPair{
				card:   card,
				vertex: vertex,
			}

			nonRotatedCardVertexPairs = append(nonRotatedCardVertexPairs, cardVertexPair)
		}
	}

	for _, cardVertexPair := range nonRotatedCardVertexPairs {

		// CAUTION: HACK - This logic should be moved somewhere else to apply rotations directly to the cards themselves maybe?
		var rotationDegree = convertDirectionToDegree(rotation)
		var rotatedPlayPair = rotateTilePatterns(float64(location.X), float64(location.Y), float64(cardVertexPair.vertex.X), float64(cardVertexPair.vertex.Y), rotationDegree)

		var x = int(math.RoundToEven(rotatedPlayPair.X))
		var y = int(math.RoundToEven(rotatedPlayPair.Y))

		rotatedCardVertexPair := &CardVertexRotationPair{
			cardVertexPair: &CardVertexPair{
				vertex: &Vertex{
					X: x,
					Y: y,
				},
				card: cardVertexPair.card,
			},
			rotation: rotation,
		}

		rotatedCardVertexPairs = append(rotatedCardVertexPairs, rotatedCardVertexPair)
	}

	return rotatedCardVertexPairs
}

func GenerateCardVertexPairs(location *gridNode, entity *deviant.Entity, entities []*deviant.EntitiesRow, rotation deviant.EntityRotationNames) []*CardVertexRotationPair {
	nonRotatedCardVertexPairs := []*CardVertexPair{}
	rotatedCardVertexPairs := []*CardVertexRotationPair{}

	for _, card := range entity.Hand.Cards {
		if card.Cost <= entity.Ap-int32(location.apCost) {
			for _, play := range card.Action.Pattern {
				offsetVertex := &Vertex{
					X: 0,
					Y: 0,
				}

				for _, offSet := range play.Offset {
					offsetVertex = generateOffSetVertex(offsetVertex, offSet)
				}

				playPatternVertexes := generatePatternVertex(location, offsetVertex, play.Direction, play.Distance)

				for _, vertex := range playPatternVertexes {
					cardVertexPair := &CardVertexPair{
						card:   card,
						vertex: vertex,
					}

					nonRotatedCardVertexPairs = append(nonRotatedCardVertexPairs, cardVertexPair)
				}
			}
		}
	}

	for _, cardVertexPair := range nonRotatedCardVertexPairs {

		// CAUTION: HACK - This logic should be moved somewhere else to apply rotations directly to the cards themselves maybe?
		var rotationDegree = convertDirectionToDegree(rotation)
		var rotatedPlayPair = rotateTilePatterns(float64(location.X), float64(location.Y), float64(cardVertexPair.vertex.X), float64(cardVertexPair.vertex.Y), rotationDegree)

		var x = int(math.RoundToEven(rotatedPlayPair.X))
		var y = int(math.RoundToEven(rotatedPlayPair.Y))

		rotatedCardVertexPair := &CardVertexRotationPair{
			cardVertexPair: &CardVertexPair{
				vertex: &Vertex{
					X: x,
					Y: y,
				},
				card: cardVertexPair.card,
			},
			rotation: rotation,
		}

		rotatedCardVertexPairs = append(rotatedCardVertexPairs, rotatedCardVertexPair)
	}

	return rotatedCardVertexPairs
}

// Generate a list of all plays at all locations with avaliable AP.
func GenerateAllLocationMoveCombinations(entity *deviant.Entity, entities []*deviant.EntitiesRow, encounter *deviant.Encounter) []*CardVertexRotationPair {
	cardVertexRotationPairs := []*CardVertexRotationPair{}

	entityRotations := []deviant.EntityRotationNames{
		deviant.EntityRotationNames_NORTH, deviant.EntityRotationNames_SOUTH, deviant.EntityRotationNames_EAST, deviant.EntityRotationNames_WEST,
	}
	validMoveVertexes := GenerateValidMoveVertexes(entity, entities, encounter)

	for _, moveVertex := range validMoveVertexes {
		for _, rotation := range entityRotations {
			generatedPairs := GenerateCardVertexPairs(moveVertex, entity, entities, rotation)

			for _, generatedPair := range generatedPairs {
				// Include the origin point along with the rotation and card for this move.
				generatedPair.origin = &Vertex{
					X:      int(moveVertex.X),
					Y:      int(moveVertex.Y),
					apCost: moveVertex.apCost,
				}

				cardVertexRotationPairs = append(cardVertexRotationPairs, generatedPair)
			}
		}
	}

	return cardVertexRotationPairs
}

// Filter list of all plays down to plays which hit an enemy

func FilterCardPlaysToHits(entities []*deviant.EntitiesRow, encounter *deviant.Encounter, alignment deviant.Alignment) []*CardVertexRotationPair {

	locationMoveCombinationsThatHit := []*CardVertexRotationPair{}
	entityLocationPairs := GenerateEntityLocationPairs(alignment, entities)
	allLocationMoveCombinations := GenerateAllLocationMoveCombinations(encounter.ActiveEntity, entities, encounter)

	for _, locationMoveCombination := range allLocationMoveCombinations {
		for _, entityLocationPair := range entityLocationPairs {
			if locationMoveCombination.cardVertexPair.vertex.X == entityLocationPair.vertex.X && locationMoveCombination.cardVertexPair.vertex.Y == entityLocationPair.vertex.Y {
				locationMoveCombinationsThatHit = append(locationMoveCombinationsThatHit, locationMoveCombination)
			}
		}
	}

	return locationMoveCombinationsThatHit
}

// Determine HP of each enemy after each hit

func SortCardPlaysByDamageInflicted(cardVertexRotationPairs []*CardVertexRotationPair, entities *deviant.Entities) []*CardVertexRotationPair {

	// This logic isn't actually the correct damage number but it is
	for _, cardVertexRotationPair := range cardVertexRotationPairs {
		effectiveDamage := entities.Entities[cardVertexRotationPair.cardVertexPair.vertex.X].Entities[cardVertexRotationPair.cardVertexPair.vertex.Y].MaxHp - (entities.Entities[cardVertexRotationPair.cardVertexPair.vertex.X].Entities[cardVertexRotationPair.cardVertexPair.vertex.Y].Hp - cardVertexRotationPair.cardVertexPair.card.Damage)

		// Remove overkill damage
		if effectiveDamage < 0 {
			effectiveDamage = 0
		}

		for _, otherCardVertexRotationPairs := range cardVertexRotationPairs {
			if cardVertexRotationPair.cardVertexPair.card.InstanceId == otherCardVertexRotationPairs.cardVertexPair.card.InstanceId && otherCardVertexRotationPairs.origin.X == cardVertexRotationPair.origin.X && otherCardVertexRotationPairs.origin.Y == cardVertexRotationPair.origin.Y && otherCardVertexRotationPairs.rotation == cardVertexRotationPair.rotation && (cardVertexRotationPair.cardVertexPair.vertex.X != otherCardVertexRotationPairs.cardVertexPair.vertex.X || otherCardVertexRotationPairs.cardVertexPair.vertex.Y != cardVertexRotationPair.cardVertexPair.vertex.Y) {
				newDamage := entities.Entities[otherCardVertexRotationPairs.cardVertexPair.vertex.X].Entities[otherCardVertexRotationPairs.cardVertexPair.vertex.Y].MaxHp - (entities.Entities[otherCardVertexRotationPairs.cardVertexPair.vertex.X].Entities[otherCardVertexRotationPairs.cardVertexPair.vertex.Y].Hp - otherCardVertexRotationPairs.cardVertexPair.card.Damage)

				// Remove overkill damage
				if newDamage < 0 {
					newDamage = 0
				}

				effectiveDamage += newDamage
			}
		}

		cardVertexRotationPair.damage = int(effectiveDamage)
	}

	// Sort the slice by highest damage
	sort.SliceStable(cardVertexRotationPairs, func(i, j int) bool { return cardVertexRotationPairs[i].damage > cardVertexRotationPairs[j].damage })

	return cardVertexRotationPairs
}

// Sort list by lowest enemy HP
func GetPlayThatDealsTheMostDamageToTheLowestHealthTargets(cardVertexRotationPairs []*CardVertexRotationPair, entityVertexPairs []*EntityVertexPair) *CardVertexRotationPair {
	sortedEntityVertexes := entityVertexPairs
	sort.SliceStable(sortedEntityVertexes, func(i, j int) bool { return sortedEntityVertexes[i].entity.Hp < sortedEntityVertexes[j].entity.Hp })

	for _, lowestHealthEnemy := range sortedEntityVertexes {
		for _, cardVertexRotationPair := range cardVertexRotationPairs {
			if cardVertexRotationPair.cardVertexPair.vertex.X == lowestHealthEnemy.vertex.X && cardVertexRotationPair.cardVertexPair.vertex.Y == lowestHealthEnemy.vertex.Y {
				cardMovePlayThatDealsTheMostDamageToTheLowestHealthTargets := cardVertexRotationPair
				return cardMovePlayThatDealsTheMostDamageToTheLowestHealthTargets
			}
		}
	}

	return nil
}

func GenerateMoveAction(cardVertexRotationPair *CardVertexRotationPair, encounter *deviant.Encounter) *deviant.EncounterRequest {
	entityMoveAction := &deviant.EntityMoveAction{
		StartXPosition: int32(GetEntityVertex(encounter.ActiveEntity, encounter.Board.Entities.Entities).X),
		StartYPosition: int32(GetEntityVertex(encounter.ActiveEntity, encounter.Board.Entities.Entities).Y),
		FinalXPosition: int32(cardVertexRotationPair.origin.X),
		FinalYPosition: int32(cardVertexRotationPair.origin.Y),
	}

	encounterRequest := &deviant.EncounterRequest{
		PlayerId:         encounter.ActiveEntity.Id,
		EntityActionName: deviant.EntityActionNames_MOVE,
		EntityMoveAction: entityMoveAction,
	}

	return encounterRequest
}

func GenerateTargetAction(cardVertexRotationPair *CardVertexRotationPair, encounter *deviant.Encounter) *deviant.EncounterRequest {
	targettedTiles := []*deviant.Tile{}

	gridNode := &gridNode{
		apCost: cardVertexRotationPair.cardVertexPair.vertex.apCost,
		X:      int32(cardVertexRotationPair.origin.X),
		Y:      int32(cardVertexRotationPair.origin.Y),
	}

	generatedPairs := GenerateCardVertexPair(cardVertexRotationPair.cardVertexPair.card, gridNode, encounter.ActiveEntity, encounter.Board.Entities.Entities, cardVertexRotationPair.rotation)

	for _, generatedPair := range generatedPairs {
		// Include the origin point along with the rotation and card for this move.
		tile := &deviant.Tile{
			X:  int32(generatedPair.cardVertexPair.vertex.X),
			Y:  int32(generatedPair.cardVertexPair.vertex.Y),
			Id: "select_0002",
		}

		targettedTiles = append(targettedTiles, tile)
	}

	// Update the Server With Newly Highlighted Overlay Tiles
	encounterOverlayTilesRequest := &deviant.EncounterRequest{}
	encounterOverlayTilesRequest.EntityTargetAction = &deviant.EntityTargetAction{}
	encounterOverlayTilesRequest.EntityTargetAction.Id = encounter.ActiveEntity.Id
	encounterOverlayTilesRequest.EntityTargetAction.Tiles = []*deviant.Tile{}
	for _, tile := range targettedTiles {
		encounterOverlayTilesRequest.EntityTargetAction.Tiles = append(encounterOverlayTilesRequest.EntityTargetAction.Tiles, tile)
	}

	return encounterOverlayTilesRequest
}

func GenerateClearTargetAction(encounter *deviant.Encounter) *deviant.EncounterRequest {
	// Update the Server With Newly Highlighted Overlay Tiles
	encounterOverlayTilesRequest := &deviant.EncounterRequest{}
	encounterOverlayTilesRequest.EntityTargetAction = &deviant.EntityTargetAction{}
	encounterOverlayTilesRequest.EntityTargetAction.Id = encounter.ActiveEntity.Id
	encounterOverlayTilesRequest.EntityTargetAction.Tiles = []*deviant.Tile{}
	return encounterOverlayTilesRequest
}

func GeneratePlayAction(cardVertexRotationPair *CardVertexRotationPair, encounter *deviant.Encounter) *deviant.EncounterRequest {
	plays := []*deviant.Play{}

	gridNode := &gridNode{
		apCost: cardVertexRotationPair.cardVertexPair.vertex.apCost,
		X:      int32(cardVertexRotationPair.origin.X),
		Y:      int32(cardVertexRotationPair.origin.Y),
	}

	generatedPairs := GenerateCardVertexPair(cardVertexRotationPair.cardVertexPair.card, gridNode, encounter.ActiveEntity, encounter.Board.Entities.Entities, cardVertexRotationPair.rotation)

	for _, generatedPair := range generatedPairs {
		// Include the origin point along with the rotation and card for this move.
		play := &deviant.Play{
			X:  int32(generatedPair.cardVertexPair.vertex.X),
			Y:  int32(generatedPair.cardVertexPair.vertex.Y),
			Id: cardVertexRotationPair.cardVertexPair.card.InstanceId,
		}

		plays = append(plays, play)
	}

	// Update the Server With Newly Highlighted Overlay Tiles
	encounterPlayActionRequest := &deviant.EncounterRequest{}
	encounterPlayActionRequest.EntityPlayAction = &deviant.EntityPlayAction{}
	encounterPlayActionRequest.EntityPlayAction.CardId = cardVertexRotationPair.cardVertexPair.card.InstanceId
	encounterPlayActionRequest.EntityPlayAction.Plays = []*deviant.Play{}

	for _, play := range plays {
		encounterPlayActionRequest.EntityPlayAction.Plays = append(encounterPlayActionRequest.EntityPlayAction.Plays, play)
	}

	return encounterPlayActionRequest
}

func GenerateEndTurnAction(encounter *deviant.Encounter) *deviant.EncounterRequest {
	encounterRequest := &deviant.EncounterRequest{
		PlayerId:         encounter.ActiveEntity.Id,
		EntityActionName: deviant.EntityActionNames_CHANGE_PHASE,
	}

	return encounterRequest
}

func GenerateClosestMove(alignment deviant.Alignment, encounter *deviant.Encounter) *deviant.EncounterRequest {

	manhattenPairs := []*manhattenPair{}
	entityLocations := GenerateEntityLocationPairs(alignment, encounter.Board.Entities.Entities)
	validMoves := GenerateValidMoveVertexes(encounter.ActiveEntity, encounter.Board.Entities.Entities, encounter)

	for _, validMove := range validMoves {
		for _, location := range entityLocations {

			distance := int(math.Abs(float64(int32(location.vertex.X)-validMove.X)) + math.Abs(float64(int32(location.vertex.Y)-validMove.Y)))
			newManhattenPair := &manhattenPair{
				X:        int32(validMove.X),
				Y:        int32(validMove.Y),
				distance: distance,
			}
			log.Output(0, fmt.Sprintf("%v", newManhattenPair))

			manhattenPairs = append(manhattenPairs, newManhattenPair)
		}
	}

	sort.SliceStable(manhattenPairs, func(i, j int) bool { return manhattenPairs[i].distance < manhattenPairs[j].distance })

	entityMoveAction := &deviant.EntityMoveAction{
		StartXPosition: int32(GetEntityVertex(encounter.ActiveEntity, encounter.Board.Entities.Entities).X),
		StartYPosition: int32(GetEntityVertex(encounter.ActiveEntity, encounter.Board.Entities.Entities).Y),
		FinalXPosition: int32(manhattenPairs[0].X),
		FinalYPosition: int32(manhattenPairs[0].Y),
	}

	encounterRequest := &deviant.EncounterRequest{
		PlayerId:         encounter.ActiveEntity.Id,
		EntityActionName: deviant.EntityActionNames_MOVE,
		EntityMoveAction: entityMoveAction,
	}

	return encounterRequest
}

func TakeTurn(encounterResponse *deviant.EncounterResponse) []*deviant.EncounterRequest {

	encounterRequests := []*deviant.EncounterRequest{}

	allHittingMoveCombinations := FilterCardPlaysToHits(encounterResponse.Encounter.Board.Entities.Entities, encounterResponse.Encounter, deviant.Alignment_FRIENDLY)
	bestMovesInDamageOrder := SortCardPlaysByDamageInflicted(allHittingMoveCombinations, encounterResponse.Encounter.Board.Entities)
	entityLocationVertexPairs := GenerateEntityLocationPairs(deviant.Alignment_FRIENDLY, encounterResponse.Encounter.Board.Entities.Entities)

	theBestPlay := GetPlayThatDealsTheMostDamageToTheLowestHealthTargets(bestMovesInDamageOrder, entityLocationVertexPairs)

	if theBestPlay != nil {
		moveEncounterRequest := GenerateMoveAction(theBestPlay, encounterResponse.Encounter)

		moveEncounterRequest.PlayerId = "0002"
		targetEncounterRequest := GenerateTargetAction(theBestPlay, encounterResponse.Encounter)

		targetEncounterRequest.PlayerId = "0002"
		playEncounterRequest := GeneratePlayAction(theBestPlay, encounterResponse.Encounter)
		playEncounterRequest.PlayerId = "0002"

		clearTargetAction := GenerateClearTargetAction(encounterResponse.Encounter)
		clearTargetAction.PlayerId = "0002"

		encounterRequests = append(encounterRequests, moveEncounterRequest)
		encounterRequests = append(encounterRequests, targetEncounterRequest)
		encounterRequests = append(encounterRequests, playEncounterRequest)
		encounterRequests = append(encounterRequests, clearTargetAction)
	} else {
		moveEncounterRequest := GenerateClosestMove(deviant.Alignment_FRIENDLY, encounterResponse.Encounter)
		moveEncounterRequest.PlayerId = "0002"
		encounterRequests = append(encounterRequests, moveEncounterRequest)
	}

	endTurnEncounterRequest := GenerateEndTurnAction(encounterResponse.Encounter)
	endTurnEncounterRequest.PlayerId = "0002"
	encounterRequests = append(encounterRequests, endTurnEncounterRequest)

	return encounterRequests
}
