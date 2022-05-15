package api

import (
	"database/sql"
	"net/http"

	db "shortURL/db/sqlc"
	"shortURL/util"

	"github.com/gin-gonic/gin"
)

type createShortURLRequest struct {
	OriginUrl string `json:"originUrl" binding:"required,url"`
}

// 建立短連結
func (server *Server) createShortURL(ctx *gin.Context) {
	var req createShortURLRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateURLParams{
		OriginUrl: req.OriginUrl,
		ShortUrl:  util.RandomString(6),
	}

	account, err := server.store.CreateURL(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type getRedirectRequest struct {
	ShortUrl string `uri:"short_url" binding:"required,min=6"`
}

// 取得導向的長連結
func (server *Server) getRedirect(ctx *gin.Context) {
	var req getRedirectRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	url, err := server.store.GetURL(ctx, req.ShortUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, url.OriginUrl)
}
