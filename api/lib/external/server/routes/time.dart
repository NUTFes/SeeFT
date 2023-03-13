import 'package:shelf_router/shelf_router.dart';

import '../../../interface/controller/time_controller.dart';

class TimeApi {
  final TimeController controller;

  TimeApi(this.controller);

  Router get router {
    final router = Router();

    router.get('/', controller.getTimes);

    return router;
  }
}
