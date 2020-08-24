# 主版本
VERSION ?= $(shell git describe --tags --always --dirty)

# build 序列
BUILD_NO ?= $(shell git show -s --format=%H)

# Go文件列表
GOFILES := $(shell find . ! -path "./vendor/*" -name "*.go")

# 支持的操作系统列表
GOOSES := linux windows darwin plan9
# 支持的CPU架构
GOARCHES := 386 amd64

# 目标输出目录
DIST_FOLDER := dist

# 目标可执行程序源码所在的目录
COMMAND_ROOT := cmd

# 构建附加选项
BUILD_OPTS := -ldflags "-s -w"

# 单元测试附加选项
TEST_OPTS := -v

# 基准测试附加选项
BENCHMARK_OPTS := -cpu 1,2,3,4,5,6,7,8

# 目标可执行程序列表
COMMAND_LIST := scalpel suture

# sonar 相关报告输出路径（包括：单元测试报告输出，单元测试覆盖率报告，golint 报告，golangci-lint 报告）
REPORT_FOLDER := sonar

# sonar report
TEST_REPORT := ${REPORT_FOLDER}/test.report 
COVER_REPORT := ${REPORT_FOLDER}/cover.report
GOLANGCI_LINT_REPORT := ${REPORT_FOLDER}/golangci-lint.xml 
GOLINT_REPORT := ${REPORT_FOLDER}/golint.report 

.PHONY: build format test benchmark sonar all clean

.DEFAULT: build 

# 构建目标
build: $(GOFILES)
	@for command in ${COMMAND_LIST} ; do 																		\
		go build ${BUILD_OPTS} -o ${DIST_FOLDER}/$${command}/$${command} ./${COMMAND_ROOT}/$${command} ;		\
	done																										\

# 格式化
format:
	@for f in ${GOFILES} ; do 																					\
		gofmt -w $${f};																							\
	done																										\

# 单元测试
test: 
	go test ${TEST_OPTS} ./...

# 基准测试
benchmark:
	go test -bench . -run ^$$ ${BENCHMARK_OPTS}  ./...

# sonar
sonar: 
	mkdir -p ${REPORT_FOLDER}
	go test -json ./... > ${TEST_REPORT}
	go test -coverprofile=${COVER_REPORT} ./... 
	golangci-lint run --out-format checkstyle  ./... > ${GOLANGCI_LINT_REPORT}
	golint ./... > ${GOLINT_REPORT}
	sonar-scanner

# 构建所有支持的操作系统和架构的目标文件
all:
	@for command in ${COMMAND_LIST} ; do 																	\
		for os in ${GOOSES} ; do																			\
			for arch in ${GOARCHES} ; do 																	\
				if [ "$${os}" = "windows" ] ;then															\
					GOOS=$${os} GOARCH=$${arch}  															\
					go build ${BUILD_OPTS} 																	\
					-o ${DIST_FOLDER}/$${command}/$${os}_$${arch}/$${command}.exe 							\
					./${COMMAND_ROOT}/$${command} ;															\
				else																						\
					GOOS=$${os} GOARCH=$${arch}  															\
					go build ${BUILD_OPTS} 																	\
					-o ${DIST_FOLDER}/$${command}/$${os}_$${arch}/$${command} 								\
					./${COMMAND_ROOT}/$${command} ;															\
				fi																							\
			done																							\
		done																								\
	done																									\

# 清理
clean:
	-rm -rf $(DIST_FOLDER)/*
	-rm -f ${TEST_REPORT}
	-rm -f ${COVER_REPORT}
	-rm -f ${GOLANGCI_LINT_REPORT}
	-rm -f ${GOLINT_REPORT}
