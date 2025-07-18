import requests
import threading


class Caller:
    def __init__(self, url, count, alias):
        self.url = url
        self.count = count
        self.alias = alias

    def _call(self):
        response = requests.get(self.url)
        return response

    def _call_loop(self):
        print(f'[Caller] Starting caller "{self.alias}"')
        for i in range(self.count):
            print(f"[{self.alias} #{i+1}] Calling {self.url}")
            _ = self._call()
        print(f"[Caller {self.alias}] Done.")

    def do_loop(self):
        t = threading.Thread(target=self._call_loop)
        t.start()
        return t


class Spawner:
    def __init__(self, callers):
        self.callers = callers

    def run(self):
        threads = [(c.alias, c.do_loop()) for c in self.callers]
        for a, t in threads:
            t.join()
            print(f"[Spawner] Joined thread for caller {a}")


caller_a1 = Caller("http://localhost:8080/function/openfaas-fn/function-a", 100, "A1")
caller_a2 = Caller("http://localhost:8080/function/openfaas-fn/function-a", 100, "A2")
caller_b1 = Caller("http://localhost:8080/function/openfaas-fn/function-b", 100, "B1")
caller_b2 = Caller("http://localhost:8080/function/openfaas-fn/function-b", 100, "B2")
caller_c1 = Caller("http://localhost:8080/function/openfaas-fn/function-c", 100, "C1")
caller_c2 = Caller("http://localhost:8080/function/openfaas-fn/function-c", 100, "C2")

spawner = Spawner([caller_a1, caller_a2, caller_b1, caller_b2, caller_c1, caller_c2])
spawner.run()
