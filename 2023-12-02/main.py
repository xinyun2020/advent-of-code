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

def find_top_3_indices(arrays):
    # Create a list of tuples (index, sum) and sort it in descending order based on the sum
    sorted_indices = sorted(enumerate(map(sum, arrays)), key=lambda x: x[1], reverse=True)

    # Take the top 3 indices
    top_3_indices = sorted_indices[:3]

    return top_3_indices

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
    top_3_indices = find_top_3_indices(arrays)

    print("List of Lists:", arrays)
    print("Top 3 Indices with the Highest Sum:", top_3_indices)
    print("Top 3 Indices total sum:", sum(item[1] for item in top_3_indices))

