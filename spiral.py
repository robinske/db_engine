def spiral(n, i):
	if n == i:
		return
	x = i
	y = i

	while x < n:
		print(x, i)
		x += 1
	while y <= n:
		print(n, y)
		y += 1
	x = n - 1
	y = n - 1
	while x > i:
		print (x, n)
		x = x - 1
	while y >= i:
		print (i, y)
		y = y - 1
	spiral(n-1, i+1)

spiral(5, 1)