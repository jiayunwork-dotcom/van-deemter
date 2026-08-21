# van-deemter

van-deemter 是色谱柱塔板高度与 van Deemter 曲线核算的 Go 命令行内核：输入涡流扩散 A、纵向扩散 B、传质阻力 C 与流动相线速 u（A、B、C、u 均须为正，u≤0 或任一系数非正一律报错），按 H(u)=A+B/u+C·u 输出塔板高度 H；给定柱长 L>0 时按 N=L/H(u) 输出理论板数。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run . eval example/packed.json --u 0.02     # H≈3.44e-05 m，N≈7267，满足与否见输出
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
