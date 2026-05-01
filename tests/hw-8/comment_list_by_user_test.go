package hw_8

import (
	"net/http"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
	"route256/tests/app/domain"
)

func (s *Suite) TestCommentListByOneUser_Success(t provider.T) {
	t.Title("Успешное получение списка комментариев, которые оставил один юзер по SKU")

	const (
		userID       = int64(40)
		sku          = int64(40)
		commentText1 = "тестовый комментарий 1"
		commentText2 = "тестовый комментарий 2"
	)

	var (
		commentIDs []int64
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID", userID, "sku", sku)

	t.WithNewStep("Создание комментария 1", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID, sku, commentText1)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID1", commentID)

	t.WithNewStep("Получение первого созданного комментария", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListByUser(s.ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Len(commentList, 1)
		sCtx.Require().NotNil(commentList[0])

		resultComment := commentList[0]
		sCtx.Require().Equal(commentID, resultComment.ID)
		sCtx.Require().Equal(sku, resultComment.SKU)
		sCtx.Require().Equal(commentText1, resultComment.Comment)
	})

	t.WithNewStep("Создание комментария 2", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID, sku, commentText2)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID2", commentID)

	t.WithNewStep("Получение созданных комментариев в обратном хронологическом порядке", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListByUser(s.ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		sCtx.Require().Len(commentList, 2)
		sCtx.Require().NotNil(commentList[0])
		sCtx.Require().NotNil(commentList[1])

		// commentIDs содержит в себе id комментариев в прямом хронологическом порядке
		sCtx.Require().Equal(commentIDs[1], commentList[0].ID)
		sCtx.Require().Equal(sku, commentList[0].SKU)
		sCtx.Require().Equal(commentText2, commentList[0].Comment)

		sCtx.Require().Equal(commentIDs[0], commentList[1].ID)
		sCtx.Require().Equal(sku, commentList[1].SKU)
		sCtx.Require().Equal(commentText1, commentList[1].Comment)
	})
}

func (s *Suite) TestCommentListByOneUserManySKU_Success(t provider.T) {
	t.Title("Успешное получение списка комментариев, которые оставил один юзер к нескольким SKU по user_id")

	const (
		userID       = int64(41)
		sku1         = int64(410)
		sku2         = int64(411)
		commentText1 = "тестовый комментарий 1"
		commentText2 = "тестовый комментарий 2"
	)

	var (
		commentIDs []int64
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID", userID, "sku", sku1)

	t.WithNewStep("Создание комментария 1", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID, sku1, commentText1)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID1", commentID)

	t.WithNewStep("Получение первого созданного комментария", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListByUser(s.ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Len(commentList, 1)
		sCtx.Require().NotNil(commentList[0])

		resultComment := commentList[0]
		sCtx.Require().Equal(commentID, resultComment.ID)
		sCtx.Require().Equal(sku1, resultComment.SKU)
		sCtx.Require().Equal(commentText1, resultComment.Comment)
	})

	t.WithNewStep("Создание комментария 2", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID, sku2, commentText2)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID2", commentID)

	t.WithNewStep("Получение созданных комментариев в обратном хронологическом порядке", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListByUser(s.ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		sCtx.Require().Len(commentList, 2)
		sCtx.Require().NotNil(commentList[0])
		sCtx.Require().NotNil(commentList[1])

		// commentIDs содержит в себе id комментариев в прямом хронологическом порядке
		sCtx.Require().Equal(commentIDs[1], commentList[0].ID)
		sCtx.Require().Equal(sku2, commentList[0].SKU)
		sCtx.Require().Equal(commentText2, commentList[0].Comment)

		sCtx.Require().Equal(commentIDs[0], commentList[1].ID)
		sCtx.Require().Equal(sku1, commentList[1].SKU)
		sCtx.Require().Equal(commentText1, commentList[1].Comment)
	})
}

func (s *Suite) TestCommentListByManyUser_Success(t provider.T) {
	t.Title("Успешное получение списка комментариев, которые оставили несколько юзеров по user_id")

	const (
		userID1      = int64(410)
		userID2      = int64(411)
		sku          = int64(40)
		commentText1 = "тестовый комментарий 1"
		commentText2 = "тестовый комментарий 2"
	)

	var (
		commentIDs []int64
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID", userID1, "sku", sku)

	t.WithNewStep("Создание комментария 1", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID1, sku, commentText1)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID1", commentID)

	t.WithNewStep("Получение первого созданного комментария", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListByUser(s.ctx, sCtx, userID1)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Len(commentList, 1)
		sCtx.Require().NotNil(commentList[0])

		resultComment := commentList[0]
		sCtx.Require().Equal(commentID, resultComment.ID)
		sCtx.Require().Equal(sku, resultComment.SKU)
		sCtx.Require().Equal(commentText1, resultComment.Comment)
	})

	t.WithNewStep("Создание комментария 2", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID2, sku, commentText2)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID2", commentID)

	t.WithNewStep("Получение созданного комментария первым юзером", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListByUser(s.ctx, sCtx, userID1)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Len(commentList, 1)
		sCtx.Require().NotNil(commentList[0])

		resultComment := commentList[0]
		sCtx.Require().Equal(commentIDs[0], resultComment.ID)
		sCtx.Require().Equal(sku, resultComment.SKU)
		sCtx.Require().Equal(commentText1, resultComment.Comment)
	})
}

func (s *Suite) TestCommentListByUser_Negative(t provider.T) {
	t.Title("Негативные кейсы получения списка комментариев по user_id")

	const (
		userID = int64(44)
	)

	var (
		statusCode int
	)

	t.WithNewStep("Получение списка комментариев по user_id, у которого еще нет комментариев", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListByUser(s.ctx, sCtx, userID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		sCtx.Require().Len(commentList, 0)
	})

	t.WithNewStep("Получение списка комментариев по user_id без указания user_id", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListByUser(s.ctx, sCtx, 0)
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		sCtx.Require().Len(commentList, 0)
	})
}
