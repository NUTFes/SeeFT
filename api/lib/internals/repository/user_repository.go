package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type userRepository struct {
	client   db.Client
	crud abstract.Crud
}

type UserRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	FindByStudentNumber(context.Context, string) (*sql.Row)
	Create(context.Context, string, string, string, string, string, string, string, string, string) error
	Update(context.Context, string, string, string, string, string, string, string, string, string, string) error
	Delete(context.Context, string) error
	FindNewRecord(context.Context) (*sql.Row, error)
}

func NewUserRepository(c db.Client, ac abstract.Crud) UserRepository {
	return &userRepository{c, ac}
}

// 全件取得
func (ur *userRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM users"
	return ur.crud.Read(c, query)
}

// 1件取得
func (ur *userRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM users WHERE id = " + id
	return ur.crud.ReadByID(c, query)
}

// 学籍番号から取得
func (ur *userRepository) FindByStudentNumber(c context.Context, studentNumber string) (*sql.Row) {
	query := "SELECT * FROM users WHERE student_number = " + studentNumber 
	row := ur.client.DB().QueryRowContext(c, query)
	fmt.Printf("\x1b[36m%s\n", query)
	return row
}


// 作成
func (ur *userRepository) Create(c context.Context, name string, mail string, gradeID string, departmentID string, bureauID string, roleID string, studentNumber string, tel string, password string) error {
	query := `
		INSERT INTO
			users (name, mail, grade_id, department_id, bureau_id, role_id, student_number, tel, password)
		VALUES ('` + name + "', '" + mail + "', " + gradeID + ", " + departmentID + ", " +  bureauID + ", " + roleID + ", " + studentNumber + ", '" + tel +  "', '" + password +"')"
	return ur.crud.UpdateDB(c, query)
}

// 編集
func (ur *userRepository) Update(c context.Context, id string, name string, mail string, gradeID string, departmentID string, bureauID string, roleID string, studentNumber string, tel string, password string) error {
	query := `
		UPDATE
			users
		SET
			name = '` + name +
		"', mail = " + mail +
		", grade_id = " + gradeID +
		", departmentid = " + departmentID +
		", bureau_id = " + bureauID +
		", role_id = " + roleID +
		", student_number = " + studentNumber +
		", tel = " + tel +
		", password = " + password +
		" WHERE id = " + id
	return ur.crud.UpdateDB(c, query)
}

// 削除
func (ur *userRepository) Delete(c context.Context, id string) error {
	query := "DELETE FROM users WHERE id = " + id
	return ur.crud.UpdateDB(c, query)
}

func (ur *userRepository) FindNewRecord(c context.Context) (*sql.Row, error) {
	query := "SELECT * FROM users ORDER BY id DESC LIMIT 1"
	return ur.crud.ReadByID(c, query)
}

// import '../../usecase/repository/user_repository.dart';
// import './external/database.dart';
// import '../../entity/entity.dart';

// class UserRepositoryImpl implements UserRepository {
//   Database database;

//   UserRepositoryImpl(this.database);

//   @override
//   Future<List<User>> getUsers(ctx) async {
//     String sql = '''
// SELECT * FROM users;
// ''';

//     List<Map<String, dynamic>> data = await database.finds(ctx, sql);
//     List<User> list = [];

//     for (Map<String, dynamic> d in data) {
//       User user = _convertUser(d);

//       if (user.isDeleted) {
//         continue;
//       }

//       list.add(user);
//     }

//     return list;
//   }

//   @override
//   Future<User> getUser(ctx, int id) async {
//     String sql = '''
// SELECT * FROM users WHERE id=$id;
// ''';
//     Map<String, dynamic> data = await database.find(ctx, sql);

//     return _convertUser(data);
//   }

//   @override
//   Future<User> insertUser(ctx, User req) async {
//     String sql = '''
// INSERT INTO users (name, bureau_id, grade_id)
//     VALUES ("${req.name}", "${req.bureauId}", "${req.gradeId}") 
//     returning *;
// ''';

//     var data = await database.insert(ctx, sql);
//     if (data["id"] == 0) {
//       print("error at UserRepository.insertUser");
//     }

//     return _convertUser(data);
//   }

//   @override
//   Future<User> updateUser(ctx, User user) async {
//     String sql = '''
// UPDATE users SET name="${user.name}" WHERE id=${user.id};
// ''';
//     String returningSQL = '''
// SELECT * FROM users WHERE id=${user.id};
// ''';

//     var data = await database.update(ctx, sql, returningSQL);

//     return _convertUser(data);
//   }

//   @override
//   Future<User> deleteUser(ctx, User user) async {
//     String sql = '''
// UPDATE users SET deleted_at=NOW() WHERE id=${user.id};
// ''';
//     String returningSQL = '''
// SELECT * FROM users WHERE id=${user.id};
// ''';

//     var data = await database.update(ctx, sql, returningSQL);

//     return _convertUser(data);
//   }

//   @override
//   Future<User> getUserByMail(ctx, User req) async {
//     String sql = '''
// SELECT id, name, mail, bureau_id, grade_id, created_at, updated_at, deleted_at
// FROM users WHERE mail="${req.mail}";
// ''';

//     Map<String, dynamic> data = await database.find(ctx, sql);
//     print(data);

//     return _convertUser(data);
//   }

//   User _convertUser(Map<String, dynamic> data) {
//     return User(
//       id: data['id'],
//       name: data['name'],
//       mail: data['mail'],
//       bureauId: data['bureau_id'],
//       gradeId: data['grade_id'],
//       createdAt: data['created_at'],
//       updatedAt: data['updated_at'],
//       deletedAt: data['deleted_at'],
//     );
//   }
// }
