import requests
import time
import threading
import os
import io


SEQ_URL = "http://localhost:8080/function/openfaas-fn/prime-numbers-caller-sequential/entrypoint?mode=seq"
NODEP_URL = "http://localhost:8080/function/openfaas-fn/prime-numbers-caller-sequential-nodep/entrypoint?mode=seq"


def do_request(url: str, file_name: str, name: str):
    with open(file_name, "w") as file:
        r = requests.get(url)
        print(f"{name}: {r.text}")
        file.write(os.popen("kubectl get podscales -n openfaas-fn -o json").read())


def run_requests(url: str, count: int, output_dir: str):
    if not os.path.isdir(output_dir):
        os.mkdir(output_dir)
    for i in range(count):
        threading.Thread(
            target=do_request,
            args=(url, f"{output_dir}/{i}.json", f"{output_dir} #{i}"),
        ).start()
        time.sleep(1)


threads = [
    t
    for i, t in enumerate(
        [
            threading.Thread(target=run_requests, args=(SEQ_URL, 900, "./out_dep")),
            threading.Thread(target=run_requests, args=(NODEP_URL, 900, "./out_nodep")),
        ]
    )
    if t.start() or print(f"Started thread #{i+1}") or True
]

for t in threads:
    t.join()
print("done")
