package melody

type envelope struct {
	t      int
	msg    []byte
	filter filterFunc
}

type envelope_to_self struct {
	t     int
	msg   []byte
	afunc AfterRunFunc
	sess  *Session
}
