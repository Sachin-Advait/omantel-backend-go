package scheduler

func LoadIntrinsicValues() {

	scheduleProvide.LoadPlansFromDB(true)
	scheduleProvide.LoadTargetFromDB(true)
	scheduleProvide.LoadMaxPayFromDB(true)
	scheduleProvide.LoadMaxExtraPayFromDB(true)
	scheduleProvide.LoadThresholdCheckFromDB(true)
	scheduleProvide.LoadValueOfPointFromDB(true)
	scheduleProvide.LoadVOCFromDB(true)

	PrintIntrinsicValues()
}
