package ceres

import (
	"CERES/src/utils"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type ExecutionResult struct {
	answer   string
	passOver map[string]any
}

type Executable interface {
	Execute(map[string]any) (ExecutionResult, error)
	load(string) error
	save(string) error
}

type ExecutionChain []Executable

func (ec *ExecutionChain) Execute(m map[string]any) (ExecutionResult, error) {
	if ec == nil {
		panic("Cannot execute nil Execution Chain")
	}
	var previous = ExecutionResult{answer: "", passOver: nil}
	var err error = nil

	for i := range *ec {
		previous, err = (*ec)[i].Execute(utils.Merge(m, previous.passOver))
		if err != nil {
			return previous, err
		}
	}
	return previous, nil
}

func (ec *ExecutionChain) load(Fname string) error {
	if ec == nil {
		panic("Cannot load into nil Execution Chain")
	}
	return nil
}

func (ec *ExecutionChain) save(Fname string) error {
	if ec == nil {
		panic("Cannot save as nil Execution Chain")
	}

	f, err := os.Open(Fname)
	if err != nil {
		return err
	}

	var fileLines []string = make([]string, 0, 8)
	for {
		bs, err := bufio.NewReader(f).ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		fileLines = append(fileLines, string(bs))
	}

	*(ec) = ExecutionChain(make(ExecutionChain, len(fileLines)))
	for i := range *ec {
		fields := strings.Split(fileLines[i], "|")
		if len(fields) != 2 {
			return fmt.Errorf("cannot parse execution chain at  %s : Line %v is not parseable (\"%s\" has wrong number of fields %v != 2)",
				Fname, i, fileLines[i], len(fields))
		}
		(*ec)[i], err = newExecutable(fields[0])
		if err != nil {
			return nil
		}
		err = (*ec)[i].load(fields[1])
		if err != nil {
			return nil
		}
	}

	return nil
}

func newExecutable(dtype string) (Executable, error) {
	switch dtype {
	case "ExecutableChain", "Chain":
		return new(ExecutionChain), nil
	}

	return nil, fmt.Errorf("cannot process dtype \"%s\"", dtype)
}
