package uistrategy

import (
	"context"
	"fmt"
	"path"

	log "github.com/dnitsch/simplelog"
	"github.com/dnitsch/uistrategy/internal/config"
	"github.com/dnitsch/uistrategy/internal/util"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Element struct {
	CSSSelector *string `yaml:"cssSelector,omitempty" json:"cssSelector,omitempty"`
	XPath       *string `yaml:"xPath,omitempty" json:"xPath,omitempty"`
	Value       *string `yaml:"value,omitempty" json:"value,omitempty"`
	Must        bool    `yaml:"must" json:"must"`
}

type Auth struct {
	Username        Element `yaml:"username" json:"username"`
	Password        Element `yaml:"password" json:"password"`
	ConfirmPassword Element `yaml:"confirmPassword,omitempty json:"confirmPassword,omitempty`
	RequireConfirm  bool    `yaml:"requireConfirm,omitempty" json:"requireConfirm,omitempty"`
	Navigate        string  `yaml:"navigate" json:"navigate"`
	Submit          Element `yaml:"submit" json:"submit"`
}

type LoggedInPage struct {
	page    *rod.Page
	browser *rod.Browser
	log     log.Loggeriface
}

type WebConfig struct {
	Headless bool `yaml:"headless" json:"headless"`
	// if enabled it will store session data on disk
	// when used in CI, if you also want this enabled
	// you should also cache the default location of where the cache is:
	// ~/.uistratetegy-data
	PersistSessionOnDisk bool `yaml:"persist" json:"persist"`
}

// BaseConfig is the base config object
// each web session will have its own go routine to run the entire session
// Auth -> LoggedInPage ->[]Actions
type BaseConfig struct {
	BaseUrl string `yaml:"baseUrl" json:"baseUrl"`
	// Timeout will initialises a copy of the page with a context Timeout
	Timeout           int        `yaml:"timeout" json:"timeout"`
	WebConfig         *WebConfig `yaml:"webConfig,omitempty" json:"webConfig,omitempty"`
	ReuseBrowserCache bool       `yaml:"reuseBrowserCache" json:"reuseBrowserCache"`
}

type UiStrategyConf struct {
	BaseConfig
	Auth    Auth         `yaml:"auth" json:"auth"`
	Actions []ViewAction `yaml:"actions" json:"actions"`
}

// ViewAction defines a single action to do
// e.g. look up item, input text, click/swipe
// can include Assertion that action successfully occured
type ViewAction struct {
	navigate       string          `yaml:"-" json:"-"`
	Name           string          `yaml:"name" json:"name"`
	Navigate       string          `yaml:"navigate" json:"navigate"`
	ElementActions []ElementAction `yaml:"elementActions" json:"elementActions"`
}

type ElementAction struct {
	Name       string  `yaml:"name" json:"name"`
	Timeout    int     `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Element    Element `yaml:"element" json:"element"`
	InputText  *string `yaml:"inputText,omitempty" json:"inputText,omitempty"`
	ClickSwipe bool    `yaml:"clickSwipe" json:"clickSwipe"`
	Assert     any     `yaml:"assert,omitempty" json:"assert,omitempty"`
}

// WithNavigate correctly formats the Navigate URL to include the full baseUrl
func (a *ViewAction) WithNavigate(baseUrl string) *ViewAction {
	a.navigate = fmt.Sprintf("%s%s", baseUrl, a.Navigate)
	return a
}

// Web is the single instance struct
type Web struct {
	datadir  *string
	launcher *launcher.Launcher
	browser  *rod.Browser
	log      log.Loggeriface
	config   BaseConfig
}

// New returns an initialised instance of Web struct
// with the provided BaseConfig
func New(conf BaseConfig) *Web {
	ddir := path.Join(util.HomeDir(), fmt.Sprintf(".%s-data", config.SELF_NAME))

	l := launcher.New().
		Set("proxy-server", "localhost:8888").
		Headless(false).
		Devtools(false).
		Leakless(true)

	// url := l.UserDataDir(ddir).MustLaunch()
	url := l.MustLaunch()

	browser := rod.New().
		ControlURL(url).
		MustConnect().NoDefaultDevice()

	return &Web{
		datadir:  &ddir,
		launcher: l,
		browser:  browser,
		config:   conf,
	}
}

func (w *Web) WithLogger(l log.Logger) *Web {
	w.log = l
	return w
}

// Drive runs a single UIStrategy in the same logged in session
func (web *Web) Drive(ctx context.Context, auth Auth, action []ViewAction) []error {
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
	for _, v := range action {
		v = *v.WithNavigate(web.config.BaseUrl)
		if e := page.PerformActions(v); e != nil {
			errs = append(errs, err)
		}
	}
	// logOut
	return errs
}

// DoAuth performs the required Authentication
// in the browser and returns a authed Page
func (web *Web) DoAuth(auth Auth) (*LoggedInPage, error) {

	util.WriteDataDir(*web.datadir)
	web.log.Debug("begin auth")
	page := web.browser.MustPage(web.config.BaseUrl + auth.Navigate).MustWaitLoad()
	// page.MustHas("span", "Login")
	// .MustWaitInvisible()
	lp := &LoggedInPage{page, web.browser, web.log}
	// wait := lp.page.MustWaitRequestIdle()
	// wait()
	// determine which selector is available special case for AuthHandler
	mustSelector(lp, auth.Username).MustInput(*auth.Username.Value)
	mustSelector(lp, auth.Password).MustInput(*auth.Password.Value)
	mustSelector(lp, auth.Submit).MustClick().MustWaitInvisible()
	page.MustWaitLoad()
	// lp.page.MustWaitLoad()

	web.log.Debug("end auth")
	return lp, nil
}

// PerformAction handles a single action
func (p *LoggedInPage) PerformActions(action ViewAction) error {
	if err := p.page.Navigate(action.navigate); err != nil {
		return err
	}
	p.page.MustWaitLoad()
	p.log.Debugf("navigated to: %s", action.navigate)
	for _, a := range action.ElementActions {
		p.log.Debugf("starting to perform action: %s", a.Name)
		// if Assert is specified do not perform action only assert on pageObjects
		// extend screenshots here
		actionedElement, err := p.DetermineActionType(a, p.DetermineActionElement(a))
		if err != nil {
			p.log.Debugf("action: %s, errored with %#+v", a.Name, err)
			return err
		}
		actionedElement.MustWaitLoad()
		p.log.Debugf("completed action: %s", a.Name)
	}
	return nil
}

// DetermineActionType returns the rod.Element with correct action
func (p *LoggedInPage) DetermineActionElement(action ElementAction) *rod.Element {
	return determinActionElement(p, action.Element)
}

func determinActionElement(lp *LoggedInPage, elem Element) *rod.Element {
	if elem.Must {
		return mustSelector(lp, elem)
	}
	return maySelector(lp, elem)
}

func mustSelector(lp *LoggedInPage, elem Element) *rod.Element {
	lp.log.Debugf("mustElement: %#v", elem)
	if elem.XPath != nil {
		return lp.page.MustElementX(*elem.XPath)
	}
	return lp.page.MustElement(*elem.CSSSelector)
}

func maySelector(lp *LoggedInPage, elem Element) *rod.Element {
	lp.log.Debugf("mayElement: %#v", elem)
	if elem.XPath != nil {
		elx, err := lp.page.ElementX(*elem.XPath)
		if err != nil {
			lp.log.Debugf("element: %#v", elem)
			lp.log.Errorf("error %+v", err)
		}
		return elx
	}

	el, err := lp.page.Element(*elem.CSSSelector)
	if err != nil {
		lp.log.Debugf("element: %#v", elem)
		lp.log.Errorf("error %+v", err)
	}
	return el
}

// DetermineActionType returns the rod.Element with correct action
// either Click/Swipe or Input
// when Input is selected - ensure you have specified the input HTML element
// as the enclosing elements may not always allow for input...
func (p *LoggedInPage) DetermineActionType(action ElementAction, elem *rod.Element) (*rod.Element, error) {
	// if Value is present on the actionElement then always give preference
	if action.Element.Value != nil {
		p.log.Debugf("action includes Value on actionElement: ")
		return elem.MustInput(*action.Element.Value), nil
	}
	// if Value is missing then click element
	// simplified for now
	return elem.MustClick(), nil
	// // TODO: expand this into a more switch statement type implementation
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
