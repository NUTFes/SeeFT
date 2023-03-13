import 'package:shelf/shelf.dart';

class StatusResponse {
  final corsHeaders = {
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
    'Access-Control-Allow-Headers': 'Origin, Content-Type'
  };

  Response responseOK(dynamic body) {
    return Response(
      200,
      body: body,
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    );
  }

  Response responseBadRequest(dynamic body) {
    return Response(
      400,
      body: body,
      headers: {
        'Content-Type': 'application/json',
        ...corsHeaders,
      },
    );
  }

  Response responseNotFound(dynamic body) {
    return Response(
      404,
      body: body,
      headers: {
        ...corsHeaders,
      },
    );
  }

  Response responseServerError(dynamic body) {
    return Response(
      500,
      body: body,
      headers: {
        ...corsHeaders,
      },
    );
  }
}
