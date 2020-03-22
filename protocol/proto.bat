set outPath=../proto
set fileArray=(WSBaseReqProto WSBaseResProto WSMessageResProto WSUserResProto)

# 将.proto文件生成java类
for %%i in %fileArray% do (
    echo generate cli protocol go code: %%i.proto
    protoc --go_out=%outPath% ./%%i.proto
)

pause