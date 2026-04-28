package hw_8

import (
	"net/http"
	"strings"
	"time"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
	"route256/tests/app/domain"
)

func (s *Suite) TestEditComment_Success(t provider.T) {

	t.Title("успешное создание комментария")

	const (
		sku            = int64(20)
		userID         = int64(20)
		commentText    = "тестовый комментарий"
		newCommentText = "новый текстовый комментарий"
	)

	var (
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID", userID, "sku", sku)

	t.WithNewStep("Создание комментария", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID, sku, commentText)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID", commentID)

	t.WithNewStep("Получение созданного комментария", func(sCtx provider.StepCtx) {
		comment := &domain.Comment{}
		comment, statusCode = s.commentsClient.GetCommentByID(s.ctx, sCtx, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Equal(commentID, comment.ID)
		sCtx.Require().Equal(userID, comment.UserID)
		sCtx.Require().Equal(commentText, comment.Comment)
	})

	t.WithNewStep("Редактирование комментария", func(sCtx provider.StepCtx) {
		statusCode = s.commentsClient.EditComment(s.ctx, sCtx, userID, commentID, newCommentText)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewStep("Проверка корректного редактирования комментария", func(sCtx provider.StepCtx) {
		comment := &domain.Comment{}
		comment, statusCode = s.commentsClient.GetCommentByID(s.ctx, sCtx, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Equal(newCommentText, comment.Comment)
	})
}

func (s *Suite) TestEditComment_EditOnlyYourOwnComments(t provider.T) {

	t.Title("пользователь может редактировать только свои комментарии")

	const (
		sku            = int64(21)
		userID1        = int64(21)
		userID2        = int64(4321)
		commentText    = "тестовый комментарий"
		newCommentText = "новый текстовый комментарий"
	)

	var (
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID", userID1, "sku", sku)

	t.WithNewStep("Создание комментария", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID1, sku, commentText)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("userID", userID2, "commentID", commentID)

	t.WithNewStep("Редактирование чужого комментария", func(sCtx provider.StepCtx) {
		statusCode = s.commentsClient.EditComment(s.ctx, sCtx, userID2, commentID, newCommentText)
		assert.StatusCode(sCtx, http.StatusForbidden, statusCode)
	})
}

func (s *Suite) TestEditComment_EditAfterEditInterval(t provider.T) {
	t.Title("пользователь может редактировать комментарии только в течение n времени после создания")

	const (
		sku            = int64(22)
		userID         = int64(22)
		commentText    = "тестовый комментарий"
		newCommentText = "новый текстовый комментарий"
	)

	var (
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID", userID, "sku", sku)

	t.WithNewStep("Создание комментария", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID, sku, commentText)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID", commentID)

	t.WithNewStep("Редактирование комментария после истечения editInterval", func(sCtx provider.StepCtx) {
		t.Logf("ждем истечения editInterval=%vs", s.editInterval.Seconds())
		time.Sleep(s.editInterval)

		statusCode = s.commentsClient.EditComment(s.ctx, sCtx, userID, commentID, newCommentText)

		// GRPC status FailedPrecondition is actually a 400 HTTP code
		// assert.StatusCode(sCtx, http.StatusPreconditionFailed, statusCode)
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})
}

func (s *Suite) TestEditComment_Negative(t provider.T) {

	t.Title("успешное создание комментария")

	const (
		sku               = int64(23)
		userID1           = int64(230)
		userID2           = int64(231)
		NotFoundCommentID = int64(1)
		commentText       = "тестовый комментарий"
		newCommentText    = "новый текстовый комментарий"
	)

	var (
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID", userID1, "sku", sku)

	t.WithNewStep("Создание комментария", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID1, sku, commentText)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID", commentID)

	t.WithNewStep("Получение созданного комментария", func(sCtx provider.StepCtx) {
		comment := &domain.Comment{}
		comment, statusCode = s.commentsClient.GetCommentByID(s.ctx, sCtx, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Equal(commentID, comment.ID)
		sCtx.Require().Equal(userID1, comment.UserID)
		sCtx.Require().Equal(commentText, comment.Comment)
	})

	t.WithNewStep("Редактирование комментария без указания user_id", func(sCtx provider.StepCtx) {
		statusCode = s.commentsClient.EditComment(s.ctx, sCtx, 0, commentID, newCommentText)
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})
	t.WithNewStep("Редактирование комментария без указания comment_id", func(sCtx provider.StepCtx) {
		statusCode = s.commentsClient.EditComment(s.ctx, sCtx, userID1, 0, newCommentText)
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Редактирование комментария с запрещенным минимальным размером текста", func(sCtx provider.StepCtx) {
		statusCode = s.commentsClient.EditComment(s.ctx, sCtx, userID1, commentID, "")
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Редактирование комментария с запрещенным максимальным размером текста", func(sCtx provider.StepCtx) {
		statusCode = s.commentsClient.EditComment(s.ctx, sCtx, userID1, commentID, strings.Repeat("q", 260))
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
	})

	t.WithNewStep("Редактирование не своего комментария", func(sCtx provider.StepCtx) {
		statusCode = s.commentsClient.EditComment(s.ctx, sCtx, userID2, commentID, newCommentText)
		assert.StatusCode(sCtx, http.StatusForbidden, statusCode)
	})

	t.WithNewStep("Редактирование несуществующего комментария", func(sCtx provider.StepCtx) {
		statusCode = s.commentsClient.EditComment(s.ctx, sCtx, userID1, NotFoundCommentID, newCommentText)
		assert.StatusCode(sCtx, http.StatusNotFound, statusCode)
	})
}
