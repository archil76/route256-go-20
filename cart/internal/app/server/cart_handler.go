package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	model "route256/cart/internal/domain/model"
	"route256/cart/internal/infra/utils"

	gody "github.com/guiferpa/gody/v2"
	"github.com/guiferpa/gody/v2/rule"
)

var (
	ErrInvalidUserID = errors.New("Идентификатор пользователя должен быть натуральным числом (больше нуля)")
	ErrInvalidSKU    = errors.New("SKU должен быть натуральным числом (больше нуля)")
	ErrInvalidCount  = errors.New("Количество должно быть натуральным числом (больше нуля)")
	ErrPSFail        = errors.New("SKU должен существовать в сервисе product-service")
	ErrUnmarshalling = errors.New("Unmarshalling error")
	ErrOther         = errors.New("Ошибка сервера")
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
	if _, err := validator.Validate(addItemRequest); err != nil {
		utils.WriteErrorToResponse(writer, request, ErrOther, "", http.StatusBadRequest)
		return
	}

	if addItemRequest.Count < 1 {
		utils.WriteErrorToResponse(writer, request, ErrInvalidCount, "", http.StatusBadRequest)
		return
	}

	_, err = s.cartService.AddItem(request.Context(), userID, skuID, addItemRequest.Count)

	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			utils.WriteErrorToResponse(writer, request, ErrPSFail, "", http.StatusPreconditionFailed)
		} else {
			utils.WriteErrorToResponse(writer, request, ErrOther, "", http.StatusBadRequest)
		}
		return
	}
	utils.WriteStatusToResponse(writer, request, "", http.StatusOK)

}

func (s *Server) DeleteItem(writer http.ResponseWriter, request *http.Request) {

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

	_, err = s.cartService.DeleteItem(request.Context(), userID, skuID)
	if err != nil {
		// тут ошибки могут быть из-за невалидных ID а они проверены раньше. Поэтому просто лог и ответ ОК.
		utils.WriteErrorToLog(request, err, "")
	}

	utils.WriteStatusToResponse(writer, request, "", http.StatusOK)

}

func (s *Server) ClearCart(writer http.ResponseWriter, request *http.Request) {

	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(rawUserID)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidUserID, "", http.StatusBadRequest)
		return
	}

	_, err = s.cartService.DeleteItemByUserId(request.Context(), userID)

	if err != nil {
		// тут ошибки могут быть из-за невалидных ID а они проверены раньше. Поэтому просто лог и ответ ОК.
		utils.WriteErrorToLog(request, err, "")
	}

	utils.WriteStatusToResponse(writer, request, "", http.StatusNoContent)

}

func (s *Server) GetCart(writer http.ResponseWriter, request *http.Request) {

	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(rawUserID)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidUserID, "", http.StatusBadRequest)
		return
	}

	cart, err := s.cartService.GetItemsByUserId(request.Context(), userID)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrOther, "", http.StatusBadRequest)
		return
	}

	reportCart := ReportCart{
		UserID:     cart.UserID,
		Items:      map[model.Sku]ItemInСart{},
		TotalPrice: cart.TotalPrice}

	for sku, item := range cart.Items {
		reportCart.Items[item.SKU] = ItemInСart{
			SKU:   sku,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		}

	}

	rawResponce, err := json.Marshal(reportCart)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, err, "can't get cart. marshalling error", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(rawResponce)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, err, "can't get cart. marshalling error", http.StatusBadRequest)

	}
}
