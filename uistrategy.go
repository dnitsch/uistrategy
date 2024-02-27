package uistrategy

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
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
}

type ActionReportItem struct {
	Screenshot string   `json:"screenshot"`
	Errored    bool     `json:"errored"`
	Message    string   `json:"message"`
	Output     []string `json:"output"`
}

type ActionsReport map[string]ActionReportItem

type ViewReportItem struct {
	Message                   string         `json:"message"`
	CapturedHeaderRequestKeys map[string]any `json:"capturedHeaderReqKeys"`
	Actions                   ActionsReport  `json:"actions"`
}

type ViewReport map[string]ViewReportItem

type LoggedInPage struct {
	*Web
	page   *rod.Page
	errors UIStrategyError
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
	NoSandbox         bool   `yaml:"noSandbox" json:"noSandbox"`
}

// BaseConfig is the base config object
// each web session will have its own go routine to run the entire session
// Auth -> LoggedInPage ->[]Actions
type BaseConfig struct {
	BaseUrl         string     `yaml:"baseUrl" json:"baseUrl"`
	LauncherConfig  *WebConfig `yaml:"browserConfig,omitempty" json:"browserConfig,omitempty"`
	ContinueOnError bool       `yaml:"continueOnError" json:"continueOnError"`
	IsSinglePageApp bool       `yaml:"isSinglePageApp" json:"isSinglePageApp"`
	WriteReport     bool       `yaml:"writeReport" json:"writeReport"`
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
	message               string           `yaml:"-" json:"-"`
	capturedReqHeaders    map[string]any   `yaml:"-" json:"-"`
	Iframe                *IframeAction    `yaml:"iframe,omitempty" json:"iframe,omitempty"`
	Name                  string           `yaml:"name" json:"name"`
	Navigate              string           `yaml:"navigate" json:"navigate"`
	CaptureRequestHeaders []string         `yaml:"captureRequestHeaders" json:"captureRequestHeaders"`
	ElementActions        []*ElementAction `yaml:"elementActions" json:"elementActions"`
}

// IframeAction
type IframeAction struct {
	Selector string `yaml:"selector,omitempty" json:"selector,omitempty"`
	// WaitEval has to be in the form of a boolean return
	// e.g. `myVar !== null` or `(myVar !== null || document.title == "ready")`
	// the supplied value will be appended to an existing
	// `return document.readyState === 'complete' && ${WaitEval};`
	WaitEval string `yaml:"waitEval,omitempty" json:"waitEval,omitempty"`
}

// ElementAction
type ElementAction struct {
	Name               string  `yaml:"name" json:"name"`
	Element            Element `yaml:"element" json:"element"`
	Assert             bool    `yaml:"assert,omitempty" json:"assert,omitempty"`
	SkipOnErrorMessage string  `yaml:"skipOnErrorMessage,omitempty" json:"skipOnErrorMessage,omitempty"`
	CaptureOutput      bool    `yaml:"captureOutput,omitempty" json:"captureOutput,omitempty"`
	// Timeout in seconds to cancel the action if unable to MustClick or MustInput
	Timeout int `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	// report attrs
	message        string
	errored        bool
	screenshot     string
	capturedOutput []string
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
	output  io.Writer
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

func (w *Web) WithWriter(writer io.Writer) *Web {
	w.output = writer
	return w
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
		if webconf.NoSandbox {
			l.NoSandbox(webconf.NoSandbox)
		}
	}

	return l
}

// WithLogger
func (w *Web) WithLogger(l log.Logger) *Web {
	w.log = l
	return w
}

// UIStrategyError custom error handler
// TODO: enable mutex lock on this just in case
type UIStrategyError struct {
	errorMap []struct {
		view    string
		action  string
		message string
	}
}

func (e *UIStrategyError) setError(view, action, message string) {
	e.errorMap = append(e.errorMap, struct {
		view    string
		action  string
		message string
	}{view, action, message})
}

func (e *UIStrategyError) Error() string {
	if len(e.errorMap) > 0 {
		es := []string{"following errors occured:\n"}
		for _, v := range e.errorMap {
			es = append(es, fmt.Sprintf("\n\tin view: %s, performing action: %s, failed on: %s", v.view, v.action, v.message))
		}
		return strings.Join(es, "")
	}
	return ""
}

func (e *UIStrategyError) hasError() bool {
	return len(e.errorMap) > 0
}

// Drive runs a single UIStrategy in the same logged in session
// returns a custom error type with details of errors per action
func (web *Web) Drive(ctx context.Context, auth *Auth, allActions []*ViewAction) ([]*ViewAction, error) {
	uiErr := &UIStrategyError{}
	// and re-use same browser for all calls
	// defer web.browser.MustClose()
	defer web.browser.MustClose()

	// doAuth
	page, err := web.DoAuth(auth)
	if err != nil {
		uiErr.setError("auth", "login", err.Error())
		return allActions, uiErr
	}

	// start driving in that session
	for _, v := range allActions {
		if err := page.PerformActions(v.WithNavigate(web.config.BaseUrl)); err != nil {
			// returning on error if ContinueOnError was not specified
			return allActions, err
		}
	}
	// return errors to caller for visibility if any
	if page.errors.hasError() {
		return allActions, &page.errors
	}
	return allActions, nil
}

// PerformAction handles a single action on Navigate'd page/view of SPA
func (lp *LoggedInPage) PerformActions(action *ViewAction) error {
	actionPage, err := lp.navigateHelper(lp.page, action)
	if err != nil {
		return err
	}

	lp.log.Debugf("navigated to: %s", action.navigate)

	if action.Iframe != nil {
		iframe, err := lp.ensureIframeLoaded(actionPage, action)
		if err != nil {
			return err
		}
		actionPage = iframe
	}

	actionPage.MustWaitLoad()

	for _, ap := range action.ElementActions {
		// perform action
		lp.log.Debugf("starting to perform action: %s", ap.Name)
		// end perform action
		if skip, e := lp.handleActionError(actionPage, ap, lp.performAction(actionPage, ap)); e != nil {
			if skip {
				break
			}
			return e
		}
		lp.log.Debugf("completed action: %s", ap.Name)
	}
	return nil
}

// navigateHelper ensures page is navigated to and waited sufficient time to ensure
func (lp *LoggedInPage) navigateHelper(page *rod.Page, action *ViewAction) (*rod.Page, error) {
	lp.log.Debug(page.MustInfo().URL)

	//
	targetUrl, err := url.Parse(action.navigate)
	if err != nil {
		return nil, err
	}
	currentUrl, err := url.Parse(page.MustInfo().URL)
	if err != nil {
		return nil, err
	}
	lp.log.Debugf("targetUrl: %v", targetUrl)
	lp.log.Debugf("currentUrl: %v", currentUrl)

	// set request hijack in background
	if len(action.CaptureRequestHeaders) > 0 {
		go func(lp *LoggedInPage, action *ViewAction) {
			defer lp.captureHeaders(action)
		}(lp, action)
	}
	// classic applications need to perform full page postbacks
	// more often than SPAs
	// There is a isSPA flag - but currently not used
	if strings.EqualFold(currentUrl.Path, targetUrl.Path) {
		waitNav := page.MustWaitNavigation()
		page.MustNavigate(action.navigate)
		waitNav()
	}

	page.MustWaitIdle()
	page.MustWaitLoad()

	action.message = fmt.Sprintf("successfully navigated to: %s", action.navigate)

	return page, nil
}

func (lp *LoggedInPage) captureHeaders(action *ViewAction) func() {
	router := lp.browser.HijackRequests()
	// defer router.MustStop()
	router.MustAdd(action.navigate, func(ctx *rod.Hijack) {
		_ = ctx.LoadResponse(http.DefaultClient, true)
		headerMap := make(map[string]any)
		for _, v := range action.CaptureRequestHeaders {
			for hkey, hval := range ctx.Request.Headers() {
				if strings.EqualFold(strings.ToLower(v), strings.ToLower(hkey)) {
					headerMap[v] = hval
				}
			}
		}
		action.capturedReqHeaders = headerMap
	})

	go router.Run()

	return router.MustStop
}

// ensureIframeLoaded returns an instnace of a pointer to a rod.Page
// which is a document tree inside an iframe
func (lp *LoggedInPage) ensureIframeLoaded(page *rod.Page, action *ViewAction) (*rod.Page, error) {
	// allow for extremely slow iframe and page loads
	// search the page for iframe element and then apply selector
	if _, err := page.Search(action.Iframe.Selector); err != nil {
		lp.log.Errorf("not found element with iframe: %v", err.Error())
	}

	iframe, err := determinActionElement(lp.log, page, Element{Selector: &action.Iframe.Selector})
	if err != nil {
		return nil, err
	}
	iframe.MustWaitLoad()

	action.message = fmt.Sprintf("%s\n%s", action.message, "will perform following actions inside an iframe")

	page = iframe.MustFrame()
	// page.MustWaitLoad()

	page.MustWait(fmt.Sprintf(`() => { 
		console.log("trying to look for elements in iframe page");
		try {
			return document.readyState === 'complete' && %s;
		} catch (ex) {
			console.log("failed wait iframe eval", ex.message)
			return false
		}
	}`, action.Iframe.WaitEval))

	return page, nil
}

// handleActionError returns a skip error and error depending on config set up
func (p *LoggedInPage) handleActionError(page *rod.Page, a *ElementAction, err UIStrategyError) (bool, error) {

	if len(err.errorMap) > 0 && p.config.ContinueOnError {
		p.log.Debugf("action: %#v, errored with %v", a, err)
		p.log.Debugf("continue on error...")
		a.errored = true
		a.screenshot = p.captureAndSave(page)
		return true, nil
	}
	if len(err.errorMap) > 0 {
		return false, &err
	}
	return false, nil
}

// performAction handles finding the element and any related actions on it
// i.e. click or input
func (p *LoggedInPage) performAction(page *rod.Page, a *ElementAction) UIStrategyError {
	rodElem, err := p.DetermineActionElement(page, a)
	a.errored = false
	a.screenshot = ""
	if rodElem == nil {
		a.errored = true
		a.screenshot = p.captureAndSave(page)
		p.errors.setError(a.Name, *a.Element.Selector, "element not found")
		return p.errors
	}
	if a.CaptureOutput {
		html, err := rodElem.HTML()
		if err != nil {
			p.errors.setError(a.Name, *a.Element.Selector, fmt.Sprintf("failed to capture the output: %v", err))
		}
		a.capturedOutput = append(a.capturedOutput, html)
	}
	if err != nil {
		p.log.Debugf("action: %s, errored with %+#v", a.Name, err)
		// extend screenshots here
		a.message = fmt.Sprintf("locating element with selector: %s, errored with %+#v", *a.Element.Selector, err)
		a.errored = true
		a.screenshot = p.captureAndSave(page)
		p.errors.setError(a.Name, *a.Element.Selector, err.Error())
	}
	a.message = fmt.Sprintf("found element: %s", *a.Element.Selector)
	if err := p.DetermineActionType(a, rodElem); err != nil {
		p.log.Debugf("action: %s, errored with %v", a.Name, err)
		a.message = fmt.Sprintf("performing action on element with selector: %s, errored with %+v", *a.Element.Selector, err)
		a.errored = true
		a.screenshot = p.captureAndSave(page)
		p.errors.setError(a.Name, *a.Element.Selector, err.Error())
	}
	p.page.WaitRequestIdle(5*time.Second, []string{p.config.BaseUrl}, []string{}, []proto.NetworkResourceType{})
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

	if elem.Selector == nil {
		return nil, fmt.Errorf("action must include selector")
	}
	page.MustWaitLoad()

	type searchElemFunc func(selector string) (rod.Elements, error)
	searchfuncs := []searchElemFunc{
		func(selector string) (rod.Elements, error) {
			return page.Elements(selector)
		},
		func(selector string) (rod.Elements, error) {
			return page.ElementsX(selector)
		},
		// TODO: add more types here e.g. regex
		// func(selector string) (bool, *rod.Element, error) {return page.HasR(selector)},
	}

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
	log.Debugf("not found element using any search method")
	return nil, fmt.Errorf("element not found by selector: %v", *elem.Selector)
}

// DetermineActionType returns the rod.Element with correct action
// either Click/Swipe or Input
// when Input is selected - ensure you have specified the input/clickable HTML element
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

	if action.Timeout > 0 {
		ctx := context.Background()
		cctx, cancel := context.WithTimeout(ctx, time.Duration(action.Timeout*int(time.Second)))
		defer cancel()
		elem = elem.Context(cctx)
	}

	if elem != nil {
		if action.Assert {
			// update report with step found
			// item found not performing action
			lp.log.Debug("only assert only returning early")
			
			return nil
		}
		if action.CaptureOutput {
			lp.log.Debug("capturing output only  returning early")
			return nil
		}
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
	// action hover
	// // if
	// if action.SkipOnErrorMessage != "" && {
	// 	lp.page.Race().Search()
	// }

	elem.MustWaitLoad() // when clicked we wait for a

	return nil
}

// captureAndSave wil store the captured image under the .report/captures/*.png
// it will swallow any errors and log them out
func (lp *LoggedInPage) captureAndSave(page *rod.Page) string {
	file := fmt.Sprintf(`.report/captures/%v.png`, time.Now().UnixNano())
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
