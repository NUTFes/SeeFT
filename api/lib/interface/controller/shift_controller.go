import 'dart:convert';
import 'package:shelf/shelf.dart';

import '../../config/config.dart';
import '../../entity/entity.dart';
import '../../usecase/usecase.dart';

class ShiftController {
  final StatusResponse statusResponse;
  final ShiftUsecase shiftUsecase;

  ShiftController(
    this.statusResponse,
    this.shiftUsecase,
  );

  Future<Response> getShiftsByUser(Request request, String id) async {
    try {
      final req = User(id: int.parse(id));
      final res = await shiftUsecase.getShiftsByUser(request.context, req);

      return statusResponse.responseOK(jsonEncode(res));
    } catch (e) {
      Log.severe('shiftContoller.getShiftsByUser: ${e.toString()}');
      var json = jsonEncode({'message': e.toString()});
      return statusResponse.responseBadRequest(json);
    }
  }

  Future<Response> getShiftsByUserAndDateAndWeather(
      Request request, String userId, String dateId, String weatherId) async {
    try {
      final req = Shift(
          user: User(id: int.parse(userId)),
          date: Date(id: int.parse(dateId)),
          weather: Weather(id: int.parse(weatherId)));
      final res = await shiftUsecase.getShiftsByUserAndDateAndWeather(request.context, req);
      return statusResponse.responseOK(jsonEncode(res));
    } catch (e) {
      Log.severe('shiftController.getShiftsByUserAndDateAndWeather: ${e.toString()}');
      var json = jsonEncode({'message': e.toString()});
      return statusResponse.responseBadRequest(json);
    }
  }

  Future<Response> getShiftsByYearAndDateAndWeather(
      Request request, String yearId, String dateId, String weatherId) async {
    try {
      final req = Shift(
          year: Year(id: int.parse(yearId)),
          date: Date(id: int.parse(dateId)),
          weather: Weather(id: int.parse(weatherId)));
      final res = await shiftUsecase.getShiftsByYearAndDateAndWeather(request.context, req);
      return statusResponse.responseOK(jsonEncode(res));
    } catch (e) {
      Log.severe('shiftController.getShiftsByYearAndDateAndWeather: ${e.toString()}');
      var json = jsonEncode({'message': e.toString()});
      return statusResponse.responseBadRequest(json);
    }
  }
}
