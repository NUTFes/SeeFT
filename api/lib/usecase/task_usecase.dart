import '../entity/entity.dart';
import './repository/repository.dart';

abstract class TaskUsecase {
  Future<List<Task>> getTasks(ctx);
  Future<TaskDetail> getTaskByShift(ctx, Shift req);
}

class TaskUsecaseImpl implements TaskUsecase {
  TaskRepository taskRepository;
  ShiftRepository shiftRepository;

  TaskUsecaseImpl(this.taskRepository, this.shiftRepository);

  @override
  Future<List<Task>> getTasks(ctx) async {
    List<Task> list = await taskRepository.getTasks(ctx);

    return list;
  }

  @override
  Future<TaskDetail> getTaskByShift(ctx, Shift req) async {
    Shift shift = await shiftRepository.getShift(ctx, req);
    int count = await taskRepository.countUserFromTask(ctx, shift);
    TaskDetail task;
    if (count > 1) {
      task = await taskRepository.getTaskByShift(ctx, shift);
    } else if (count == 1) {
      task = await taskRepository.getTask(ctx, shift);
    } else {
      throw Exception('not find.');
    }
    return task;
  }
}
