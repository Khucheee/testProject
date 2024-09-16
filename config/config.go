package config

import "os"

var Kuber string
var KafkaHost string
var KafkaPort string
var KafkaTopic string
var KafkaLogTopic string
var PostgresHost string
var PostgresPort string
var PostgresDatabaseName string
var PostgresPassword string //храним с секрете
var PostgresUser string     //храним в секрете
var RedisHost string
var RedisPort string
var RedisPassword string
var LogstashPort string
var LogstashHost string
var ElasticsearchPort string
var ElasticsearchHost string
var KibanaHost string
var KibanaPort string
var WorkersCount string
var RepositoryRetries string
var KafkaEnabled bool

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
	if KafkaLogTopic = os.Getenv("kafkaLogTopic"); KafkaLogTopic == "" {
		KafkaLogTopic = "log_topic"
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
	if LogstashHost = os.Getenv("logstashHost"); LogstashHost == "" {
		LogstashHost = "logstash"
	}
	if LogstashPort = os.Getenv("logstashPort"); LogstashPort == "" {
		LogstashPort = "5044"
	}
	if ElasticsearchHost = os.Getenv("elasticsearchHost"); ElasticsearchHost == "" {
		ElasticsearchHost = "elasticsearch"
	}
	if ElasticsearchPort = os.Getenv("elasticsearchPort"); ElasticsearchPort == "" {
		ElasticsearchPort = "9200"
	}
	if KibanaHost = os.Getenv("kibanaHost"); KibanaHost == "" {
		KibanaHost = "localhost"
	}
	if KibanaPort = os.Getenv("kibanaPort"); KibanaPort == "" {
		KibanaPort = "5601"
	}
	if WorkersCount = os.Getenv("workersCount"); WorkersCount == "" {
		WorkersCount = "15"
	}
	if RepositoryRetries = os.Getenv("repositoryRetries"); RepositoryRetries == "" {
		RepositoryRetries = "3"
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
