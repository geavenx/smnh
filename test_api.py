#!/bin/env python3

import os
from sys import argv

import requests
from dotenv import load_dotenv

if __name__ == "__main__":
    script_name, *argv = argv
    usage = f"Usage: {script_name} <method> <username> <limit> <period>"

    if argv == []:
        print("Error: No arguments provided")
        print(usage)
        exit(1)

    elif argv == ["help"] or argv == ["h"]:
        print(usage)
        exit(0)

    elif len(argv) < 4:
        print("Error: Not enough arguments provided")
        print(usage)
        exit(1)

    elif len(argv) > 4:
        print("Error: Too much arguments provided")
        print(usage)
        exit(1)

    load_dotenv()
    last_fm_api_key = os.getenv("LAST_FM_API_KEY")
    last_fm_base_url = os.getenv("LAST_FM_BASE_URL")

    params = {
        "method": argv[0],
        "user": argv[1],
        "api_key": last_fm_api_key,
        "limit": argv[2],
        "period": argv[3],
    }

    qs = "&".join([f"{k}={v}" for k, v in params.items()])
    final_url = f"{last_fm_base_url}?{qs}"

    resp = requests.get(final_url)
    response_filename = "response.xml"
    with open(response_filename, "w") as f:
        f.write(resp.text)

        print(f"Wrote response to {response_filename}")
