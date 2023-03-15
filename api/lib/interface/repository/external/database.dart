abstract class Database {
  Future<List<Map<String, dynamic>>> finds(ctx, String sql);
  Future<Map<String, dynamic>> find(ctx, String sql);
  Future<Map<String, dynamic>> insert(ctx, String sql);
  Future<Map<String, dynamic>> update(ctx, String sql, String getSQL);
}
