// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package internal

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Prediction struct {
	date    time.Time
	cutoff  int
	invites int
}

func Predict(crs *CRS, score int) (string, time.Time) {
	var userITA time.Time
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
				continue
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
		}

		prediction = append(prediction, newPrediction)
		date = newDate
	}
	minDiff := math.MaxInt32
	sBuilder := strings.Builder{}
	for _, items := range prediction {
		sBuilder.WriteString(
			fmt.Sprintf("%v\t\t\t%d\t\t\t\t%d\n", items.date.Format("01-02-2006"), items.invites, items.cutoff))
		diff := absDiff(score, items.cutoff)
		if diff < minDiff && items.cutoff <= score {
			userITA = items.date
			minDiff = diff
		}
	}
	return sBuilder.String(), userITA
}

func absDiff(a, b int) int {
	diff := a - b
	if diff < 0 {
		return diff * -1
	}
	return diff
}
