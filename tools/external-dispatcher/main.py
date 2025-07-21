from flask import Flask, request, jsonify, abort, Response
import requests
import time
import re

app = Flask(__name__)

# Base port for OpenFaaS functions (often 8080)
OPENFAAS_FUNCTION_PORT = 8080
LOCAL_PORT = 8090


@app.route(
    "/function/<namespace>/<function_name>",
    methods=["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"],
)
@app.route(
    "/function/<namespace>/<function_name>/<path:sub_path>",
    methods=["GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"],
)
def proxy_openfaas_function(namespace, function_name, sub_path=""):
    """
    Proxies requests to specific OpenFaaS functions within the cluster.
    Expected incoming path format: /function/<namespace>/<function_name>/<potential_sub_path>
    """
    incoming_path = request.path  # e.g., /function/openfaas-fn/my-func/foo/bar
    print(f"[{time.time()}] Dispatcher received request for path: {incoming_path}")
    start_time = time.perf_counter()

    # Construct the internal cluster service hostname
    # e.g., "function-b.openfaas-fn.svc.cluster.local"
    # internal_service_hostname = f"{function_name}.{namespace}.svc.cluster.local"
    target_url = f"http://localhost:{OPENFAAS_FUNCTION_PORT}/function/{namespace}/{function_name}"

    # Construct the full target URL within the cluster
    # OpenFaaS functions typically serve their main route at '/' within their pod
    # So, we append the 'sub_path' directly if it exists.
    if sub_path:
        target_url += f"/{sub_path}"

    # Add original query parameters to the target URL
    if request.query_string:
        target_url += f"?{request.query_string.decode('utf-8')}"

    print(f"[{time.time()}] Proxying to internal URL: {target_url}")

    try:
        # Reconstruct headers for the proxied request
        # Remove 'Host' header as it should be set by requests to the target's host
        # Remove 'Transfer-Encoding' to avoid issues, let requests handle it
        headers = {
            name: value
            for name, value in request.headers
            if name.lower() not in ["host", "transfer-encoding"]
        }

        # Use requests.request to handle any HTTP method
        resp = requests.request(
            method=request.method,
            url=target_url,
            headers=headers,
            data=request.get_data(),  # Raw request body
            allow_redirects=False,  # Control redirects explicitly if needed
            stream=True,  # Efficiently stream response body
        )

        end_time = time.perf_counter()
        response_time_ms = (end_time - start_time) * 1000
        print(
            f"[{time.time()}] Request to {target_url} completed in {response_time_ms:.2f}ms. Status: {resp.status_code}"
        )

        # *** HERE IS WHERE YOU INTEGRATE WITH NEPTUNE ***
        # You now have:
        # - `response_time_ms`: The total time for the round trip through the dispatcher.
        # - `namespace`: The OpenFaaS namespace (e.g., 'openfaas-fn').
        # - `function_name`: The name of the function (e.g., 'function-b').
        # - `incoming_path`: The full original path received by the dispatcher.
        # - `status_code`: The HTTP status code from the function.
        # - `request.method`: The HTTP method used.
        # Use this data to send metrics/logs to Neptune.

        # Prepare the response to send back to the original caller (e.g., Function A)
        response_headers = [
            (name, value)
            for name, value in resp.headers.items()
            if name.lower()
            not in ["content-encoding", "content-length", "transfer-encoding"]
        ]
        # Flask can handle content-encoding and content-length, or you might strip them if streaming.

        # Stream the content back
        def generate():
            for chunk in resp.iter_content(chunk_size=8192):
                yield chunk

        return Response(generate(), status=resp.status_code, headers=response_headers)

    except requests.exceptions.ConnectionError as e:
        print(f"[{time.time()}] Connection error proxying to {target_url}: {e}")
        return (
            jsonify({"error": f"Could not connect to target service: {e}"}),
            503,
        )  # Service Unavailable

    except requests.exceptions.RequestException as e:
        print(f"[{time.time()}] General error proxying request to {target_url}: {e}")
        return (
            jsonify({"error": f"Proxy request failed: {e}"}),
            500,
        )  # Internal Server Error


if __name__ == "__main__":
    print(f"Starting local OpenFaaS-specific dispatcher on port {LOCAL_PORT}")
    print(
        f"All requests to /function/<namespace>/<function_name>/<...> will be proxied."
    )
    # Listen on 0.0.0.0 to be accessible from within the k3d cluster via host.k3d.internal
    app.run(
        host="0.0.0.0", port=LOCAL_PORT, debug=True
    )  # debug=True is good for dev, remove for prod
