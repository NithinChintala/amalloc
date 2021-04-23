import sys
import time

CURS_UP = "\033[F"
CR = "\r"

CLEAR = "\033[2J"
ORIG = "\033[0;0H"
CLEAR_ORIG = CLEAR + ORIG
val = 16

def data(level):
    byte = "0" * 8
    empty = " " * 8
    if level not in [16, 8, 4, 2]:
        sys.exit(1)
    
    out = []
    tmp = []
    while level != 1:
        tmp = [empty] * (16 - level)
        tmp.extend([byte] * level)
        out.append(f"{level:-2} : {' '.join(tmp)}")

        level = level // 2

    # return " ".join(["0" * 8] * level)
    return "\n".join(out)

def ask():
    global val
    sys.stdout.write(f"{CLEAR_ORIG}{data(int(val))}\n\n>>> ")
    val = input()
    if val == "q":
        sys.exit(0)

while True:
    ask()
