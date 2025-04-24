# EPS_Clean
Clean the data from EPS

> Given that EPS is a Chinese platform, and most of the user is chinese speaker, so this document is written by chinese.

# EPS_Clean
一个基于Python实现的命令行工具，来帮你把从EPS下载下来的数据转换成utf-8格式的CSV数据。

**以下内容由Claude撰写**

## 简介

EPS_Clean 主要功能是处理从EPS（经济数据聚合平台）下载的CSV文件，解决编码问题，并提取有用数据。工具可以：

- 读取GB2312编码的CSV文件
- 删除最后三行数据（一些版权行）
- 将结果以UTF-8编码保存，便于跨平台使用或导入Stata等

## 安装方法

### 下载预编译版本（推荐）

您可以直接从GitHub Releases页面下载已经打包好的可执行文件：

1. 访问 [GitHub Releases页面](https://github.com/sepinetam/eps_clean/releases)
2. 下载适合您操作系统的版本
3. 将文件放置到以下位置可全局使用：
   - **macOS**: 将`epsclean`复制到`/usr/local/bin/`并添加执行权限
     ```bash
     sudo cp epsclean /usr/local/bin/
     sudo chmod +x /usr/local/bin/epsclean
     ```
   - **Windows**: 将`epsclean.exe`复制到`C:\Windows\System32\`或其他PATH目录

> **注意**：目前仅提供macOS版本的预编译文件，Windows版本待社区贡献。

### 通过pip安装

```bash
pip install eps-clean
```

### 手动安装

```bash
git clone https://github.com/sepinetam/eps_clean.git
cd eps_clean
pip install -e .
```

### 更新到最新版本

```bash
pip install --upgrade eps-clean
```

### 卸载

```bash
pip uninstall eps-clean
```

### 虚拟环境注意事项

如果您在虚拟环境中安装了 EPS_Clean，请注意以下几点：

1. 工具只在该虚拟环境激活时可用
2. 要在任何环境下都能使用，可以考虑以下方法：
   - 在系统全局环境中安装（不推荐）
   - 创建一个简单的脚本，激活虚拟环境并调用该工具
   - 使用 PyInstaller 创建独立可执行文件（参见贡献指南）

## 使用方法

### 基本用法

```bash
epsclean filename.csv
```

这将处理`filename.csv`文件，删除最后三行，并以UTF-8编码覆盖原文件。

### 保存到新文件

```bash
epsclean filename.csv newfilename.csv
```

这将处理`filename.csv`文件，删除最后三行，并以UTF-8编码保存到`newfilename.csv`。

### 指定输入文件编码

```bash
epsclean --encoding gb2312 filename.csv
```

默认情况下，工具假设输入文件使用GB2312编码（EPS导出的标准编码）。如果您的文件使用其他编码，可以通过`--encoding`参数指定。

### 帮助信息

```bash
epsclean --help
```

显示完整的使用说明和可用选项。

## 示例

```bash
# 处理单个文件，覆盖原文件
epsclean filename.csv

# 处理文件并保存到新文件
epsclean filename.csv filename_cleaned.csv

# 处理使用UTF-8编码的文件
epsclean --encoding utf8 filename.csv filename_cleaned.csv

# 处理使用gbk编码的文件
epsclean --encoding gbk filename.csv filename_cleaned.csv
```

## 常见问题

**Q: 为什么我的文件处理后出现乱码？**  
A: 可能是输入文件的编码不是GB2312。尝试使用`--encoding`参数指定正确的编码。

**Q: 我需要删除不同数量的行怎么办？**  
A: 目前版本固定删除最后三行。如有特殊需求，欢迎提交Issue或PR。

## 贡献指南

欢迎对EPS_Clean做出贡献！您可以通过以下方式参与：

1. 提交Issue：报告bug或提出新功能建议
2. 提交Pull Request：改进代码或文档

### 创建独立可执行文件

如果您想要创建不依赖Python环境的独立可执行文件，可以使用PyInstaller：

```bash
# 安装PyInstaller
pip install pyinstaller

# 打包成独立可执行文件
pyinstaller --onefile eps_clean/main.py --name epsclean

# 打包后的文件位于dist目录中
```

### Windows版本贡献

目前该项目由macOS用户维护，我们非常欢迎Windows用户的贡献：

1. 如果您使用Windows系统并愿意提供打包好的Windows可执行文件，请创建PR
2. 您也可以帮助测试在Windows环境下的兼容性问题
3. 欢迎提供Windows特定的安装和使用说明

### 平台兼容性

本工具应该能在所有主流平台（Windows、macOS、Linux）上运行，但我们主要在macOS上进行测试。如果您在其他平台上遇到问题，请提交Issue。

项目地址：https://github.com/sepinetam/eps_clean

## 许可证

本项目采用[GNU AFFERO GENERAL PUBLIC LICENSE许可证](LICENSE)。详情请参阅LICENSE文件。
