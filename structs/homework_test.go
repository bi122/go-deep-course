package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.nameArr[:], name[:42])
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	if x < maxNegativeCoordinate {
		x = maxNegativeCoordinate
	}
	if x > maxPositiveCoordinate {
		x = maxPositiveCoordinate
	}
	if y < maxNegativeCoordinate {
		y = maxNegativeCoordinate
	}
	if y > maxPositiveCoordinate {
		y = maxPositiveCoordinate
	}
	if z < maxNegativeCoordinate {
		z = maxNegativeCoordinate
	}
	if z > maxPositiveCoordinate {
		z = maxPositiveCoordinate
	}
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	if gold < 0 {
		return func(person *GamePerson) {}
	}
	if gold > maxGold {
		gold = maxGold
	}
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	if mana < 0 {
		return func(person *GamePerson) {}
	}
	if mana > maxMana {
		mana = maxMana
	}
	return func(person *GamePerson) {
		var mask byte = 0b1100_0000
		person.manaAndHealthArr[0] = byte(mana)
		person.manaAndHealthArr[1] = (byte(mana>>8) << 6) | (person.manaAndHealthArr[1] &^ mask)
	}
}

func WithHealth(health int) func(*GamePerson) {
	if health < 0 {
		return func(person *GamePerson) {}
	}
	if health > maxHealth {
		health = maxHealth
	}
	return func(person *GamePerson) {
		person.manaAndHealthArr[2] = byte(health)
		person.manaAndHealthArr[1] = person.manaAndHealthArr[1]&^healthMask | byte(health>>8)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	if respect < 0 {
		return func(person *GamePerson) {}
	}
	if respect > maxRespect {
		respect = maxRespect
	}
	return func(person *GamePerson) {
		person.respectAndStrengthArr = person.respectAndStrengthArr&^lastHalfMask | (byte(respect) << 4)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	if strength < 0 {
		return func(person *GamePerson) {}
	}
	if strength > maxStrength {
		strength = maxStrength
	}
	return func(person *GamePerson) {
		person.respectAndStrengthArr = person.respectAndStrengthArr&^firstHalfMask | byte(strength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	if experience < 0 {
		return func(person *GamePerson) {}
	}
	if experience > maxExperience {
		experience = maxExperience
	}
	return func(person *GamePerson) {
		person.experienceAndLevelArr = person.experienceAndLevelArr&^lastHalfMask | (byte(experience) << 4)
	}
}

func WithLevel(level int) func(*GamePerson) {
	if level < 0 {
		return func(person *GamePerson) {}
	}
	if level > maxLevel {
		level = maxLevel
	}
	return func(person *GamePerson) {
		person.experienceAndLevelArr = person.experienceAndLevelArr&^firstHalfMask | byte(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.propertyBitmap |= houseBitMask
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.propertyBitmap |= gunBitMask
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.propertyBitmap |= familyBitMask
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.propertyBitmap = person.propertyBitmap&^firstHalfMask | byte(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	houseBitMask  byte = 0b0001_0000
	gunBitMask    byte = 0b0010_0000
	familyBitMask byte = 0b0100_0000

	maxPositiveCoordinate                            = 2_000_000_000
	maxNegativeCoordinate                            = -2_000_000_000
	maxGold                                          = 2_000_000_000
	maxMana, maxHealth                               = 1000, 1000
	maxLevel, maxExperience, maxStrength, maxRespect = 10, 10, 10, 10
	firstHalfMask                                    = 0b0000_1111
	lastHalfMask                                     = 0b1111_0000
	healthMask                                       = 0b0000_0011
	healthAntiMask                                   = 0b1111_1100
)

// GamePerson struct tag format: `my-json:"${MethodName}:${JSONProperty},..."
type GamePerson struct {
	x int32 `my-json:"X:x"`
	y int32 `my-json:"Y:y"`
	z int32 `my-json:"Z:z"`

	gold    uint32   `my-json:"Gold:gold"`
	nameArr [42]byte `my-json:"Name:name"`

	// mana 0-9 bits; not used 10-13 bits; health 14-23 bits
	manaAndHealthArr [3]byte `my-json:"Mana:mana,Health:health"`
	// respect 0-3 bits; strength 4-7 bits
	respectAndStrengthArr uint8 `my-json:"Respect:respect,Strength:strength"`
	// experience 0-3 bits; level 4-7 bits
	experienceAndLevelArr uint8 `my-json:"Experience:experience,Level:level"`
	// not used 0th bit;have_family 2nd bit; have_gun 3rd bit; have_house 4th bit; type 4-7 bits;
	propertyBitmap uint8 `my-json:"HasHouse:has_house;HasGun:has_gun,HasFamily:has_family;Type:type_id"`
}

func (p *GamePerson) MarshalJSON() ([]byte, error) {
	v := reflect.ValueOf(*p)

	methodMap := make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		jsonTag, ok := v.Type().Field(i).Tag.Lookup("my-json")
		if !ok {
			continue
		}
		for _, tagToken := range strings.Split(jsonTag, ",") {
			tokens := strings.Split(tagToken, ":")
			if len(tokens) != 2 {
				continue
			}
			methodMap[tokens[0]] = tokens[1]
		}
	}

	v1 := reflect.ValueOf(p)
	buf := new(bytes.Buffer)
	buf.WriteString("{")
	needComma := false
	for methodName, propertyName := range methodMap {
		res := v1.MethodByName(methodName).Call([]reflect.Value{})
		if needComma {
			buf.WriteString(",")
		}
		switch res[0].Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			buf.WriteString("\"" + propertyName + "\":" + strconv.Itoa(int(res[0].Int())))
			needComma = true
		case reflect.String:
			buf.WriteString("\"" + propertyName + "\":" + strconv.Quote(res[0].String()))
			needComma = true
		case reflect.Bool:
			buf.WriteString("\"" + propertyName + "\":" + strconv.FormatBool(res[0].Bool()))
			needComma = true
		default:
			slog.Info("unhandled default case", "method", methodName)
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func NewGamePerson(options ...Option) GamePerson {
	gp := GamePerson{
		nameArr:               [42]byte{},
		manaAndHealthArr:      [3]byte{0, 0, 0},
		respectAndStrengthArr: byte(0),
		experienceAndLevelArr: byte(0),
		propertyBitmap:        byte(0),
	}
	for _, option := range options {
		option(&gp)
	}
	return gp
}

func (p *GamePerson) Name() string {
	return string(p.nameArr[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.manaAndHealthArr[1]>>6)<<8 + int(p.manaAndHealthArr[0])
}

func (p *GamePerson) Health() int {
	return int(p.manaAndHealthArr[1]&^healthAntiMask)<<8 + int(p.manaAndHealthArr[2])
}

func (p *GamePerson) Respect() int {
	return int(p.respectAndStrengthArr&^firstHalfMask) >> 4
}

func (p *GamePerson) Strength() int {
	return int(p.respectAndStrengthArr &^ lastHalfMask)
}

func (p *GamePerson) Experience() int {
	return int(p.experienceAndLevelArr&^firstHalfMask) >> 4
}

func (p *GamePerson) Level() int {
	return int(p.experienceAndLevelArr &^ lastHalfMask)
}

func (p *GamePerson) HasHouse() bool {
	return p.propertyBitmap&houseBitMask != 0
}

func (p *GamePerson) HasGun() bool {
	return p.propertyBitmap&gunBitMask != 0
}

func (p *GamePerson) HasFamily() bool {
	return p.propertyBitmap&familyBitMask != 0
}

func (p *GamePerson) Type() int {
	return int(p.propertyBitmap&^lastHalfMask) >> 4
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const expX, expY, expZ = -2_000_000_000, 2_000_000_000, z
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const expGold = 2_000_000_000
	const mana = 10
	const health = 1000
	const respect = 3
	const strength = 7
	const experience = 3
	const level = 7

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, expX, person.X())
	assert.Equal(t, expY, person.Y())
	assert.Equal(t, expZ, person.Z())
	assert.Equal(t, expGold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())

	expJSON := "{\"y\":2000000000,\"z\":0,\"gold\":2000000000,\"mana\":10,\"respect\":3,\"strength\":7,\"experience\":3,\"name\":\"aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc\",\"health\":1000,\"level\":7,\"x\":-2000000000}"
	data, err := json.Marshal(&person)
	assert.NoError(t, err)
	assert.JSONEq(t, expJSON, string(data))
}
