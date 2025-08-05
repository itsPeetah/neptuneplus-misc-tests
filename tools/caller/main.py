import requests
import threading
import time


class Saver:
    def __init__(self, filename):
        self.filename = filename
        self.lines = []

    def record(self, alias, count, time):
        self.lines.append(f"{alias};{count};{time}\n")

    def save(self):
        with open(self.filename, "w") as f:
            f.write("alias;count;time\n")
            f.writelines(self.lines)


class Caller:
    def __init__(self, url, count, alias, saver):
        self.url = url
        self.count = count
        self.alias = alias
        self.saver = saver

    def _call(self):
        response = requests.get(self.url)
        return response

    def _call_loop(self):
        print(f'[Caller] Starting caller "{self.alias}"')
        for i in range(self.count):
            time.sleep(1)
            print(f"[{self.alias} #{i+1}] Calling {self.url}")
            start = time.time()
            _ = self._call()
            end = time.time()
            self.saver.record(self.alias, i + 1, end - start)
        print(f"[Caller {self.alias}] Done.")

    def do_loop(self):
        t = threading.Thread(target=self._call_loop)
        t.start()
        return t


class Spawner:
    def __init__(self, callers):
        self.callers = callers

    def run(self):
        print("[Spawner] Starting...")
        threads = [(c.alias, c.do_loop()) for c in self.callers]
        print("[Spawner] Started.")
        for a, t in threads:
            t.join()
            print(f"[Spawner] Joined thread for caller {a}")
        print("[Spawner] Finished.")


saver = Saver("times.csv")

caller_a1 = Caller(
    "http://localhost:8080/function/openfaas-fn/function-a/handle", 100, "A1", saver
)
caller_a2 = Caller(
    "http://localhost:8080/function/openfaas-fn/function-a/handle", 100, "A2", saver
)
caller_b1 = Caller(
    "http://localhost:8080/function/openfaas-fn/function-b/handle", 100, "B1", saver
)
caller_b2 = Caller(
    "http://localhost:8080/function/openfaas-fn/function-b/handle", 100, "B2", saver
)
caller_c1 = Caller(
    "http://localhost:8080/function/openfaas-fn/function-c/handle", 100, "C1", saver
)
caller_c2 = Caller(
    "http://localhost:8080/function/openfaas-fn/function-c/handle", 100, "C2", saver
)

spawner = Spawner([caller_a1, caller_a2, caller_b1, caller_b2, caller_c1, caller_c2])
spawner.run()

saver.save()
