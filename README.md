# van-deemter

van-deemter 是色谱柱塔板高度与 van Deemter 曲线核算的 Go 命令行内核：输入涡流扩散 A、纵向扩散 B、传质阻力 C 与流动相线速 u（A、B、C、u 均须为正，u≤0 或任一系数非正一律报错），按 H(u)=A+B/u+C·u 输出塔板高度 H；给定柱长 L>0 时按 N=L/H(u) 输出理论板数。`--scan` 直接给出最优线速 u_opt=√(B/C) 与最少塔板高 H_min=A+2√(BC)，并用 u 网格数值抽查与中心差分在 u_opt 处验证 dH/du≈0。给定需要的理论板数 N_req 时，反求最大允许塔板高 H_max=L/N_req，再判断当前 u 下的 H(u) 是否满足；同一柱长上 u 偏离最优则 N 下降，只增大 B 则 u_opt 上升，只增大 C 则 u_opt 下降且 H_min 上升。算例同时给出保留因子 k' 与分离因子 α 时，按简化分辨率 Rs=(√N/4)·((α−1)/α)·(k'/(1+k')) 输出 Rs；不给 k' 就不假装有 Rs。边界：A、B、C、L、u 必须为正，n_req 必须为正，k' 与 α 必须成对给出（只给一个报错），违反一律 stderr 明文报错并以非零退出码结束。本仓只做 van Deemter 塔板核算，不做进样序列排班，不做 McCabe–Thiele 精馏逐板。

- 输入：一个 JSON 算例文件。`example/packed.json` 字段：`a`（涡流扩散，m）、`b`（纵向扩散，m²/s）、`c`（传质阻力，s）、`length`（柱长，m），可选 `n_req`（需要达到的理论板数）、`kprime`（保留因子）、`alpha`（分离因子）。手算量级：A=2.4e-5 m、B=1.6e-8 m²/s、C=4.8e-4 s 时 u_opt=√(1.6e-8/4.8e-4)≈5.77e-3 m/s，H_min=2.4e-5+2√(7.68e-12)≈2.95e-5 m。
- 输出：`eval --u <m/s>` 打印 H、N（给了 n_req 时追加 H_max 与是否满足，给了 kprime/alpha 时追加 Rs）；`--scan` 打印 u_opt、H_min 及网格抽查与差分导数结果。
- 边界：仅处理单个 packed 柱算例；u_opt 与 H_min 用闭式公式，网格只做数值抽查不替代公式；`--scan` 与 `--u` 同时给出时以 `--scan` 为准；两个都不给时按用法错误退出。

## 钉死的约定

- **塔板高度**：H(u)=A+B/u+C·u，A、B、C、u 必须为正；最优处 dH/du=−B/u²+C=0，闭式 u_opt=√(B/C)、H_min=A+2√(BC)，H(u_opt)=H_min（`internal/hmodel`）。
- **最优一致性**：u_opt 两侧各取一点 H 都 ≥ H_min；中心差分 (H(u+δ)−H(u−δ))/(2δ) 在 u_opt 处以导数量级的 1e-6 相对容差收敛到零。
- **单调性**：只增大 B → u_opt 上升；只增大 C → u_opt 下降、H_min 上升；网格最小 H 与闭式 H_min 相对差 ≤1e-3（`internal/hmodel` 交叉测试）。
- **理论板数**：N=L/H(u)，L>0；给定 N_req 时 H_max=L/N_req，H(u)≤H_max 即满足；反解 H(u)=H_max 得满足要求的 u 区间（二次方程两根，`internal/column`）。
- **分辨率**：k' 与 α 同时给定时 Rs=(√N/4)·((α−1)/α)·(k'/(1+k'))，其中 N 取当前 u 下的理论板数；两者缺一即 error，不给则输出不含 Rs。

## 构建 / 运行 / 测试

```text
go build ./...                       # 编译（纯标准库）
go test ./...                        # 全部测试（hmodel / column / evalx / main）
go run . eval example/packed.json --u 0.02     # H≈3.44e-05 m，N≈7267，满足与否见输出
go run . eval example/packed.json --scan       # u_opt≈0.005774 m/s，H_min≈2.9543e-05 m
```

命令反例（应 stderr 报错并非零退出）：

```text
go run . eval example/packed.json --u 0        # 线速必须为正
go run . eval example/packed.json --u -0.1     # 线速必须为正
go run . eval example/packed.json --scan       # 先把 example/packed.json 的 c 改成 -4.8e-4，系数必须为正
go run . eval example/packed.json              # 缺 --u 且缺 --scan，用法错误
```
