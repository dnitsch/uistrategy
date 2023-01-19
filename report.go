package uistrategy

import (
	"encoding/json"
	"os"
)

func (web *Web) buildReport(allActions []ViewAction) {

	vrs := []ViewReport{}
	for _, v := range allActions {
		vr := ViewReport{
			Name:    v.Name,
			Message: v.message,
		}
		for _, a := range v.ElementActions {
			va := ActionReport{
				Name:       a.Name,
				Message:    a.message,
				Screenshot: a.screenshot,
				Errored:    a.errored,
			}
			vr.Actions = append(vr.Actions, va)
		}
		vrs = append(vrs, vr)
	}

	web.flushReport(vrs)
}

func (web *Web) flushReport(report []ViewReport) error {
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
