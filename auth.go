package uistrategy

import (
	"fmt"

	"github.com/go-rod/rod"
)

type Auth struct {
	Username        Element  `yaml:"username" json:"username"`
	Password        Element  `yaml:"password" json:"password"`
	ConfirmPassword Element  `yaml:"confirmPassword,omitempty" json:"confirmPassword,omitempty"`
	RequireConfirm  bool     `yaml:"requireConfirm,omitempty" json:"requireConfirm,omitempty"`
	Navigate        string   `yaml:"navigate" json:"navigate"`
	IdpManaged      bool     `yaml:"idpManaged" json:"idpManaged"`
	IdpSelector     *Element `yaml:"idpSelector,omitempty" json:"idpSelector,omitempty"`
	MfaSelector     *Element `yaml:"mfaSelector,omitempty" json:"mfaSelector,omitempty"`
	IdpUrl          string   `yaml:"idpUrl" json:"idpUrl"`
	Submit          Element  `yaml:"submit" json:"submit"`
}

// DoAuth performs the required Authentication
// in the browser and returns a authed Page
func (web *Web) DoAuth(auth *Auth) (*LoggedInPage, error) {

	if auth != nil {
		if auth.IdpManaged {
			return web.doIdpAuth(*auth)
		}
		return web.doLocalAuth(*auth)
	}
	page := web.browser.MustPage(web.config.BaseUrl).MustWaitLoad()

	lp := &LoggedInPage{web, page, UIStrategyError{}}
	return lp, nil
}

// doLocalAuth will drive a local username and password login form
func (web *Web) doLocalAuth(auth Auth) (*LoggedInPage, error) {

	page := web.browser.MustPage(web.config.BaseUrl + auth.Navigate).MustWaitLoad()
	lp := &LoggedInPage{web, page, UIStrategyError{}}

	web.log.Debug("begin auth")
	if err := web.sharedLoginForm(page, auth); err != nil {
		return nil, err
	}
	page.MustWaitLoad()
	web.log.Debug("end auth")
	return lp, nil
}

// doIdpAuth will drive either an SP or IdP initiated login
// SP initiated will be simpler as you can omit the idpUrl
// the flow will follow redirects
func (web *Web) doIdpAuth(auth Auth) (*LoggedInPage, error) {
	if auth.IdpSelector == nil {
		return nil, fmt.Errorf("idpSelector must be specified")
	}

	page := web.browser.MustPage(web.config.BaseUrl + auth.Navigate).MustWaitLoad()
	lp := &LoggedInPage{web, page, UIStrategyError{}}

	web.log.Debug("begin auth")

	idpSelect, err := determinActionElement(lp.log, page, *auth.IdpSelector)
	if err != nil {
		web.log.Errorf("unable to find IdpSelector field, by selector: %v, error: %v", auth.IdpSelector.Selector, err.Error())
		return nil, err
	}
	idpSelect.MustClick()
	page.MustWaitLoad()

	if err := web.sharedLoginForm(page, auth); err != nil {
		return nil, err
	}
	page.MustWaitLoad()
	waitNav := page.MustWaitNavigation()
	// if MFA needs to be triggered and when not done automatically
	if auth.MfaSelector != nil {
		mfaSelect, err := determinActionElement(lp.log, page, *auth.MfaSelector)
		if err != nil {
			web.log.Errorf("unable to find mfaSelector field, by selector: %v", *auth.MfaSelector.Selector)
			return nil, err
		}
		mfaSelect.MustClick()
		if err := rod.Try(func() { mfaSelect.WaitInvisible() }); err != nil {
			// if err := mfaSelect.WaitInvisible()
			lp.log.Debugf("error waiting on invisible: %v", err)
		}
	}
	waitNav()
	page.MustWaitLoad()
	web.log.Debug("end auth")
	return lp, nil
}

func (web *Web) sharedLoginForm(page *rod.Page, auth Auth) error {
	uname, err := determinActionElement(web.log, page, auth.Username)
	if err != nil {
		web.log.Errorf("unable to find username field, by selector: %v", *auth.Username.Selector)
		return err
	}
	uname.MustInput(*auth.Username.Value)
	passwd, err := determinActionElement(web.log, page, auth.Password)
	if err != nil {
		web.log.Errorf("unable to find password field, by selector: %v", *auth.Username.Selector)
		return err
	}
	passwd.MustInput(*auth.Password.Value)
	submit, err := determinActionElement(web.log, page, auth.Submit)
	if err != nil {
		web.log.Errorf("unable to find password field, by selector: %v", *auth.Username.Selector)
		return err
	}
	// in case the auth page is not the same as the app it will crash
	submit.MustClick()
	if err := submit.WaitInvisible(); err != nil {
		web.log.Debugf("error waiting on submit invisible: %v", err)
	}
	return nil
}
