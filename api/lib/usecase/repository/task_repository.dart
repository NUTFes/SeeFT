import '../../entity/entity.dart';

abstract class TaskRepository {
  Future<List<Task>> getTasks(ctx);
  Future<TaskDetail> getTask(ctx, Shift req);
  Future<TaskDetail> getTaskByShift(ctx, Shift req);
  Future<int> countUserFromTask(ctx, Shift req);
}
