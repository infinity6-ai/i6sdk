"""Quickstart example demonstrating i6sdkpy usage with zero external dependencies."""

import pathlib
import sys

# Ensure local i6sdkpy package is importable when running example directly
sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent.parent / "src"))

import i6sdkpy


def main() -> None:
    # 1. Initialize SDK client (reads I6_SDK_DECSUITE_URL from environment by default)
    # Persisting credentials to ./gen/sdkpy.tokens.db and caching SQLite bundles to ./gen/cache
    client = i6sdkpy.create_client(
        serv_urls={"decsuite": "https://decsuite.sandbox.rs.infinity6.ai"},
        storage="./gen/sdkpy.tokens.db",
        cache_dir="./gen/cache",
    )

    print("Client initialized with serv_urls:", client.config().serv_urls)

    me = client.auth().me()
    print("me:", me)

    if not me:
        print("\nNo active user session. Initiating OAuth device authorization...")
        flow = client.auth().device_code_start()
        verify_url = flow.get("verification_uri_complete") or flow.get("verification_uri")
        user_code = flow.get("user_code")

        print(f"Please visit: {verify_url}")
        if user_code:
            print(f"Confirmation code: {user_code}")
        print("Waiting for authorization in browser...")

        interval = float(flow.get("interval_sec") or 5)
        timeout = float(flow.get("expires_in_sec") or 300)
        client.auth().device_code_poll(
            device_code=flow.get("device_code", ""),
            interval=interval,
            timeout=timeout,
        )

        me = client.auth().me()
        print("Successfully authenticated as:", me)

    if not me:
        raise Exception('no user')

    user_id = me.get("user_id")
    if not user_id:
        raise Exception("no user id")

    resp = client.apis().dsz_list_datasets(req={"params": {"userId": user_id}})
    if resp.ok:
        data = resp.parse()
        print("Successfully queried datasets:", data)
        datasets = data.get("datasets")
        print("datasets", datasets)
        if datasets:
            print("datasets[0].id", datasets[0].get("id"))
            print("datasets[0].name", datasets[0].get("grants"))
    else:
        print("Request failed:", resp.error)


    # client.infer().run({
    #     "dataset": "dectest",
    #     "model_name": "relevance-fashion-fbt",
    #     "partitions": {"store_id": "1"},
    #     "params": {}
    # })

    result = client.funcs().product_relevance_fashion_dash_display_metrics_status_v1({
        "partitions": {"store_id": "1"},
        "input": {
            "dataset": "dectest",
            "params": {
                "min_year_month": "2026-08",
                "max_year_month": "2027-12",                   
            }
        },
    })
    print("product_relevance_fashion_dash_display_metrics_status_v1", result)

    # 3. Use bitmask utilities and Grants enum
    from i6sdkpy import Grants, bitmapz

    # Combine permissions: Read + Write
    my_grants = bitmapz.join(Grants.READ, Grants.WRITE)
    print(f"Combined bitmask: {my_grants}")
    print(f"Has write permission: {bitmapz.contains_all(my_grants, Grants.WRITE)}")


if __name__ == "__main__":
    main()
