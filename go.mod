module go-mysql

go 1.13

replace vitess.io/vitess => ../vitess

require (
	github.com/go-sql-driver/mysql v1.5.0
	vitess.io/vitess v2.1.1+incompatible
)
