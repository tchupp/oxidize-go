package blockhttp

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/tclchiam/oxidize-go/blockchain"
	"github.com/tclchiam/oxidize-go/blockchain/entity"
	"github.com/tclchiam/oxidize-go/encoding"
	"github.com/tclchiam/oxidize-go/server/httpserver"
)

type E map[string]interface{}

func HandleBlockByIndex(bc blockchain.Blockchain) echo.HandlerFunc {
	return func(c echo.Context) error {
		rawIndex := c.Param("index")
		blockIndex, err := strconv.Atoi(rawIndex)
		if err != nil {
			return c.JSON(http.StatusBadRequest, E{"message": "invalid index: " + rawIndex})
		}

		block, err := bc.BlockByIndex(uint64(blockIndex))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, E{"message": "error finding block"})
		}
		if block == nil {
			return c.JSON(http.StatusNotFound, E{"message": "block with index not found: " + rawIndex})
		}

		wireBlock := encoding.ToWireBlock(block)
		return c.JSON(http.StatusOK, wireBlock)
	}
}

func HandleBlockByHash(bc blockchain.Blockchain) echo.HandlerFunc {
	return func(c echo.Context) error {
		blockHash, err := entity.NewHashFromString(c.Param("hash"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, E{"message": "invalid hash: " + c.Param("hash")})
		}

		block, err := bc.BlockByHash(blockHash)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, E{"message": "error finding block"})
		}
		if block == nil {
			return c.JSON(http.StatusNotFound, E{"message": "block with hash not found: " + c.Param("hash")})
		}

		wireBlock := encoding.ToWireBlock(block)
		return c.JSON(http.StatusOK, wireBlock)
	}
}

func RegisterBlockServer(server *httpserver.Server, bc blockchain.Blockchain) {

	server.GET("/blocks/", func(c echo.Context) error { return c.Redirect(http.StatusFound, "/blocks/index") })
	server.GET("/blocks/index", func(c echo.Context) error { return c.Redirect(http.StatusFound, "/blocks/index/0") })
	server.GET("/blocks/index/:index", HandleBlockByIndex(bc))
	server.GET("/blocks/hash/:hash", HandleBlockByHash(bc))
}
