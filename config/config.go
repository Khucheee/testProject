package config

import "os"

var Kuber string
var KafkaHost string
var KafkaPort string
var KafkaTopic string
var PostgresHost string
var PostgresPort string
var PostgresDatabaseName string
var PostgresPassword string //храним с секрете
var PostgresUser string     //храним в секрете

func SetConfig() {
	Kuber = os.Getenv("kuber")

	if KafkaHost = os.Getenv("kafkahost"); KafkaHost == "" {
		KafkaHost = "localhost"
	}
	if KafkaPort = os.Getenv("kafkaport"); KafkaPort == "" {
		KafkaPort = "9092"
	}
	if KafkaTopic = os.Getenv("kafkatopic"); KafkaTopic == "" {
		KafkaTopic = "json_topic"
	}
	if PostgresHost = os.Getenv("postgreshost"); PostgresHost == "" {
		PostgresHost = "localhost"
	}
	if PostgresPort = os.Getenv("postgresport"); PostgresPort == "" {
		PostgresPort = "5432"
	}

	//в случае с локальным запуском пароль и логин от базы можно забирать с помощью флагов командной строки
	//но тогда все равно придется указывать дефолтные значения на случай если аргументов нет
	if PostgresDatabaseName = os.Getenv("postgresdatabasename"); PostgresDatabaseName == "" {
		PostgresDatabaseName = "postgres"
	}
	if PostgresPassword = os.Getenv("postgrespassword"); PostgresPassword == "" {
		PostgresPassword = "root"
	}
	if PostgresUser = os.Getenv("postgresuser"); PostgresUser == "" {
		PostgresUser = "postgres"
	}
}
