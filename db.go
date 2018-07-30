package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Usuario es la estructura para el ORM - tabla
type Usuario struct {
	gorm.Model
	Nombre   string
	Password string
}

//Servicio es la estructura de inventario MariaDB para el proyecto de MiniRDS
type Servicio struct {
	gorm.Model
	Usuario  Usuario
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

// Dbinit inicializa la base de datos
func Dbinit(elhost string, labase string) {
	db, err := gorm.Open("mysql", "gouser:gopasswd@tcp("+elhost+")/"+labase+"?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("Error de base de datos")
		panic(err)
	}
	defer db.Close()
	db.AutoMigrate(&Usuario{})
	db.AutoMigrate(&Servicio{})
	fmt.Println("-------------")
	fmt.Println("Tables created")
}

//InsertaUsuario is a Funcion para insertar campos en Base de datos
func InsertaUsuario(elhost string, labase string, nombre string, password string) {
	db, err := gorm.Open("mysql", "gouser:gopasswd@tcp("+elhost+")/"+labase+"?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("Error de base de datos")
		panic(err)
	}
	defer db.Close()
	db.Create(&Usuario{Nombre: nombre, Password: password})
	fmt.Println("-------------")
	fmt.Println("Usuario insertado")
}

//InsertaServicio is a Funcion para insertar campos en Base de datos
func InsertaServicio(elhost string, labase string, nombre string,
	host string, port string, database string, user string, password string) {
	usuario := EncuentraUsuario(elhost, labase, nombre)
	db, err := gorm.Open("mysql", "gouser:gopasswd@tcp("+elhost+")/"+labase+"?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("Error de base de datos")
		panic(err)
	}
	defer db.Close()
	db.Create(&Servicio{Usuario: usuario, Host: host, Port: port,
		Database: database, User: user, Password: password})
	fmt.Println("-------------")
	fmt.Println("Usuario insertado")
}

// EncuentraUsuario es una función para hacer queries simples
func EncuentraUsuario(elhost string, labase string, busca string) Usuario {
	var usuario Usuario
	db, err := gorm.Open("mysql",
		"gouser:gopasswd@tcp("+elhost+")/"+labase+"?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("Error de base de datos")
		panic(err)
	}
	defer db.Close()
	//sal := db.Where("Name LIKE ?", "%"+busca+"%").
	//	Or("Sku LIKE ?", "%"+busca+"%").Find(&Usuario)
	db.Where("`Nombre` LIKE ? ", "%"+busca+"%").Find(&usuario)
	fmt.Println("")
	return usuario
}

// EncuentraServicio es una función para hacer queries simples
func EncuentraServicio(elhost string, labase string, busca string) Servicio {
	var servicio Servicio
	db, err := gorm.Open("mysql",
		"gouser:gopasswd@tcp("+elhost+")/"+labase+"?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println("Error de base de datos")
		panic(err)
	}
	defer db.Close()
	//sal := db.Where("Name LIKE ?", "%"+busca+"%").
	//	Or("Sku LIKE ?", "%"+busca+"%").Find(&Usuario)
	db.Where("`Host` LIKE ? OR `Database` LIKE ? OR `User` LIKE ?", "%"+busca+"%",
		"%"+busca+"%", "%"+busca+"%").Find(&servicio)
	fmt.Println("")
	return (servicio)
}
