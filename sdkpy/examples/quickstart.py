"""Quickstart example demonstrating i6sdk usage with zero external dependencies."""

import i6sdk


def main() -> None:
    # 1. Initialize SDK client (reads I6_SDK_DECSUITE_URL from environment by default)
    client = i6sdk.create_client(
        serv_urls={"decsuite": "https://decsuite.sandbox.rs.infinity6.ai"},
    )

    print("Client initialized with serv_urls:", client.config().serv_urls)

    # 2. Call an API endpoint (e.g. list datasets)
    resp = client.apis().dsz_list_datasets()
    if resp.ok:
        data = resp.parse()
        print("Successfully queried datasets:", data)
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
