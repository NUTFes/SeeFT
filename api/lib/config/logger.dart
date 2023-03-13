import 'package:logging/logging.dart';

final logger = Logger('SeeFT-API-Logger');

class Log {
  static setupDev() {
    Logger.root.level = Level.CONFIG;
    Logger.root.onRecord.listen((record) {
      print('${record.level.name}: ${record.time}: ${record.message}');
    });
  }

  static setupProd() {
    Logger.root.level = Level.WARNING;
    Logger.root.onRecord.listen((record) {
      print('${record.level.name}: ${record.time}: ${record.message}');
    });
  }

  static severe(dynamic text) {
    logger.severe(text);
  }

  static warning(dynamic text) {
    logger.warning(text);
  }

  static info(dynamic text) {
    logger.info(text);
  }

  static config(dynamic text) {
    logger.config(text);
  }
}
