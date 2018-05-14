tabtoy ^
--mode=exportorv2 ^
--go_out=task.go ^
--json_out=config.json ^
--combinename=Config ^
--lan=zh_cn ^
Globals.xlsx
task.xlsx

@IF %ERRORLEVEL% NEQ 0 pause