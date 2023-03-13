import 'dart:convert';
import 'package:shelf/shelf.dart';

import '../../config/config.dart';
import '../../entity/entity.dart';
import '../../usecase/usecase.dart';

class TimeController {
  final StatusResponse statusResponse;
  final TimeUsecase timeUsecase;

  TimeController(
    this.statusResponse,
    this.timeUsecase,
  );

  Future<Response> getTimes(Request request) async {
    try {
      List<Time> times = await timeUsecase.getTimes(request.context);
      return statusResponse.responseOK(jsonEncode(times));
    } catch (e) {
      Log.severe('timeContoller.getTimes: ${e.toString()}');
      var json = jsonEncode({'message': e.toString()});
      return statusResponse.responseBadRequest(json);
    }
  }
}
