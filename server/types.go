package server

import (
	"fmt"
	"github.com/dofusdude/api/gen"
	"github.com/dofusdude/api/utils"
	"github.com/hashicorp/go-memdb"
	"log"
)

type ApiImageUrls struct {
	Icon string `json:"icon"`
	Sd   string `json:"sd,omitempty"`
	Hq   string `json:"hq,omitempty"`
	Hd   string `json:"hd,omitempty"`
}

func RenderImageUrls(urls []string) ApiImageUrls {
	if len(urls) == 0 {
		return ApiImageUrls{}
	}

	var res ApiImageUrls
	res.Icon = urls[0]
	if len(urls) > 1 {
		res.Sd = urls[1]
	}
	if len(urls) > 2 {
		res.Hq = urls[2]
	}
	if len(urls) > 3 {
		res.Hd = urls[3]
	}

	return res
}

type ApiEffect struct {
	MinInt       int           `json:"int_minimum"`
	MaxInt       int           `json:"int_maximum"`
	Type         ApiEffectType `json:"type"`
	IgnoreMinInt bool          `json:"ignore_int_min"`
	IgnoreMaxInt bool          `json:"ignore_int_max"`
	Formatted    string        `json:"formatted"`
}

func RenderEffects(effects *[]gen.MappedMultilangEffect, lang string) []ApiEffect {
	var retEffects []ApiEffect
	for _, effect := range *effects {
		retEffects = append(retEffects, ApiEffect{
			MinInt:       effect.Min,
			MaxInt:       effect.Max,
			IgnoreMinInt: effect.IsMeta || effect.MinMaxIrrelevant == -2,
			IgnoreMaxInt: effect.IsMeta || effect.MinMaxIrrelevant <= -1,
			Type: ApiEffectType{
				Name:     effect.Type[lang],
				Id:       effect.ElementId,
				IsMeta:   effect.IsMeta,
				IsActive: effect.Active,
			},
			Formatted: effect.Templated[lang],
		})
	}

	if len(retEffects) > 0 {
		return retEffects
	}

	return nil
}

type ApiCondition struct {
	Operator string           `json:"operator"`
	IntValue int              `json:"int_value"`
	Element  ApiConditionType `json:"element"`
}

type APIResource struct {
	Id          int            `json:"ankama_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        ApiType        `json:"type"`
	Level       int            `json:"level"`
	Pods        int            `json:"pods"`
	ImageUrls   ApiImageUrls   `json:"image_urls,omitempty"`
	Effects     []ApiEffect    `json:"effects,omitempty"`
	Conditions  []ApiCondition `json:"conditions,omitempty"`
	Recipe      []APIRecipe    `json:"recipe,omitempty"`
}

func RenderResource(item *gen.MappedMultilangItem, lang string) APIResource {
	resource := APIResource{
		Id:   item.AnkamaId,
		Name: item.Name[lang],
		Type: ApiType{
			Name: item.Type.Name[lang],
			Id:   item.Type.ItemTypeId,
		},
		Description: item.Description[lang],
		Level:       item.Level,
		Pods:        item.Pods,
		ImageUrls:   RenderImageUrls(utils.ImageUrls(item.IconId, "item")),
		Recipe:      nil,
	}

	conditions := RenderConditions(&item.Conditions, lang)
	if len(conditions) == 0 {
		resource.Conditions = nil
	} else {
		resource.Conditions = conditions
	}

	effects := RenderEffects(&item.Effects, lang)
	if len(effects) == 0 {
		resource.Effects = nil
	} else {
		resource.Effects = effects
	}

	return resource
}

type APIEquipment struct {
	Id          int                `json:"ankama_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Type        ApiType            `json:"type"`
	IsWeapon    bool               `json:"is_weapon"`
	Level       int                `json:"level"`
	Pods        int                `json:"pods"`
	ImageUrls   ApiImageUrls       `json:"image_urls,omitempty"`
	Effects     []ApiEffect        `json:"effects,omitempty"`
	Conditions  []ApiCondition     `json:"conditions,omitempty"`
	Recipe      []APIRecipe        `json:"recipe,omitempty"`
	ParentSet   *APISetReverseLink `json:"parent_set,omitempty"`
}

func RenderEquipment(item *gen.MappedMultilangItem, lang string) APIEquipment {
	var setLink *APISetReverseLink = nil
	if item.HasParentSet {
		setLink = &APISetReverseLink{
			Id:   item.ParentSet.Id,
			Name: item.ParentSet.Name[lang],
		}
	}

	equip := APIEquipment{
		Id:   item.AnkamaId,
		Name: item.Name[lang],
		Type: ApiType{
			Name: item.Type.Name[lang],
			Id:   item.Type.ItemTypeId,
		},
		Description: item.Description[lang],
		Level:       item.Level,
		Pods:        item.Pods,
		ImageUrls:   RenderImageUrls(utils.ImageUrls(item.IconId, "item")),
		IsWeapon:    false,
		Recipe:      nil,
		ParentSet:   setLink,
	}

	conditions := RenderConditions(&item.Conditions, lang)
	if len(conditions) == 0 {
		equip.Conditions = nil
	} else {
		equip.Conditions = conditions
	}

	effects := RenderEffects(&item.Effects, lang)
	if len(effects) == 0 {
		equip.Effects = nil
	} else {
		equip.Effects = effects
	}

	return equip
}

type APIRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type APISetReverseLink struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type APIWeapon struct {
	Id                     int                `json:"ankama_id"`
	Name                   string             `json:"name"`
	Description            string             `json:"description"`
	Type                   ApiType            `json:"type"`
	IsWeapon               bool               `json:"is_weapon"`
	Level                  int                `json:"level"`
	Pods                   int                `json:"pods"`
	ImageUrls              ApiImageUrls       `json:"image_urls,omitempty"`
	Effects                []ApiEffect        `json:"effects,omitempty"`
	Conditions             []ApiCondition     `json:"conditions,omitempty"`
	CriticalHitProbability int                `json:"critical_hit_probability"`
	CriticalHitBonus       int                `json:"critical_hit_bonus"`
	TwoHanded              bool               `json:"is_two_handed"`
	MaxCastPerTurn         int                `json:"max_cast_per_turn"`
	ApCost                 int                `json:"ap_cost"`
	Range                  APIRange           `json:"range"`
	Recipe                 []APIRecipe        `json:"recipe,omitempty"`
	ParentSet              *APISetReverseLink `json:"parent_set,omitempty"`
}

func RenderWeapon(item *gen.MappedMultilangItem, lang string) APIWeapon {
	var setLink *APISetReverseLink = nil
	if item.HasParentSet {
		setLink = &APISetReverseLink{
			Id:   item.ParentSet.Id,
			Name: item.ParentSet.Name[lang],
		}
	}

	weapon := APIWeapon{
		Id:   item.AnkamaId,
		Name: item.Name[lang],
		Type: ApiType{
			Name: item.Type.Name[lang],
			Id:   item.Type.ItemTypeId,
		},
		Description:            item.Description[lang],
		Level:                  item.Level,
		Pods:                   item.Pods,
		ImageUrls:              RenderImageUrls(utils.ImageUrls(item.IconId, "item")),
		Recipe:                 nil,
		CriticalHitBonus:       item.CriticalHitBonus,
		CriticalHitProbability: item.CriticalHitProbability,
		TwoHanded:              item.TwoHanded,
		MaxCastPerTurn:         item.MaxCastPerTurn,
		ApCost:                 item.ApCost,
		Range: APIRange{
			Min: item.MinRange,
			Max: item.Range,
		},
		IsWeapon:  true,
		ParentSet: setLink,
	}

	conditions := RenderConditions(&item.Conditions, lang)
	if len(conditions) == 0 {
		weapon.Conditions = nil
	} else {
		weapon.Conditions = conditions
	}

	effects := RenderEffects(&item.Effects, lang)
	if len(effects) == 0 {
		weapon.Effects = nil
	} else {
		weapon.Effects = effects
	}

	return weapon
}

type MappedMultilangCondition struct {
	Element   string            `json:"element"`
	Operator  string            `json:"operator"`
	Value     int               `json:"value"`
	Templated map[string]string `json:"templated"`
}

func RenderConditions(conditions *[]gen.MappedMultilangCondition, lang string) []ApiCondition {
	var retConditions []ApiCondition
	for _, condition := range *conditions {
		retConditions = append(retConditions, ApiCondition{
			Operator: condition.Operator,
			IntValue: condition.Value,
			Element: ApiConditionType{
				Name: condition.Templated[lang],
				Id:   condition.ElementId,
			},
		})
	}

	if len(retConditions) > 0 {
		return retConditions
	}

	return nil
}

type ApiType struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type ApiConditionType struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type ApiEffectType struct {
	Name     string `json:"name"`
	Id       int    `json:"id"`
	IsMeta   bool   `json:"is_meta"`
	IsActive bool   `json:"is_active"`
}

type APIListItem struct {
	Id        int          `json:"ankama_id"`
	Name      string       `json:"name"`
	Type      ApiType      `json:"type"`
	Level     int          `json:"level"`
	ImageUrls ApiImageUrls `json:"image_urls,omitempty"`

	// extra fields
	Description *string        `json:"description,omitempty"`
	Recipe      []APIRecipe    `json:"recipe,omitempty"`
	Conditions  []ApiCondition `json:"conditions,omitempty"`
	Effects     []ApiEffect    `json:"effects,omitempty"`

	// extra equipment
	IsWeapon  *bool              `json:"is_weapon,omitempty"`
	Pods      *int               `json:"pods,omitempty"`
	ParentSet *APISetReverseLink `json:"parent_set,omitempty"`

	// extra weapon
	CriticalHitProbability *int      `json:"critical_hit_probability,omitempty"`
	CriticalHitBonus       *int      `json:"critical_hit_bonus,omitempty"`
	TwoHanded              *bool     `json:"is_two_handed,omitempty"`
	MaxCastPerTurn         *int      `json:"max_cast_per_turn,omitempty"`
	ApCost                 *int      `json:"ap_cost,omitempty"`
	Range                  *APIRange `json:"range,omitempty"`
}

func RenderItemListEntry(item *gen.MappedMultilangItem, lang string) APIListItem {
	return APIListItem{
		Id:   item.AnkamaId,
		Name: item.Name[lang],
		Type: ApiType{
			Name: item.Type.Name[lang],
			Id:   item.Type.ItemTypeId,
		},
		Level:     item.Level,
		ImageUrls: RenderImageUrls(utils.ImageUrls(item.IconId, "item")),
	}
}

type APIListTypedItem struct {
	Id          int          `json:"ankama_id"`
	Name        string       `json:"name"`
	Type        ApiType      `json:"type"`
	ItemSubtype string       `json:"item_subtype"`
	Level       int          `json:"level"`
	ImageUrls   ApiImageUrls `json:"image_urls,omitempty"`
}

func RenderTypedItemListEntry(item *gen.MappedMultilangItem, lang string) APIListTypedItem {
	return APIListTypedItem{
		Id:   item.AnkamaId,
		Name: item.Name[lang],
		Type: ApiType{
			Name: item.Type.Name[lang],
			Id:   item.Type.ItemTypeId,
		},
		ItemSubtype: utils.CategoryIdApiMapping(item.Type.CategoryId),
		Level:       item.Level,
		ImageUrls:   RenderImageUrls(utils.ImageUrls(item.IconId, "item")),
	}
}

type APIListMount struct {
	Id         int          `json:"ankama_id"`
	Name       string       `json:"name"`
	FamilyName string       `json:"family_name"`
	ImageUrls  ApiImageUrls `json:"image_urls,omitempty"`

	// extra fields
	Effects []ApiEffect `json:"effects,omitempty"`
}

func RenderMountListEntry(mount *gen.MappedMultilangMount, lang string) APIListMount {
	return APIListMount{
		Id:         mount.AnkamaId,
		Name:       mount.Name[lang],
		ImageUrls:  RenderImageUrls(utils.ImageUrls(mount.AnkamaId, "mount")),
		FamilyName: mount.FamilyName[lang],
	}
}

type APIRecipe struct {
	AnkamaId int    `json:"item_ankama_id"`
	ItemType string `json:"item_subtype"`
	Quantity int    `json:"quantity"`
}

func RenderRecipe(recipe gen.MappedMultilangRecipe, db *memdb.MemDB) []APIRecipe {
	if len(recipe.Entries) == 0 {
		return nil
	}

	txn := db.Txn(false)
	defer txn.Abort()

	var apiRecipes []APIRecipe
	for _, entry := range recipe.Entries {
		raw, err := txn.First(fmt.Sprintf("%s-%s", utils.CurrentRedBlueVersionStr(Version.MemDb), "all_items"), "id", entry.ItemId)
		if err != nil {
			log.Println(err)
			return nil
		}
		item := raw.(*gen.MappedMultilangItem)

		apiRecipes = append(apiRecipes, APIRecipe{
			AnkamaId: entry.ItemId,
			Quantity: entry.Quantity,
			ItemType: utils.CategoryIdApiMapping(item.Type.CategoryId),
		})
	}
	return apiRecipes
}

type APIPageItem struct {
	Links utils.PaginationLinks `json:"_links,omitempty"`
	Items []APIListItem         `json:"items"`
}

type APIPageMount struct {
	Links utils.PaginationLinks `json:"_links,omitempty"`
	Items []APIListMount        `json:"mounts"`
}

type APIPageSet struct {
	Links utils.PaginationLinks `json:"_links,omitempty"`
	Items []APIListSet          `json:"sets"`
}

type APIMount struct {
	Id         int          `json:"ankama_id"`
	Name       string       `json:"name"`
	FamilyName string       `json:"family_name"`
	ImageUrls  ApiImageUrls `json:"image_urls,omitempty"`
	Effects    []ApiEffect  `json:"effects,omitempty"`
}

func RenderMount(mount *gen.MappedMultilangMount, lang string) APIMount {
	resMount := APIMount{
		Id:         mount.AnkamaId,
		Name:       mount.Name[lang],
		FamilyName: mount.FamilyName[lang],
		ImageUrls:  RenderImageUrls(utils.ImageUrls(mount.AnkamaId, "mount")),
	}

	effects := RenderEffects(&mount.Effects, lang)
	if len(effects) == 0 {
		resMount.Effects = nil
	} else {
		resMount.Effects = effects
	}

	return resMount
}

type APIListSet struct {
	Id    int    `json:"ankama_id"`
	Name  string `json:"name"`
	Items int    `json:"items"`
	Level int    `json:"level"`

	// extra fields
	Effects [][]ApiEffect `json:"effects,omitempty"`
	ItemIds []int         `json:"equipment_ids,omitempty"`
}

func RenderSetListEntry(set *gen.MappedMultilangSet, lang string) APIListSet {
	return APIListSet{
		Id:    set.AnkamaId,
		Name:  set.Name[lang],
		Items: len(set.ItemIds),
		Level: set.Level,
	}
}

type APISet struct {
	AnkamaId int           `json:"ankama_id"`
	Name     string        `json:"name"`
	ItemIds  []int         `json:"equipment_ids"`
	Effects  [][]ApiEffect `json:"effects,omitempty"`
	Level    int           `json:"highest_equipment_level"`
}

func RenderSet(set *gen.MappedMultilangSet, lang string) APISet {
	var effects [][]ApiEffect
	for _, effect := range set.Effects {
		effects = append(effects, RenderEffects(&effect, lang))
	}

	resSet := APISet{
		AnkamaId: set.AnkamaId,
		Name:     set.Name[lang],
		ItemIds:  set.ItemIds,
		Effects:  effects,
		Level:    set.Level,
	}

	if len(effects) == 0 {
		resSet.Effects = nil
	} else {
		resSet.Effects = effects
	}

	return resSet
}
