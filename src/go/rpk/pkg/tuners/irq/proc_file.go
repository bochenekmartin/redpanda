// Copyright 2020 Vectorized, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package irq

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/redpanda-data/redpanda/src/go/rpk/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type ProcFile interface {
	GetIRQProcFileLinesMap() (map[int]string, error)
}

func NewProcFile(fs afero.Fs) ProcFile {
	return &procFile{
		fs: fs,
	}
}

type procFile struct {
	fs afero.Fs
}

func (procFile *procFile) GetIRQProcFileLinesMap() (map[int]string, error) {
	log.Debugf("Reading '/proc/interrupts' file...")
	lines, err := utils.ReadFileLines(procFile.fs, "/proc/interrupts")
	if err != nil {
		return nil, err
	}
	linesByIRQ := make(map[int]string)
	irqPattern := regexp.MustCompile("^\\s*\\d+:.*$")
	for _, line := range lines {
		if !irqPattern.MatchString(line) {
			continue
		}
		irq, err := strconv.Atoi(strings.TrimSpace(strings.Split(line, ":")[0]))
		if err != nil {
			return nil, err
		}
		linesByIRQ[irq] = line
	}
	for irq, line := range linesByIRQ {
		log.Tracef("IRQ -> /proc/interrupts %d - %s", irq, line)
	}
	return linesByIRQ, nil
}
