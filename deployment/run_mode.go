package deployment

type RunMode uint8

const Production, Development, Test = 1, 2, 3

var CurrentRunMode RunMode

func IsTest() bool {
	return CurrentRunMode == Test
}
