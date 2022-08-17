nohup ./flypd -config=./flypd_config.toml -raftcluster=1@1@http://localhost:8811@127.0.0.1:10050@voter flyfish > /dev/null 2>pd_err.log &
nohup ./flykv -config=./flykv_config.toml -id=1 flyfish > /dev/null 2>kv_err.log &
nohup ./flygate -config=./flygate_config.toml flyfish > /dev/null 2>gate_err.log &