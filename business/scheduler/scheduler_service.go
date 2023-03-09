package scheduler

type Scheduler interface {
	LoadPlansFromDB(printLog bool) error
	LoadTargetFromDB(printLog bool) error
	LoadMaxPayFromDB(printLog bool) error
	LoadMaxExtraPayFromDB(printLog bool) error
	LoadThresholdCheckFromDB(printLog bool) error
	LoadValueOfPointFromDB(printLog bool) error
	LoadVOCFromDB(printLog bool) error
}
