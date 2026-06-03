---
name: karpathy-coding-guidelines
description: Karpathy-inspired coding guidelines for AI assistants. Use when writing, reviewing, or modifying code to ensure simplicity, precision, and goal-driven execution.
---

# Karpathy-Inspired Coding Guidelines

> 源自 [Andrej Karpathy 的观察](https://x.com/karpathy/status/2015883857489522876) 和 [forrestchang/andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills)
> 用途：做新项目时将此文件拷贝到项目根目录，指导 AI 助手减少常见编码错误。

**权衡：** 这些指南倾向于谨慎而非速度。对于琐碎任务，自行判断。

---

## 1. Think Before Coding（编码前思考）

**不要假设。不要隐藏困惑。呈现权衡。**

实现之前：
- **明确说明假设** — 如果不确定，询问而不是猜测
- **呈现多种解释** — 当存在歧义时，不要默默选择
- **适时提出异议** — 如果存在更简单的方法，说出来
- **困惑时停下来** — 指出不清楚的地方并要求澄清

---

## 2. Simplicity First（简洁优先）

**用最少的代码解决问题。不要过度推测。**

- 不要添加要求之外的功能
- 不要为一次性代码创建抽象
- 不要添加未要求的"灵活性"或"可配置性"
- 不要为不可能发生的场景做错误处理
- 如果 200 行代码可以写成 50 行，重写它

**检验标准：** 资深工程师会觉得这过于复杂吗？如果是，简化。

---

## 3. Surgical Changes（精准修改）

**只碰必须碰的。只清理自己造成的混乱。**

编辑现有代码时：
- 不要"改进"相邻的代码、注释或格式
- 不要重构没坏的东西
- 匹配现有风格，即使你更倾向于不同的写法
- 如果注意到无关的死代码，提一下 —— 不要删除它

当你的改动产生孤儿代码时：
- 删除因你的改动而变得无用的导入/变量/函数
- 不要删除预先存在的死代码，除非被要求

**检验标准：** 每一行修改都应该能直接追溯到用户的请求。

---

## 4. Goal-Driven Execution（目标驱动执行）

**定义成功标准。循环验证直到达成。**

将指令式任务转化为可验证的目标：

| 不要这样做... | 转化为... |
|-------------|---------|
| "添加验证" | "为无效输入编写测试，然后让它们通过" |
| "修复 bug" | "编写重现 bug 的测试，然后让它通过" |
| "重构 X" | "确保重构前后测试都能通过" |

对于多步骤任务，说明一个简短的计划：

```
1. [步骤] → 验证: [检查]
2. [步骤] → 验证: [检查]
3. [步骤] → 验证: [检查]
```

---

## 核心洞察

> "LLM 非常擅长循环执行直到达成特定目标……不要告诉它该做什么，给它成功标准，然后看着它完成。"
> — Andrej Karpathy
