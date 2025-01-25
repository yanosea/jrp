package jrp

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	jrpApp "github.com/yanosea/jrp/v2/app/application/jrp"
	wnjpnApp "github.com/yanosea/jrp/v2/app/application/wnjpn"
	"github.com/yanosea/jrp/v2/app/infrastructure/database"
	"github.com/yanosea/jrp/v2/app/infrastructure/wnjpn/query_service"
	"github.com/yanosea/jrp/v2/app/presentation/api/jrp/formatter"

	"github.com/yanosea/jrp/v2/pkg/proxy"
)

var (
	format = "json"
)

// BindGetJrpHandler binds the getJrp handler to the server.
func BindGetJrpHandler(g proxy.Group) {
	g.GET("/jrp", getJrp)
}

// getJrp is a handler that returns a random Japanese phrase.
func getJrp(c echo.Context) error {
	connManager := database.GetConnectionManager()
	if connManager == nil {
		log.Error("Connection manager is not initialized...")
		return c.NoContent(http.StatusInternalServerError)
	}

	if _, err := connManager.GetConnection(database.WNJpnDB); err != nil {
		log.Error("Failed to get a connection to the database...")
		return c.NoContent(http.StatusInternalServerError)
	}

	wordQueryService := query_service.NewWordQueryService()
	fwuc := wnjpnApp.NewFetchWordsUseCase(wordQueryService)

	fwoDtos, err := fwuc.Run(
		c.Request().Context(),
		"jpn",
		[]string{"a", "v", "n"},
	)
	if err != nil {
		log.Error("Failed to fetch words...")
		return c.NoContent(http.StatusInternalServerError)
	}

	var gjiDtos []*jrpApp.GenerateJrpUseCaseInputDto
	for _, fwoDto := range fwoDtos {
		gjiDto := &jrpApp.GenerateJrpUseCaseInputDto{
			WordID: fwoDto.WordID,
			Lang:   fwoDto.Lang,
			Lemma:  fwoDto.Lemma,
			Pron:   fwoDto.Pron,
			Pos:    fwoDto.Pos,
		}
		gjiDtos = append(gjiDtos, gjiDto)
	}

	gjuc := jrpApp.NewGenerateJrpUseCase()
	gjoDto := gjuc.RunWithRandom(gjiDtos)

	f, err := formatter.NewFormatter(format)
	if err != nil {
		log.Error("Failed to create a new formatter...")
		return c.NoContent(http.StatusInternalServerError)
	}

	body, err := f.Format(gjoDto)
	if body == nil || err != nil {
		log.Error("Failed to format the output...")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSONBlob(http.StatusOK, body)
}
