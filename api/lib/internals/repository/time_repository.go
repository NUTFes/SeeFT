package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type timeRepository struct {
	client db.Client
	crud   abstract.Crud
}

type TimeRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
}

func NewTimeRepository(c db.Client, ac abstract.Crud) TimeRepository {
	return &timeRepository{c, ac}
}

// 全件取得
func (b *timeRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM times"
	return b.crud.Read(c, query)
}

// 1件取得
func (b *timeRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM times WHERE id =" + id
	return b.crud.ReadByID(c, query)
}

// import '../../usecase/repository/time_repository.dart';
// import '../../entity/entity.dart';
// import './external/database.dart';

// class TimeRepositoryImpl implements TimeRepository {
//   Database database;

//   TimeRepositoryImpl(this.database);

//   @override
//   Future<List<Time>> getTimes(ctx) async {
//     String sql = '''
// SELECT * FROM times;
// ''';

//     List<Map<String, dynamic>> data = await database.finds(ctx, sql);
//     List<Time> list = [];

//     for (var d in data) {
//       Time time = Time(
//         id: d['id'],
//         time: d['time'],
//         createdAt: d['created_at'],
//         updatedAt: d['updated_at'],
//         deletedAt: d['deleted_at'],
//       );

//       if (time.isDeleted) {
//         continue;
//       }

//       list.add(time);
//     }

//     return list;
//   }
// }
