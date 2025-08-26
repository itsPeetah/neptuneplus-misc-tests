import requests


CALL_COUNT = 20


def main():
    with open("profile-list.txt", "r") as f_in:
        text = [line.strip() for line in f_in.readlines()]
        services = zip(
            [text[i] for i in range(0, len(text), 2)],
            [text[i] for i in range(1, len(text), 2)],
        )
        for service_name, service_uri in services:

            with open(f"results-{service_name}.txt", "w") as f_out:
                pass


if __name__ == "__main__":
    main()
