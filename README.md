# sometimes

A simple script language

## 模块结构
- [ast](https://github.com/0x5459/sometimes/tree/main/ast) 抽象语法树
- [lexer](https://github.com/0x5459/sometimes/tree/main/lexer) 词法解析
- [parser](https://github.com/0x5459/sometimes/tree/main/parser) 语法解析
- [hir](https://github.com/0x5459/sometimes/tree/main/hir) High level Intermediate Representation
- [vm](https://github.com/0x5459/sometimes/tree/main/vm) 字节码虚拟机

## Example
```
// 计算第 n 项斐波拉契数列小程序
fn main() {
	let a=1, b=1, t=0;
	let i=2, n=7;
	loop (i<n) {
		t = a;
		a = b;
		b = t + a;
		i += 1;
	};
	print(b);
}
```
