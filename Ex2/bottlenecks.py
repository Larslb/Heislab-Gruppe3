from threading import Thread, Lock
i = 0
mtx = Lock()

def threadFunc1():	

	global i
	for j in range(0,1000001):
		mtx.acquire()
		i += 1
		mtx.release()

def threadFunc2():

	global i
	
	for j in range(0,1000000):
		mtx.acquire()
		i -= 1
		mtx.release()

def main():
	global i
	thread1 = Thread(target = threadFunc1(), args = (),)
	thread2 = Thread(target = threadFunc2(), args = (),)
	
	thread1.start()
	thread2.start()
	
	thread1.join()
	thread2.join()
	print(i)
	

main()
