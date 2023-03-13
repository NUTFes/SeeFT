import '../../entity/entity.dart';

abstract class BureauRepository {
  Future<List<Bureau>> getBureaus(ctx);
}
