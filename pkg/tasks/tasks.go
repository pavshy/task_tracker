package tasks

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/pavshy/task_tracker/pkg/history"
)

var (
	TodayTasks = make(Tasks) // map[taskName]task
)

type Tasks map[string]Task

type Task struct {
	Name  string
	Spent time.Duration
}

func Listen() error {
	histContent, err := history.Load()
	if err != nil {
		fmt.Printf("error loading history: %v\n", err)
	}
	TodayTasks = ParseHistory(histContent)

	fmt.Println("Enter time and task")
	for {
		in := bufio.NewReader(os.Stdin)

		ln, err := in.ReadString('\n')
		if err != nil {
			return fmt.Errorf("scanline error: %w", err)
		}
		ln = strings.TrimSpace(ln)
		switch ln {
		case "q", "quit", "exit":
			return fmt.Errorf("exiting")
		case "":
			cls()
			report, err := FormReport(TodayTasks)
			if err != nil {
				return fmt.Errorf("form reprot error: %w", err)
			}
			fmt.Println(report)
		default:
			spltLn := strings.SplitN(ln, " ", 2)
			if len(spltLn) != 2 {
				fmt.Println("wrong number of arguments in:", ln)
				break
			}
			spentTime := spltLn[0]
			taskName := spltLn[1]
			spentT, err := time.ParseDuration(spentTime)
			if err != nil {
				fmt.Println("error parsing duration:", err)
				break
			}
			oldTask, ok := TodayTasks[taskName]
			if !ok {
				TodayTasks[taskName] = Task{
					Name:  taskName,
					Spent: spentT,
				}
			} else {
				TodayTasks[taskName] = Task{
					Name:  oldTask.Name,
					Spent: oldTask.Spent + spentT,
				}
			}
			report, err := FormReport(TodayTasks)
			if err != nil {
				return fmt.Errorf("error forming report: %w", err)
			}
			err = history.Save(report)
			if err != nil {
				return fmt.Errorf("error saving report: %w", err)
			}
			cls()
			fmt.Printf(report)
		}
	}
}

func FormReport(tasks Tasks) (string, error) {
	buf := new(bytes.Buffer)
	_, err := fmt.Fprintf(buf, "Today:\n")
	if err != nil {
		return "", fmt.Errorf("cannot write today's history to file: %w", err)
	}
	if len(tasks) == 0 {
		_, err := fmt.Fprintf(buf, "nothing today\n")
		if err != nil {
			return "", fmt.Errorf("cannot write today's history to file: %w", err)
		}
	}
	for _, task := range tasks {
		var timeStr string
		hrs := int(task.Spent.Hours())
		mins := int(task.Spent.Minutes()) % 60
		if hrs > 0 {
			timeStr = fmt.Sprintf("%dч%dм", hrs, mins)
		} else {
			timeStr = fmt.Sprintf("%dм", mins)
		}
		_, err := fmt.Fprintf(buf, "%s %s\n", timeStr, task.Name)
		if err != nil {
			return "", fmt.Errorf("cannot write today's history to file: %w", err)
		}
	}
	return buf.String(), nil
}

func ParseHistory(hist string) Tasks {
	var tasks = make(Tasks)
linesLoop:
	for _, line := range strings.Split(hist, "\n") {
		line = strings.TrimSpace(line)
		split := strings.SplitN(line, " ", 2)
		if len(split) != 2 {
			continue
		}
		spentStr := split[0]
		var hrs, mins int
		spentSplit := strings.Split(spentStr, "ч")
		switch len(spentSplit) {
		case 2:
			var err error
			hrs, err = strconv.Atoi(strings.TrimSpace(spentSplit[0]))
			if err != nil {
				continue
			}
			mins, err = strconv.Atoi(strings.Trim(spentSplit[1], "\n\t м"))
			if err != nil {
				continue
			}
		case 1:
			var err error
			mins, err = strconv.Atoi(strings.Trim(spentSplit[0], "\n\t м"))
			if err != nil {
				continue
			}
		default:
			continue linesLoop
		}
		spent := time.Duration(hrs)*time.Hour + time.Duration(mins)*time.Minute
		taskName := split[1]
		tasks[taskName] = Task{
			Name:  taskName,
			Spent: spent,
		}
	}
	return tasks
}

func cls() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
