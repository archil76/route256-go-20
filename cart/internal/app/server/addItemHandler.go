package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/logger"
	"route256/cart/internal/infra/utils"

	gody "github.com/guiferpa/gody/v2"
	rule "github.com/guiferpa/gody/v2/rule"
)

func (s *Server) AddItem(writer http.ResponseWriter, request *http.Request) {
	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(rawUserID)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidUserID, "", http.StatusBadRequest)
		return
	}

	rawSkuID := request.PathValue("sku_id")
	skuID, err := utils.PrepareID(rawSkuID)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidSKU, "", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidSKU, "", http.StatusBadRequest)
		return
	}

	var addItemRequest AddItemRequest

	err = json.Unmarshal(body, &addItemRequest)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrUnmarshalling, "", http.StatusBadRequest)
		return
	}

	validator := gody.NewValidator()
	err = validator.AddRules(rule.Min)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidCount, "", http.StatusBadRequest)
		return
	}
	if _, err = validator.Validate(addItemRequest); err != nil {
		utils.WriteErrorToResponse(writer, request, ErrOther, "", http.StatusBadRequest)
		return
	}

	if addItemRequest.Count < 1 {
		utils.WriteErrorToResponse(writer, request, ErrInvalidCount, "", http.StatusBadRequest)
		return
	}

	_, err = s.cartService.AddItem(request.Context(), userID, skuID, uint32(addItemRequest.Count))
	if err != nil {
		logger.Errorw("controller", err)
		if errors.Is(err, model.ErrProductNotFound) {
			utils.WriteErrorToResponse(writer, request, ErrPSFail, "", http.StatusPreconditionFailed)
		} else {
			utils.WriteErrorToResponse(writer, request, fmt.Errorf("3: %v;", err), "0", http.StatusBadRequest)
		}
		return
	}
	utils.WriteStatusToResponse(writer, request, "", http.StatusOK)
}
