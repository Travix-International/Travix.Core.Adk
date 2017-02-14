package appix

// GlobalArgs defines the basic set of arguments which can be used by all verbs.
type GlobalArgs struct {
	Verbose   bool
	TargetEnv string
	Timeout   int
}
