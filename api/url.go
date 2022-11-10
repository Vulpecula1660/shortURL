package api

import (
	"database/sql"
	"fmt"
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

	// 產生短網址
	var shortUrl string
	retry := 1

	for retry < 5 {
		shortUrl = util.RandomString(6)

		// 設置布隆過濾器
		exist, err := server.redis.SetBloom(ctx, shortUrl)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if exist {
			break
		}
	}

	arg := db.CreateURLParams{
		OriginUrl: req.OriginUrl,
		ShortUrl:  shortUrl,
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

	// 檢查布隆過濾器
	exist, err := server.redis.ExistBloom(ctx, req.ShortUrl)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !exist {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("布隆過濾器內無資料")))
		return
	}

	// redis 取資料
	redisUrl, haveData, err := server.redis.GetData(ctx, req.ShortUrl)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	if haveData {
		ctx.Redirect(http.StatusMovedPermanently, redisUrl.OriginUrl)
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

	// 放入 redis
	err = server.redis.SetData(ctx, req.ShortUrl, url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, url.OriginUrl)
}
