import 'package:shelf_router/shelf_router.dart';

import '../../../interface/controller/bureau_controller.dart';

class BureauApi {
  final BureauController controller;

  BureauApi(this.controller);

  Router get router {
    final router = Router();

    router.get('/', controller.getBureaus);
//    router.get('/<id>', controller.getBureau);

    return router;
  }
}
