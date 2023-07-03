package mysql

type AutoMigrator interface {
	Init()
}

var migrates = []AutoMigrator{}

func Register(m AutoMigrator) {
	migrates = append(migrates, m)
}
