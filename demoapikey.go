package streamcommons

import "time"

// APIKeyDemo is the demo API-key
const APIKeyDemo = "demo"

// DemoAPIKeyAllowedStart is the start date of range of time that data is allowed to be fetched using demo apikey
const DemoAPIKeyAllowedStart = 1577836800 * time.Second

// DemoAPIKeyAllowedEnd is the end date of range of time that data is allowed to be fetched using demo apikey
const DemoAPIKeyAllowedEnd = (1577836800 + 60*60) * time.Second
