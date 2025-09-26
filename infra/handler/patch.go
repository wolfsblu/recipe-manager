package handler

import (
	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/ogen-go/ogen/json"
)

func (h *RecipeHandler) mergePatch(original any, patch any) ([]byte, error) {
	originalJSON, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}

	patchJSON, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	return jsonpatch.MergePatch(originalJSON, patchJSON)
}
