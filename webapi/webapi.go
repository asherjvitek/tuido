package webapi

import (
	"net/http"
	"strconv"

	"tuido/data"

	"github.com/gin-gonic/gin"
)

var (
	provider data.Provider = data.SqliteProvider{}
)

func getBoards(c *gin.Context) {
	boards, err := provider.Boards()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

    c.JSON(http.StatusOK, boards)
}

func createBoard(c *gin.Context) {

	var board data.Board
	err := c.BindJSON(&board)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = provider.InsertBoard(&board)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, board)
}

func updateBoard(c *gin.Context) {

	var board data.Board
	err := c.BindJSON(&board)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = provider.UpdateBoard(board)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, board)
}

func deleteBoard(c *gin.Context) {

	boardId := c.Params.ByName("boardId")
	bId, err := strconv.Atoi(boardId)

	err = provider.DeleteBoard(bId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func getLists(c *gin.Context) {
	boardId := c.Params.ByName("boardId")

	bId, err := strconv.Atoi(boardId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	lists, err := provider.Lists(bId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    c.JSON(http.StatusOK, lists)
}

func Run() {
    router := gin.Default()
    router.GET("/boards", getBoards)
	// router.GET("/boards/:boardId", getBoard)
    router.POST("/boards", createBoard)
    router.PUT("/boards", updateBoard)
    router.DELETE("/boards", deleteBoard)

	router.GET("/lists/:boardId", getLists)

    router.Run("localhost:8082")
}
