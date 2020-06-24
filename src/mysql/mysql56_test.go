/*
 * Xenon
 *
 * Copyright 2018 The Xenon Authors.
 * Code is licensed under the GPLv3.
 *
 */

package mysql

import (
	"testing"

	"config"
	"xbase/xlog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestMysql56Handler(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.PANIC))
	conf := config.DefaultMysqlConfig()
	conf.Version = "mysql56"

	mysql := NewMysql(conf, log)
	want := new(Mysql56)
	got := mysql.mysqlHandler
	assert.Equal(t, want, got)
}

func TestMysql56SetSemiWaitSlaveCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// log
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	conf := config.DefaultMysqlConfig()
	conf.Version = "mysql56"
	mysql56 := NewMysql(conf, log)
	mysql56.db = db

	queryList := []string{
		"SET GLOBAL rpl_semi_sync_master_wait_for_slave_count = 2",
	}

	mock.ExpectExec(queryList[0]).WillReturnResult(sqlmock.NewResult(1, 1))
	err = mysql56.SetSemiWaitSlaveCount(2)
	assert.Nil(t, err)
}

func TestMysql56ChangeUserPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// log
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	conf := config.DefaultMysqlConfig()
	conf.Version = "mysql56"
	mysql56 := NewMysql(conf, log)
	mysql56.db = db

	queryList := []string{
		"SET PASSWORD FOR `usr`@`127.0.0.1` = PASSWORD('pwd')",
	}

	mock.ExpectExec(queryList[0]).WillReturnResult(sqlmock.NewResult(1, 1))
	err = mysql56.ChangeUserPasswd("usr", "127.0.0.1", "pwd")
	assert.Nil(t, err)
}

func TestMysql56CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// log
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	conf := config.DefaultMysqlConfig()
	conf.Version = "mysql56"
	mysql56 := NewMysql(conf, log)
	mysql56.db = db

	// ssl is NO
	query := "GRANT USAGE ON *.* TO `xx`@`192.168.0.%` IDENTIFIED BY 'xxx'"
	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	err = mysql56.CreateUser("xx", "192.168.0.%", "xxx", "NO")
	assert.Nil(t, err)

	// ssl is YES
	query = "GRANT USAGE ON *.* TO `xx`@`192.168.0.%` IDENTIFIED BY 'xxx' REQUIRE SSL"
	mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	err = mysql56.CreateUser("xx", "192.168.0.%", "xxx", "YES")
	assert.Nil(t, err)
}

func TestMysql56CreateUserWithPrivileges(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// log
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	conf := config.DefaultMysqlConfig()
	conf.Version = "mysql56"
	mysql56 := NewMysql(conf, log)
	mysql56.db = db

	queryList := []string{
		"GRANT USAGE ON *.* TO `xx`@`127.0.0.1` IDENTIFIED BY 'pwd'",
		"GRANT USAGE ON *.* TO `xx`@`127.0.0.1` IDENTIFIED BY 'pwd' REQUIRE SSL",
		"GRANT ALTER , ALTER ROUTINE ON test.* TO `xx`@`127.0.0.1`",
	}

	// ssl is NO
	mock.ExpectExec(queryList[0]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(queryList[2]).WillReturnResult(sqlmock.NewResult(1, 1))
	err = mysql56.CreateUserWithPrivileges("xx", "pwd", "test", "*", "127.0.0.1", "ALTER , ALTER ROUTINE", "NO")
	assert.Nil(t, err)

	// ssl is YES
	mock.ExpectExec(queryList[1]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(queryList[2]).WillReturnResult(sqlmock.NewResult(1, 1))
	err = mysql56.CreateUserWithPrivileges("xx", "pwd", "test", "*", "127.0.0.1", "ALTER , ALTER ROUTINE", "YES")
	assert.Nil(t, err)
}

func TestMysql56GrantAllPrivileges(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	// log
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	conf := config.DefaultMysqlConfig()
	conf.Version = "mysql56"
	mysql56 := NewMysql(conf, log)
	mysql56.db = db

	queryList := []string{
		"GRANT USAGE ON *.* TO `xx`@`192.168.0.%` IDENTIFIED BY 'pwd'",
		"GRANT USAGE ON *.* TO `xx`@`192.168.0.%` IDENTIFIED BY 'pwd' REQUIRE SSL",
		"GRANT ALL ON *.* TO `xx`@`192.168.0.%` WITH GRANT OPTION",
	}

	// ssl is NO
	mock.ExpectExec(queryList[0]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(queryList[2]).WillReturnResult(sqlmock.NewResult(1, 1))
	err = mysql56.GrantAllPrivileges("xx", "192.168.0.%", "pwd", "NO")
	assert.Nil(t, err)

	// ssl is YES
	mock.ExpectExec(queryList[1]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(queryList[2]).WillReturnResult(sqlmock.NewResult(1, 1))
	err = mysql56.GrantAllPrivileges("xx", "192.168.0.%", "pwd", "YES")
	assert.Nil(t, err)
}
