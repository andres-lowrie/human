#!/usr/bin/env python3
import json
import os
import subprocess
import sys
import uuid


def log(msg, color="white"):
    prefixes = {
        "black": "\33[30m",
        "red": "\33[31m",
        "green": "\33[32m",
        "yellow": "\33[33m",
        "blue": "\33[34m",
        "magenta": "\33[35m",
        "cyan": "\33[36m",
        "white": "\33[37m",
    }
    suffix = "\33[0m"
    print(prefixes[color] + msg + suffix)


def good(msg):
    log(msg, "green")


def bad(msg):
    log(msg, "red")


tokens = {
    "%%human%%": "path-to-bin",
}


def replace_tokens(s):
    return "".join([s.replace(t, tokens[t]) for t in tokens])


def parse_spec(fp):
    out = subprocess.run(["yj", fp], capture_output=True, check=True)
    return json.loads(out.stdout)


def run(cmd):
    # @TODO make own return type in case I want to change the implementation details
    proc = subprocess.run([cmd], shell=True, capture_output=True)
    return proc


def create_script_file(cmd, tmpdi):
    tmpfile = f"{tmpdir}/{uuid.uuid4()}"
    fh = open(tmpfile, "w")
    cmd = "#!/usr/bin/env bash\n\n" + cmd
    fh.write(cmd)
    fh.close()
    os.chmod(tmpfile, 0o700)
    return tmpfile


def rm_script(script):
    os.remove(script)


def build(src_path, output_path):
    subprocess.check_output(["go", "build", "-o", output_path], cwd=src_path)


if __name__ == "__main__":
    # each file should be a yaml file
    files = sys.argv[1:]
    tmpdir = os.environ["E2E_TMP_DIR"]
    path_to_bin = f"{tmpdir}/human"
    path_to_src = "/".join(tmpdir.split("/")[0:-2])

    build(path_to_src, path_to_bin)
    tokens["%%human%%"] = path_to_bin

    for f in files:
        suite = parse_spec(f)["suite"]

        for case in suite["cases"]:
            if "setup" in case:
                print("handle setup")

            if "test" in case:
                test = case["test"]
                cmd = replace_tokens(test)
                script = create_script_file(cmd, tmpdir)
                out = run(script)

                # @TODO add option to only log failures?
                if out.returncode != 0:
                    bad(f'Failed: {case["name"]}')
                    print("\tStdOut", out.stdout)
                    print("\tStdErr", out.stderr)
                else:
                    good(f'{case["name"]}')

            if "cleanup" in case:
                print("handle cleanup")

            rm_script(script)
