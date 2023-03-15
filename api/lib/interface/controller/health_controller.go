import 'package:shelf/shelf.dart';

class HealthController {
  getHealth(Request request) {
    return Response.ok('Get request');
  }
}
