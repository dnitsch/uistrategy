package uistrategy

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"time"

	log "github.com/dnitsch/simplelog"
	"github.com/dnitsch/uistrategy/internal/util"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type Element struct {
	// Selector can be a CSSStyle selector or XPath
	Selector *string `yaml:"selector,omitempty" json:"selector,omitempty"`
	Value    *string `yaml:"value,omitempty" json:"value,omitempty"`
	Must     bool    `yaml:"must" json:"must"`
	Timeout  int
}

type ActionReport struct {
	Name       string `json:"name"`
	Screenshot string `json:"screenshot"`
	Errored    bool   `json:"errored"`
	Message    string `json:"message"`
}

type ViewReport struct {
	Name    string         `json:"name"`
	Message string         `json:"message"`
	Actions []ActionReport `json:"actions"`
}

type LoggedInPage struct {
	*Web
	page   *rod.Page
	errors []error
	// report []ViewReport
}

type WebConfig struct {
	Headless bool `yaml:"headless" json:"headless"`
	// if enabled it will store session data on disk
	// when used in CI, if you also want this enabled
	// you should also cache the default location of where the cache is:
	// ~/.uistratetegy-data
	PersistSessionOnDisk bool `yaml:"persist" json:"persist"`
	// Timeout will initialises a copy of the page with a context Timeout
	Timeout           int    `yaml:"timeout" json:"timeout"`
	BrowserPathExec   string `yaml:"execPath" json:"execPath"`
	UserMode          bool   `yaml:"userMode" json:"userMode"`
	DataDir           string `yaml:"dataDir" json:"dataDir"`
	ReuseBrowserCache bool   `yaml:"reuseBrowserCache" json:"reuseBrowserCache"`
}

// BaseConfig is the base config object
// each web session will have its own go routine to run the entire session
// Auth -> LoggedInPage ->[]Actions
type BaseConfig struct {
	BaseUrl         string     `yaml:"baseUrl" json:"baseUrl"`
	LauncherConfig  *WebConfig `yaml:"browserConfig,omitempty" json:"browserConfig,omitempty"`
	ContinueOnError bool       `yaml:"continueOnError" json:"continueOnError"`
}

type UiStrategyConf struct {
	Setup BaseConfig `yaml:"setup" json:"setup"`
	// Auth is optional
	// should be omitted for apps that do not require a login
	Auth    *Auth         `yaml:"auth,omitempty" json:"auth,omitempty"`
	Actions []*ViewAction `yaml:"actions" json:"actions"`
}

// ViewAction defines a single action to do
// e.g. look up item, input text, click/swipe
// can include Assertion that action successfully occured
type ViewAction struct {
	navigate string `yaml:"-" json:"-"`
	// report attr
	message        string           `yaml:"-" json:"-"`
	Iframe         *string          `yaml:"iframe,omitempty" json:"iframe,omitempty"`
	Name           string           `yaml:"name" json:"name"`
	Navigate       string           `yaml:"navigate" json:"navigate"`
	ElementActions []*ElementAction `yaml:"elementActions" json:"elementActions"`
}

type ElementAction struct {
	Name    string  `yaml:"name" json:"name"`
	Element Element `yaml:"element" json:"element"`
	Assert  bool    `yaml:"assert,omitempty" json:"assert,omitempty"`
	// report attrs
	message    string
	errored    bool
	screenshot string
	// TODO: currently unused
	// Timeout int     `yaml:"timeout" json:"timeout"`
	// InputText  *string `yaml:"inputText,omitempty" json:"inputText,omitempty"`
	// ClickSwipe bool    `yaml:"clickSwipe" json:"clickSwipe"`
}

// WithNavigate correctly formats the Navigate URL to include the full baseUrl
func (a *ViewAction) WithNavigate(baseUrl string) *ViewAction {
	a.navigate = fmt.Sprintf("%s%s", baseUrl, a.Navigate)
	return a
}

// Web is the single instance struct
type Web struct {
	browser *rod.Browser
	log     log.Loggeriface
	config  BaseConfig
}

// New returns an initialised instance of Web struct
// with the provided BaseConfig
func New(conf BaseConfig) *Web {
	_ = util.InitDirDeps()
	url := newLauncher(conf.LauncherConfig).MustLaunch()
	browser := rod.New().
		ControlURL(url).
		MustConnect().NoDefaultDevice()

	return &Web{
		browser: browser,
		config:  conf,
	}
}

// newLauncher returns a launcher with specified properties
func newLauncher(webconf *WebConfig) *launcher.Launcher {
	// ddir := path.Join(util.HomeDir(), fmt.Sprintf(".%s-data", config.SELF_NAME))

	l := launcher.New()

	l.Leakless(true).Devtools(false).Headless(false)

	if webconf != nil {
		if webconf.UserMode {
			l = launcher.NewUserMode()
		}
		if len(webconf.BrowserPathExec) > 0 {
			if l != nil {
				l.Bin(webconf.BrowserPathExec)
			} else {
				l = launcher.New().Bin(webconf.BrowserPathExec)
			}
		}
		if len(webconf.DataDir) > 0 {
			l.UserDataDir(webconf.DataDir)
		}
		if webconf.Headless {
			l.Headless(true)
		}
	}

	return l
}

// WithLogger
func (w *Web) WithLogger(l log.Logger) *Web {
	w.log = l
	return w
}

// Drive runs a single UIStrategy in the same logged in session
func (web *Web) Drive(ctx context.Context, auth *Auth, allActions []*ViewAction) []error {
	var errs []error
	// and re-use same browser for all calls
	// defer web.browser.MustClose()
	defer web.browser.MustClose()

	// doAuth
	page, err := web.DoAuth(auth)
	if err != nil {
		return []error{err}
	}

	// start driving in that session
	for _, v := range allActions {
		v = v.WithNavigate(web.config.BaseUrl)
		if e := page.PerformActions(v); e != nil {
			errs = append(errs, e)
		}
	}
	// send to report builder here
	web.buildReport(allActions)
	// logOut
	return errs
}

// PerformAction handles a single action on Navigate'd page/view of SPA
func (lp *LoggedInPage) PerformActions(action *ViewAction) error {

	actionPage := lp.page
	waitNav := actionPage.MustWaitNavigation()
	if err := actionPage.Navigate(action.navigate); err != nil {
		return err
	}
	waitNav()
	// lp.page.MustWaitIdle()
	actionPage.MustWaitLoad()

	action.message = fmt.Sprintf("successfully navigated to: %s", action.navigate)
	lp.log.Debugf("navigated to: %s", action.navigate)

	if action.Iframe != nil {
		iframe, err := determinActionElement(lp.log, actionPage, Element{Selector: action.Iframe})
		if err != nil {
			return err
		}
		iframe.MustWaitLoad()
		action.message = fmt.Sprintf("%s\n%s", action.message, "will perform following actions inside an iframe")
		lp.page = iframe.MustFrame()
	}

	// lp.page.MustWaitIdle()
	if err := rod.Try(func() { actionPage.WaitLoad() }); err != nil {
		lp.log.Debugf("err on page load: %s", err.Error())
	}

	for _, a := range action.ElementActions {
		// perform action
		lp.log.Debugf("starting to perform action: %s", a.Name)
		// end perform action
		if skip, e := lp.handleActionError(actionPage, a, lp.performAction(actionPage, a)); e != nil {
			if skip {
				break
			}
			return e
		}
		lp.log.Debugf("completed action: %s", a.Name)
	}
	return nil
}

// handleActionError returns a skip error and error depending on config set up
func (p *LoggedInPage) handleActionError(page *rod.Page, a *ElementAction, err []error) (bool, error) {

	if len(err) > 0 && p.config.ContinueOnError {
		p.log.Debugf("action: %#v, errored with %v", a, err)
		p.log.Debugf("continue on error...")
		a.errored = true
		a.screenshot = p.captureAndSave(page)
		return true, nil
	}
	if len(err) > 0 {
		return false, fmt.Errorf("%+v", err)
	}
	return false, nil
}

// performAction handles finding the element and any related actions on it
// i.e. click or input
func (p *LoggedInPage) performAction(page *rod.Page, a *ElementAction) []error {
	rodElem, err := p.DetermineActionElement(page, a)
	a.errored = false
	a.screenshot = ""
	if err != nil {
		p.log.Debugf("action: %s, errored with %+#v", a.Name, err)
		// extend screenshots here
		a.message = fmt.Sprintf("locating element with selector: %s, errored with %+#v", *a.Element.Selector, err)
		a.errored = true
		a.screenshot = p.captureAndSave(page)
		p.errors = append(p.errors, err)
	}
	a.message = fmt.Sprintf("found element: %s", *a.Element.Selector)
	if err := p.DetermineActionType(a, rodElem); err != nil {
		p.log.Debugf("action: %s, errored with %v", a.Name, err)
		a.message = fmt.Sprintf("performing action on element with selector: %s, errored with %+v", *a.Element.Selector, err)
		a.errored = true
		a.screenshot = p.captureAndSave(page)
		p.errors = append(p.errors, err)
	}

	// also add results to Report outcome
	return p.errors
}

// DetermineActionType returns the rod.Element with correct action
func (lp *LoggedInPage) DetermineActionElement(page *rod.Page, action *ElementAction) (*rod.Element, error) {
	return determinActionElement(lp.log, page, action.Element)
}

// determinActionElement
func determinActionElement(log log.Loggeriface, page *rod.Page, elem Element) (*rod.Element, error) {
	log.Debugf("looking for element: %+v", elem)
	// when timeout is properly implemented
	// we need to wrap it in Try as it will panic on timeout
	// err := rod.Try(func() {
	// })
	if elem.Selector == nil {
		//
		return nil, fmt.Errorf("action must include selector")
	}

	type searchElemFunc func(selector string) (rod.Elements, error)
	searchfuncs := []searchElemFunc{
		func(selector string) (rod.Elements, error) {
			return page.Elements(selector)
		},
		func(selector string) (rod.Elements, error) {
			return page.ElementsX(selector)
		},
		// TODO: add more types here e.g. regex
		// func(selector string) (bool, *rod.Element, error) {
		// 	return page.HasR(selector)
		// },
	}
	// NOTE: shove this in known length channel slice and range over that so that it's done in parallel
	for k, searchEl := range searchfuncs {
		felems, err := searchEl(*elem.Selector)
		if err != nil {
			log.Debugf("error: %v occured when looking for element: %v, using method: %v", err.Error(), *elem.Selector, k)

		}
		if !felems.Empty() {
			// update report with success for step
			log.Debugf("found element using method: %v", k)
			return felems.First(), nil
		}
		log.Debugf("not found element using method: %v", k)
	}
	// update report with error for step
	log.Debugf("not found element using any search method")
	return nil, fmt.Errorf("element not found by selector: %v", *elem.Selector)
}

// DetermineActionType returns the rod.Element with correct action
// either Click/Swipe or Input
// when Input is selected - ensure you have specified the input HTML element
// as the enclosing elements may not always allow for input...
func (lp *LoggedInPage) DetermineActionType(action *ElementAction, elem *rod.Element) error {
	if elem == nil {
		if action.Assert {
			// TODO: custom errors here
			return fmt.Errorf("assert set to true. unable to perform action: %+v. element not found", action)
		}
		// update report with step miss
		lp.log.Debugf("element not found but ignoring error as assert is set to false")
		return nil
	}
	if elem != nil && action.Assert {
		// update report with step found
		// item found not performing action
		return nil
	}
	// if Value is present on the actionElement then always give preference
	if action.Element.Value != nil {
		lp.log.Debugf("action includes value on actionElement: %v", *action.Element.Value)
		elem.MustSelectAllText().MustInput("").MustInput(*action.Element.Value)
		return nil
	}

	// TODO: expand this into a more switch statement type implementation
	// allow - double tap/click, swipe, etc..
	elem.MustClick()
	elem.MustWaitLoad() // when clicked we wait for a
	// lp.page.MustWaitLoad()

	return nil
}

// captureAndSave wil store the captured image under the .uistrategy/captures/*.png
// it will swallow any errors and log them out
func (lp *LoggedInPage) captureAndSave(page *rod.Page) string {
	file := fmt.Sprintf(`.uistrategy/captures/%v.png`, time.Now().UnixNano())
	b, err := page.Screenshot(true, &proto.PageCaptureScreenshot{Format: "png", Clip: nil, FromSurface: true, Quality: util.Int(100)})
	if err != nil {
		lp.log.Debugf("failed to capture page: %+v", page)
	}

	r := bytes.NewReader(b)
	i, _ := png.Decode((r))
	// helper create dir if not exists
	w, err := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		lp.log.Debugf("failed to open write location for screenshot: %v", err)
	}

	if err := png.Encode(w, i); err != nil {
		lp.log.Debugf("failed to write screenshot: %v", err)
	}
	return file
}

// // DoRegistration performs the required registration
// // currently unused but will be a special dispensation
// // for when the UI run of actions will require a registration of users
// func (web *Web) DoRegistration(auth Auth) (*LoggedInPage, error) {

// 	util.WriteDataDir(*web.datadir)

// 	page := web.browser.MustPage(auth.Navigate)
// 	lp := &LoggedInPage{page, web.browser, web.log}
// 	// determine which selector is available special case for AuthHandler
// 	determinActionElement(lp, auth.Username).MustInput(*auth.Username.Value)
// 	determinActionElement(lp, auth.Password).MustInput(*auth.Password.Value)
// 	if auth.RequireConfirm {
// 		determinActionElement(lp, auth.ConfirmPassword).MustInput(*auth.ConfirmPassword.Value)
// 	}
// 	determinActionElement(lp, auth.Submit).MustClick()
// 	return lp, nil
// }
