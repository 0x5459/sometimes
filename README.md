# sometimes

A simple script language


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
