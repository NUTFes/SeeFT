import 'package:shelf_router/shelf_router.dart';

import '../../../interface/controller/controller.dart';

class AuthApi {
  final AuthController controller;

  AuthApi(this.controller);

  Router get router {
    final router = Router();

    router.get('/signin/<mail>', controller.signIn);
//    router.get('/<id>', controller.getBureau);

    return router;
  }
}
