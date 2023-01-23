FROM stellarisw/alpine as build

# 为我们的镜像设置必要的环境变量
ENV TZ=Asia/Shanghai

WORKDIR /server

# 将代码复制到容器中
COPY ./douyin-chat-rpc-sys /server/douyin-chat-rpc-sys

# 设置时区
RUN chmod 777 ./douyin-chat-rpc-sys

# 启动容器时运行的命令
CMD ["./douyin-chat-rpc-sys"]
