package context

type ContextHeader struct {
	Authorization      string `header:"Authorization"` // "Authorization": ""
	SuperAuthorization string `header:"Super-Authorization"`
	XUserId            string `header:"X-User-Id"`         // "X-User-Id": ""
	XUserType          string `header:"X-UserType"`        // "X-UserType": ""
	XSourceID          string `header:"X-SourceID"`        // "X-SourceID": ""
	XDeviceID          string `header:"X-DeviceID"`        // "X-DeviceID": "xxxxxxxxx"
	XDeviceType        string `header:"X-DeviceType"`      // "X-DeviceType": "Samsung Galaxy 11 / iPhone 12 / OnePlus"
	XRequestID         string `header:"X-Request-Id"`      // "X-Request-Id": "0bc86576911d7468f7bbd7d55fb2b72d"
	XCorrelationID     string `header:"X-Correlation-Id"`  // "X-Correlation-Id": ""
	XClientLocalIP     string `header:"X-ClientLocalIP"`   // "X-ClientLocalIP": ""
	XClientPublicIP    string `header:"X-ClientPublicIP"`  // "X-ClientPublicIP": ""
	XMACAddress        string `header:"X-MACAddress"`      // "X-MACAddress": ""
	XSystemInfo        string `header:"X-SystemInfo"`      // "X-SystemInfo": ""
	XLocation          string `header:"X-Location"`        // "X-Location": ""
	XAppID             string `header:"X-AppID"`           // "X-AppID": "XXXXXXX"
	XProductVersion    string `header:"X-ProductVersion"`  // "X-ProductVersion": ""
	XOperatingSystem   string `header:"X-OperatingSystem"` // "X-OperatingSystem": "Windows"
	XPlatform          string `header:"X-Platform"`        // "X-Platform": "Android / iOS / IE / Chrome / Safari"
	XAppVersion        string `header:"X-App-version"`     // "X-App-version": "39.0.1 / Chrome 90.04 / EDGE 11"
	XOSVersion         string `header:"X-OS-version"`      // "X-OS-version": "Android 10 / iOS 10 / Windows 10"
	XAPIVersion        string `header:"X-Api-Version"`     // "X-Api-Version": "1"
	IsAdmin            bool
}
