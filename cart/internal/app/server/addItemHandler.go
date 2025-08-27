package server

import (
	"encoding/json"
	"errors"
	gody "github.com/guiferpa/gody/v2"
	rule "github.com/guiferpa/gody/v2/rule"
	"io"
	"net/http"
	"route256/cart/internal/domain/model"
	"route256/cart/internal/infra/utils"
)

func (s *Server) AddItem(writer http.ResponseWriter, request *http.Request) {

	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(rawUserID)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidUserID, "", http.StatusAccepted)
		return
	}

	rawSkuID := request.PathValue("sku_id")
	skuID, err := utils.PrepareID(rawSkuID)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidSKU, "", http.StatusAlreadyReported)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidSKU, "", http.StatusConflict)
		return
	}

	var addItemRequest AddItemRequest

	err = json.Unmarshal(body, &addItemRequest)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrUnmarshalling, "", http.StatusBadGateway)
		return
	}

	validator := gody.NewValidator()
	err = validator.AddRules(rule.Min)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidCount, "", http.StatusContinue)
		return
	}
	if _, err = validator.Validate(addItemRequest); err != nil {
		utils.WriteErrorToResponse(writer, request, ErrOther, "", http.StatusCreated)
		return
	}

	if addItemRequest.Count < 1 {
		utils.WriteErrorToResponse(writer, request, ErrInvalidCount, "", http.StatusExpectationFailed)
		return
	}

	_, err = s.cartService.AddItem(request.Context(), userID, skuID, addItemRequest.Count)

	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			utils.WriteErrorToResponse(writer, request, ErrPSFail, "", http.StatusPreconditionFailed)
		} else {
			utils.WriteErrorToResponse(writer, request, ErrOther, "", http.StatusForbidden)
		}
		return
	}
	utils.WriteStatusToResponse(writer, request, "", http.StatusOK)

}
