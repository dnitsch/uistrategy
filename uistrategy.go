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
	Value       string  `yaml:"value,omitempty" json:"value,omitempty"`
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

// func (lp *LoggedInPage) WithLogger(l log.Logger) *LoggedInPage {
// 	lp.log = l
// 	return lp
// }

type UiStrategyConf struct {
	BaseUrl   string `yaml:"baseUrl" json:"baseUrl"`
	Timeout   int    `yaml:"timeout" json:"timeout"`
	WebConfig struct {
		Headless bool `yaml:"headless" json:"headless"`
		// if enabled it will store session data on disk
		// when used in CI, if you also want this enabled
		// you should also cache the default location of where the cache is:
		// ~/.uistratetegy-data
		PersistSessionOnDisk bool `yaml:"persist" json:"persist"`
	} `yaml:"webConfig,omitempty" json:"webConfig,omitempty"`
	ReuseBrowserCache bool       `yaml:"reuseBrowserCache" json:"reuseBrowserCache"`
	Auth              Auth       `yaml:"auth" json:"auth"`
	Actions           []UIAction `yaml:"actions" json:"actions"`
}

type UIAction struct {
	Name string `yaml:"name" json:"name"`
	Action
}

// Action defines a single action to do
// e.g. look up item, input text, click/swipe
// can include Assertion that action successfully occured
type Action struct {
	navigate   string  `yaml:"-" json:"-"`
	Timeout    int     `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Navigate   string  `yaml:"navigate" json:"navigate"`
	Element    Element `yaml:"element" json:"element"`
	InputText  *string `yaml:"inputText,omitempty" json:"inputText,omitempty"`
	ClickSwipe bool    `yaml:"clickSwipe" json:"clickSwipe"`
	Assert     any     `yaml:"assert,omitempty" json:"assert,omitempty"`
}

func (a *Action) WithNavigate(baseUrl string) *Action {
	a.navigate = fmt.Sprintf("%s%s", baseUrl, a.Navigate)
	return a
}

type Web struct {
	datadir  *string
	launcher *launcher.Launcher
	browser  *rod.Browser
	log      log.Loggeriface
}

// New returns an initialised instance of Web struct
func New() *Web {
	ddir := path.Join(util.HomeDir(), fmt.Sprintf(".%s-data", config.SELF_NAME))

	l := launcher.New().
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
	}
}

func (w *Web) WithLogger(l log.Logger) *Web {
	w.log = l
	return w
}

// ActionPerform wraps around a single action
func (web *Web) ActionsPerform(ctx context.Context, auth Auth, action []Action, config UiStrategyConf) []error {
	var errs []error
	// and re-use same browser for all calls
	// defer web.browser.MustClose()

	// doAuth
	page, err := web.DoAuth(auth)
	defer web.browser.MustClose()
	if err != nil {
		return []error{err}
	}

	// start driving in that session
	for _, v := range action {
		v = *v.WithNavigate(config.BaseUrl)
		if e := page.PerformAction(v); e != nil {
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

	page := web.browser.MustPage(auth.Navigate)
	lp := &LoggedInPage{page, web.browser, web.log}
	// determine which selector is available special case for AuthHandler
	determinActionElement(lp, auth.Username).MustInput(auth.Username.Value)
	determinActionElement(lp, auth.Password).MustInput(auth.Password.Value)
	determinActionElement(lp, auth.Submit).MustClick()
	wait := lp.page.MustWaitRequestIdle()
	wait()
	return lp, nil
}

// DoRegistration performs the required registration
// currently unused but will be a special dispensation
// for when the UI run of actions will require a registration of users
func (web *Web) DoRegistration(auth Auth) (*LoggedInPage, error) {

	util.WriteDataDir(*web.datadir)

	page := web.browser.MustPage(auth.Navigate)
	lp := &LoggedInPage{page, web.browser, web.log}
	// determine which selector is available special case for AuthHandler
	determinActionElement(lp, auth.Username).MustInput(auth.Username.Value)
	determinActionElement(lp, auth.Password).MustInput(auth.Password.Value)
	if auth.RequireConfirm {
		determinActionElement(lp, auth.ConfirmPassword).MustInput(auth.ConfirmPassword.Value)
	}
	determinActionElement(lp, auth.Submit).MustClick()
	wait := lp.page.MustWaitRequestIdle()
	wait()
	return lp, nil
}

// PerformAction handles a single action
func (p *LoggedInPage) PerformAction(action Action) error {
	if err := p.page.Navigate(action.navigate); err != nil {
		return err
	}
	// if Assert is specified do not perform action only assert on pageObjects
	// extend screenshots here
	if _, err := p.DetermineActionType(action, p.DetermineActionElement(action)); err != nil {
		return err
	}
	return nil
}

// DetermineActionType returns the rod.Element with correct action
func (p *LoggedInPage) DetermineActionElement(action Action) *rod.Element {
	return determinActionElement(p, action.Element)
}

func determinActionElement(lp *LoggedInPage, elem Element) *rod.Element {
	if elem.Must {
		return mustSelector(lp, elem)
	}
	return maySelector(lp, elem)
}

func mustSelector(lp *LoggedInPage, elem Element) *rod.Element {
	if elem.XPath != nil {
		return lp.page.MustElementX(*elem.XPath)
	}
	return lp.page.MustElement(*elem.CSSSelector)
}

func maySelector(lp *LoggedInPage, elem Element) *rod.Element {
	if elem.XPath != nil {
		if el, err := lp.page.ElementX(*elem.XPath); el != nil {
			fmt.Printf("error %+v", err)
			return el
		}
	}
	if el, err := lp.page.Element(*elem.CSSSelector); el != nil {
		fmt.Printf("error %+v", err)
		return el
	}
	return nil
}

// DetermineActionType returns the rod.Element with correct action
// either Click/Swipe or Input
// when Input is selected - ensure you have specified the input HTML element
// as the enclosing elements may not always allow for input...
func (p *LoggedInPage) DetermineActionType(action Action, elem *rod.Element) (*rod.Element, error) {
	if action.ClickSwipe {
		return elem.MustClick(), nil
	}
	if action.InputText != nil {
		return elem.MustInput(*action.InputText), nil
	}
	return nil, fmt.Errorf("must specify either click/swipe action or input text must not be <nil>")
}
