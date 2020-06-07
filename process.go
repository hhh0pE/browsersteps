package browsersteps

import "github.com/cucumber/godog"

func (b *BrowserSteps) buildProcessSteps(s *godog.Suite) {
	s.Step(`^I accept alert$`, func() error { return b.GetWebDriver().AcceptAlert() })
	s.Step(`^I dismiss alert$`, func() error { return b.GetWebDriver().DismissAlert() })
}
