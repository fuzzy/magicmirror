[![Build Status](https://travis-ci.org/splatpm/subhuman.svg?branch=master)](https://travis-ci.org/splatpm/subhuman)
## Subhuman (Humanizing library for Go)

The current 'go test' output is as follows:

	TestHumanSize       : Expected        32B and got        32B SUCCESS
	TestHumanSize       : Expected     3.16KB and got     3.16KB SUCCESS
	TestHumanSize       : Expected   315.66KB and got   315.66KB SUCCESS
	TestHumanSize       : Expected    30.83MB and got    30.83MB SUCCESS
	TestHumanSize       : Expected     3.01GB and got     3.01GB SUCCESS
	TestHumanTimeColon  : Expected   01:01:07 and got   01:01:07 SUCCESS
	TestHumanTimeColon  : Expected   00:00:32 and got   00:00:32 SUCCESS
	TestHumanTimeColon  : Expected   00:01:04 and got   00:01:04 SUCCESS
	TestHumanTimeConcise: Expected  00h00m32s and got  00h00m32s SUCCESS
	TestHumanTimeConcise: Expected  00h01m04s and got  00h01m04s SUCCESS
	TestHumanTimeConcise: Expected  01h01m07s and got  01h01m07s SUCCESS

#### Usage:

###### HumanSize(int64) string

fmt.Println(HumanSize(3232323232))
// Returns: "3.01GB"

###### HumanTimeColon(int64) string

fmt.Println(HumanTimeColon(32))
// Returns: "00:00:32"

###### HumanTimeConcise(int64) string

fmt.Println(HumanTimeConcise(32))
// Returns: "00h00m32s"

