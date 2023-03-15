package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/drivers/db"
	"github.com/NUTFes/SeeFT/api/lib/externals/repository/abstract"
)

type bureauRepository struct {
	client db.Client
	crud   abstract.Crud
}

type BureauRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	Create(context.Context, string) error
	Update(context.Context, string, string) error
	Destroy(context.Context, string) error
	FindLatestRecord(context.Context) (*sql.Row, error)
}

func NewBureauRepository(c db.Client, ac abstract.Crud) BureauRepository {
	return &bureauRepository{c, ac}
}

// 全件取得
func (b *bureauRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM bureaus"
	return b.crud.Read(c, query)
}

// 1件取得
func (b *bureauRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM bureaus WHERE id =" + id
	return b.crud.ReadByID(c, query)
}

// 作成
func (b *bureauRepository) Create(c context.Context, name string) error {
	query := "INSERT INTO bureaus (name) VALUES ('" + name + "')"
	return b.crud.UpdateDB(c, query)
}

// 編集
func (b *bureauRepository) Update(c context.Context, id string, name string) error {
	query := "UPDATE bureaus SET name = '" + name + "' WHERE id = " + id
	return b.crud.UpdateDB(c, query)
}

// 削除
func (b *bureauRepository) Destroy(c context.Context, id string) error {
	query := "DELETE FROM bureaus WHERE id =" + id
	return b.crud.UpdateDB(c, query)
}

// 最新のbureauを取得する
func (b *bureauRepository) FindLatestRecord(c context.Context) (*sql.Row, error) {
	query := `
		SELECT
			*
		FROM
			bureaus
		ORDER BY
			id
		DESC LIMIT 1
	`
	return b.crud.ReadByID(c, query)
}


// import '../../usecase/repository/bureau_repository.dart';
// import '../../entity/entity.dart';
// import './external/database.dart';

// class BureauRepositoryImpl implements BureauRepository {
//   Database database;

//   BureauRepositoryImpl(this.database);

//   @override
//   Future<List<Bureau>> getBureaus(ctx) async {
//     String sql = '''
// SELECT * FROM bureaus;
// ''';

//     List<Map<String, dynamic>> data = await database.finds(ctx, sql);
//     List<Bureau> list = [];

//     for (var d in data) {
//       Bureau bureau = Bureau(
//         id: d['id'],
//         bureau: d['bureau'],
//         color: Color(int.parse(d['color'], radix: 16)),
//         createdAt: d['created_at'],
//         updatedAt: d['updated_at'],
//         deletedAt: d['deleted_at'],
//       );

//       if (bureau.isDeleted) {
//         continue;
//       }

//       list.add(bureau);
//     }

//     return list;
//   }
// }
