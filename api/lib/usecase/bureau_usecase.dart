import '../entity/entity.dart';
import './repository/repository.dart';

abstract class BureauUsecase {
  Future<List<Bureau>> getBureaus(ctx);
}

class BureauUsecaseImpl implements BureauUsecase {
  BureauRepository bureauRepository;

  BureauUsecaseImpl(this.bureauRepository);

  @override
  Future<List<Bureau>> getBureaus(ctx) async {
    List<Bureau> list = await bureauRepository.getBureaus(ctx);

    return list;
  }
}
