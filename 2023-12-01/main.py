def process_input(input_text):
    lines = input_text.strip().split('\n')
    result = []
    current_subarray = []

    for line in lines:
        if not line.strip():  # Check if the line is empty
            if current_subarray:
                result.append(current_subarray)
                current_subarray = []
        else:
            current_subarray.extend(map(int, line.strip().split()))

    if current_subarray:
        result.append(current_subarray)

    return result

def find_max_sum_index(arrays):
    max_sum = float('-inf')
    max_sum_index = -1

    for i, subarray in enumerate(arrays):
        current_sum = sum(subarray)
        if current_sum > max_sum:
            max_sum = current_sum
            max_sum_index = i

    return max_sum_index, max_sum

if __name__ == "__main__":
    input_text = """1000
    2000
    3000

    4000

    5000
    6000

    7000
    8000
    9000

    10000"""

    arrays = process_input(input_text)
    max_sum_index, max_sum = find_max_sum_index(arrays)

    print("List of Lists:", arrays)
    print("Index with the Highest Sum:", max_sum_index, " with Sum:", max_sum)

