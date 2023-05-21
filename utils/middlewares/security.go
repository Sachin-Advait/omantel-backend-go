package middlewares

import (
	"errors"
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

	expectedHostName := applicationConfig.GetStringSlice(constant.ExpectedHostName)
	expectedReferer := applicationConfig.GetStringSlice(constant.ExpectedReferer)

	if len(expectedHostName) == 0 || len(expectedReferer) == 0 {
		logger.FileLogger.Fatal().Err(err).Msg("SecurityCheck: blank security config, env - " + flags.Env())
	}

	return func(c *gin.Context) {

		var foundHostName bool
		for _, hostname := range expectedHostName {
			if c.Request.Host == hostname {
				foundHostName = true
				break
			}
		}

		if !foundHostName {
			logger.FileLogger.CtxError(c).Interface("expected", expectedHostName).Err(errors.New("found-" + c.Request.Host)).Msg("SecurityCheck: invalid hostname")
			fmt.Println("host error: ", c.Request.Host)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var foundExpectedReferer bool
		for _, referer := range expectedReferer {
			if strings.Contains(c.Request.Referer(), referer) {
				foundExpectedReferer = true
				break
			}
		}

		if !foundExpectedReferer {
			logger.FileLogger.CtxError(c).Interface("expected", expectedReferer).Err(errors.New("found-" + c.Request.Referer())).Msg("SecurityCheck: invalid referer")
			fmt.Println("referer error: ", c.Request.Referer())
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Next()
	}
}
