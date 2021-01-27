package runner

// quick and dirty observer implementation
// to know when an Input causes newline

type InputDetector struct {
	update func()
}

type InputNotifier struct {
	observers []InputDetector
}

// register a new detector to listen to the subject
func (s *InputNotifier) register(i InputDetector) {
	s.observers = append(s.observers, i)
}

// remove all listeners (easy to implement and all i needed really)
func (s *InputNotifier) unregisterAll() {
	s.observers = []InputDetector{}
}

// notify all detectors of a change
func (s InputNotifier) notifyInput() {
	for _, detector := range s.observers {
		detector.detected()
	}
}

// the detected func runs the passed update func from initialization,
// in order to have access to necessary data
func (d *InputDetector) detected() {
	d.update()
}
