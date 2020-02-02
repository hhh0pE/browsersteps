package browsersteps

import (
	"errors"
	"fmt"
	"github.com/DATA-DOG/godog/gherkin"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/tebeka/selenium"
)

/*BrowserSteps represents a WebDriver context to run the Scenarios*/
type BrowserSteps struct {
	wd             selenium.WebDriver
	Capabilities   selenium.Capabilities
	DefaultURL     string
	URL            *url.URL
	ScreenshotPath string
	Timeout        time.Duration
	PingDuration   time.Duration
}

/*SetBaseURL sets the absolute URL used to complete relative URLs*/
func (b *BrowserSteps) SetBaseURL(url *url.URL) error {
	if !url.IsAbs() {
		return errors.New("BaseURL must be absolute")
	}
	b.URL = url
	return nil
}

func (b *BrowserSteps) iWriteTo(text, selector, by string) error {
	return RunWithTimeout(func() error {
		// Click the element
		element, err := b.GetWebDriver().FindElement(by, selector)
		if err != nil {
			return err
		}

		err = element.Clear()
		if err != nil {
			return err
		}
		return element.SendKeys(text)
	}, b.Timeout, b.PingDuration)

}

func (b *BrowserSteps) iClick(selector, by string) error {
	return RunWithTimeout(func() error {
		// Submit the element
		element, err := b.GetWebDriver().FindElement(by, selector)
		if err != nil {
			return err
		}
		return element.Click()
	}, b.Timeout, b.PingDuration)
}

func (b *BrowserSteps) iSubmit(selector, by string) error {
	return RunWithTimeout(func() error {
		// Submit the element
		element, err := b.GetWebDriver().FindElement(by, selector)
		if err != nil {
			return err
		}
		return element.Submit()
	}, b.Timeout, b.PingDuration)
}

func (b *BrowserSteps) iMoveTo(selector, by string) error {
	return RunWithTimeout(func() error {
		element, err := b.GetWebDriver().FindElement(by, selector)
		if err != nil {
			return err
		}
		return element.MoveTo(0, 0)
	}, b.Timeout, b.PingDuration)

}

//BeforeScenario is executed before each scenario
func (b *BrowserSteps) BeforeScenario(a interface{}) {
	var err error
	b.wd, err = selenium.NewRemote(b.Capabilities, b.DefaultURL)
	if err != nil {
		log.Panic(err)
	}
}

//AfterScenario is executed after each scenario
func (b *BrowserSteps) AfterScenario(a interface{}, err error) {
	if err != nil && b.ScreenshotPath != "" {
		var filename = "FAILED STEP.png"
		if gerkinDef, ok := a.(*gherkin.Scenario); ok {
			filename = "FAILED STEP -- " + gerkinDef.Name + ".png"
		}

		buff, err := b.GetWebDriver().Screenshot()
		if err != nil {
			fmt.Printf("Error %+v\n", err)
		}

		if _, err := os.Stat(b.ScreenshotPath); os.IsNotExist(err) {
			os.MkdirAll(b.ScreenshotPath, 0755)
		}
		pathname := filepath.Join(b.ScreenshotPath, filename)
		if write_err := ioutil.WriteFile(pathname, buff, 0644); write_err != nil {
			fmt.Errorf("AfterScenario saving screenshot error: %s", write_err.Error())
		}
	}
	b.GetWebDriver().Quit()
}

func (b *BrowserSteps) buildSteps(s *godog.Suite) {
	b.buildNavigationSteps(s)
	b.buildAssertionSteps(s)
	b.buildProcessSteps(s)

	s.Step(`^I am a anonymous user$`, func() error { return b.GetWebDriver().DeleteAllCookies() })

	s.Step(`^I write "([^"]*)" to "([^"]*)" `+ByOption+`$`, b.iWriteTo)
	s.Step(`^I click "([^"]*)" `+ByOption+`$`, b.iClick)
	s.Step(`^I submit "([^"]*)" `+ByOption+`$`, b.iSubmit)

	s.Step(`^I move to "([^"]*)" `+ByOption+`$`, b.iMoveTo)

}

//NewBrowserSteps starts a new BrowserSteps instance.
func NewBrowserSteps(s *godog.Suite, cap selenium.Capabilities, defaultURL string) *BrowserSteps {
	bs := &BrowserSteps{Capabilities: cap, DefaultURL: defaultURL, ScreenshotPath: os.Getenv("SCREENSHOT_PATH")}
	bs.Timeout = time.Second * 4
	bs.PingDuration = time.Second / 5
	bs.buildSteps(s)

	s.BeforeScenario(bs.BeforeScenario)
	s.AfterScenario(bs.AfterScenario)

	return bs
}
