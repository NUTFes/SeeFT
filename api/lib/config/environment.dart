import 'dart:io' show Platform;

class Environment {
  String _applicationEnv = '';
  late _Server _server;
  late _Database _database;
  Environment() {
    Map<String, String> env = Platform.environment;
    _server = _Server();
    _database = _Database();
    _applicationEnv = env['APPLICATION_ENV'] ?? 'develop';
  }

  get applicationEnv => _applicationEnv;
  get server => _server;
  get database => _database;
}

class _Server {
  String _address = '';
  String _port = '';

  _Server() {
    Map<String, String> env = Platform.environment;
    _address = env['SERVER_ADDRESS'] ?? '0.0.0.0';
    _port = env['SERVER_PORT'] ?? '3000';
  }

  get address => _address;
  get port => _port;
}

class _Database {
  String _host = '';
  String _user = '';
  String _password = '';
  String _name = '';

  _Database() {
    Map<String, String> env = Platform.environment;
    _host = env['DATABASE_HOST'] ?? '';
    _user = env['DATABASE_USER'] ?? '';
    _password = env['DATABASE_PASSWORD'] ?? '';
    _name = env['DATABASE_NAME'] ?? 'seeft';
  }

  get host => _host;
  get user => _user;
  get password => _password;
  get name => _name;
}
