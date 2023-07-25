package uistrategy

import (
	"encoding/json"
)

func (web *Web) BuildReport(allActions []*ViewAction) ViewReport {

	vrs := make(ViewReport)
	for _, v := range allActions {
		actions := make(ActionsReport)
		vrs[v.Name] = ViewReportItem{
			Message:                   v.message,
			CapturedHeaderRequestKeys: v.capturedReqHeaders,
			Actions:                   actions,
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
	if web.config.WriteReport {
		web.flushReport(vrs)
	}
	return vrs
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
