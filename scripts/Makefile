# proj root folder
PROJROOT := $(abspath $(shell pwd)/../)
# keys path
KEYPATH := $(PROJROOT)/keys
APPPATH := $(PROJROOT)/src/go/notion


SQLSERVER_PASS := deliciousMango1

notion_env:
	./setup_notion.sh

google_env:
	export GOOGLE_APPLICATION_CREDENTIALS=$(KEYPATH)/animated-bay-338209-e3cea0254898.json

sql_install:
	sudo docker pull mcr.microsoft.com/mssql/server:2019-latest

sql_run:
	docker-compose up -d

sql_create:
	sqlcmd -S localhost,1433 -U SA -Q "CREATE DATABASE TestDB;"

sql_listDB:
	sqlcmd -S localhost,1433 -U SA -P "tmpPass1" -Q "SELECT Name FROM sys.Databases;"

sql_createUser:
	sqlcmd -S localhost,1433 -U SA -P "tmpPass1" -Q 'CREATE LOGIN autoio WITH PASSWORD = "autoPass1";'
	sqlcmd -S localhost,1433 -U SA -P "tmpPass1" -Q "CREATE USER autoio FOR LOGIN autoio;"
sql_grant:
	sqlcmd -S localhost,1433 -U SA -P "tmpPass1" -Q "GRANT CREATE ON TestDB.* TO 'autoio'@'34.206.38.217';"

build_app:
	cd $(APPPATH); go build -o app .; mv $(APPPATH)/app $(PROJROOT)/bin/app

test_all:
	./test.sh

test_generate:
	./test_generate.sh

test_insert:
	./test_insert.sh

test_update:
	./test_update.sh

clean:
	./clean.sh