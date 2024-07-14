# Todoer

`Todoer` is a command-line tool written in Go that processes a markdown file containing a to-do list. It separates completed and uncompleted tasks, and moves completed tasks to a specific date section in the file. This helps in organizing tasks based on their completion date.

## Features

- Parses a markdown file for tasks.
- Separates completed and uncompleted tasks.
- Moves completed tasks to a specific date section in the file.
- Configurable to move tasks to either the current day or the previous day based on the current time.

## Usage

```
todoer [options] [file]
```

### Options

- `-h`, `--help`: Show the help message.

### Example

```
todoer tasks.md
```

## Task Format

Tasks in the markdown file should follow this format:

```
- [ ] Task description
- [x] Completed task description
```

## Installation

1. Clone the repository:

    ```
    git clone https://github.com/danielsrojo/todoer.git
    ```

2. Change to the project directory:

    ```
    cd todoer
    ```

3. Build the project:

    ```
    go build -o todoer
    ```

4. Move the executable to a directory in your `PATH`:

    ```
    mv todoer /usr/local/bin/
    ```

## How It Works

1. The tool reads the specified markdown file.
2. It parses tasks, separating them into completed and uncompleted tasks.
3. Completed tasks are moved to a section labeled with the current date or the previous day if run before a specific time (9 AM).
4. The updated content is written back to the file.

## Example File

### Before

```
# ToDo List

- [ ] Task 1
- [ ] Task 2
- [x] Task 3
- [x] Task 4

# Monday, 1 January
- [x] Task 0
```

### After

```
# ToDo List

- [ ] Task 1
- [ ] Task 2

# Tuesday, 2 January
- [x] Task 3
- [x] Task 4

# Monday, 1 January
- [x] Task 0
```
