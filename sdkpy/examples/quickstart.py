"""Quickstart example demonstrating i6sdk usage with zero external dependencies."""

import pathlib
import sys

# Ensure local i6sdk package is importable when running example directly
sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent.parent / "src"))

import i6sdk


def main() -> None:
    # 1. Initialize SDK client (reads I6_SDK_DECSUITE_URL from environment by default)
    # Persisting credentials to ./gen/sdkpy.tokens.json
    client = i6sdk.create_client(
        serv_urls={"decsuite": "https://decsuite.sandbox.rs.infinity6.ai"},
        storage="./gen/sdkpy.tokens.json",
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

        client.auth().device_code_poll(
            device_code=flow["device_code"],
            interval=float(flow.get("interval_sec", 5)),
            timeout=float(flow.get("expires_in_sec", 300)),
        )

        me = client.auth().me()
        print("Successfully authenticated as:", me)

    if not me:
        raise Exception('no user')

    user_id = me.get("user_id")
    if not user_id:
        raise Exception("no user id")

    resp = client.apis().dsz_list_datasets(req={"params": {"userId": str(user_id)}})
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


    # 3. Use internal bitmask utilities
    from i6sdk.internal.bitmapz import contains_all, join

    # Combine permissions: Read (8) + Write (4)
    my_grants = join(4, 8)
    print(f"Combined bitmask: {my_grants}")
    print(f"Has write permission: {contains_all(my_grants, 4)}")


if __name__ == "__main__":
    main()
