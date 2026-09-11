# Go 论文查重

本项目实现一个只依赖 Go 标准库的命令行论文文本相似度计算程序。

## 使用

```text
go build -o main.exe .
main.exe <原文文件绝对路径> <抄袭版文件绝对路径> <答案文件绝对路径>
```

程序读取前两个文件并把结果写入第三个文件。答案是 0 到 1 之间、保留两位小数的浮点数，例如 `0.67`。输入文件必须是有效 UTF-8；文件开头的 UTF-8 BOM 会被忽略。输出路径不能与任一输入文件相同。

## 算法

程序保留 Unicode 字母和数字，删除空白及标点，并将字母转换为小写。对于长度至少为两个字符的文本，统计相邻字符二元组的多重集频次，使用 Dice 相似度：

```text
2 * sum(min(originalCount, candidateCount)) /
    (originalBigramCount + candidateBigramCount)
```

当任一文本不足两个字符时，改用字符频次的 Dice 相似度。两篇规范化文本都为空时视为完全相同，得分为 `1.00`；只有一篇为空时得分为 `0.00`。

算法时间复杂度为 O(n)，空间复杂度为 O(n)，其中 n 是两篇文本的字符总数。该方法衡量字面相似度，不能识别同义词或语义改写。

## 验证

```text
go test ./...
go test ./... -cover
go test -run '^$' -bench . -benchmem
go vet ./...
```
