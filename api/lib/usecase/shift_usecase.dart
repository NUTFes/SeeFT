import '../entity/entity.dart';
import './repository/repository.dart';

abstract class ShiftUsecase {
  Future<List<Shift>> getShiftsByUser(ctx, User req);
  Future<List<Shift>> getShiftsByUserAndDateAndWeather(ctx, Shift req);
  Future<List<Shift>> getShiftsByYearAndDateAndWeather(ctx, Shift req);
}

class ShiftUsecaseImpl implements ShiftUsecase {
  ShiftRepository shiftRepository;

  ShiftUsecaseImpl(this.shiftRepository);

  @override
  Future<List<Shift>> getShiftsByUser(ctx, User req) async {
    List<Shift> list = await shiftRepository.getShiftsByUser(ctx, req);
    return list;
  }

  @override
  Future<List<Shift>> getShiftsByUserAndDateAndWeather(ctx, Shift req) async {
    List<Shift> list = await shiftRepository.getShiftsByUserAndDateAndWeather(ctx, req);
    return list;
  }

  @override
  Future<List<Shift>> getShiftsByYearAndDateAndWeather(ctx, Shift req) async {
    List<Shift> list = await shiftRepository.getShiftsByYearAndDateAndWeather(ctx, req);
    return list;
  }
}
