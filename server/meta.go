package server

import (
	"encoding/json"
	"github.com/dofusdude/api/gen"
	"github.com/dofusdude/api/utils"
	"net/http"
)

func ListEffectConditionElements(w http.ResponseWriter, r *http.Request) {
	txn := Db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("effect-condition-elements", "id")
	if err != nil || it == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var effects []string
	for obj := it.Next(); obj != nil; obj = it.Next() {
		effectElement := obj.(*gen.EffectConditionDbEntry)
		effects = append(effects, effectElement.Name)
	}

	utils.WriteCacheHeader(&w)
	err = json.NewEncoder(w).Encode(effects)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
