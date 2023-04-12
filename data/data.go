package data

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:123456@tcp(127.0.0.1:3310)/isolation_test?charset=utf8mb4&parseTime=True&loc=Local"

type Account struct {
	ID      uint `gorm:"primaryKey"`
	Balance int
}

// Transfer transfers money between two accounts.
func Transfer(db *gorm.DB, fromID, toID uint, amount int) error {
	var fromAccount, toAccount Account
	if err := db.First(&fromAccount, fromID).Error; err != nil {
		return err
	}
	if err := db.First(&toAccount, toID).Error; err != nil {
		return err
	}
	if fromAccount.Balance < amount {
		return errors.New("insufficient balance")
	}
	fromAccount.Balance -= amount
	toAccount.Balance += amount
	if err := db.Save(&fromAccount).Error; err != nil {
		return err
	}
	if err := db.Save(&toAccount).Error; err != nil {
		return err
	}
	return nil
}

// 数据库构建Account表，并初始化数据
func Init() {
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&Account{})
	db.Create(&Account{ID: 1, Balance: 1000})
	db.Create(&Account{ID: 2, Balance: 1000})
}

// 读未提交（Read Uncommitted）
func ReadUncommitted() {
	db, _ := gorm.Open(mysql.Open(dsn+"&transaction_isolation='READ-UNCOMMITTED'"), &gorm.Config{})

	tx := db.Begin()
	// tx.Exec("SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED")

	err := Transfer(tx, 1, 2, 100)
	if err != nil {
		fmt.Println("Transfer failed:", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

// 读已提交（Read Committed）
func ReadCommitted() {
	db, _ := gorm.Open(mysql.Open(dsn+"&transaction_isolation='READ-COMMITTED'"), &gorm.Config{})

	tx := db.Begin()
	// tx.Exec("SET TRANSACTION ISOLATION LEVEL READ COMMITTED")

	err := Transfer(tx, 1, 2, 100)
	if err != nil {
		fmt.Println("Transfer failed:", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

// 可重复读（Repeatable Read）
func RepeatableRead() {
	db, _ := gorm.Open(mysql.Open(dsn+"&transaction_isolation='REPEATABLE-READ'"), &gorm.Config{})

	tx := db.Begin()
	// tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")

	err := Transfer(tx, 1, 2, 100)
	if err != nil {
		fmt.Println("Transfer failed:", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

// 串行化（Serializable）
func Serializable() {
	db, _ := gorm.Open(mysql.Open(dsn+"&transaction_isolation='SERIALIZABLE'"), &gorm.Config{})

	tx := db.Begin()
	// tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")

	err := Transfer(tx, 1, 2, 100)
	if err != nil {
		fmt.Println("Transfer failed:", err)
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
