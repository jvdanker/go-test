package main

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Run(ctxt, screenshot("https://www.govt.nz/"))
	if err != nil {
		log.Fatal(err)
	}

	var res []byte
	af := chromedp.ActionFunc(func(ctxt context.Context, h cdp.Executor) error {
		var err error

		root, err := dom.GetDocument().Do(ctxt, h)
		body, err := dom.QuerySelector(root.NodeID, "body").Do(ctxt, h)
		bm, err := dom.GetBoxModel().WithNodeID(body).Do(ctxt, h)
		emulation.SetDeviceMetricsOverride(1400, bm.Height, 1, false).Do(ctxt, h)
		//emulation.SetVisibleSize

		res, err = page.CaptureScreenshot().WithFromSurface(true).Do(ctxt, h)

		return err
	})
	err = c.Run(ctxt, af)

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("image.png", res, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func screenshot(urlstr string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		//chromedp.Sleep(2 * time.Second),
		chromedp.WaitVisible(".content", chromedp.ByQuery),
		//chromedp.WaitNotVisible(`div.v-middle > div.la-ball-clip-rotate`, chromedp.ByQuery),
		//chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByQuery),
	}
}
