// Package Timer is timed contain, base on Time wheel algorithm.
package lib

import (
	"container/list"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Timer timed container.
type Timer struct {
	// slotNumber total number of slots.
	slotNumber int
	// slot every slot store a series of timed task.
	// [
	// 		slot1-list: task1-->task2-->task3, ...,
	// 		slot2-list: task1-->task2-->task3, ...
	// ]
	slot []*list.List
	// currentSlot which slot is currently pointed to.
	currentSlot int
	// tick heart time tick.
	tick *time.Ticker
	// si time to turn to the next slot.
	si time.Duration
	// stopChannel stopChannel all timed task.
	stopChannel chan struct{}
	// TaskID global task ID.
	taskID *TaskID
	// taskInfo store all task information, include the position in slot and the key of task.
	taskInfo map[TaskID]int
}

type task struct {
	// timeInterval after the timeInterval to run.
	timeInterval time.Duration
	// persistent whether is persistent.
	persistent bool
	// ring indicates in which circle to run.
	ring int
	// args the argument of callback.
	args interface{}
	// callback timed task.
	callback func(...interface{})
	// ID current task ID.
	ID TaskID
	// ctx cancel the task.
	ctx context.Context
}

// TaskID represent a task ID which used to delete.
type TaskID int

// NewTimer return a timer contain.
func NewTimer(slotNumber int, si time.Duration) *Timer {
	t := &Timer{}
	t.slotNumber = slotNumber
	t.si = si
	t.tick = time.NewTicker(si)
	t.stopChannel = make(chan struct{})
	t.taskID = new(TaskID)
	t.taskInfo = make(map[TaskID]int)
	for i := 0; i < t.slotNumber; i++ {
		t.slot = append(t.slot, list.New())
	}

	return t
}

// Start start to run.
func (t *Timer) Start() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGALRM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		// handle.
		case <-t.tick.C:
			t.run(ctx)
			if t.currentSlot == (t.slotNumber - 1) {
				t.currentSlot = 0
			} else {
				t.currentSlot++
			}
		// stopChannel task.
		case <-t.stopChannel:
			t.tick.Stop()
			return
		// sigalrm signal.
		case <-ch:
		}
	}
}

// Stop stopChannel task.
func (t *Timer) Stop() {
	t.stopChannel <- struct{}{}
}

// Add add a timed task.
func (t *Timer) Add(timeInterval time.Duration, fn func(v ...interface{}), args interface{}, persistent bool) TaskID {
	ring, position := t.getNewTaskPositionRing(timeInterval)
	*t.taskID++

	task := &task{
		timeInterval: timeInterval,
		persistent:   persistent,
		args:         args,
		callback:     fn,
		ID:           *t.taskID,
		ring:         ring,
	}

	t.slot[position].PushBack(task)

	// store task information.
	t.taskInfo[*t.taskID] = position

	return task.ID
}

// taskPersistent when a task need persistent to run, the task may be set in a new slot position,
// in order to delete this task, the new task position information should be update in taskInfo.
func (t *Timer) taskPersistent(timeInterval time.Duration, fn func(...interface{}), args interface{}, persistent bool, oldTaskID TaskID) {
	ring, position := t.getNewTaskPositionRing(timeInterval)
	newTask := &task{
		timeInterval: timeInterval,
		persistent:   persistent,
		args:         args,
		callback:     fn,
		ID:           oldTaskID,
		ring:         ring,
	}
	t.slot[position].PushBack(newTask)

	oldPosition, ok := t.taskInfo[oldTaskID]

	// update the new position of the task.
	t.taskInfo[oldTaskID] = position

	if ok {
		// delete the old task in slot.
		for e := t.slot[oldPosition].Front(); e != nil; e = e.Next() {
			if e.Value.(*task).ID == oldTaskID {
				t.slot[oldPosition].Remove(e)
				break
			}
		}
	}
}

// getNewTaskPositionRing return the position of new task locate in slot and which ring to run.
func (t *Timer) getNewTaskPositionRing(timeInterval time.Duration) (ring, position int) {
	if timeInterval <= 0 {
		panic("bad timeInterval")
	}

	siToSecond := t.si.Seconds()
	timeIntervalSecond := timeInterval.Seconds()
	ring = int(timeIntervalSecond/siToSecond) / t.slotNumber
	position = (t.currentSlot + int(timeIntervalSecond/siToSecond)) % t.slotNumber
	return
}

// Del delete a timed task.
func (t *Timer) Del(id TaskID) bool {
	if id == 0 {
		return false
	}

	position, ok := t.taskInfo[id]
	if !ok {
		return false
	}

	var tt *task
	for e := t.slot[position].Front(); e != nil; e = e.Next() {
		tt = e.Value.(*task)
		if tt.ID == id {
			t.slot[position].Remove(e)
			delete(t.taskInfo, id)
			return true
		}
	}
	return false
}

// run run timed task.
func (t *Timer) run(ctx context.Context) {
	if t.slot[t.currentSlot].Len() > 0 {
		for e := t.slot[t.currentSlot].Front(); e != nil; e = e.Next() {
			currentTask := e.Value.(*task)
			if currentTask.ring > 0 {
				currentTask.ring--
				continue
			}
			go func(currentTask *task) {
				for {
					select {
					default:
						currentTask.callback(currentTask.args)
						return
					case <-ctx.Done():
						return
					}
				}
			}(currentTask)

			if currentTask.persistent {
				t.taskPersistent(currentTask.timeInterval, currentTask.callback, currentTask.args, currentTask.persistent, currentTask.ID)
			} else {
				t.slot[t.currentSlot].Remove(e)
				delete(t.taskInfo, currentTask.ID)
			}
		}
	}
}
