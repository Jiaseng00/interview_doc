package main

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	var db *gorm.DB //模拟数据库链接
	err := updateBalanceWithPessimisticLock(db, 1, 50, "deposit")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = updateBalanceWithOptimisticLock(db, 1, 100, "deposit")
	if err != nil {
		fmt.Println(err)
		return
	}
}

// 悲观锁
func updateBalanceWithPessimisticLock(db *gorm.DB, userID int, amount float64, action string) error {
	type wallet struct {
		ID      int     `gorm:"id"`
		UserId  int     `gorm:"user_id"`
		Balance float64 `gorm:"balance"`
	}
	var userWallet wallet
	tx := db.Begin()

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", userID).Find(&userWallet).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	switch action {
	case "deposit":
		userWallet.Balance += amount
	case "withdrawal":
		if (userWallet.Balance - amount) >= 0 {
			userWallet.Balance -= amount
		} else {
			tx.Rollback()
			return errors.New("insufficient balance ")
		}
	default:
		tx.Rollback()
		return errors.New("invalid action,only \"deposit\" and \"withdrawal\" are supported")
	}

	err = tx.Model(&wallet{}).Where("user_id = ?", userID).Updates(&userWallet).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 乐观锁
func updateBalanceWithOptimisticLock(db *gorm.DB, userID int, amount float64, action string) error {
	type wallet struct {
		ID      int     `gorm:"id"`
		UserId  int     `gorm:"user_id"`
		Balance float64 `gorm:"balance"`
		Version int     `gorm:"version"`
	}
	var userWallet wallet
	tx := db.Begin()

	err := tx.Model(&wallet{}).Where("user_id = ?", userID).Find(&userWallet).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	switch action {
	case "deposit":
		userWallet.Balance += amount
	case "withdrawal":
		if (userWallet.Balance - amount) >= 0 {
			userWallet.Balance -= amount
		} else {
			tx.Rollback()
			return errors.New("insufficient balance ")
		}
	default:
		tx.Rollback()
		return errors.New("invalid action,only \"deposit\" and \"withdrawal\" are supported")
	}

	err = tx.Model(&wallet{}).Where("user_id = ? AND version = ?", userID, userWallet.Version).
		Updates(&wallet{Balance: userWallet.Balance, Version: userWallet.Version + 1}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
