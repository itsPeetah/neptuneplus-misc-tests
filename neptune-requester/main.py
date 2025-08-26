import time
import requests
import threading
import random


def run(alias: str, count: int, endpoint: str, range_min: int, range_max: int):
    for i in range(count):
        res = requests.get(endpoint + str(random.randint(range_min, range_max + 1)))
        print(f"[{alias}] R{i+1} : {res.text}")
        time.sleep(1.5)


def multirun(threads: int, count: int, endpoint: str, range_min: int, range_max: int):
    ts: list[threading.Thread] = []
    for i in range(threads):
        t = threading.Thread(
            target=run,
            args=[
                f"{endpoint.strip('/').split('/')[-2]}_{i}",
                count,
                endpoint,
                range_min,
                range_max,
            ],
        )
        t.start()
    for t in ts:
        t.join()


t1 = threading.Thread(
    target=multirun,
    args=[
        4,
        1000,
        "http://localhost:8080/function/openfaas-fn/prime-numbers/prime/",
        100000,
        100000,
    ],
)

t2 = threading.Thread(
    target=multirun,
    args=[
        4,
        1000,
        "http://localhost:8080/function/openfaas-fn/crime-numbers/prime/",
        100000,
        100000,
    ],
)


t3 = threading.Thread(
    target=multirun,
    args=[
        4,
        100,
        "http://localhost:8080/function/openfaas-fn/grime-numbers/prime/",
        100000,
        100000,
    ],
)

t1.start()
t2.start()
t3.start()

t1.join()
t2.join()
t3.join()
