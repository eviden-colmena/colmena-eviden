/*
Copyright © 2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This work has been implemented within the context of COLMENA project.
*/

package model

import (
	"colmena/sla-management-svc/app/common"
	"colmena/sla-management-svc/app/common/cfg"
	"colmena/sla-management-svc/app/common/expressions"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/lithammer/shortuuid/v4"
)

/**
 * Transforms an SLA Model to a OutputSLA model
 */
func SLAModelToOutputSLA(qos SLA) (OutputSLA, error) {

	output_model := OutputSLA{
		ServiceId: qos.Name,
		SLAId:     qos.Id,
		Kpis: []OutputSLAKpi{
			{
				RoleId:          qos.Details.Guarantees[0].Name,
				Query:           qos.Details.Guarantees[0].Constraint,
				Value:           0, // TODO res, //result.LastValues,
				Level:           qos.Assessment.Level,
				Threshold:       "", //qos.Details.Guarantees[0].Query,
				TotalViolations: qos.Assessment.TotalViolations,
			},
		},
	}

	return output_model, nil
}

/**
 * Transforms a list of SLA Models to a list of OutputSLA models
 */
func SLAModelsToOutputSLAs(l SLAs) ([]OutputSLA, error) {
	var lout []OutputSLA

	for _, m := range l {
		mout, err := SLAModelToOutputSLA(m)
		if err == nil {
			lout = append(lout, mout)
		} else {
			return lout, err
		}
	}

	return lout, nil
}

/**
 * Transforms the input to an SLA Model
 */
func InputSLAModelToSLAModel(c *gin.Context) ([]SLA, error) {
	var input InputSLA
	var slas []SLA

	err := c.ShouldBindJSON(&input)
	if err != nil {
		return slas, err
	}

	// InputSLA ==> SLA(s) managed by the app
	// KPIs
	if len(input.Kpis) > 0 {
		slas1 := listToSLAModel(input, "", input.Kpis)
		if len(slas1) > 0 {
			slas = append(slas, slas1...)
		}
	}

	// Roles definition
	if len(input.Roles) > 0 {
		for _, r := range input.Roles {
			if len(r.Kpis) > 0 {
				slas2 := listToSLAModel(input, r.Id, r.Kpis)
				if len(slas2) > 0 {
					slas = append(slas, slas2...)
				}
			}
		}
	}

	return slas, nil
}

// listToSLAModel
func listToSLAModel(input InputSLA, roleId string, l []InputSLARoleKPI) []SLA {
	var slas []SLA

	x := common.GetIntEnv(cfg.ASSESSMENT_X, DEFAULT_ASSESSMENT_X)
	y := common.GetIntEnv(cfg.ASSESSMENT_Y, DEFAULT_ASSESSMENT_Y)
	z := common.GetIntEnv(cfg.ASSESSMENT_Z, DEFAULT_ASSESSMENT_Z)

	for _, kpi := range l {
		uid := uuid.New()
		sla := SLA{}

		sla.Name = input.ServiceId.Value
		sla.Id = input.ServiceId.Value + "-" + uid

		// assessment
		sla.Assessment.TotalExecutions = 0
		sla.Assessment.TotalViolations = 0
		sla.Assessment.X = x
		sla.Assessment.XCounter = 0
		sla.Assessment.Y = y
		sla.Assessment.YCounter = 0
		sla.Assessment.Z = z
		sla.Assessment.ZCounter = 0
		sla.Assessment.Level = ASSESSMENT_LEVEL_UNKNOWN // Broken, Critical, Met, Desired, Unstable, Unknown

		// constraint expression
		expr, err := expressions.CheckAndParseConstraint(kpi.Query)
		if err != nil {
			expr = kpi.Query
		}

		// guarantees
		sla.Details.Guarantees = make([]Guarantee, 1) // TODO for each KPI => 1 Guarantee
		sla.Details.Guarantees[0].Name = roleId
		sla.Details.Guarantees[0].Constraint = kpi.Query
		sla.Details.Guarantees[0].Query = expr
		sla.Details.Guarantees[0].Scope = strings.ReplaceAll(kpi.Scope, " ", "")
		sla.Details.Guarantees[0].ScopeTemplate = strings.ReplaceAll(kpi.Scope, " ", "")

		if len(sla.Details.Guarantees[0].Constraint) > 0 && len(sla.Details.Guarantees[0].Scope) > 0 {
			sla.State = PAUSED
		} else if len(sla.Details.Guarantees[0].Constraint) > 0 {
			sla.State = STARTED
		} else {
			sla.State = INVALID
		}

		slas = append(slas, sla)
	}

	return slas
}
