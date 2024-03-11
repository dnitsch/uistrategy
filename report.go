package uistrategy

import (
	"encoding/json"
)

func (web *Web) BuildReport(allActions []*ViewAction) ViewReport {
	vr := make(ViewReport)

	if !web.config.WriteReport {
		return vr
	}

	for _, action := range allActions {
		actions := make(ActionsReport)
		vr[action.Name] = ViewReportItem{
			Message:                   action.message,
			CapturedHeaderRequestKeys: action.capturedReqHeaders,
			Actions:                   actions,
		}
		for _, ap := range action.ElementActions {
			vr[action.Name].Actions[ap.Name] = ActionReportItem{
				Message:    ap.message,
				Screenshot: ap.screenshot,
				Errored:    ap.errored,
				Output:     ap.capturedOutput,
			}
		}
	}

	if err := web.flushReport(vr); err != nil {
		web.log.Errorf("failed to write report with: %s", err)
	}
	return vr
}

func (web *Web) flushReport(report ViewReport) error {
	b, err := json.Marshal(report)
	if err != nil {
		return err
	}
	if _, err := web.output.Write(b); err != nil {
		web.log.Errorf("failed to write report: %v", err)
		return err
	}
	return nil
}
