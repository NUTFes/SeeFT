import 'dart:async';
import 'package:mysql1/mysql1.dart';

import '../../interface/repository/external/database.dart';
import '../../config/config.dart';

class Mysql implements Database {
  final MySqlConnection conn;

  Mysql(this.conn);

  static Future<MySqlConnection> connect(Environment env) async {
    final ConnectionSettings settings = ConnectionSettings(
      host: env.database.host,
      user: env.database.user,
      password: env.database.password,
      db: env.database.name,
    );
    return await MySqlConnection.connect(settings);
  }

  @override
  Future<List<Map<String, dynamic>>> finds(ctx, sql) async {
    try {
      final Results res = await conn.query(sql);
      final List results = res.toList();
      final List<Map<String, dynamic>> list = [];

      if (results.isEmpty) {
        throw Exception("Bad state: No element");
      }

      for (var result in results) {
        list.add(result.fields);
      }

      return list;
    } catch (e) {
      throw Exception(e);
    }
  }

  @override
  Future<Map<String, dynamic>> find(ctx, sql) async {
    try {
      final Results result = await conn.query(sql);
      var data = result.first.fields;
      return data;
    } catch (e) {
      throw Exception(e);
    }
  }

  @override
  Future<Map<String, dynamic>> insert(ctx, sql) async {
    try {
      var result = await conn.query(sql);
      var data = result.first.fields;
      return data;
    } catch (e) {
      throw Exception(e);
    }
  }

  @override
  Future<Map<String, dynamic>> update(ctx, sql, getSQL) async {
    try {
      await conn.query(sql);
      var result = await conn.query(getSQL);
      var data = result.first.fields;
      return data;
    } catch (e) {
      throw Exception(e);
    }
  }
}
