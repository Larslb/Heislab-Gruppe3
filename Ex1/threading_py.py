from threading import Thread
i = 0

def threadFunc1():
	global i
	for j in range(0,1000000):
		i += 1
		

def threadFunc2():
	global i
	for j in range(0,1000000):
		i -= 1

def main():
	global i
	thread1 = Thread(target = threadFunc1(), args = (),)
	thread2 = Thread(target = threadFunc2(), args = (),)
	
	thread1.start()
	thread2.start()
	
	thread1.join()
	thread2.join()
	

main()
