import requests
import time
import threading

f = open("output.txt", "+w")


def req():
    global f
    r = requests.get(
        "http://localhost:8080/function/openfaas-fn/prime-numbers/prime/50000"
    )
    f.write(r.text + "\n")


for _ in range(1000):
    print(time.time())

    d = 0.1
    [
        (t.start(), time.sleep(d))
        for t in [threading.Thread(target=req, daemon=True) for __ in range(10)]
    ],
