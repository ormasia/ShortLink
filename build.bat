chcp 65001
@echo off
:loop
@echo off&amp;color 0A
cls
echo,
@REM echo 请选择要编译的系统环境：
@REM echo,
@REM echo 1. Windows_amd64
@REM echo 2. linux_amd64


echo 编译Linux版本64位
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o target/api ./apigateway/main.go
go build -o target/short ./shortlinkcore/main.go
go build -o target/user ./userservice/main.go