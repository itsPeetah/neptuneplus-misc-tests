# %%
import requests
import time
import threading
import os
import io

# %%
URL_TO_HIT = "http://localhost:8080/function/openfaas-fn/prime-numbers/prime/50000"

# %% [markdown]
# # Requests


# %%
def make_ramp(
    start_count: int,
    end_count: int,
    steps: int,
    start_padding: int = 0,
    end_padding: int = 0,
) -> list[int]:
    delta = end_count - start_count
    return (
        [start_count for _ in range(start_padding)]
        + [start_count + int((i / steps) * delta) for i in range(steps + 1)]
        + [end_count for _ in range(end_padding)]
    )


# %%
def do_request():
    r = requests.get(URL_TO_HIT)


def run_requests(ramp: list[int]):
    for reqs in ramp:
        for i in range(reqs):
            threading.Thread(target=do_request).start()
        time.sleep(1)


# %% [markdown]
# # Analysis


# %%
def get_pods(output_file: io.TextIOWrapper, t: int):
    txt = os.popen("kubectl get pods -n openfaas-fn").read().strip().split("\n")[1:]
    for line in txt:
        output_file.write(str(t) + " " + line.strip().replace("pod/", "") + "\n")


def monitor_pods(file_name: str, duration: int):
    with open(file_name, "w") as f:
        for x in range(duration):
            get_pods(f, x + 1)
            time.sleep(1)


# %% [markdown]
# # Main

# %%
t_anal = threading.Thread(target=monitor_pods, args=("pod_names.txt", 500))
t_reqs = threading.Thread(
    target=run_requests, kwargs={"ramp": make_ramp(20, 20, 1000, 10000, 200)}
)

print("Starting requests")
t_reqs.start()
print("Starting analysis")
t_anal.start()
print("Started")
t_reqs.join()
print("Joined request thread")
t_anal.join()
print("Joined analysis thread")
print("Done.")


# %%
