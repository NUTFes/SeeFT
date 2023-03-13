import 'package:shelf_router/shelf_router.dart';

import '../../../interface/controller/user_controller.dart';

class UserApi {
  final UserController controller;

  UserApi(this.controller);

  Router get router {
    final router = Router();

    router.get('/', controller.getUsers);
    router.get('/<id>', controller.getUser);
    router.post('/', controller.insertUser);
    router.put('/<id>', controller.updateUser);
    router.delete('/<id>', controller.deleteUser);

    return router;
  }
}
