import '../../entity/entity.dart';

abstract class ShiftRepository {
  Future<Shift> getShift(ctx, Shift req);
  Future<List<Shift>> getShiftsByUser(ctx, User req);
  Future<List<Shift>> getShiftsByUserAndDateAndWeather(ctx, Shift req);
  Future<List<Shift>> getShiftsByYearAndDateAndWeather(ctx, Shift req);
}
