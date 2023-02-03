package uistrategy

import (
	"encoding/json"
	"os"
)

func (web *Web) buildReport(allActions []*ViewAction) {

	vrs := make(ViewReport)
	for _, v := range allActions {
		actions := make(ActionsReport)
		vrs[v.Name] = ViewReportItem{
			Message: v.message,
			Actions: actions,
		}
		for _, ap := range v.ElementActions {
			vrs[v.Name].Actions[ap.Name] = ActionReportItem{
				Message:    ap.message,
				Screenshot: ap.screenshot,
				Errored:    ap.errored,
				Output:     ap.capturedOutput,
			}
		}
	}

	web.flushReport(vrs)
}

func (web *Web) flushReport(report ViewReport) error {
	file := `.report/report.json`

	w, err := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		web.log.Debugf("unable to get a writer: %v", err)
		return err
	}

	b, err := json.Marshal(report)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		web.log.Errorf("failed to write report: %v", err)
		return err
	}
	return nil
}
