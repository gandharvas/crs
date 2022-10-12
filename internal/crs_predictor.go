// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package internal

import (
	"time"
)

type Prediction struct {
	date    time.Time
	cutoff  int
	invites int
}

func Predict(crs *CRS, score int) ([]Prediction, time.Time) {
	date := crs.GetPreviousDrawDate()
	cutOff := crs.GetPreviousDrawCutoff()
	stepSize := crs.GetPreviousDrawStepsize()
	totalInvites := crs.GetPreviousTotalInvitesSent()
	var prediction []Prediction
	previousdraw := Prediction{}
	tempInvites := totalInvites
	distriIdx := 0
	for idx, distribution := range crs.backlog {
		if distribution.lowerBound == 451 && distribution.upperBound == 500 {
			continue
		}

		if distribution.upperBound < cutOff {
			break
		}
		previousdraw.cutoff = cutOff
		previousdraw.date = date
		previousdraw.invites = totalInvites
		if tempInvites > 0 {
			if distribution.candidates < tempInvites {
				tempInvites -= distribution.candidates
				crs.backlog[idx].candidates = 0
				distriIdx++
			} else {
				crs.backlog[idx].candidates -= tempInvites
				tempInvites = 0
			}
		} else {
			break
		}

	}
	prediction = append(prediction, previousdraw)

	for {
		if distriIdx > 13 {
			break
		}
		newPrediction := Prediction{}
		newDate := date.AddDate(0, 0, 14)
		totalInvites += stepSize

		newPrediction.date = newDate
		newPrediction.invites = totalInvites
		tempTotalInvites := totalInvites
		for tempTotalInvites > 0 {
			if crs.backlog[distriIdx].upperBound == 500 && crs.backlog[distriIdx].lowerBound == 451 {
				distriIdx++
			}
			distriBacklog := crs.backlog[distriIdx].candidates
			if distriBacklog < tempTotalInvites {
				tempTotalInvites -= distriBacklog
				crs.backlog[distriIdx].candidates = 0
				distriIdx++
			} else {
				distri := crs.backlog[distriIdx]
				candidatePerPoint := distri.candidates / (distri.upperBound - distri.lowerBound)
				totalPointsToReduce := tempTotalInvites / candidatePerPoint
				crs.backlog[distriIdx].candidates -= tempTotalInvites
				tempTotalInvites = 0
				newPrediction.cutoff = distri.upperBound - totalPointsToReduce
				break
			}

			if crs.backlog[distriIdx].upperBound >= score && crs.backlog[distriIdx].lowerBound <= score {
				return prediction, newPrediction.date
			}
		}

		prediction = append(prediction, newPrediction)
		date = newDate
	}
	return prediction, time.Now()
}
