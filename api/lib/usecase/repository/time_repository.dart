import '../../entity/entity.dart';

abstract class TimeRepository {
  Future<List<Time>> getTimes(ctx);
}
