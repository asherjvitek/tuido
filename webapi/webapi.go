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

func insertBoard(c *gin.Context) {
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

	c.JSON(http.StatusOK, nil)
}

func deleteBoard(c *gin.Context) {
	boardId := c.Params.ByName("boardId")
	id, err := strconv.Atoi(boardId)

	err = provider.DeleteBoard(id)

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

func insertList(c *gin.Context) {
	var list data.List
	err := c.BindJSON(&list)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = provider.InsertList(&list)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, list)
}

func updateList(c *gin.Context) {
	var list data.List
	err := c.BindJSON(&list)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = provider.UpdateList(list)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func deleteList(c *gin.Context) {
	listId := c.Params.ByName("listId")
	id, err := strconv.Atoi(listId)

	err = provider.DeleteList(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func insertItem(c *gin.Context) {
	var item data.Item
	err := c.BindJSON(&item)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = provider.InsertItem(&item)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func updateItem(c *gin.Context) {
	var item data.Item
	err := c.BindJSON(&item)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = provider.UpdateItem(item)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func deleteItem(c *gin.Context) {
	itemId := c.Params.ByName("itemId")
	id, err := strconv.Atoi(itemId)

	err = provider.DeleteItem(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func Run() {
	router := gin.Default()
	router.GET("/boards", getBoards)
	router.POST("/boards", insertBoard)
	router.PUT("/boards", updateBoard)
	router.DELETE("/boards/:boardId", deleteBoard)

	router.GET("/lists/:boardId", getLists)
	router.POST("/lists", insertList)
	router.PUT("/lists", updateList)
	router.DELETE("/lists/:listId", deleteList)

	router.POST("/items", insertItem)
	router.PUT("/items", updateItem)
	router.DELETE("/items/:itemId", deleteItem)

	router.Run("localhost:8082")
}
