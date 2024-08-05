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
var RedisHost string
var RedisPort string
var RedisPassword string

//если приложение поднято в кубере, то будут использоваться переменные окружения из configmap и secrets

func SetConfig() {
	Kuber = os.Getenv("kuber")

	if KafkaHost = os.Getenv("kafkaHost"); KafkaHost == "" {
		KafkaHost = "localhost"
	}
	if KafkaPort = os.Getenv("kafkaPort"); KafkaPort == "" {
		KafkaPort = "9092"
	}
	if KafkaTopic = os.Getenv("kafkaTopic"); KafkaTopic == "" {
		KafkaTopic = "json_topic"
	}
	if PostgresHost = os.Getenv("postgresHost"); PostgresHost == "" {
		PostgresHost = "localhost"
	}
	if PostgresPort = os.Getenv("postgresPort"); PostgresPort == "" {
		PostgresPort = "5432"
	}

	if RedisHost = os.Getenv("redisHost"); RedisHost == "" {
		RedisHost = "localhost"
	}
	if RedisPort = os.Getenv("redisPort"); RedisPort == "" {
		RedisPort = "6379"
	}

	//в случае с локальным запуском пароль и логин от базы можно забирать с помощью флагов командной строки
	//но тогда все равно придется указывать дефолтные значения на случай если аргументов нет
	if PostgresDatabaseName = os.Getenv("postgresDatabaseName"); PostgresDatabaseName == "" {
		PostgresDatabaseName = "postgres"
	}
	if PostgresPassword = os.Getenv("postgresPassword"); PostgresPassword == "" {
		PostgresPassword = "root"
	}
	if PostgresUser = os.Getenv("postgresUser"); PostgresUser == "" {
		PostgresUser = "postgres"
	}
	if RedisPassword = os.Getenv("redisPassword"); RedisPassword == "" {
		RedisPassword = ""
	}
}
