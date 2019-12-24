package api

import (
	"enlabs"
	"enlabs/pkg/account"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const sourceTypeHeader = "Source-Type"

type httpServer struct {
	am  account.Manager
	log *logrus.Entry
}

func (hs *httpServer) Run(addr string) error {
	router := gin.Default()
	router.GET("/balance", hs.getBalance)
	router.POST("/transaction", hs.addTransaction)
	return router.Run(addr)
}

//HTTPServer http server
type HTTPServer interface {
	Run(addr string) error
}

func (hs *httpServer) addTransaction(g *gin.Context) {
	sourceType := g.GetHeader(sourceTypeHeader)
	var req addTransactionRequest
	if err := g.BindJSON(&req); err != nil {
		_ = g.AbortWithError(http.StatusBadRequest, err)
		return
	}
	tran, mapErr := req.MapToTransaction(enlabs.Source(sourceType))
	if mapErr != nil {
		_ = g.AbortWithError(http.StatusBadRequest, mapErr)
		return
	}
	if err := hs.am.AddTransaction(tran); err != nil {
		hs.log.WithError(err).Errorf("can't add transaction %s", req.ID)
		g.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	hs.log.WithField("tran_id", req.ID).Info("transaction recorded")
}

func (hs *httpServer) getBalance(g *gin.Context) {
	balance, getBalanceErr := hs.am.GetBalance()
	if getBalanceErr != nil {
		hs.log.WithError(getBalanceErr).Error("can't get balance")
		g.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	g.JSON(http.StatusOK, balance)
}

//NewHTTPServer initialize http server
func NewHTTPServer(am account.Manager, log *logrus.Entry) HTTPServer {
	return &httpServer{
		am:  am,
		log: log,
	}
}
