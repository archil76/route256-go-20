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

func (s *Server) AddItem(writer http.ResponseWriter, request *http.Request) {

	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(writer, request, rawUserID)

	if err != nil {
		return
	}

	rawSkuID := request.PathValue("sku_id")
	skuID, err := utils.PrepareID(writer, request, rawSkuID)

	if err != nil {
		return
	}

	body, err := io.ReadAll(request.Body)

	var addItemRequest AddItemRequest

	err = json.Unmarshal(body, &addItemRequest)
	if err != nil {
		err = utils.WriteErrorToResponse(writer, request, err, "unmarshalling error", http.StatusBadRequest)
		return
	}

	validator := gody.NewValidator()
	err = validator.AddRules(rule.Min)
	if err != nil {
		_ = utils.WriteErrorToResponse(writer, request, err, "", http.StatusBadRequest)
		return
	}
	if _, err := validator.Validate(addItemRequest); err != nil {
		_ = utils.WriteErrorToResponse(writer, request, err, "", http.StatusBadRequest)
		return
	}

	if addItemRequest.Count < 1 {
		_ = utils.WriteErrorToResponse(writer, request, errors.New("count of items must be greater than zero"), "", http.StatusBadRequest)
		return
	}

	_, err = s.cartService.AddItem(request.Context(), userID, skuID, addItemRequest.Count)

	if err != nil {
		err = utils.WriteErrorToResponse(writer, request, err, "unable to add cart", http.StatusPreconditionFailed)

		return
	}
	utils.WriteStatusToResponse(writer, request, "", http.StatusOK)

}

func (s *Server) DeleteItem(writer http.ResponseWriter, request *http.Request) {

	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(writer, request, rawUserID)

	if err != nil {

		utils.WriteStatusToResponse(writer, request, "", http.StatusOK)
		return
	}

	rawSkuID := request.PathValue("sku_id")
	skuID, err := utils.PrepareID(writer, request, rawSkuID)

	if err != nil {
		utils.WriteStatusToResponse(writer, request, "", http.StatusOK)
		return
	}

	_, _ = s.cartService.DeleteItem(request.Context(), userID, skuID)

	utils.WriteStatusToResponse(writer, request, "", http.StatusOK)

}

func (s *Server) ClearCart(writer http.ResponseWriter, request *http.Request) {
	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(writer, request, rawUserID)

	if err != nil {
		return
	}

	_, err = s.cartService.DeleteItemByUserId(request.Context(), userID)

	utils.WriteStatusToResponse(writer, request, "", http.StatusNoContent)

}

func (s *Server) GetCart(writer http.ResponseWriter, request *http.Request) {

	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(writer, request, rawUserID)

	if err != nil {
		return
	}

	cart, err := s.cartService.GetItemsByUserId(request.Context(), userID)

	if err != nil {
		err = utils.WriteErrorToResponse(writer, request, err, "", http.StatusOK)
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
		_ = utils.WriteErrorToResponse(writer, request, err, "can't get cart. marshalling error", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(rawResponce)

	if err != nil {
		_ = utils.WriteErrorToResponse(writer, request, err, "can't get cart. marshalling error", http.StatusBadRequest)

	}
}
