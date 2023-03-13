import '../external/server/server.dart';
import '../external/server/service.dart';
import '../external/mysql/mysql.dart';
import '../interface/controller/controller.dart';
import '../interface/repository/repository.dart';
import '../usecase/usecase.dart';
import '../config/config.dart';

initializeServer(Environment env) async {
  final statusResponse = StatusResponse();
  final conn = await Mysql.connect(env);
  final database = Mysql(conn);

  final userRepository = UserRepositoryImpl(database);
  final timeRepositoty = TimeRepositoryImpl(database);
  final bureauRepository = BureauRepositoryImpl(database);
  final shiftRepository = ShiftRepositoryImpl(database);
  final taskRepository = TaskRepositoryImpl(database);

  final userUsecase = UserUsecaseImpl(userRepository);
  final timeUsecase = TimeUsecaseImpl(timeRepositoty);
  final bureauUsecase = BureauUsecaseImpl(bureauRepository);
  final shiftUsecase = ShiftUsecaseImpl(shiftRepository);
  final authUsecase = AuthUsecaseImpl(userRepository);
  final taskUsecase = TaskUsecaseImpl(taskRepository, shiftRepository);

  final healthController = HealthController();
  final userController = UserController(statusResponse, userUsecase);
  final timeController = TimeController(statusResponse, timeUsecase);
  final bureauController = BureauController(statusResponse, bureauUsecase);
  final shiftController = ShiftController(statusResponse, shiftUsecase);
  final authController = AuthController(statusResponse, authUsecase);
  final taskController = TaskController(statusResponse, taskUsecase);

  final service = Service(
    healthController,
    userController,
    timeController,
    bureauController,
    shiftController,
    authController,
    taskController,
  );
  final server = Server(service, env);

  return server;
}
