# go-h5-os-app

go和h5代码，组建系统级APP。基于Lorca。用 Loca 创建窗口

# 实现思路

## 用 Lorca 创建窗口

我了解到 Go 的如下库可以实现窗口：

1. [lorca](https://github.com/zserge/lorca) - 调用系统现有的 Chrome/Edge 实现简单的窗口，UI 通过 HTML/CSS/JS 实现
2. [webview](https://github.com/webview/webview) - 比 lorca 功能更强，实现 UI 的思路差不多
3. [fyne](https://github.com/fyne-io/fyne) - 使用 Canvas 绘制的 UI 框架，性能不错
4. [qt](https://github.com/therecipe/qt) - 更复杂更强大的 UI 框架
