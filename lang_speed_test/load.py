import time

def load(N):
	x = []
	i = 0
	s = time.time()
	while i<N:
		x.append("document_%07d" %i)
		i+=1
	e = time.time()

	print e-s


load(10000000)


	
		
	
