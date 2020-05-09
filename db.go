package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "service_man:canson123@/greenbook")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func IsUnique(email string) bool {
	db, err := OpenDB()
	if err != nil {
		return false
	}
	defer db.Close()
	var num int
	row := db.QueryRow("call isExistsEmail(?)", email)
	err = row.Scan(&num)
	if num == 1 || err != nil {
		return false
	}
	return true
}

func Registration(user User) (int64, error) {
	var id int64
	db, err := OpenDB()
	if err != nil {
		return id, err
	}
	defer db.Close()
	row := db.QueryRow("call registration(?, ?, ?, ?, ?, ?, ?, ?)", user.FirstName, user.LastName, user.Age, user.City, user.Sex, user.Interests, user.Email, user.Password)
	err = row.Scan(&id)
	if err != nil {
		log.Println(err)
		return id, err
	}
	return id, nil
}

func GetPrivatePage(id int64) (User, error) {
	user := User{}
	user.Id = id
	db, err := OpenDB()
	if err != nil {
		return user, err
	}
	defer db.Close()
	row := db.QueryRow("call getPrivatePage(?)", id)
	err = row.Scan(&user.FirstName, &user.LastName, &user.Age, &user.Sex, &user.City, &user.Interests)
	if err != nil {
		return user, err
	}
	return user, err
}

func GetCommonPage(userId, commonId int64) (CommonUser, error) {
	commonUser := CommonUser{}
	var num int
	db, err := OpenDB()
	if err != nil {
		return commonUser, err
	}
	defer db.Close()
	row := db.QueryRow("call getCommonPage(?, ?)", userId, commonId)
	err = row.Scan(&commonUser.FirstName, &commonUser.LastName, &commonUser.Age, &commonUser.Sex, &commonUser.City, &commonUser.Interests, &num)
	if err != nil {
		return commonUser, err
	}
	if num > 0 {
		commonUser.IsFriend = true
	} else {
		commonUser.IsFriend = false
	}
	return commonUser, nil
}

func GetFriends(userId int64) ([]Friend, error) {
	var friends []Friend
	db, err := OpenDB()
	if err != nil {
		return friends, err
	}
	defer db.Close()
	rows, err := db.Query("call getFriends(?)", userId)
	if err != nil {
		return friends, err
	}
	defer rows.Close()
	friend := Friend{}
	for rows.Next() {
		err = rows.Scan(&friend.Id, &friend.FirstName, &friend.LastName)
		if err != nil {
			return friends, err
		}
		friends = append(friends, friend)
	}
	return friends, err
}

func AddFriend(userId, friendId int64) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("call addFriend(?, ?)", userId, friendId)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFriend(userId, friendId int64) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("call removeFriend(?, ?)", userId, friendId)
	if err != nil {
		return err
	}
	return nil
}

func LoginPass(login Login) bool {
	db, err := OpenDB()
	if err != nil {
		return false
	}
	defer db.Close()
	row := db.QueryRow("call login(?)", login.Email)
	var hashPass string
	err = row.Scan(&hashPass)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(login.Password))
	if err != nil {
		return false
	}
	return true
}
