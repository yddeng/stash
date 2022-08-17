proto:
	make gen_cs;make gen_ss;make gen_rpc;make gen_cmd
gen_ss:
	cd protocol/;make gen_ss;cd ../
gen_cs:
	cd protocol/;make gen_cs;cd ../
gen_rpc:
	cd rpc;go run gen_rpc.go;cd ../
gen_cmd:
	cd protocol/cmdEnum/gen;go run cmd_gen.go;cd ../../../

build_center:
	test -d bin || mkdir -p bin
	cd bin;go build ../center/main/center.go;cd ../
build_dir:
	test -d bin || mkdir -p bin	
	cd bin;go build ../node/node_dir/main/node_dir.go;cd ../
build_login:
	test -d bin || mkdir -p bin	
	cd bin;go build ../node/node_login/main/node_login.go;cd ../
build_game:
	test -d bin || mkdir -p bin	
	cd bin;go build ../node/node_game/main/node_game.go;cd ../
build_world:
	test -d bin || mkdir -p bin
	cd bin;go build ../node/node_world/main/node_world.go;cd ../
build_map:
	test -d bin || mkdir -p bin	
	cd bin;go build ../node/node_map/main/node_map.go;cd ../
build_gate:
	test -d bin || mkdir -p bin	
	cd bin;go build ../node/node_gate/main/node_gate.go;cd ../
build_webservice:
	test -d bin || mkdir -p bin
	go build -o bin/node_webservice node/node_webservice/main/node_webservice.go
build_team:
	test -d bin || mkdir -p bin
	go build -o bin/node_team node/node_team/main/node_team.go
build_rank:
	test -d bin || mkdir -p bin
	go build -o bin/node_rank node/node_rank/main/node_rank.go
exe:
	make build_center
	make build_dir
	make build_login
	make build_game
#	make build_world
#	make build_map
	make build_gate
	make build_webservice
#	make build_team
	make build_rank


build_flyfish:
	test -d bin || mkdir -p bin
	cd bin;go build ../flyfish/main/flyfishd.go;cd ../


sourcePath=/Users/yidongdeng/svn/TheInitial3_Dev/TheInitial3_Dev_Config# 配置,注意前后不能有空格。 可通过 make table sourcePath=... 的方式变量赋值
table:
	cd node/table/excel;make all excelPath=${sourcePath}/Excel


