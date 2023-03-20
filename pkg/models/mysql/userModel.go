package mysql

import (
	"database/sql"
	"mApi/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) InsertUser(name string, age int64) (int, error) {
	stmt := `insert into socialDB.users (name, age) VALUES (?, ?)`

	result, err := m.DB.Exec(stmt, name, age)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *UserModel) MakeFriends(user_id, target_id int64) error {
	stmt := `insert into socialDB.friends (user_id, friend_id) values (?, ?)`

	_, err := m.DB.Exec(stmt, user_id, target_id)
	if err != nil {
		return err
	}

	_, err = m.DB.Exec(stmt, target_id, user_id)
	if err != nil {
		return nil
	}

	return nil
}

func (m *UserModel) DeleteUser(user_id int64) error {
	stmt := `delete from socialDB.users where user_id = ?`

	_, err := m.DB.Exec(stmt, user_id)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) GetFriends(user_id int64) (*[]models.UserForDB, error) {
	stmt := `select f.friend_id, u.name, u.age from socialDB.users u join friends f on f.friend_id = u.user_id where f.user_id = ?;
`

	row, err := m.DB.Query(stmt, user_id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var results []models.UserForDB
	for row.Next() {
		userTemp := models.UserForDB{}

		err := row.Scan(&userTemp.User_ID, &userTemp.Name, &userTemp.Age)
		if err != nil {
			return nil, err
		}
		//age := json.Number(userTemp.Age)

		user := models.UserForDB{User_ID: userTemp.User_ID, Name: userTemp.Name, Age: userTemp.Age}
		results = append(results, user)
	}
	return &results, nil
}

func (m *UserModel) UpdateAge(user_id int, age int64) error {
	stmt := `update socialDB.users set age = ? where user_id = ?`

	_, err := m.DB.Exec(stmt, age, user_id)
	if err != nil {
		return err
	}

	return nil
}
