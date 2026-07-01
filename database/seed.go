package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to hash password: " + err.Error())
	}
	return string(hash)
}

func Seed(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	if err := seedDefaultAdmin(db); err != nil {
		return err
	}

	log.Println("Database seeding completed!")
	return nil
}

func seedDefaultAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&User{}).Where("role = ?", "admin").Count(&count)

	if count > 0 {
		log.Println("Admin user already exists, skipping...")
		return nil
	}

	admin := User{
		Name:     "Admin",
		Email:    "admin@learningsession.local",
		Password: hashPassword("admin123"),
		Role:     "admin",
	}

	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	log.Println("Admin user created:")
	log.Println("   Email: admin@learningsession.local")
	log.Println("   Password: admin123")

	return nil
}
