import '../entity/entity.dart';
import './repository/repository.dart';

abstract class TimeUsecase {
  Future<List<Time>> getTimes(ctx);
}

class TimeUsecaseImpl implements TimeUsecase {
  TimeRepository timeRepository;

  TimeUsecaseImpl(this.timeRepository);

  @override
  Future<List<Time>> getTimes(ctx) async {
    List<Time> list = await timeRepository.getTimes(ctx);

    return list;
  }
}
