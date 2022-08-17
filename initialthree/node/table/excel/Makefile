excelPath = ""
all:
	cd EnumType/CharacterUsualAttribute;go run gen.go ${excelPath};cd ../../
	cd EnumType/CharacterBattleAttribute;go run gen.go ${excelPath};cd ../../
	cd EnumType/CommonEnums;go run gen.go ${excelPath};cd ../../
	go run main/gen_table.go ${excelPath}
	go run main/gen_check.go ${excelPath}
	cd check;go run main.go ${excelPath}