# TodoWrite

## Tool Description
Creates or updates a structured todo list used for task tracking.  
Supports full replacement mode or merge mode keyed by todo `id`.

Behavior notes:
- At most one todo should be `in_progress` at a time.
- Cancel tasks that become unnecessary.
- Common practice is setting the first task to `in_progress`.

## Tool Param
### merge
* description: If true, merge incoming todos into existing list by `id`; if false, replace the existing todo list entirely.
* type: boolean
* required: yes

### todos
* description: Array of todo items to create/update.
* type: array<object>
* required: yes
* minimum_items: 2

#### todos[].id
* description: Unique identifier for the todo item.
* type: string
* required: yes

#### todos[].content
* description: Task description/content.
* type: string
* required: yes

#### todos[].status
* description: Current status for the task.
* type: string
* required: yes
* enum: ["pending", "in_progress", "completed", "cancelled"]
