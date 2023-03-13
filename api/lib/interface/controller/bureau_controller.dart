import 'dart:convert';
import 'package:shelf/shelf.dart';

import '../../config/config.dart';
import '../../entity/entity.dart';
import '../../usecase/usecase.dart';

class BureauController {
  final StatusResponse statusResponse;
  final BureauUsecase bureauUsecase;

  BureauController(
    this.statusResponse,
    this.bureauUsecase,
  );

  getBureaus(Request request) async {
    try {
      List<Bureau> bureaus = await bureauUsecase.getBureaus(request.context);
      return statusResponse.responseOK(jsonEncode(bureaus));
    } catch (e) {
      Log.severe('BureauController.getBureaus: ${e.toString()}');
      var json = jsonEncode({'message': e.toString()});
      return statusResponse.responseBadRequest(json);
    }
  }
}
