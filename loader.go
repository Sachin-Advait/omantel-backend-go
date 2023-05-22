package main

import (
	"time"

	"gitlab.com/rendezvousPoint/incentive_calculator/business/admin"
	adminService "gitlab.com/rendezvousPoint/incentive_calculator/business/admin/impl"
	"gitlab.com/rendezvousPoint/incentive_calculator/business/health"
	healthService "gitlab.com/rendezvousPoint/incentive_calculator/business/health/impl"
	"gitlab.com/rendezvousPoint/incentive_calculator/business/scheduler"
	schedulerService "gitlab.com/rendezvousPoint/incentive_calculator/business/scheduler/impl"
	"gitlab.com/rendezvousPoint/incentive_calculator/business/user"
	userService "gitlab.com/rendezvousPoint/incentive_calculator/business/user/impl"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	databaseImpl "gitlab.com/rendezvousPoint/incentive_calculator/repository/database/impl"
	fcmImpl "gitlab.com/rendezvousPoint/incentive_calculator/repository/fcm/impl"
	cacheImpl "gitlab.com/rendezvousPoint/incentive_calculator/repository/redis/impl"
	soapImpl "gitlab.com/rendezvousPoint/incentive_calculator/repository/soap/impl"
	cache "gitlab.com/rendezvousPoint/incentive_calculator/utils/cache"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	database "gitlab.com/rendezvousPoint/incentive_calculator/utils/database"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/fcm"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/http"
	jwthelper "gitlab.com/rendezvousPoint/incentive_calculator/utils/jwtHelper"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

type Components struct {
	HealthProvider health.HealthCheckProvider
	AdminProvider  admin.AdminProvider
	UserProvider   user.UserProvider
	Scheduler      scheduler.Scheduler
}

func LoadApplicationComponents() Components {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("LoadApp: error loading application config file, env - " + flags.Env())
	}

	jwthelper.JwtSecretKey = applicationConfig.GetString(constant.JwtKey)
	jwthelper.AdminJwtSecretKey = applicationConfig.GetString(constant.AdminJwtKey)
	if jwthelper.JwtSecretKey == "" || jwthelper.AdminJwtSecretKey == "" || jwthelper.JwtSecretKey == jwthelper.AdminJwtSecretKey {
		logger.FileLogger.Fatal().Err(err).Msg("jwt secret key invalid, env - " + flags.Env())
	}

	models.ShortCode = applicationConfig.GetString(constant.ShortCodeKey)
	models.RequestorChannelId = applicationConfig.GetString(constant.RequestorChannelIdKey)
	models.CaptchaSitesEnableKey = applicationConfig.GetString(constant.CaptchaSitesEnableKey)
	models.CaptchaSecretKey = applicationConfig.GetString(constant.CaptchaSecretKey)
	models.VerifyCaptchaGoogleApi = applicationConfig.GetString(constant.VerifyCaptchaGoogleApi)
	if models.CaptchaSitesEnableKey == "" || models.CaptchaSecretKey == "" || models.VerifyCaptchaGoogleApi == "" {
		logger.FileLogger.Fatal().Err(err).Msg("captcha keys not found, env - " + flags.Env())
	}

	adminConfig := &configs.AdminConfig{
		FCMOfferTopic:        applicationConfig.GetString(constant.FCMOfferTopicKey),
		FCMTrainingTopic:     applicationConfig.GetString(constant.FCMTrainingTopicKey),
		FCMNotificationTopic: applicationConfig.GetString(constant.FCMNotificationTopicKey),
	}

	userConfig := &configs.UserConfig{
		OTPUrl:                 applicationConfig.GetString(constant.OTPEndpoint),
		FCMOfferTopic:          applicationConfig.GetString(constant.FCMOfferTopicKey),
		FCMTrainingTopic:       applicationConfig.GetString(constant.FCMTrainingTopicKey),
		FCMNotificationTopic:   applicationConfig.GetString(constant.FCMNotificationTopicKey),
		ReportCacheTTL:         time.Duration(applicationConfig.GetInt(constant.ReportCacheTTLInHours)) * time.Hour,
		SessionExpiryTime:      time.Duration(applicationConfig.GetInt(constant.SessionExpiryTimeMinutes)) * time.Minute,
		RefreshTokenExpiryTime: time.Duration(applicationConfig.GetInt(constant.RefreshTokenExpiryHours)) * time.Hour,
	}

	dbClient := database.Get()
	redisClient := cache.Get()
	httpClient := http.Get()
	fcmClient := fcm.Get()

	dbService := databaseImpl.GetDBServiceInstance(dbClient)
	cacheService := cacheImpl.GetCacheServiceInstance(redisClient)
	fcmService := fcmImpl.GetFCMServiceInstance(fcmClient)
	soapService := soapImpl.GetSOAPServiceInstance(httpClient)

	return Components{
		HealthProvider: healthService.GetHealthService(dbService, cacheService),
		AdminProvider:  adminService.GetAdminService(adminConfig, dbService, cacheService, fcmService),
		UserProvider:   userService.GetUserService(userConfig, dbService, cacheService, fcmService, soapService),
		Scheduler:      schedulerService.GetSchedulerService(dbService),
	}
}
