package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Task struct {
	Description string
	Completed   bool
	Aborted     bool
}

type Day struct {
	Date  string
	Tasks []Task
}

func main() {

	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage:\n\ttodoer [options] [file]")
		return
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage:\n\ttodoer [options] [file]")
		fmt.Println()
		fmt.Println("Example:\ntodoer ~/todo.md")
		return
	}

	// Get file file path from command-line arguments
	filePath := os.Args[1]
	fileInfo, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("No file found at specified path:", filePath)
		return
	} else if err != nil {
		fmt.Printf("Error accessing file: %v\n", err)
		return
	} else if fileInfo.IsDir() {
		fmt.Println("Specified path is a directory, not a file.")
		return
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Split todo list's tasks into completed/uncompleted
	tasks, err := getTasks(file)
	if err != nil {
		fmt.Printf("error getting tasks: %v\n", err)
		return
	}
	completedTasks, uncompletedTasks, err := splitTasks(tasks)
	if err != nil {
		fmt.Printf("error splitting tasks: %v\n", err)
		return
	}

	// Reset file's offset
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Printf("error seeking to offset: %v\n", err)
		return
	}

	// Parse tasks from days before today
	days, err := parseDays(file)
	if err != nil {
		fmt.Printf("error parsing days: %v\n", err)
		return
	}

	// Completion day
	completionDay := Day{
		Date:  chooseDay(),
		Tasks: completedTasks,
	}

	// Add completed tasks to the beginning of the days slice
	if len(completedTasks) > 0 {
		days = append([]Day{completionDay}, days...)
	}

	// Format the output
	output := formatTasks(uncompletedTasks) + "---\n" + formatDays(days)

	err = writeToFile(filePath, output)
	if err != nil {
		fmt.Printf("error writing to file %v\n", err)
		return
	}

	fmt.Printf("File %s updated successfully\n", filePath)
}

func isTask(line string) bool {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "- ") {
		return true
	} else if strings.HasPrefix(line, "~~") {
		return true
	}

	return false
}

func getTasks(file *os.File) ([]Task, error) {
	var tasks []Task
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !isTask(line) {
			break
		}
		task, err := parseTask(line)
		if err != nil {
			fmt.Printf("error parsing task: %v\n", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func splitTasks(tasks []Task) ([]Task, []Task, error) {

	var completedTasks, uncompletedTasks []Task
	for _, task := range tasks {
		if task.Completed || task.Aborted {
			completedTasks = append(completedTasks, task)
		} else {
			uncompletedTasks = append(uncompletedTasks, task)
		}
	}

	return completedTasks, uncompletedTasks, nil
}

func parseTask(s string) (Task, error) {

	aborted := false
	// Remove aborted task notation
	if strings.HasPrefix(s, "~~- [") {
		s = strings.TrimPrefix(s, "~~")
		if strings.HasSuffix(s, "~~") {
			s = strings.TrimSuffix(s, "~~")
		}
		aborted = true
	}

	// Check if string has markdown task format
	if !strings.HasPrefix(s, "- [") {
		return Task{}, fmt.Errorf("Invalid task format: %s", s)
	}

	// Extract description and completed status from the string
	var completed bool
	switch s[3] {
	case ' ':
		completed = false
	case 'x':
		completed = true
	default:
		return Task{}, fmt.Errorf("Invalid task format: %s", s)
	}

	return Task{
		Description: strings.TrimSpace(s[5:]),
		Completed:   completed,
		Aborted:     aborted,
	}, nil
}

func parseDays(file *os.File) ([]Day, error) {

	// Read file line by line
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Parse days from the lines
	var days []Day
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			date := strings.TrimPrefix(line, "# ")
			tasks, err := extractTasks(lines[i+1:])
			if err != nil {
				fmt.Printf("error getting tasks from day \"%v\": %v\n", date, err)
				continue
			}
			days = append(days, Day{
				Date:  date,
				Tasks: tasks,
			})
		}
	}
	return days, nil
}

func extractTasks(lines []string) ([]Task, error) {
	var tasks []Task
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "---") {
			break
		}
		if !isTask(line) {
			continue
		}
		task, err := parseTask(line)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func formatTasks(tasks []Task) string {

	var sb strings.Builder
	for _, task := range tasks {
		completedMark := ' '
		if task.Completed {
			completedMark = 'x'
		}
		abortedMark := ""
		if task.Aborted {
			abortedMark = "~~"
		}
		sb.WriteString(fmt.Sprintf("%s- [%c] %s%s\n", abortedMark, completedMark, task.Description, abortedMark))
	}
	return sb.String()
}

func formatDays(days []Day) string {

	var formattedDays string
	for _, day := range days {
		formattedDays += fmt.Sprintf("\n# %s\n", day.Date) + formatTasks(day.Tasks)
	}

	return formattedDays
}

func writeToFile(filePath string, content string) error {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func chooseDay() string {
	date := time.Now()
	year, month, day := date.Date()
	comparisonTime := time.Date(year, month, day, 9, 0, 0, 0, date.Location())
	if date.Before(comparisonTime) {
		fmt.Println("It's before 6 PM today. Using yesterday's completed tasks.")
		date = date.AddDate(0, 0, -1)
	}
	return date.Format("Monday, 2 January")
}
