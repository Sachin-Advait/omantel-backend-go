package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func SecurityCheck() gin.HandlerFunc {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("SecurityCheck: error loading application config file, env - " + flags.Env())
	}

	expectedHostName := applicationConfig.GetString(constant.ExpectedHostName)
	expectedReferer := applicationConfig.GetString(constant.ExpectedReferer)

	return func(c *gin.Context) {

		if c.Request.Host == "" || c.Request.Host != expectedHostName {
			fmt.Println("host error: ", c.Request.Host)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if !strings.Contains(c.Request.Referer(), expectedReferer) {
			fmt.Println("referer error: ", c.Request.Host)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Next()
	}
}
