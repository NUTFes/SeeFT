package usecase

import (
  "context"

	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/pkg/errors"
)

type taskUseCase struct {
  rep rep.TaskRepository
}

type TaskUseCase interface {
  GetTasks(context.Context) ([]entity.Task, error)
  GetTaskByID(context.Context, string) (entity.Task, error)
  GetTasksByShift(context.Context, string) ([]entity.Task, error)
}

func NewTaskUseCase(rep rep.TaskRepository) TaskUseCase {
  return &taskUseCase{rep}
}

func (b *taskUseCase) GetTasks(c context.Context) ([]entity.Task, error) {
  task := entity.Task{}
  var tasks []entity.Task

  rows, err := b.rep.All(c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

  for rows.Next() {
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Place,
			&task.Url,
      &task.Superviser,
	  &task.Color,
      &task.Notes,
      &task.YearID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}
    
		tasks = append(tasks, task)
	}
  
	return tasks, nil
}

func (b *taskUseCase) GetTaskByID(c context.Context, id string) (entity.Task, error) {
  var task entity.Task
  row, err := b.rep.Find(c, id)
	err = row.Scan(
		&task.ID,
		&task.Name,
		&task.Place,
		&task.Url,
    &task.Superviser,
	&task.Color,
    &task.Notes,
    &task.YearID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return task, err
	}
		
	return task, nil
}

func (b *taskUseCase) GetTasksByShift(c context.Context, shift string) ([]entity.Task, error) {
  task := entity.Task{}
  var tasks []entity.Task

  rows, err := b.rep.Shift(c, shift)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

  for rows.Next() {
		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.Place,
			&task.Url,
      &task.Superviser,
	  &task.Color,
      &task.Notes,
      &task.YearID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// import '../entity/entity.dart';
// import './repository/repository.dart';

// abstract class TaskUsecase {
//   Future<List<Task>> getTasks(ctx);
//   Future<TaskDetail> getTaskByShift(ctx, Shift req);
// }

// class TaskUsecaseImpl implements TaskUsecase {
//   TaskRepository taskRepository;
//   ShiftRepository shiftRepository;

//   TaskUsecaseImpl(this.taskRepository, this.shiftRepository);

//   @override
//   Future<List<Task>> getTasks(ctx) async {
//     List<Task> list = await taskRepository.getTasks(ctx);

//     return list;
//   }

//   @override
//   Future<TaskDetail> getTaskByShift(ctx, Shift req) async {
//     Shift shift = await shiftRepository.getShift(ctx, req);
//     int count = await taskRepository.countUserFromTask(ctx, shift);
//     TaskDetail task;
//     if (count > 1) {
//       task = await taskRepository.getTaskByShift(ctx, shift);
//     } else if (count == 1) {
//       task = await taskRepository.getTask(ctx, shift);
//     } else {
//       throw Exception('not find.');
//     }
//     return task;
//   }
// }
