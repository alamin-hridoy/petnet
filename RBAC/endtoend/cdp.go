package endtoend

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

var showBrowser = flag.Bool("show-browser", true, "show the browser window for chromedp")

func (ut util) login(ctx context.Context, url, email, pass string, withNavi bool) string {
	const (
		loginBtnSel = `//*[@id="contact-us"]/div/form/button`
		emailSel    = `//*[@id="email"]`
		passSel     = `//*[@id="password"]`
	)

	if withNavi {
		ut.t.Log("Navigating to homepage")
		if err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(loginBtnSel),
		); err != nil {
			ut.t.Fatalf("navigating to homepage: %v", err)
		}
	}

	ut.t.Log("Logging in")
	var code string
	if err := chromedp.Run(ctx,
		chromedp.SetValue(emailSel, email),
		chromedp.SetValue(passSel, pass),
		chromedp.Click(loginBtnSel),
		getCodeFromRedirect(ctx, ut.t, &code),
	); err != nil {
		ut.t.Fatalf("logging in: %v", err)
	}
	return code
}

func getCodeFromRedirect(ctx context.Context, t *testing.T, code *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(c context.Context) error {
			ctxWithTimeout, canc := context.WithTimeout(ctx, cdpTimeout)
			defer canc()
			lCtx, lCanc := context.WithCancel(ctxWithTimeout)
			defer lCanc()
			done := make(chan struct{})
			chromedp.ListenTarget(lCtx, func(ev interface{}) {
				if a, ok := ev.(*network.EventRequestWillBeSent); ok {
					u, err := url.Parse(a.DocumentURL)
					if err != nil {
						t.Fatal("parsing url: ", err)
					}
					q := u.Query()
					if c := q.Get("code"); c != "" {
						*code = c
						close(done)
						lCanc()
					}
				}
			})
			for range done {
				lCanc()
			}
			return nil
		}),
	}
}

func (ut util) inviteSignup(ctx context.Context, t *testing.T, siteURL, pass string, withNavi bool) {
	const (
		fnSel          = `/html/body/main/div/div/div/form/div[1]/input`
		lnSel          = `/html/body/main/div/div/div/form/div[2]/input`
		passSel        = `/html/body/main/div/div/div/form/div[3]/div/input`
		confirmPassSel = `/html/body/main/div/div/div/form/div[4]/div/input`
		signupBtnSel   = `/html/body/main/div/div/div/form/div[5]/button`
	)

	if withNavi {
		t.Log("Navigating to sign up page")
		if err := chromedp.Run(ctx,
			chromedp.Navigate(siteURL),
			chromedp.WaitVisible(signupBtnSel),
		); err != nil {
			t.Fatalf("navigating to homepage: %v", err)
		}
	}

	t.Log("Signing up")
	if err := chromedp.Run(ctx,
		chromedp.SetValue(fnSel, "first"),
		chromedp.SetValue(lnSel, "last"),
		chromedp.SetValue(passSel, pass),
		chromedp.SetValue(confirmPassSel, pass),
		chromedp.Click(signupBtnSel),
		waitPageLoad(ctx, t),
	); err != nil {
		t.Fatalf("signing up: %v", err)
	}
}

func waitPageLoad(ctx context.Context, t *testing.T) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(c context.Context) error {
			ctxWithTimeout, canc := context.WithTimeout(ctx, cdpTimeout)
			defer canc()
			lCtx, lCanc := context.WithCancel(ctxWithTimeout)
			defer lCanc()
			done := make(chan struct{})
			chromedp.ListenTarget(lCtx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					close(done)
					lCanc()
				}
			})
			for range done {
				lCanc()
			}
			return nil
		}),
	}
}

// testBrowser creates a new Chrome browser, returning its chromedp context. The
// context will be cancelled when the test finishes.
func testBrowser(t *testing.T, opts ...chromedp.ContextOption) context.Context {
	allocOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", !*showBrowser),
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), allocOpts...)

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		append([]chromedp.ContextOption{
			chromedp.WithLogf(t.Logf),
		}, opts...)...,
	)
	t.Cleanup(cancel)
	// Start the browser, to ensure that the caller receives a browser
	// that's ready and already running. This also allows the caller to set
	// up timeouts without affecting the lifetime of the browser process.
	// todo implement runWithScreenShot()
	if err := chromedp.Run(ctx); err != nil {
		t.Fatal(err)
	}
	chromedp.ListenTarget(ctx, func(ev interface{}) { eventListener(t, ev) })
	return ctx
}

func eventListener(t *testing.T, ev interface{}) {
	switch ev := ev.(type) {
	case *runtime.EventConsoleAPICalled:
		switch ev.Type {
		case runtime.APITypeError, runtime.APITypeWarning:
		default:
			// hide log/info/debug, as they are too verbose
			return
		}
		var b strings.Builder
		fmt.Fprintf(&b, "console.%s(", ev.Type)
		for i, arg := range ev.Args {
			if i > 0 {
				fmt.Fprintf(&b, ", ")
			}
			// Value is a json.RawMessage, so it's already readable
			fmt.Fprintf(&b, "%s", arg.Value)
		}
		fmt.Fprintf(&b, ")")
		t.Logf("%s", &b)
	}
}
