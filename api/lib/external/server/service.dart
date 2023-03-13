import 'package:shelf_router/shelf_router.dart';
import 'package:shelf/shelf.dart';

import 'routes/route.dart';
import '../../interface/controller/controller.dart';

class Service {
  final HealthController healthController;
  final AuthController authController;
  final UserController userController;
  final TimeController timeController;
  final BureauController bureauController;
  final ShiftController shiftController;
  final TaskController taskController;

  Service(
    this.healthController,
    this.userController,
    this.timeController,
    this.bureauController,
    this.shiftController,
    this.authController,
    this.taskController,
  );

  Handler get handler {
    final router = Router();

    router.get('/', (Request request) {
      return Response.ok('Get request');
    });

    router.get('/healthz', healthController.getHealth);

    router.mount('/auth', AuthApi(authController).router);

    router.mount('/users', UserApi(userController).router);

    router.mount('/times', TimeApi(timeController).router);

    router.mount('/bureaus', BureauApi(bureauController).router);

    router.mount('/shifts', ShiftApi(shiftController).router);

    router.mount('/tasks', TaskApi(taskController).router);

    router.all('/<ignored|.*>', (Request request) {
      return Response.notFound('Page not found');
    });

    return router;
  }
}
