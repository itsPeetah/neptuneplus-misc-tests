This is just a simple app made of serverless functions to test the work I'm doing in k8s for my thesis

App structure:

Entrypoint A

Entrypoint B

Entrypoint C

Function D

Waiter 2 seconds

Random Waiter (0-5seconds)

A —--> 2 seconds (1 call) , Random Waiter 1 call sequential
B —--> Random waiter (1 call) + 2 seconds (1 call) parallel , 2 seconds 1 call sequential
C —--> (D —---> Random waiter 2 call parallel), 2 seconds sequential

The code is pretty shit but it's not like it's going somewhere lol
