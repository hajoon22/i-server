c: main.o utils.o icmp.o
	gcc -o c main.o utils.o icmp.o
main.o: client/main.c
	gcc -c client/main.c -o main.o
utils.o: client/utils/utils.c client/utils/utils.h
	gcc -c client/utils/utils.c -o utils.o
icmp.o: client/icmp/icmp.c client/icmp/icmp.h
	gcc -c client/icmp/icmp.c -o icmp.o

serv:
	cd server;go build -o ../serv .
