package install

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/code"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/proposal/tablesqls"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type initExecuteRequest struct {
	Language  string `form:"language" `  // 语言包
	RedisAddr string `form:"redis_addr"` // 连接地址，例如：127.0.0.1:6379
	RedisPass string `form:"redis_pass"` // 连接密码
	RedisDb   string `form:"redis_db"`   // 连接 db

	MySQLAddr string `form:"mysql_addr"`
	MySQLUser string `form:"mysql_user"`
	MySQLPass string `form:"mysql_pass"`
	MySQLName string `form:"mysql_name"`

	PgSQLAddr string `form:"pgsql_addr"`
	PgSQLUser string `form:"pgsql_user"`
	PgSQLPass string `form:"pgsql_pass"`
	PgSQLName string `form:"pgsql_name"`
	PgSQLPort string `form:"pgsql_port"`
}

func (h *handler) Execute() core.HandlerFunc {

	installTableList := map[string]map[string]string{
		"authorized": {
			"table_sql":        tablesqls.CreateAuthorizedTableSql(),
			"table_pgsql":      tablesqls.CreateAuthorizedTablePGSql(),
			"table_data_sql":   tablesqls.CreateAuthorizedTableDataSql(),
			"table_data_pgsql": tablesqls.CreateAuthorizedTableDataPGSql(),
		},
		"authorized_api": {
			"table_sql":        tablesqls.CreateAuthorizedAPITableSql(),
			"table_pgsql":      tablesqls.CreateAuthorizedAPITablePGSql(),
			"table_data_sql":   tablesqls.CreateAuthorizedAPITableDataSql(),
			"table_data_pgsql": tablesqls.CreateAuthorizedAPITableDataPGSql(),
		},
		"admin": {
			"table_sql":        tablesqls.CreateAdminTableSql(),
			"table_pgsql":      tablesqls.CreateAdminTablePGSql(),
			"table_data_sql":   tablesqls.CreateAdminTableDataSql(),
			"table_data_pgsql": tablesqls.CreateAdminTableDataPGSql(),
		},
		"admin_menu": {
			"table_sql":        tablesqls.CreateAdminMenuTableSql(),
			"table_pgsql":      tablesqls.CreateAdminMenuTablePGSql(),
			"table_data_sql":   tablesqls.CreateAdminMenuTableDataSql(),
			"table_data_pgsql": tablesqls.CreateAdminMenuTableDataPGSql(),
		},
		"menu": {
			"table_sql":        tablesqls.CreateMenuTableSql(),
			"table_pgsql":      tablesqls.CreateMenuTablePGSql(),
			"table_data_sql":   tablesqls.CreateMenuTableDataSql(),
			"table_data_pgsql": tablesqls.CreateMenuTableDataPGSql(),
		},
		"menu_action": {
			"table_sql":        tablesqls.CreateMenuActionTableSql(),
			"table_pgsql":      tablesqls.CreateMenuActionTablePGSql(),
			"table_data_sql":   tablesqls.CreateMenuActionTableDataSql(),
			"table_data_pgsql": tablesqls.CreateMenuActionTableDataPGSql(),
		},
		"cron_task": {
			"table_sql":        tablesqls.CreateCronTaskTableSql(),
			"table_pgsql":      tablesqls.CreateCronTaskTablePGSql(),
			"table_data_sql":   "",
			"table_data_pgsql": "",
		},
	}

	return func(ctx core.Context) {
		req := new(initExecuteRequest)
		if err := ctx.ShouldBindForm(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		// region 验证 version
		versionStr := runtime.Version()
		version := cast.ToFloat32(versionStr[2:6])
		if version < configs.MinGoVersion {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.GoVersionError,
				code.Text(code.GoVersionError)),
			)
			return
		}
		// endregion

		// region 验证 Redis 配置
		cfg := configs.Get()
		redisClient := redis.NewClient(&redis.Options{
			Addr:         req.RedisAddr,
			Password:     req.RedisPass,
			DB:           cast.ToInt(req.RedisDb),
			MaxRetries:   cfg.Redis.MaxRetries,
			PoolSize:     cfg.Redis.PoolSize,
			MinIdleConns: cfg.Redis.MinIdleConns,
		})

		if err := redisClient.Ping().Err(); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.RedisConnectError,
				code.Text(code.RedisConnectError)).WithError(err),
			)
			return
		}

		defer redisClient.Close()

		outPutString := "已检测 Redis 配置可用。\n"
		// endregion

		// region 验证 MySQL 配置
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
			req.MySQLUser,
			req.MySQLPass,
			req.MySQLAddr,
			req.MySQLName,
			true,
			"Local")

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			//Logger: logger.Default.LogMode(logger.Info), // 日志配置
		})

		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.MySQLConnectError,
				code.Text(code.MySQLConnectError)).WithError(err),
			)
			return
		}

		db.Set("gorm:table_options", "CHARSET=utf8mb4")

		dbClient, _ := db.DB()
		defer dbClient.Close()

		outPutString += "已检测 MySQL 配置可用。\n"
		// endregion

		// region 验证 PgsqlSQL 配置
		pgDsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			req.PgSQLAddr,
			req.PgSQLPort,
			req.PgSQLUser,
			req.PgSQLName,
			"disable",
			req.PgSQLPass,
		)
		//gorm.Open("postgres", dataSource)
		pgDB, err := gorm.Open(postgres.Open(pgDsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			//Logger: logger.Default.LogMode(logger.Info), // 日志配置
		})

		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.MySQLConnectError,
				code.Text(code.PgSQlConnectError)).WithError(err),
			)
			return
		}

		db.Set("gorm:table_options", "CHARSET=utf8mb4")

		pgDBClient, _ := pgDB.DB()
		defer pgDBClient.Close()

		outPutString += "已检测 PgSQL 配置可用。\n"

		// region 写入配置文件
		viper.Set("language.local", req.Language)

		viper.Set("redis.addr", req.RedisAddr)
		viper.Set("redis.pass", req.RedisPass)
		viper.Set("redis.db", req.RedisDb)

		viper.Set("mysql.read.addr", req.MySQLAddr)
		viper.Set("mysql.read.user", req.MySQLUser)
		viper.Set("mysql.read.pass", req.MySQLPass)
		viper.Set("mysql.read.name", req.MySQLName)

		viper.Set("mysql.write.addr", req.MySQLAddr)
		viper.Set("mysql.write.user", req.MySQLUser)
		viper.Set("mysql.write.pass", req.MySQLPass)
		viper.Set("mysql.write.name", req.MySQLName)

		viper.Set("pgsql.read.addr", req.PgSQLAddr)
		viper.Set("pgsql.read.user", req.PgSQLUser)
		viper.Set("pgsql.read.pass", req.PgSQLPass)
		viper.Set("pgsql.read.name", req.PgSQLName)
		viper.Set("pgsql.read.port", req.PgSQLPort)

		viper.Set("pgsql.write.addr", req.PgSQLAddr)
		viper.Set("pgsql.write.user", req.PgSQLUser)
		viper.Set("pgsql.write.pass", req.PgSQLPass)
		viper.Set("pgsql.write.name", req.PgSQLName)
		viper.Set("pgsql.write.port", req.PgSQLPort)

		if viper.WriteConfig() != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.WriteConfigError,
				code.Text(code.WriteConfigError)).WithError(err),
			)
			return
		}

		outPutString += "语言包 " + req.Language + " 配置成功。\n"
		outPutString += "配置项 Redis、MySQL、PgSQL 配置成功。\n"
		// endregion

		// region 初始化表结构 + 默认数据
		for k, v := range installTableList {
			if v["table_sql"] != "" {
				// region 初始化表结构
				if err = db.Exec(v["table_sql"]).Error; err != nil {
					ctx.AbortWithError(core.Error(
						http.StatusBadRequest,
						code.MySQLExecError,
						code.Text(code.MySQLExecError)+" "+err.Error()).WithError(err),
					)
					return
				}

				outPutString += "初始化 MySQL 数据表：" + k + " 成功。\n"
				// endregion

				// region 初始化默认数据
				if v["table_data_sql"] != "" {
					if err = db.Exec(v["table_data_sql"]).Error; err != nil {
						ctx.AbortWithError(core.Error(
							http.StatusBadRequest,
							code.MySQLExecError,
							code.Text(code.MySQLExecError)+" "+err.Error()).WithError(err),
						)
						return
					}

					outPutString += "初始化 MySQL 数据表：" + k + " 默认数据成功。\n"
				}
				// endregion
			}
			if v["table_pgsql"] != "" {
				// region 初始化表结构
				if err = pgDB.Exec(v["table_pgsql"]).Error; err != nil {
					ctx.AbortWithError(core.Error(
						http.StatusBadRequest,
						code.MySQLExecError,
						code.Text(code.MySQLExecError)+" "+err.Error()).WithError(err),
					)
					return
				}

				outPutString += "初始化 Postgre 数据表：" + k + " 成功。\n"
				// endregion

				// region 初始化默认数据
				if v["table_data_sql"] != "" {
					if err = pgDB.Exec(v["table_data_pgsql"]).Error; err != nil {
						ctx.AbortWithError(core.Error(
							http.StatusBadRequest,
							code.MySQLExecError,
							code.Text(code.MySQLExecError)+" "+err.Error()).WithError(err),
						)
						return
					}

					outPutString += "初始化 Postgre 数据表：" + k + " 默认数据成功。\n"
				}
				// endregion
			}

		}
		// endregion

		// region 生成 install 完成标识
		f, err := os.Create(configs.ProjectInstallMark)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.MySQLExecError,
				code.Text(code.MySQLExecError)+" "+err.Error()).WithError(err),
			)
			return
		}
		defer f.Close()
		// endregion

		ctx.Payload(outPutString)
	}
}
