import 'package:shelf_router/shelf_router.dart';
import '../../../interface/controller/controller.dart';

class TaskApi {
  final TaskController controller;

  TaskApi(this.controller);

  Router get router {
    final router = Router();

    router.get('/', controller.getTasks);
    router.get('/shifts/<id>', controller.getTaskByShift);

    return router;
  }
}
